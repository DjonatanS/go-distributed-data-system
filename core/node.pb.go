// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v4.25.3
// source: 25_distribuited_system/core/node.proto

package core

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Action        string                 `protobuf:"bytes,1,opt,name=action,proto3" json:"action,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Request) Reset() {
	*x = Request{}
	mi := &file__25_distribuited_system_core_node_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file__25_distribuited_system_core_node_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file__25_distribuited_system_core_node_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          string                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Response) Reset() {
	*x = Response{}
	mi := &file__25_distribuited_system_core_node_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file__25_distribuited_system_core_node_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file__25_distribuited_system_core_node_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File__25_distribuited_system_core_node_proto protoreflect.FileDescriptor

const file__25_distribuited_system_core_node_proto_rawDesc = "" +
	"\n" +
	"&25_distribuited_system/core/node.proto\x12\x04core\"!\n" +
	"\aRequest\x12\x16\n" +
	"\x06action\x18\x01 \x01(\tR\x06action\"\x1e\n" +
	"\bResponse\x12\x12\n" +
	"\x04data\x18\x01 \x01(\tR\x04data2o\n" +
	"\vNodeService\x12/\n" +
	"\fReportStatus\x12\r.core.Request\x1a\x0e.core.Response\"\x00\x12/\n" +
	"\n" +
	"AssignTask\x12\r.core.Request\x1a\x0e.core.Response\"\x000\x01B\bZ\x06.;coreb\x06proto3"

var (
	file__25_distribuited_system_core_node_proto_rawDescOnce sync.Once
	file__25_distribuited_system_core_node_proto_rawDescData []byte
)

func file__25_distribuited_system_core_node_proto_rawDescGZIP() []byte {
	file__25_distribuited_system_core_node_proto_rawDescOnce.Do(func() {
		file__25_distribuited_system_core_node_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file__25_distribuited_system_core_node_proto_rawDesc), len(file__25_distribuited_system_core_node_proto_rawDesc)))
	})
	return file__25_distribuited_system_core_node_proto_rawDescData
}

var file__25_distribuited_system_core_node_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file__25_distribuited_system_core_node_proto_goTypes = []any{
	(*Request)(nil),  // 0: core.Request
	(*Response)(nil), // 1: core.Response
}
var file__25_distribuited_system_core_node_proto_depIdxs = []int32{
	0, // 0: core.NodeService.ReportStatus:input_type -> core.Request
	0, // 1: core.NodeService.AssignTask:input_type -> core.Request
	1, // 2: core.NodeService.ReportStatus:output_type -> core.Response
	1, // 3: core.NodeService.AssignTask:output_type -> core.Response
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file__25_distribuited_system_core_node_proto_init() }
func file__25_distribuited_system_core_node_proto_init() {
	if File__25_distribuited_system_core_node_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file__25_distribuited_system_core_node_proto_rawDesc), len(file__25_distribuited_system_core_node_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file__25_distribuited_system_core_node_proto_goTypes,
		DependencyIndexes: file__25_distribuited_system_core_node_proto_depIdxs,
		MessageInfos:      file__25_distribuited_system_core_node_proto_msgTypes,
	}.Build()
	File__25_distribuited_system_core_node_proto = out.File
	file__25_distribuited_system_core_node_proto_goTypes = nil
	file__25_distribuited_system_core_node_proto_depIdxs = nil
}
