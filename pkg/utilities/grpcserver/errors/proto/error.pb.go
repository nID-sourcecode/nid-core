// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        v3.12.4
// source: error.proto

package grpcerr

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ErrorType int32

const (
	ErrorType_UNSPECIFIED         ErrorType = 0
	ErrorType_CANCELED            ErrorType = 1
	ErrorType_UNKNOWN             ErrorType = 2
	ErrorType_INVALID_ARGUMENT    ErrorType = 3
	ErrorType_DEADLINE_EXCEEDED   ErrorType = 4
	ErrorType_NOT_FOUND           ErrorType = 5
	ErrorType_ALREADY_EXISTS      ErrorType = 6
	ErrorType_PERMISSION_DENIED   ErrorType = 7
	ErrorType_RESOURCE_EXHAUSTED  ErrorType = 8
	ErrorType_FAILED_PRECONDITION ErrorType = 9
	ErrorType_ABORTED             ErrorType = 10
	ErrorType_OUT_OF_RANGE        ErrorType = 11
	ErrorType_UNIMPLEMENTED       ErrorType = 12
	ErrorType_INTERNAL            ErrorType = 13
	ErrorType_UNAVAILABLE         ErrorType = 14
	ErrorType_DATA_LOSS           ErrorType = 15
	ErrorType_UNAUTHENTICATED     ErrorType = 16
)

// Enum value maps for ErrorType.
var (
	ErrorType_name = map[int32]string{
		0:  "UNSPECIFIED",
		1:  "CANCELED",
		2:  "UNKNOWN",
		3:  "INVALID_ARGUMENT",
		4:  "DEADLINE_EXCEEDED",
		5:  "NOT_FOUND",
		6:  "ALREADY_EXISTS",
		7:  "PERMISSION_DENIED",
		8:  "RESOURCE_EXHAUSTED",
		9:  "FAILED_PRECONDITION",
		10: "ABORTED",
		11: "OUT_OF_RANGE",
		12: "UNIMPLEMENTED",
		13: "INTERNAL",
		14: "UNAVAILABLE",
		15: "DATA_LOSS",
		16: "UNAUTHENTICATED",
	}
	ErrorType_value = map[string]int32{
		"UNSPECIFIED":         0,
		"CANCELED":            1,
		"UNKNOWN":             2,
		"INVALID_ARGUMENT":    3,
		"DEADLINE_EXCEEDED":   4,
		"NOT_FOUND":           5,
		"ALREADY_EXISTS":      6,
		"PERMISSION_DENIED":   7,
		"RESOURCE_EXHAUSTED":  8,
		"FAILED_PRECONDITION": 9,
		"ABORTED":             10,
		"OUT_OF_RANGE":        11,
		"UNIMPLEMENTED":       12,
		"INTERNAL":            13,
		"UNAVAILABLE":         14,
		"DATA_LOSS":           15,
		"UNAUTHENTICATED":     16,
	}
)

func (x ErrorType) Enum() *ErrorType {
	p := new(ErrorType)
	*p = x
	return p
}

func (x ErrorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorType) Descriptor() protoreflect.EnumDescriptor {
	return file_error_proto_enumTypes[0].Descriptor()
}

func (ErrorType) Type() protoreflect.EnumType {
	return &file_error_proto_enumTypes[0]
}

func (x ErrorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorType.Descriptor instead.
func (ErrorType) EnumDescriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{0}
}

type ErrDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrType ErrorType `protobuf:"varint,1,opt,name=err_type,json=errType,proto3,enum=grpcerr.ErrorType" json:"err_type,omitempty"`
}

func (x *ErrDetails) Reset() {
	*x = ErrDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_error_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrDetails) ProtoMessage() {}

func (x *ErrDetails) ProtoReflect() protoreflect.Message {
	mi := &file_error_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrDetails.ProtoReflect.Descriptor instead.
func (*ErrDetails) Descriptor() ([]byte, []int) {
	return file_error_proto_rawDescGZIP(), []int{0}
}

func (x *ErrDetails) GetErrType() ErrorType {
	if x != nil {
		return x.ErrType
	}
	return ErrorType_UNSPECIFIED
}

var File_error_proto protoreflect.FileDescriptor

var file_error_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x67,
	0x72, 0x70, 0x63, 0x65, 0x72, 0x72, 0x22, 0x3b, 0x0a, 0x0a, 0x45, 0x72, 0x72, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x12, 0x2d, 0x0a, 0x08, 0x65, 0x72, 0x72, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x65, 0x72, 0x72,
	0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x65, 0x72, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x2a, 0xc4, 0x02, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x45, 0x44, 0x10, 0x01,
	0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x02, 0x12, 0x14, 0x0a,
	0x10, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x41, 0x52, 0x47, 0x55, 0x4d, 0x45, 0x4e,
	0x54, 0x10, 0x03, 0x12, 0x15, 0x0a, 0x11, 0x44, 0x45, 0x41, 0x44, 0x4c, 0x49, 0x4e, 0x45, 0x5f,
	0x45, 0x58, 0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10, 0x04, 0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x4f,
	0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x05, 0x12, 0x12, 0x0a, 0x0e, 0x41, 0x4c, 0x52,
	0x45, 0x41, 0x44, 0x59, 0x5f, 0x45, 0x58, 0x49, 0x53, 0x54, 0x53, 0x10, 0x06, 0x12, 0x15, 0x0a,
	0x11, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x4e, 0x49,
	0x45, 0x44, 0x10, 0x07, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45,
	0x5f, 0x45, 0x58, 0x48, 0x41, 0x55, 0x53, 0x54, 0x45, 0x44, 0x10, 0x08, 0x12, 0x17, 0x0a, 0x13,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x5f, 0x50, 0x52, 0x45, 0x43, 0x4f, 0x4e, 0x44, 0x49, 0x54,
	0x49, 0x4f, 0x4e, 0x10, 0x09, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x42, 0x4f, 0x52, 0x54, 0x45, 0x44,
	0x10, 0x0a, 0x12, 0x10, 0x0a, 0x0c, 0x4f, 0x55, 0x54, 0x5f, 0x4f, 0x46, 0x5f, 0x52, 0x41, 0x4e,
	0x47, 0x45, 0x10, 0x0b, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x49, 0x4d, 0x50, 0x4c, 0x45, 0x4d,
	0x45, 0x4e, 0x54, 0x45, 0x44, 0x10, 0x0c, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x4e, 0x54, 0x45, 0x52,
	0x4e, 0x41, 0x4c, 0x10, 0x0d, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x41, 0x56, 0x41, 0x49, 0x4c,
	0x41, 0x42, 0x4c, 0x45, 0x10, 0x0e, 0x12, 0x0d, 0x0a, 0x09, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x4c,
	0x4f, 0x53, 0x53, 0x10, 0x0f, 0x12, 0x13, 0x0a, 0x0f, 0x55, 0x4e, 0x41, 0x55, 0x54, 0x48, 0x45,
	0x4e, 0x54, 0x49, 0x43, 0x41, 0x54, 0x45, 0x44, 0x10, 0x10, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_error_proto_rawDescOnce sync.Once
	file_error_proto_rawDescData = file_error_proto_rawDesc
)

func file_error_proto_rawDescGZIP() []byte {
	file_error_proto_rawDescOnce.Do(func() {
		file_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_error_proto_rawDescData)
	})
	return file_error_proto_rawDescData
}

var (
	file_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
	file_error_proto_msgTypes  = make([]protoimpl.MessageInfo, 1)
	file_error_proto_goTypes   = []interface{}{
		(ErrorType)(0),     // 0: grpcerr.ErrorType
		(*ErrDetails)(nil), // 1: grpcerr.ErrDetails
	}
)

var file_error_proto_depIdxs = []int32{
	0, // 0: grpcerr.ErrDetails.err_type:type_name -> grpcerr.ErrorType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_error_proto_init() }
func file_error_proto_init() {
	if File_error_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_error_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrDetails); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_error_proto_goTypes,
		DependencyIndexes: file_error_proto_depIdxs,
		EnumInfos:         file_error_proto_enumTypes,
		MessageInfos:      file_error_proto_msgTypes,
	}.Build()
	File_error_proto = out.File
	file_error_proto_rawDesc = nil
	file_error_proto_goTypes = nil
	file_error_proto_depIdxs = nil
}