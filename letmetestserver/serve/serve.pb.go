// Code generated by protoc-gen-gogo.
// source: serve.proto
// DO NOT EDIT!

/*
Package serve is a generated protocol buffer package.

It is generated from these files:
	serve.proto

It has these top-level messages:
	Artist
	Song
	Album
	EndLess
	Tree
*/
package serve

import proto "github.com/gogo/protobuf/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type Instrument int32

const (
	Instrument_Voice  Instrument = 0
	Instrument_Guitar Instrument = 1
	Instrument_Drum   Instrument = 2
)

var Instrument_name = map[int32]string{
	0: "Voice",
	1: "Guitar",
	2: "Drum",
}
var Instrument_value = map[string]int32{
	"Voice":  0,
	"Guitar": 1,
	"Drum":   2,
}

func (x Instrument) String() string {
	return proto.EnumName(Instrument_name, int32(x))
}

type Genre int32

const (
	Genre_Pop          Genre = 0
	Genre_Rock         Genre = 1
	Genre_Jazz         Genre = 2
	Genre_NintendoCore Genre = 3
	Genre_Indie        Genre = 4
	Genre_Punk         Genre = 5
	Genre_Dance        Genre = 6
)

var Genre_name = map[int32]string{
	0: "Pop",
	1: "Rock",
	2: "Jazz",
	3: "NintendoCore",
	4: "Indie",
	5: "Punk",
	6: "Dance",
}
var Genre_value = map[string]int32{
	"Pop":          0,
	"Rock":         1,
	"Jazz":         2,
	"NintendoCore": 3,
	"Indie":        4,
	"Punk":         5,
	"Dance":        6,
}

func (x Genre) String() string {
	return proto.EnumName(Genre_name, int32(x))
}

type Artist struct {
	Name string     `protobuf:"bytes,1,opt,proto3" json:"Name,omitempty"`
	Role Instrument `protobuf:"varint,2,opt,proto3,enum=serve.Instrument" json:"Role,omitempty"`
}

func (m *Artist) Reset()         { *m = Artist{} }
func (m *Artist) String() string { return proto.CompactTextString(m) }
func (*Artist) ProtoMessage()    {}

type Song struct {
	Name     string    `protobuf:"bytes,1,opt,proto3" json:"Name,omitempty"`
	Track    uint64    `protobuf:"varint,2,opt,proto3" json:"Track,omitempty"`
	Duration float64   `protobuf:"fixed64,3,opt,proto3" json:"Duration,omitempty"`
	Composer []*Artist `protobuf:"bytes,4,rep" json:"Composer,omitempty"`
}

func (m *Song) Reset()         { *m = Song{} }
func (m *Song) String() string { return proto.CompactTextString(m) }
func (*Song) ProtoMessage()    {}

func (m *Song) GetComposer() []*Artist {
	if m != nil {
		return m.Composer
	}
	return nil
}

type Album struct {
	Name     string   `protobuf:"bytes,1,opt,proto3" json:"Name,omitempty"`
	Song     []*Song  `protobuf:"bytes,2,rep" json:"Song,omitempty"`
	Genre    Genre    `protobuf:"varint,3,opt,proto3,enum=serve.Genre" json:"Genre,omitempty"`
	Year     string   `protobuf:"bytes,4,opt,proto3" json:"Year,omitempty"`
	Producer []string `protobuf:"bytes,5,rep" json:"Producer,omitempty"`
	Mediocre bool     `protobuf:"varint,6,opt,proto3" json:"Mediocre,omitempty"`
	Rated    bool     `protobuf:"varint,7,opt,proto3" json:"Rated,omitempty"`
	Epilogue string   `protobuf:"bytes,8,opt,proto3" json:"Epilogue,omitempty"`
	Likes    []bool   `protobuf:"varint,9,rep" json:"Likes,omitempty"`
}

func (m *Album) Reset()         { *m = Album{} }
func (m *Album) String() string { return proto.CompactTextString(m) }
func (*Album) ProtoMessage()    {}

func (m *Album) GetSong() []*Song {
	if m != nil {
		return m.Song
	}
	return nil
}

type EndLess struct {
	Tree *Tree `protobuf:"bytes,1,opt" json:"Tree,omitempty"`
}

func (m *EndLess) Reset()         { *m = EndLess{} }
func (m *EndLess) String() string { return proto.CompactTextString(m) }
func (*EndLess) ProtoMessage()    {}

func (m *EndLess) GetTree() *Tree {
	if m != nil {
		return m.Tree
	}
	return nil
}

type Tree struct {
	Value string `protobuf:"bytes,1,opt,proto3" json:"Value,omitempty"`
	Left  *Tree  `protobuf:"bytes,2,opt" json:"Left,omitempty"`
	Right *Tree  `protobuf:"bytes,3,opt" json:"Right,omitempty"`
}

func (m *Tree) Reset()         { *m = Tree{} }
func (m *Tree) String() string { return proto.CompactTextString(m) }
func (*Tree) ProtoMessage()    {}

func (m *Tree) GetLeft() *Tree {
	if m != nil {
		return m.Left
	}
	return nil
}

func (m *Tree) GetRight() *Tree {
	if m != nil {
		return m.Right
	}
	return nil
}

func init() {
	proto.RegisterEnum("serve.Instrument", Instrument_name, Instrument_value)
	proto.RegisterEnum("serve.Genre", Genre_name, Genre_value)
}

// Client API for Label service

type LabelClient interface {
	Produce(ctx context.Context, in *Album, opts ...grpc.CallOption) (*Album, error)
	Loop(ctx context.Context, in *EndLess, opts ...grpc.CallOption) (*EndLess, error)
}

type labelClient struct {
	cc *grpc.ClientConn
}

func NewLabelClient(cc *grpc.ClientConn) LabelClient {
	return &labelClient{cc}
}

func (c *labelClient) Produce(ctx context.Context, in *Album, opts ...grpc.CallOption) (*Album, error) {
	out := new(Album)
	err := grpc.Invoke(ctx, "/serve.Label/Produce", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *labelClient) Loop(ctx context.Context, in *EndLess, opts ...grpc.CallOption) (*EndLess, error) {
	out := new(EndLess)
	err := grpc.Invoke(ctx, "/serve.Label/Loop", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Label service

type LabelServer interface {
	Produce(context.Context, *Album) (*Album, error)
	Loop(context.Context, *EndLess) (*EndLess, error)
}

func RegisterLabelServer(s *grpc.Server, srv LabelServer) {
	s.RegisterService(&_Label_serviceDesc, srv)
}

func _Label_Produce_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(Album)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(LabelServer).Produce(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Label_Loop_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(EndLess)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(LabelServer).Loop(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Label_serviceDesc = grpc.ServiceDesc{
	ServiceName: "serve.Label",
	HandlerType: (*LabelServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Produce",
			Handler:    _Label_Produce_Handler,
		},
		{
			MethodName: "Loop",
			Handler:    _Label_Loop_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
