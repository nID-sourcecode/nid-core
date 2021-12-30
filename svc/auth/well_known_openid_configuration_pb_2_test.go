// Code generated by protoc-gen-go. Edited for testing purposes.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: auth_test_2.proto

package main

import (
	"context"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/protoc-gen-go/descriptor"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
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

type Empty2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty2) Reset() {
	*x = Empty2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_test_2_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty2) ProtoMessage() {}

func (x *Empty2) ProtoReflect() protoreflect.Message {
	mi := &file_auth_test_2_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty2.ProtoReflect.Descriptor instead.
func (*Empty2) Descriptor() ([]byte, []int) {
	return file_auth_test_2_proto_rawDescGZIP(), []int{0}
}

var File_auth_test_2_proto protoreflect.FileDescriptor

var file_auth_test_2_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x32, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x08, 0x0a, 0x06, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x32, 0xca, 0x04, 0x0a, 0x09, 0x41, 0x75,
	0x74, 0x68, 0x54, 0x65, 0x73, 0x74, 0x32, 0x12, 0x38, 0x0a, 0x09, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x65, 0x12, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x1a, 0x07, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x12, 0x0b,
	0x2f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x32, 0xd0, 0xac, 0xc9, 0xba, 0x02,
	0x01, 0x12, 0x30, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x07, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x32, 0x1a, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x22, 0x15, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x09, 0x12, 0x07, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0xd0, 0xac, 0xc9,
	0xba, 0x02, 0x02, 0x12, 0x2e, 0x0a, 0x04, 0x4a, 0x77, 0x6b, 0x73, 0x12, 0x07, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x32, 0x1a, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x22, 0x14, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x08, 0x12, 0x06, 0x2f, 0x6a, 0x77, 0x6b, 0x73, 0x32, 0xd0, 0xac, 0xc9,
	0xba, 0x02, 0x03, 0x12, 0x3e, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x1a, 0x07, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x32, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0xd0, 0xac, 0xc9,
	0xba, 0x02, 0x04, 0x12, 0x3c, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x6f,
	0x63, 0x73, 0x12, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x1a, 0x07, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x32, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12, 0x0d, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x64, 0x6f, 0x63, 0x73, 0x32, 0xd0, 0xac, 0xc9, 0xba, 0x02,
	0x05, 0x12, 0x36, 0x0a, 0x08, 0x4f, 0x70, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x07, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x1a, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x22,
	0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x6f, 0x70, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x32, 0xd0, 0xac, 0xc9, 0xba, 0x02, 0x06, 0x12, 0x35, 0x0a, 0x05, 0x4f, 0x70, 0x54,
	0x6f, 0x73, 0x12, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x1a, 0x07, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x32, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x12, 0x0c, 0x2f, 0x6f,
	0x70, 0x5f, 0x74, 0x6f, 0x73, 0x5f, 0x75, 0x72, 0x69, 0x32, 0xd0, 0xac, 0xc9, 0xba, 0x02, 0x07,
	0x12, 0x32, 0x0a, 0x06, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x12, 0x07, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x32, 0x1a, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x22, 0x16, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0a, 0x12, 0x08, 0x2f, 0x72, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x32, 0xd0, 0xac,
	0xc9, 0xba, 0x02, 0x08, 0x12, 0x36, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x1a, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x32, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x69, 0x6e, 0x66, 0x6f, 0x32, 0xd0, 0xac, 0xc9, 0xba, 0x02, 0x0a, 0x12, 0x48, 0x0a, 0x12,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x66, 0x72, 0x61,
	0x6d, 0x65, 0x12, 0x07, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x1a, 0x07, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x32, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x12, 0x12, 0x2f, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x73, 0x65, 0x73, 0x73, 0x5f, 0x69, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x32,
	0xd0, 0xac, 0xc9, 0xba, 0x02, 0x0b, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x3b, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_test_2_proto_rawDescOnce sync.Once
	file_auth_test_2_proto_rawDescData = file_auth_test_2_proto_rawDesc
)

func file_auth_test_2_proto_rawDescGZIP() []byte {
	file_auth_test_2_proto_rawDescOnce.Do(func() {
		file_auth_test_2_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_test_2_proto_rawDescData)
	})
	return file_auth_test_2_proto_rawDescData
}

var (
	file_auth_test_2_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
	file_auth_test_2_proto_goTypes  = []interface{}{
		(*Empty2)(nil), // 0: Empty2
	}
)

var file_auth_test_2_proto_depIdxs = []int32{
	0,  // 0: AuthTest2.Authorize:input_type -> Empty2
	0,  // 1: AuthTest2.Token:input_type -> Empty2
	0,  // 2: AuthTest2.Jwks:input_type -> Empty2
	0,  // 3: AuthTest2.Registration:input_type -> Empty2
	0,  // 4: AuthTest2.ServiceDocs:input_type -> Empty2
	0,  // 5: AuthTest2.OpPolicy:input_type -> Empty2
	0,  // 6: AuthTest2.OpTos:input_type -> Empty2
	0,  // 7: AuthTest2.Revoke:input_type -> Empty2
	0,  // 8: AuthTest2.UserInfo:input_type -> Empty2
	0,  // 9: AuthTest2.CheckSessionIframe:input_type -> Empty2
	0,  // 10: AuthTest2.Authorize:output_type -> Empty2
	0,  // 11: AuthTest2.Token:output_type -> Empty2
	0,  // 12: AuthTest2.Jwks:output_type -> Empty2
	0,  // 13: AuthTest2.Registration:output_type -> Empty2
	0,  // 14: AuthTest2.ServiceDocs:output_type -> Empty2
	0,  // 15: AuthTest2.OpPolicy:output_type -> Empty2
	0,  // 16: AuthTest2.OpTos:output_type -> Empty2
	0,  // 17: AuthTest2.Revoke:output_type -> Empty2
	0,  // 18: AuthTest2.UserInfo:output_type -> Empty2
	0,  // 19: AuthTest2.CheckSessionIframe:output_type -> Empty2
	10, // [10:20] is the sub-list for method output_type
	0,  // [0:10] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_auth_test_2_proto_init() }
func file_auth_test_2_proto_init() {
	if File_auth_test_2_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_auth_test_2_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty2); i {
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
			RawDescriptor: file_auth_test_2_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_test_2_proto_goTypes,
		DependencyIndexes: file_auth_test_2_proto_depIdxs,
		MessageInfos:      file_auth_test_2_proto_msgTypes,
	}.Build()
	File_auth_test_2_proto = out.File
	file_auth_test_2_proto_rawDesc = nil
	file_auth_test_2_proto_goTypes = nil
	file_auth_test_2_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ context.Context
	_ grpc.ClientConnInterface
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AuthTest2Client is the client API for AuthTest2 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthTest2Client interface {
	Authorize(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error)
	Token(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error)
	Jwks(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error)
	Registration(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error)
	ServiceDocs(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error)
	OpPolicy(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error)
	OpTos(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error)
	Revoke(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error)
	UserInfo(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error)
	CheckSessionIframe(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error)
}

type authTest2Client struct {
	cc grpc.ClientConnInterface
}

func NewAuthTest2Client(cc grpc.ClientConnInterface) AuthTest2Client {
	return &authTest2Client{cc}
}

func (c *authTest2Client) Authorize(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error) {
	out := new(Empty2)
	err := c.cc.Invoke(ctx, "/AuthTest2/Authorize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTest2Client) Token(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error) {
	out := new(Empty2)
	err := c.cc.Invoke(ctx, "/AuthTest2/Token", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTest2Client) Jwks(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error) {
	out := new(Empty2)
	err := c.cc.Invoke(ctx, "/AuthTest2/Jwks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTest2Client) Registration(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error) {
	out := new(Empty2)
	err := c.cc.Invoke(ctx, "/AuthTest2/Registration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTest2Client) ServiceDocs(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error) {
	out := new(Empty2)
	err := c.cc.Invoke(ctx, "/AuthTest2/ServiceDocs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTest2Client) OpPolicy(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error) {
	out := new(Empty2)
	err := c.cc.Invoke(ctx, "/AuthTest2/OpPolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTest2Client) OpTos(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error) {
	out := new(Empty2)
	err := c.cc.Invoke(ctx, "/AuthTest2/OpTos", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTest2Client) Revoke(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error) {
	out := new(Empty2)
	err := c.cc.Invoke(ctx, "/AuthTest2/Revoke", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTest2Client) UserInfo(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error) {
	out := new(Empty2)
	err := c.cc.Invoke(ctx, "/AuthTest2/UserInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTest2Client) CheckSessionIframe(ctx context.Context, in *Empty2, opts ...grpc.CallOption) (*Empty2, error) {
	out := new(Empty2)
	err := c.cc.Invoke(ctx, "/AuthTest2/CheckSessionIframe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthTest2Server is the server API for AuthTest2 service.
type AuthTest2Server interface {
	Authorize(context.Context, *Empty2) (*Empty2, error)
	Token(context.Context, *Empty2) (*Empty2, error)
	Jwks(context.Context, *Empty2) (*Empty2, error)
	Registration(context.Context, *Empty2) (*Empty2, error)
	ServiceDocs(context.Context, *Empty2) (*Empty2, error)
	OpPolicy(context.Context, *Empty2) (*Empty2, error)
	OpTos(context.Context, *Empty2) (*Empty2, error)
	Revoke(context.Context, *Empty2) (*Empty2, error)
	UserInfo(context.Context, *Empty2) (*Empty2, error)
	CheckSessionIframe(context.Context, *Empty2) (*Empty2, error)
}

// UnimplementedAuthTest2Server can be embedded to have forward compatible implementations.
type UnimplementedAuthTest2Server struct {
}

func (*UnimplementedAuthTest2Server) Authorize(context.Context, *Empty2) (*Empty2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}

func (*UnimplementedAuthTest2Server) Token(context.Context, *Empty2) (*Empty2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Token not implemented")
}

func (*UnimplementedAuthTest2Server) Jwks(context.Context, *Empty2) (*Empty2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Jwks not implemented")
}

func (*UnimplementedAuthTest2Server) Registration(context.Context, *Empty2) (*Empty2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Registration not implemented")
}

func (*UnimplementedAuthTest2Server) ServiceDocs(context.Context, *Empty2) (*Empty2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServiceDocs not implemented")
}

func (*UnimplementedAuthTest2Server) OpPolicy(context.Context, *Empty2) (*Empty2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpPolicy not implemented")
}

func (*UnimplementedAuthTest2Server) OpTos(context.Context, *Empty2) (*Empty2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpTos not implemented")
}

func (*UnimplementedAuthTest2Server) Revoke(context.Context, *Empty2) (*Empty2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Revoke not implemented")
}

func (*UnimplementedAuthTest2Server) UserInfo(context.Context, *Empty2) (*Empty2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserInfo not implemented")
}

func (*UnimplementedAuthTest2Server) CheckSessionIframe(context.Context, *Empty2) (*Empty2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckSessionIframe not implemented")
}

func RegisterAuthTest2Server(s *grpc.Server, srv AuthTest2Server) {
	s.RegisterService(&_AuthTest2_serviceDesc, srv)
}

func _AuthTest2_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTest2Server).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest2/Authorize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTest2Server).Authorize(ctx, req.(*Empty2))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest2_Token_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTest2Server).Token(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest2/Token",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTest2Server).Token(ctx, req.(*Empty2))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest2_Jwks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTest2Server).Jwks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest2/Jwks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTest2Server).Jwks(ctx, req.(*Empty2))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest2_Registration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTest2Server).Registration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest2/Registration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTest2Server).Registration(ctx, req.(*Empty2))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest2_ServiceDocs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTest2Server).ServiceDocs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest2/ServiceDocs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTest2Server).ServiceDocs(ctx, req.(*Empty2))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest2_OpPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTest2Server).OpPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest2/OpPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTest2Server).OpPolicy(ctx, req.(*Empty2))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest2_OpTos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTest2Server).OpTos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest2/OpTos",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTest2Server).OpTos(ctx, req.(*Empty2))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest2_Revoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTest2Server).Revoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest2/Revoke",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTest2Server).Revoke(ctx, req.(*Empty2))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest2_UserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTest2Server).UserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest2/UserInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTest2Server).UserInfo(ctx, req.(*Empty2))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest2_CheckSessionIframe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTest2Server).CheckSessionIframe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest2/CheckSessionIframe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTest2Server).CheckSessionIframe(ctx, req.(*Empty2))
	}
	return interceptor(ctx, in, info, handler)
}

var _AuthTest2_serviceDesc = grpc.ServiceDesc{
	ServiceName: "AuthTest2",
	HandlerType: (*AuthTest2Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authorize",
			Handler:    _AuthTest2_Authorize_Handler,
		},
		{
			MethodName: "Token",
			Handler:    _AuthTest2_Token_Handler,
		},
		{
			MethodName: "Jwks",
			Handler:    _AuthTest2_Jwks_Handler,
		},
		{
			MethodName: "Registration",
			Handler:    _AuthTest2_Registration_Handler,
		},
		{
			MethodName: "ServiceDocs",
			Handler:    _AuthTest2_ServiceDocs_Handler,
		},
		{
			MethodName: "OpPolicy",
			Handler:    _AuthTest2_OpPolicy_Handler,
		},
		{
			MethodName: "OpTos",
			Handler:    _AuthTest2_OpTos_Handler,
		},
		{
			MethodName: "Revoke",
			Handler:    _AuthTest2_Revoke_Handler,
		},
		{
			MethodName: "UserInfo",
			Handler:    _AuthTest2_UserInfo_Handler,
		},
		{
			MethodName: "CheckSessionIframe",
			Handler:    _AuthTest2_CheckSessionIframe_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth_test_2.proto",
}
