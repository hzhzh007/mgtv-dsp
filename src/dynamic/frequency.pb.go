// Code generated by protoc-gen-go.
// source: frequency.proto
// DO NOT EDIT!

/*
Package dynamic is a generated protocol buffer package.

It is generated from these files:
	frequency.proto

It has these top-level messages:
	Record
	RedisValue
*/
package dynamic

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

type FreqType int32

const (
	FreqType_FreqPerDay   FreqType = 0
	FreqType_FreqPerWeek  FreqType = 1
	FreqType_FreqPerMonth FreqType = 2
	FreqType_FreqCustom   FreqType = 3
)

var FreqType_name = map[int32]string{
	0: "FreqPerDay",
	1: "FreqPerWeek",
	2: "FreqPerMonth",
	3: "FreqCustom",
}
var FreqType_value = map[string]int32{
	"FreqPerDay":   0,
	"FreqPerWeek":  1,
	"FreqPerMonth": 2,
	"FreqCustom":   3,
}

func (x FreqType) String() string {
	return proto.EnumName(FreqType_name, int32(x))
}
func (FreqType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Record struct {
	Id      int32    `protobuf:"varint,1,opt,name=Id,json=id" json:"Id,omitempty"`
	Expire  int32    `protobuf:"varint,2,opt,name=Expire,json=expire" json:"Expire,omitempty"`
	Type    FreqType `protobuf:"varint,3,opt,name=Type,json=type,enum=dynamic.FreqType" json:"Type,omitempty"`
	Counter int32    `protobuf:"varint,4,opt,name=Counter,json=counter" json:"Counter,omitempty"`
}

func (m *Record) Reset()                    { *m = Record{} }
func (m *Record) String() string            { return proto.CompactTextString(m) }
func (*Record) ProtoMessage()               {}
func (*Record) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type RedisValue struct {
	Impression []*Record `protobuf:"bytes,1,rep,name=Impression,json=impression" json:"Impression,omitempty"`
}

func (m *RedisValue) Reset()                    { *m = RedisValue{} }
func (m *RedisValue) String() string            { return proto.CompactTextString(m) }
func (*RedisValue) ProtoMessage()               {}
func (*RedisValue) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *RedisValue) GetImpression() []*Record {
	if m != nil {
		return m.Impression
	}
	return nil
}

func init() {
	proto.RegisterType((*Record)(nil), "dynamic.Record")
	proto.RegisterType((*RedisValue)(nil), "dynamic.RedisValue")
	proto.RegisterEnum("dynamic.FreqType", FreqType_name, FreqType_value)
}

func init() { proto.RegisterFile("frequency.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 238 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x3c, 0x90, 0x41, 0x4b, 0xc4, 0x30,
	0x10, 0x85, 0x4d, 0x5b, 0x5b, 0x99, 0x95, 0x36, 0xe6, 0x20, 0x39, 0x96, 0x05, 0xa1, 0x78, 0xa8,
	0xb0, 0x9e, 0x3d, 0xad, 0x0a, 0x7b, 0x58, 0x90, 0x20, 0x7a, 0xae, 0xcd, 0x88, 0x41, 0x9b, 0xa4,
	0x69, 0x0a, 0xe6, 0xdf, 0x2f, 0xdb, 0xed, 0xf6, 0x36, 0xdf, 0x9b, 0x79, 0x6f, 0x86, 0x81, 0xe2,
	0xdb, 0x61, 0x3f, 0xa2, 0x6e, 0x43, 0x6d, 0x9d, 0xf1, 0x86, 0x65, 0x32, 0xe8, 0xa6, 0x53, 0xed,
	0xba, 0x87, 0x54, 0x60, 0x6b, 0x9c, 0x64, 0x39, 0x44, 0x3b, 0xc9, 0x49, 0x49, 0xaa, 0x4b, 0x11,
	0x29, 0xc9, 0x6e, 0x21, 0x7d, 0xf9, 0xb7, 0xca, 0x21, 0x8f, 0x26, 0x2d, 0xc5, 0x89, 0xd8, 0x1d,
	0x24, 0xef, 0xc1, 0x22, 0x8f, 0x4b, 0x52, 0xe5, 0x9b, 0x9b, 0x7a, 0x4e, 0xaa, 0x5f, 0x1d, 0xf6,
	0xc7, 0x86, 0x48, 0x7c, 0xb0, 0xc8, 0x38, 0x64, 0x5b, 0x33, 0x6a, 0x8f, 0x8e, 0x27, 0x93, 0x3f,
	0x6b, 0x4f, 0xb8, 0x7e, 0x02, 0x10, 0x28, 0xd5, 0xf0, 0xd1, 0xfc, 0x8d, 0xc8, 0x1e, 0x00, 0x76,
	0x9d, 0x75, 0x38, 0x0c, 0xca, 0x68, 0x4e, 0xca, 0xb8, 0x5a, 0x6d, 0x8a, 0x25, 0xf4, 0x74, 0x9b,
	0x00, 0xb5, 0x8c, 0xdc, 0xef, 0xe1, 0xea, 0xbc, 0x8a, 0xe5, 0x00, 0xc7, 0xfa, 0x0d, 0xdd, 0x73,
	0x13, 0xe8, 0x05, 0x2b, 0x60, 0x35, 0xf3, 0x27, 0xe2, 0x2f, 0x25, 0x8c, 0xc2, 0xf5, 0x2c, 0xec,
	0x8d, 0xf6, 0x3f, 0x34, 0x3a, 0x5b, 0xb6, 0xe3, 0xe0, 0x4d, 0x47, 0xe3, 0xaf, 0x74, 0x7a, 0xc8,
	0xe3, 0x21, 0x00, 0x00, 0xff, 0xff, 0x59, 0x4a, 0x0d, 0x47, 0x23, 0x01, 0x00, 0x00,
}