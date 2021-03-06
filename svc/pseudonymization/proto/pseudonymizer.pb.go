// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: pseudonymizer.proto

package pseudonymization

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
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

type ConvertRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NamespaceTo string   `protobuf:"bytes,1,opt,name=namespace_to,json=namespaceTo,proto3" json:"namespace_to,omitempty"`
	Pseudonyms  []string `protobuf:"bytes,2,rep,name=pseudonyms,proto3" json:"pseudonyms,omitempty"`
}

func (x *ConvertRequest) Reset() {
	*x = ConvertRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pseudonymizer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConvertRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConvertRequest) ProtoMessage() {}

func (x *ConvertRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pseudonymizer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConvertRequest.ProtoReflect.Descriptor instead.
func (*ConvertRequest) Descriptor() ([]byte, []int) {
	return file_pseudonymizer_proto_rawDescGZIP(), []int{0}
}

func (x *ConvertRequest) GetNamespaceTo() string {
	if x != nil {
		return x.NamespaceTo
	}
	return ""
}

func (x *ConvertRequest) GetPseudonyms() []string {
	if x != nil {
		return x.Pseudonyms
	}
	return nil
}

type ConvertResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Conversions map[string][]byte `protobuf:"bytes,1,rep,name=conversions,proto3" json:"conversions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ConvertResponse) Reset() {
	*x = ConvertResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pseudonymizer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConvertResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConvertResponse) ProtoMessage() {}

func (x *ConvertResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pseudonymizer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConvertResponse.ProtoReflect.Descriptor instead.
func (*ConvertResponse) Descriptor() ([]byte, []int) {
	return file_pseudonymizer_proto_rawDescGZIP(), []int{1}
}

func (x *ConvertResponse) GetConversions() map[string][]byte {
	if x != nil {
		return x.Conversions
	}
	return nil
}

type GenerateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount uint32 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *GenerateRequest) Reset() {
	*x = GenerateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pseudonymizer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateRequest) ProtoMessage() {}

func (x *GenerateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pseudonymizer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateRequest.ProtoReflect.Descriptor instead.
func (*GenerateRequest) Descriptor() ([]byte, []int) {
	return file_pseudonymizer_proto_rawDescGZIP(), []int{2}
}

func (x *GenerateRequest) GetAmount() uint32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type GenerateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pseudonyms []string `protobuf:"bytes,1,rep,name=pseudonyms,proto3" json:"pseudonyms,omitempty"`
}

func (x *GenerateResponse) Reset() {
	*x = GenerateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pseudonymizer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateResponse) ProtoMessage() {}

func (x *GenerateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pseudonymizer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateResponse.ProtoReflect.Descriptor instead.
func (*GenerateResponse) Descriptor() ([]byte, []int) {
	return file_pseudonymizer_proto_rawDescGZIP(), []int{3}
}

func (x *GenerateResponse) GetPseudonyms() []string {
	if x != nil {
		return x.Pseudonyms
	}
	return nil
}

var File_pseudonymizer_proto protoreflect.FileDescriptor

var file_pseudonymizer_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x70, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x53, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x5f, 0x74, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x54, 0x6f, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x73,
	0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a,
	0x70, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x73, 0x22, 0xa7, 0x01, 0x0a, 0x0f, 0x43,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x54,
	0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x70, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x3e, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0x29, 0x0a, 0x0f, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x32, 0x0a, 0x10, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e,
	0x79, 0x6d, 0x73, 0x32, 0xca, 0x01, 0x0a, 0x0d, 0x50, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79,
	0x6d, 0x69, 0x7a, 0x65, 0x72, 0x12, 0x67, 0x0a, 0x08, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x65, 0x12, 0x21, 0x2e, 0x70, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x70, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e,
	0x12, 0x0c, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x12, 0x50,
	0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x12, 0x20, 0x2e, 0x70, 0x73, 0x65, 0x75,
	0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e,
	0x76, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x70, 0x73,
	0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x14, 0x5a, 0x12, 0x2e, 0x3b, 0x70, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pseudonymizer_proto_rawDescOnce sync.Once
	file_pseudonymizer_proto_rawDescData = file_pseudonymizer_proto_rawDesc
)

func file_pseudonymizer_proto_rawDescGZIP() []byte {
	file_pseudonymizer_proto_rawDescOnce.Do(func() {
		file_pseudonymizer_proto_rawDescData = protoimpl.X.CompressGZIP(file_pseudonymizer_proto_rawDescData)
	})
	return file_pseudonymizer_proto_rawDescData
}

var file_pseudonymizer_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pseudonymizer_proto_goTypes = []interface{}{
	(*ConvertRequest)(nil),   // 0: pseudonymization.ConvertRequest
	(*ConvertResponse)(nil),  // 1: pseudonymization.ConvertResponse
	(*GenerateRequest)(nil),  // 2: pseudonymization.GenerateRequest
	(*GenerateResponse)(nil), // 3: pseudonymization.GenerateResponse
	nil,                      // 4: pseudonymization.ConvertResponse.ConversionsEntry
}
var file_pseudonymizer_proto_depIdxs = []int32{
	4, // 0: pseudonymization.ConvertResponse.conversions:type_name -> pseudonymization.ConvertResponse.ConversionsEntry
	2, // 1: pseudonymization.Pseudonymizer.Generate:input_type -> pseudonymization.GenerateRequest
	0, // 2: pseudonymization.Pseudonymizer.Convert:input_type -> pseudonymization.ConvertRequest
	3, // 3: pseudonymization.Pseudonymizer.Generate:output_type -> pseudonymization.GenerateResponse
	1, // 4: pseudonymization.Pseudonymizer.Convert:output_type -> pseudonymization.ConvertResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pseudonymizer_proto_init() }
func file_pseudonymizer_proto_init() {
	if File_pseudonymizer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pseudonymizer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConvertRequest); i {
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
		file_pseudonymizer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConvertResponse); i {
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
		file_pseudonymizer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerateRequest); i {
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
		file_pseudonymizer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerateResponse); i {
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
			RawDescriptor: file_pseudonymizer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pseudonymizer_proto_goTypes,
		DependencyIndexes: file_pseudonymizer_proto_depIdxs,
		MessageInfos:      file_pseudonymizer_proto_msgTypes,
	}.Build()
	File_pseudonymizer_proto = out.File
	file_pseudonymizer_proto_rawDesc = nil
	file_pseudonymizer_proto_goTypes = nil
	file_pseudonymizer_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PseudonymizerClient is the client API for Pseudonymizer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PseudonymizerClient interface {
	Generate(ctx context.Context, in *GenerateRequest, opts ...grpc.CallOption) (*GenerateResponse, error)
	Convert(ctx context.Context, in *ConvertRequest, opts ...grpc.CallOption) (*ConvertResponse, error)
}

type pseudonymizerClient struct {
	cc grpc.ClientConnInterface
}

func NewPseudonymizerClient(cc grpc.ClientConnInterface) PseudonymizerClient {
	return &pseudonymizerClient{cc}
}

func (c *pseudonymizerClient) Generate(ctx context.Context, in *GenerateRequest, opts ...grpc.CallOption) (*GenerateResponse, error) {
	out := new(GenerateResponse)
	err := c.cc.Invoke(ctx, "/pseudonymization.Pseudonymizer/Generate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pseudonymizerClient) Convert(ctx context.Context, in *ConvertRequest, opts ...grpc.CallOption) (*ConvertResponse, error) {
	out := new(ConvertResponse)
	err := c.cc.Invoke(ctx, "/pseudonymization.Pseudonymizer/Convert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PseudonymizerServer is the server API for Pseudonymizer service.
type PseudonymizerServer interface {
	Generate(context.Context, *GenerateRequest) (*GenerateResponse, error)
	Convert(context.Context, *ConvertRequest) (*ConvertResponse, error)
}

// UnimplementedPseudonymizerServer can be embedded to have forward compatible implementations.
type UnimplementedPseudonymizerServer struct {
}

func (*UnimplementedPseudonymizerServer) Generate(context.Context, *GenerateRequest) (*GenerateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Generate not implemented")
}
func (*UnimplementedPseudonymizerServer) Convert(context.Context, *ConvertRequest) (*ConvertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Convert not implemented")
}

func RegisterPseudonymizerServer(s *grpc.Server, srv PseudonymizerServer) {
	s.RegisterService(&_Pseudonymizer_serviceDesc, srv)
}

func _Pseudonymizer_Generate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PseudonymizerServer).Generate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pseudonymization.Pseudonymizer/Generate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PseudonymizerServer).Generate(ctx, req.(*GenerateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pseudonymizer_Convert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConvertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PseudonymizerServer).Convert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pseudonymization.Pseudonymizer/Convert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PseudonymizerServer).Convert(ctx, req.(*ConvertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Pseudonymizer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pseudonymization.Pseudonymizer",
	HandlerType: (*PseudonymizerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Generate",
			Handler:    _Pseudonymizer_Generate_Handler,
		},
		{
			MethodName: "Convert",
			Handler:    _Pseudonymizer_Convert_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pseudonymizer.proto",
}
