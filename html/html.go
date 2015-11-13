// Copyright (c) 2015, Vastech SA (PTY) LTD. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package html

import (
	"strings"

	"github.com/gogo/letmegrpc/form"
	descriptor "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

const backtick = "`"

type html struct {
	*generator.Generator
	generator.PluginImports
	ioPkg      generator.Single
	reflectPkg generator.Single
	stringsPkg generator.Single
	jsonPkg    generator.Single
	strconvPkg generator.Single
}

func New() *html {
	return &html{}
}

func (p *html) Name() string {
	return "html"
}

func (p *html) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *html) typeName(name string) string {
	return p.TypeName(p.ObjectNamed(name))
}

const errString = `w.Write([]byte("<div class=\"alert alert-danger\" role=\"alert\">" + err.Error() + "</div>"))`

func (p *html) writeError(eof string) {
	p.P(`if err != nil {`)
	p.In()
	p.P(`if err != `, p.ioPkg.Use(), `.EOF {`)
	p.In()
	p.P(errString)
	p.P(`return`)
	p.Out()
	p.P(`}`)
	p.P(eof)
	p.Out()
	p.P(`}`)
}

func (p *html) getInputType(method *descriptor.MethodDescriptorProto) *descriptor.DescriptorProto {
	fileDescriptorSet := p.AllFiles()
	inputs := strings.Split(method.GetInputType(), ".")
	packageName := inputs[1]
	messageName := inputs[2]
	msg := fileDescriptorSet.GetMessage(packageName, messageName)
	if msg == nil {
		p.Fail("could not find message ", method.GetInputType())
	}
	return msg
}

func (p *html) generateFormFunc(servName string, method *descriptor.MethodDescriptorProto) {
	inputs := strings.Split(method.GetInputType(), ".")
	packageName := inputs[1]
	messageName := inputs[2]
	s := `<div class="container"><div class="jumbotron">
	<h3>` + servName + `: ` + method.GetName() + `</h3>
	` + form.Create(method.GetName(), packageName, messageName, p.Generator) + `
	</div>`
	p.P(`var Form`, servName, "_", method.GetName(), " string = `", s, "`")
}

func (p *html) CamelCaseFieldPath(fieldPath []string) string { // TODO: move to apropriate place (package generator ??)
	var path string
	for _, field := range fieldPath {
		if len(path) > 0 {
			path += "."
		}
		path += generator.CamelCase(field)
	}
	return path
}

func (p *html) generateOneofScanner(field *descriptor.FieldDescriptorProto, fieldPath []string) {
	if !field.IsMessage() {
		// skip
		return
	}

	// path to the oneof field within the root message (input message)
	fieldPath = append(fieldPath, field.GetName())

	fieldDesc := p.ObjectNamed(field.GetTypeName()).(*generator.Descriptor)

	var oneofDecls = make(map[int32]*descriptor.OneofDescriptorProto)
	var oneofFields = make(map[int32][]*descriptor.FieldDescriptorProto)

	for i, oneof := range fieldDesc.OneofDecl {
		oneofDecls[int32(i)] = oneof
	}

	for _, subfield := range fieldDesc.GetField() {
		if subfield.OneofIndex != nil {
			oneofFields[*subfield.OneofIndex] = append(oneofFields[*subfield.OneofIndex], subfield)
		}
	}

	for i, oneofDecl := range oneofDecls {
		p.P(`case "`, strings.Join(fieldPath, "."), `.`, oneofDecl.GetName(), `":`)
		p.In()
		p.P(`switch one.Selected {`)
		for _, oneofField := range oneofFields[i] { // fields.. ?
			p.P(`case "`, fieldDesc.GetName(), `_`, generator.CamelCase(oneofField.GetName()), `":`)
			p.In()
			p.P(`msg.`, p.CamelCaseFieldPath(fieldPath), ` = &`, fieldDesc.GetName(), `{`)
			p.In()
			p.P(generator.CamelCase(oneofDecl.GetName()), `: &`, fieldDesc.GetName(), `_`, generator.CamelCase(oneofField.GetName()), `{},`)
			p.Out()
			p.P(`}`)
			p.Out()
		}
		p.P(`if msg.`, p.CamelCaseFieldPath(fieldPath), ` != nil {`)
		p.In()
		p.P(`err := encoding_json.Unmarshal(one.Value, msg.`, p.CamelCaseFieldPath(fieldPath), `)`)
		p.writeError(errString)
		p.P(`}`)
		p.Out()
		p.P(`}`)
		p.Out()
	}

	// recursive for child fields
	for _, subfield := range fieldDesc.GetField() {
		p.generateOneofScanner(subfield, fieldPath)
	}
}

func (p *html) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	httpPkg := p.NewImport("net/http")
	p.jsonPkg = p.NewImport("encoding/json")
	p.ioPkg = p.NewImport("io")
	contextPkg := p.NewImport("golang.org/x/net/context")
	p.reflectPkg = p.NewImport("reflect")
	p.stringsPkg = p.NewImport("strings")
	p.strconvPkg = p.NewImport("strconv")
	logPkg := p.NewImport("log")
	grpcPkg := p.NewImport("google.golang.org/grpc")

	p.P(`var DefaultHtmlStringer = func(req, resp interface{}) ([]byte, error) {`)
	p.In()
	p.P(`header := []byte("<p><div class=\"container\"><pre>")`)
	p.P(`data, err := `, p.jsonPkg.Use(), `.MarshalIndent(resp, "", "\t")`)
	p.P(`if err != nil {`)
	p.In()
	p.P(`return nil, err`)
	p.Out()
	p.P(`}`)
	p.P(`footer := []byte("</pre></div></p>")`)
	p.P(`return append(append(header, data...), footer...), nil`)
	p.Out()
	p.P(`}`)

	p.P(`func Serve(httpAddr, grpcAddr string, stringer func(req, resp interface{}) ([]byte, error), opts ...`, grpcPkg.Use(), `.DialOption) {`)
	p.In()
	p.P(`conn, err := `, grpcPkg.Use(), `.Dial(grpcAddr, opts...)`)
	p.P(`if err != nil {`)
	p.In()
	p.P(logPkg.Use(), `.Fatalf("Dial(%q) = %v", grpcAddr, err)`)
	p.Out()
	p.P(`}`)
	for _, s := range file.GetService() {
		origServName := s.GetName()
		servName := generator.CamelCase(origServName)
		p.P(origServName, `Client := New`, servName, `Client(conn)`)
		p.P(origServName, `Server := NewHTML`, servName, `Server(`, origServName, `Client, stringer)`)
		for _, m := range s.GetMethod() {
			p.P(httpPkg.Use(), `.HandleFunc("/`, servName, `/`, m.GetName(), `", `, origServName, `Server.`, m.GetName(), `)`)
		}
	}
	p.P(`if err := `, httpPkg.Use(), `.ListenAndServe(httpAddr, nil); err != nil {`)
	p.In()
	p.P(logPkg.Use(), `.Fatal(err)`)
	p.Out()
	p.P(`}`)
	p.Out()
	p.P(`}`)

	p.P(`// oneofFields maps outer field name to selected field and it's value.`)
	p.P(`type oneofFields map[string]struct {`)
	p.In()
	p.P(`Selected string                   `, backtick, `json:"selected"`, backtick)
	p.P(`Value    encoding_json.RawMessage `, backtick, `json:"value"`, backtick)
	p.Out()
	p.P(`}`)

	for _, s := range file.GetService() {
		origServName := s.GetName()
		servName := generator.CamelCase(origServName)
		p.P(`type html`, servName, ` struct {`)
		p.In()
		p.P(`client `, servName, `Client`)
		p.P(`stringer func(req, resp interface{}) ([]byte, error)`)
		p.Out()
		p.P(`}`)

		p.P(`func NewHTML`, servName, `Server(client `, servName, `Client, stringer func(req, resp interface{}) ([]byte, error)) *html`, servName, ` {`)
		p.In()
		p.P(`return &html`, servName, `{client, stringer}`)
		p.Out()
		p.P(`}`)

		for _, m := range s.GetMethod() {
			p.generateFormFunc(servName, m)
			p.P(``)
			p.P(`func (this *html`, servName, `) `, m.GetName(), `(w `, httpPkg.Use(), `.ResponseWriter, req *`, httpPkg.Use(), `.Request) {`)
			p.In()
			p.P("w.Write([]byte(Header(`", servName, "`,`", m.GetName(), "`)))")
			p.P(`jsonString := req.FormValue("json")`)
			p.P(`oneofsString := req.FormValue("oneofs")`)
			p.P(`someValue := false`)
			p.RecordTypeUse(m.GetInputType())
			p.P(`msg := &`, p.typeName(m.GetInputType()), `{}`)
			p.P(`if len(jsonString) > 0 {`)
			p.In()
			p.P(`err := `, p.jsonPkg.Use(), `.Unmarshal([]byte(jsonString), msg)`)
			p.writeError(errString)
			p.P(`someValue = true`)

			p.P(`if len(oneofsString) > 0 {`)
			p.In()
			p.P(`oneofs := make(oneofFields)`)
			p.P(`err := encoding_json.Unmarshal([]byte(oneofsString), &oneofs)`)
			p.writeError(errString)
			p.P(`for field, one := range oneofs {`)
			p.In()
			p.P(`switch field {`)
			d := p.ObjectNamed(m.GetInputType()).(*generator.Descriptor)
			for _, f := range d.GetField() { //++ add this to generateOneofScanner ?
				p.generateOneofScanner(f, nil)
			}
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)

			p.P(`w.Write([]byte(Form`, servName, `_`, m.GetName(), `))`)
			p.P(`if someValue {`)
			p.In()
			if !m.GetClientStreaming() {
				if !m.GetServerStreaming() {
					p.P(`reply, err := this.client.`, m.GetName(), `(`, contextPkg.Use(), `.Background(), msg)`)
					p.writeError(errString)
					p.P(`out, err := this.stringer(msg, reply)`)
					p.writeError(errString)
					p.P(`w.Write(out)`)
				} else {
					p.P(`down, err := this.client.`, m.GetName(), `(`, contextPkg.Use(), `.Background(), msg)`)
					p.writeError(errString)
					p.P(`for {`)
					p.In()
					p.P(`reply, err := down.Recv()`)
					p.writeError(`break`)
					p.P(`out, err := this.stringer(msg, reply)`)
					p.writeError(errString)
					p.P(`w.Write(out)`)
					p.P(`w.(`, httpPkg.Use(), `.Flusher).Flush()`)
					p.Out()
					p.P(`}`)
				}
			} else {
				if !m.GetServerStreaming() {
					p.P(`up, err := this.client.Upstream(`, contextPkg.Use(), `.Background())`)
					p.writeError(errString)
					p.P(`err = up.Send(msg)`)
					p.writeError(errString)
					p.P(`reply, err := up.CloseAndRecv()`)
					p.writeError(errString)
					p.P(`out, err := this.stringer(msg, reply)`)
					p.writeError(errString)
					p.P(`w.Write(out)`)
				} else {
					p.P(`bidi, err := this.client.Bidi(`, contextPkg.Use(), `.Background())`)
					p.writeError(errString)
					p.P(`err = bidi.Send(msg)`)
					p.writeError(errString)
					p.P(`reply, err := bidi.Recv()`)
					p.writeError(errString)
					p.P(`out, err := this.stringer(msg, reply)`)
					p.writeError(errString)
					p.P(`w.Write(out)`)
				}
			}
			p.Out()
			p.P(`}`)
			p.P("w.Write([]byte(Footer))")
			p.Out()
			p.P(`}`)
		}
	}

	header1 := `
	<html>
	<head>
	<title>`

	header2 := `</title>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/css/bootstrap.min.css">
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.2/jquery.min.js"></script>
	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/js/bootstrap.min.js"></script>
	</head>
	<body>
	`
	footer := `
	</body>
	</html>
	`

	p.P("var Header func(servName, methodName string) string = func(servName, methodName string) string {")
	p.In()
	p.P("return `", header1, "` + servName + `:` + methodName + `", header2, "`")
	p.Out()
	p.P(`}`)

	p.P("var Footer string = `", footer, "`")

}
