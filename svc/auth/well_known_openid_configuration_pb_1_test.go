// Code generated by protoc-gen-go. Edited for testing purposes.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: auth_test.proto

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_auth_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_auth_test_proto_rawDescGZIP(), []int{0}
}

var File_auth_test_proto protoreflect.FileDescriptor

var file_auth_test_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x07, 0x0a,
	0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0xe4, 0x04, 0x0a, 0x08, 0x41, 0x75, 0x74, 0x68, 0x54,
	0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x09, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65,
	0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x65, 0xd0, 0xac, 0xc9, 0xba, 0x02, 0x01, 0x12, 0x2d, 0x0a, 0x05, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x08, 0x12, 0x06, 0x2f, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0xd0, 0xac, 0xc9, 0xba, 0x02, 0x02, 0x12, 0x2b, 0x0a, 0x04, 0x4a, 0x77, 0x6b,
	0x73, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x13, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x07, 0x12, 0x05, 0x2f, 0x6a, 0x77, 0x6b, 0x73,
	0xd0, 0xac, 0xc9, 0xba, 0x02, 0x03, 0x12, 0x3b, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12, 0x0d,
	0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0xd0, 0xac, 0xc9,
	0xba, 0x02, 0x04, 0x12, 0x39, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x6f,
	0x63, 0x73, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x12, 0x0c, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x64, 0x6f, 0x63, 0x73, 0xd0, 0xac, 0xc9, 0xba, 0x02, 0x05, 0x12, 0x33,
	0x0a, 0x08, 0x4f, 0x70, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x0b, 0x12, 0x09, 0x2f, 0x6f, 0x70, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0xd0, 0xac, 0xc9,
	0xba, 0x02, 0x06, 0x12, 0x32, 0x0a, 0x05, 0x4f, 0x70, 0x54, 0x6f, 0x73, 0x12, 0x06, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x19, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0d, 0x12, 0x0b, 0x2f, 0x6f, 0x70, 0x5f, 0x74, 0x6f, 0x73, 0x5f, 0x75, 0x72,
	0x69, 0xd0, 0xac, 0xc9, 0xba, 0x02, 0x07, 0x12, 0x2f, 0x0a, 0x06, 0x52, 0x65, 0x76, 0x6f, 0x6b,
	0x65, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x15, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x09, 0x12, 0x07, 0x2f, 0x72, 0x65, 0x76, 0x6f,
	0x6b, 0x65, 0xd0, 0xac, 0xc9, 0xba, 0x02, 0x08, 0x12, 0x37, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x72,
	0x6f, 0x73, 0x70, 0x65, 0x63, 0x74, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x12, 0x0b,
	0x2f, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x73, 0x70, 0x65, 0x63, 0x74, 0xd0, 0xac, 0xc9, 0xba, 0x02,
	0x09, 0x12, 0x33, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x06, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x17, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x12, 0x09, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f,
	0xd0, 0xac, 0xc9, 0xba, 0x02, 0x0a, 0x12, 0x45, 0x0a, 0x12, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x06, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1f, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x13, 0x12, 0x11, 0x2f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x65, 0x73, 0x73,
	0x5f, 0x69, 0x66, 0x72, 0x61, 0x6d, 0x65, 0xd0, 0xac, 0xc9, 0xba, 0x02, 0x0b, 0x42, 0x0a, 0x5a,
	0x08, 0x2e, 0x3b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_auth_test_proto_rawDescOnce sync.Once
	file_auth_test_proto_rawDescData = file_auth_test_proto_rawDesc
)

func file_auth_test_proto_rawDescGZIP() []byte {
	file_auth_test_proto_rawDescOnce.Do(func() {
		file_auth_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_test_proto_rawDescData)
	})
	return file_auth_test_proto_rawDescData
}

var (
	file_auth_test_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
	file_auth_test_proto_goTypes  = []interface{}{
		(*Empty)(nil), // 0: Empty
	}
)

var file_auth_test_proto_depIdxs = []int32{
	0,  // 0: AuthTest.Authorize:input_type -> Empty
	0,  // 1: AuthTest.Token:input_type -> Empty
	0,  // 2: AuthTest.Jwks:input_type -> Empty
	0,  // 3: AuthTest.Registration:input_type -> Empty
	0,  // 4: AuthTest.ServiceDocs:input_type -> Empty
	0,  // 5: AuthTest.OpPolicy:input_type -> Empty
	0,  // 6: AuthTest.OpTos:input_type -> Empty
	0,  // 7: AuthTest.Revoke:input_type -> Empty
	0,  // 8: AuthTest.Introspect:input_type -> Empty
	0,  // 9: AuthTest.UserInfo:input_type -> Empty
	0,  // 10: AuthTest.CheckSessionIframe:input_type -> Empty
	0,  // 11: AuthTest.Authorize:output_type -> Empty
	0,  // 12: AuthTest.Token:output_type -> Empty
	0,  // 13: AuthTest.Jwks:output_type -> Empty
	0,  // 14: AuthTest.Registration:output_type -> Empty
	0,  // 15: AuthTest.ServiceDocs:output_type -> Empty
	0,  // 16: AuthTest.OpPolicy:output_type -> Empty
	0,  // 17: AuthTest.OpTos:output_type -> Empty
	0,  // 18: AuthTest.Revoke:output_type -> Empty
	0,  // 19: AuthTest.Introspect:output_type -> Empty
	0,  // 20: AuthTest.UserInfo:output_type -> Empty
	0,  // 21: AuthTest.CheckSessionIframe:output_type -> Empty
	11, // [11:22] is the sub-list for method output_type
	0,  // [0:11] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_auth_test_proto_init() }
func file_auth_test_proto_init() {
	if File_auth_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_auth_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_auth_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_test_proto_goTypes,
		DependencyIndexes: file_auth_test_proto_depIdxs,
		MessageInfos:      file_auth_test_proto_msgTypes,
	}.Build()
	File_auth_test_proto = out.File
	file_auth_test_proto_rawDesc = nil
	file_auth_test_proto_goTypes = nil
	file_auth_test_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ context.Context
	_ grpc.ClientConnInterface
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AuthTestClient is the client API for AuthTest service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthTestClient interface {
	Authorize(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Token(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Jwks(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Registration(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	ServiceDocs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	OpPolicy(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	OpTos(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Revoke(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Introspect(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	UserInfo(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	CheckSessionIframe(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type authTestClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthTestClient(cc grpc.ClientConnInterface) AuthTestClient {
	return &authTestClient{cc}
}

func (c *authTestClient) Authorize(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/Authorize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTestClient) Token(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/Token", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTestClient) Jwks(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/Jwks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTestClient) Registration(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/Registration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTestClient) ServiceDocs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/ServiceDocs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTestClient) OpPolicy(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/OpPolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTestClient) OpTos(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/OpTos", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTestClient) Revoke(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/Revoke", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTestClient) Introspect(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/Introspect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTestClient) UserInfo(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/UserInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authTestClient) CheckSessionIframe(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/AuthTest/CheckSessionIframe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthTestServer is the server API for AuthTest service.
type AuthTestServer interface {
	Authorize(context.Context, *Empty) (*Empty, error)
	Token(context.Context, *Empty) (*Empty, error)
	Jwks(context.Context, *Empty) (*Empty, error)
	Registration(context.Context, *Empty) (*Empty, error)
	ServiceDocs(context.Context, *Empty) (*Empty, error)
	OpPolicy(context.Context, *Empty) (*Empty, error)
	OpTos(context.Context, *Empty) (*Empty, error)
	Revoke(context.Context, *Empty) (*Empty, error)
	Introspect(context.Context, *Empty) (*Empty, error)
	UserInfo(context.Context, *Empty) (*Empty, error)
	CheckSessionIframe(context.Context, *Empty) (*Empty, error)
}

// UnimplementedAuthTestServer can be embedded to have forward compatible implementations.
type UnimplementedAuthTestServer struct {
}

func (*UnimplementedAuthTestServer) Authorize(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}

func (*UnimplementedAuthTestServer) Token(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Token not implemented")
}

func (*UnimplementedAuthTestServer) Jwks(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Jwks not implemented")
}

func (*UnimplementedAuthTestServer) Registration(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Registration not implemented")
}

func (*UnimplementedAuthTestServer) ServiceDocs(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServiceDocs not implemented")
}

func (*UnimplementedAuthTestServer) OpPolicy(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpPolicy not implemented")
}

func (*UnimplementedAuthTestServer) OpTos(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpTos not implemented")
}

func (*UnimplementedAuthTestServer) Revoke(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Revoke not implemented")
}

func (*UnimplementedAuthTestServer) Introspect(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Introspect not implemented")
}

func (*UnimplementedAuthTestServer) UserInfo(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserInfo not implemented")
}

func (*UnimplementedAuthTestServer) CheckSessionIframe(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckSessionIframe not implemented")
}

func RegisterAuthTestServer(s *grpc.Server, srv AuthTestServer) {
	s.RegisterService(&_AuthTest_serviceDesc, srv)
}

func _AuthTest_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/Authorize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).Authorize(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest_Token_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).Token(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/Token",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).Token(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest_Jwks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).Jwks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/Jwks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).Jwks(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest_Registration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).Registration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/Registration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).Registration(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest_ServiceDocs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).ServiceDocs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/ServiceDocs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).ServiceDocs(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest_OpPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).OpPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/OpPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).OpPolicy(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest_OpTos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).OpTos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/OpTos",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).OpTos(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest_Revoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).Revoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/Revoke",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).Revoke(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest_Introspect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).Introspect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/Introspect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).Introspect(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest_UserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).UserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/UserInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).UserInfo(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthTest_CheckSessionIframe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthTestServer).CheckSessionIframe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthTest/CheckSessionIframe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthTestServer).CheckSessionIframe(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _AuthTest_serviceDesc = grpc.ServiceDesc{
	ServiceName: "AuthTest",
	HandlerType: (*AuthTestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authorize",
			Handler:    _AuthTest_Authorize_Handler,
		},
		{
			MethodName: "Token",
			Handler:    _AuthTest_Token_Handler,
		},
		{
			MethodName: "Jwks",
			Handler:    _AuthTest_Jwks_Handler,
		},
		{
			MethodName: "Registration",
			Handler:    _AuthTest_Registration_Handler,
		},
		{
			MethodName: "ServiceDocs",
			Handler:    _AuthTest_ServiceDocs_Handler,
		},
		{
			MethodName: "OpPolicy",
			Handler:    _AuthTest_OpPolicy_Handler,
		},
		{
			MethodName: "OpTos",
			Handler:    _AuthTest_OpTos_Handler,
		},
		{
			MethodName: "Revoke",
			Handler:    _AuthTest_Revoke_Handler,
		},
		{
			MethodName: "Introspect",
			Handler:    _AuthTest_Introspect_Handler,
		},
		{
			MethodName: "UserInfo",
			Handler:    _AuthTest_UserInfo_Handler,
		},
		{
			MethodName: "CheckSessionIframe",
			Handler:    _AuthTest_CheckSessionIframe_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth_test.proto",
}
