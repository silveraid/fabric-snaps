// Code generated by protoc-gen-go. DO NOT EDIT.
// source: peers.proto

/*
Package membership is a generated protocol buffer package.

It is generated from these files:
	peers.proto

It has these top-level messages:
	PeerEndpoint
	PeerEndpoints
*/
package membership

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// PeerEndpoint contains the internal and external endpoint of a peer
type PeerEndpoint struct {
	Endpoint         string `protobuf:"bytes,1,opt,name=Endpoint,json=endpoint" json:"Endpoint,omitempty"`
	InternalEndpoint string `protobuf:"bytes,2,opt,name=InternalEndpoint,json=internalEndpoint" json:"InternalEndpoint,omitempty"`
	MSPid            []byte `protobuf:"bytes,3,opt,name=MSPid,json=mSPid,proto3" json:"MSPid,omitempty"`
}

func (m *PeerEndpoint) Reset()                    { *m = PeerEndpoint{} }
func (m *PeerEndpoint) String() string            { return proto.CompactTextString(m) }
func (*PeerEndpoint) ProtoMessage()               {}
func (*PeerEndpoint) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PeerEndpoint) GetEndpoint() string {
	if m != nil {
		return m.Endpoint
	}
	return ""
}

func (m *PeerEndpoint) GetInternalEndpoint() string {
	if m != nil {
		return m.InternalEndpoint
	}
	return ""
}

func (m *PeerEndpoint) GetMSPid() []byte {
	if m != nil {
		return m.MSPid
	}
	return nil
}

// PeerEndpoints contains a list of peer endpoints
type PeerEndpoints struct {
	Endpoints []*PeerEndpoint `protobuf:"bytes,1,rep,name=Endpoints,json=endpoints" json:"Endpoints,omitempty"`
}

func (m *PeerEndpoints) Reset()                    { *m = PeerEndpoints{} }
func (m *PeerEndpoints) String() string            { return proto.CompactTextString(m) }
func (*PeerEndpoints) ProtoMessage()               {}
func (*PeerEndpoints) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PeerEndpoints) GetEndpoints() []*PeerEndpoint {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

func init() {
	proto.RegisterType((*PeerEndpoint)(nil), "api.PeerEndpoint")
	proto.RegisterType((*PeerEndpoints)(nil), "api.PeerEndpoints")
}

func init() { proto.RegisterFile("peers.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 201 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0x48, 0x4d, 0x2d,
	0x2a, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e, 0x2c, 0xc8, 0x54, 0xca, 0xe1, 0xe2,
	0x09, 0x48, 0x4d, 0x2d, 0x72, 0xcd, 0x4b, 0x29, 0xc8, 0xcf, 0xcc, 0x2b, 0x11, 0x92, 0xe2, 0xe2,
	0x80, 0xb1, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x38, 0x52, 0x61, 0x72, 0x5a, 0x5c, 0x02,
	0x9e, 0x79, 0x25, 0xa9, 0x45, 0x79, 0x89, 0x39, 0x70, 0x35, 0x4c, 0x60, 0x35, 0x02, 0x99, 0x68,
	0xe2, 0x42, 0x22, 0x5c, 0xac, 0xbe, 0xc1, 0x01, 0x99, 0x29, 0x12, 0xcc, 0x0a, 0x8c, 0x1a, 0x3c,
	0x41, 0xac, 0xb9, 0x20, 0x8e, 0x92, 0x03, 0x17, 0x2f, 0xb2, 0x6d, 0xc5, 0x42, 0xfa, 0x5c, 0x9c,
	0x70, 0x8e, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0xb7, 0x91, 0xa0, 0x5e, 0x62, 0x41, 0xa6, 0x1e, 0xb2,
	0xb2, 0x20, 0x4e, 0x98, 0x13, 0x8a, 0x9d, 0x8c, 0xa3, 0x0c, 0xd3, 0x33, 0x4b, 0x32, 0x4a, 0x93,
	0xf4, 0x92, 0xf3, 0x73, 0xf5, 0x8b, 0x53, 0x93, 0x4b, 0x8b, 0x52, 0xb3, 0x53, 0x2b, 0xf5, 0xd3,
	0x12, 0x93, 0x8a, 0x32, 0x93, 0x75, 0x8b, 0xf3, 0x12, 0x0b, 0x8a, 0xf5, 0x13, 0x0b, 0x32, 0xf5,
	0xc1, 0x7e, 0x2c, 0xd6, 0x07, 0x79, 0x38, 0x89, 0x0d, 0xcc, 0x31, 0x06, 0x04, 0x00, 0x00, 0xff,
	0xff, 0x5e, 0xaa, 0x14, 0xf8, 0xff, 0x00, 0x00, 0x00,
}
