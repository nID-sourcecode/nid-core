// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: verification.proto

package proto

import (
	context "context"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type RetryPhoneRequest_PhoneNumberVerificationType int32

const (
	RetryPhoneRequest_PHONE_NUMBER_VERIFICATION_TYPE_UNSPECIFIED RetryPhoneRequest_PhoneNumberVerificationType = 0
	RetryPhoneRequest_SMS                                        RetryPhoneRequest_PhoneNumberVerificationType = 1
	RetryPhoneRequest_TTS                                        RetryPhoneRequest_PhoneNumberVerificationType = 2
)

// Enum value maps for RetryPhoneRequest_PhoneNumberVerificationType.
var (
	RetryPhoneRequest_PhoneNumberVerificationType_name = map[int32]string{
		0: "PHONE_NUMBER_VERIFICATION_TYPE_UNSPECIFIED",
		1: "SMS",
		2: "TTS",
	}
	RetryPhoneRequest_PhoneNumberVerificationType_value = map[string]int32{
		"PHONE_NUMBER_VERIFICATION_TYPE_UNSPECIFIED": 0,
		"SMS": 1,
		"TTS": 2,
	}
)

func (x RetryPhoneRequest_PhoneNumberVerificationType) Enum() *RetryPhoneRequest_PhoneNumberVerificationType {
	p := new(RetryPhoneRequest_PhoneNumberVerificationType)
	*p = x
	return p
}

func (x RetryPhoneRequest_PhoneNumberVerificationType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RetryPhoneRequest_PhoneNumberVerificationType) Descriptor() protoreflect.EnumDescriptor {
	return file_verification_proto_enumTypes[0].Descriptor()
}

func (RetryPhoneRequest_PhoneNumberVerificationType) Type() protoreflect.EnumType {
	return &file_verification_proto_enumTypes[0]
}

func (x RetryPhoneRequest_PhoneNumberVerificationType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RetryPhoneRequest_PhoneNumberVerificationType.Descriptor instead.
func (RetryPhoneRequest_PhoneNumberVerificationType) EnumDescriptor() ([]byte, []int) {
	return file_verification_proto_rawDescGZIP(), []int{3, 0}
}

type VerifyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Code string `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *VerifyRequest) Reset() {
	*x = VerifyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verification_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyRequest) ProtoMessage() {}

func (x *VerifyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_verification_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyRequest.ProtoReflect.Descriptor instead.
func (*VerifyRequest) Descriptor() ([]byte, []int) {
	return file_verification_proto_rawDescGZIP(), []int{0}
}

func (x *VerifyRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *VerifyRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type RetryVerifyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RetryVerifyRequest) Reset() {
	*x = RetryVerifyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verification_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetryVerifyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetryVerifyRequest) ProtoMessage() {}

func (x *RetryVerifyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_verification_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetryVerifyRequest.ProtoReflect.Descriptor instead.
func (*RetryVerifyRequest) Descriptor() ([]byte, []int) {
	return file_verification_proto_rawDescGZIP(), []int{1}
}

func (x *RetryVerifyRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type VerifyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *VerifyResponse) Reset() {
	*x = VerifyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verification_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyResponse) ProtoMessage() {}

func (x *VerifyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_verification_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyResponse.ProtoReflect.Descriptor instead.
func (*VerifyResponse) Descriptor() ([]byte, []int) {
	return file_verification_proto_rawDescGZIP(), []int{2}
}

func (x *VerifyResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type RetryPhoneRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               string                                        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	VerificationType RetryPhoneRequest_PhoneNumberVerificationType `protobuf:"varint,2,opt,name=verification_type,json=verificationType,proto3,enum=wallet.RetryPhoneRequest_PhoneNumberVerificationType" json:"verification_type,omitempty"`
}

func (x *RetryPhoneRequest) Reset() {
	*x = RetryPhoneRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verification_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetryPhoneRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetryPhoneRequest) ProtoMessage() {}

func (x *RetryPhoneRequest) ProtoReflect() protoreflect.Message {
	mi := &file_verification_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetryPhoneRequest.ProtoReflect.Descriptor instead.
func (*RetryPhoneRequest) Descriptor() ([]byte, []int) {
	return file_verification_proto_rawDescGZIP(), []int{3}
}

func (x *RetryPhoneRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RetryPhoneRequest) GetVerificationType() RetryPhoneRequest_PhoneNumberVerificationType {
	if x != nil {
		return x.VerificationType
	}
	return RetryPhoneRequest_PHONE_NUMBER_VERIFICATION_TYPE_UNSPECIFIED
}

var File_verification_proto protoreflect.FileDescriptor

var file_verification_proto_rawDesc = []byte{
	0x0a, 0x12, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3d, 0x0a, 0x0d, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x22, 0x2e, 0x0a, 0x12, 0x52, 0x65, 0x74, 0x72, 0x79, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x20, 0x0a, 0x0e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0xf2, 0x01, 0x0a, 0x11, 0x52, 0x65, 0x74, 0x72, 0x79, 0x50, 0x68,
	0x6f, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x62, 0x0a, 0x11, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x35, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x79, 0x50, 0x68,
	0x6f, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x68, 0x6f, 0x6e, 0x65,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x10, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x22, 0x5f, 0x0a, 0x1b, 0x50, 0x68, 0x6f, 0x6e,
	0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2e, 0x0a, 0x2a, 0x50, 0x48, 0x4f, 0x4e, 0x45,
	0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x49, 0x43, 0x41,
	0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x53, 0x4d, 0x53, 0x10, 0x01,
	0x12, 0x07, 0x0a, 0x03, 0x54, 0x54, 0x53, 0x10, 0x02, 0x32, 0x9f, 0x03, 0x0a, 0x0c, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x56, 0x0a, 0x0b, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x15, 0x2e, 0x77, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12,
	0x22, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x2d, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x12, 0x66, 0x0a, 0x10, 0x52, 0x65, 0x74, 0x72, 0x79, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e,
	0x52, 0x65, 0x74, 0x72, 0x79, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x18, 0x22, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x74, 0x72, 0x79, 0x2d, 0x76, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x2d, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x62, 0x0a, 0x11, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x15, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e,
	0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x22, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x74, 0x72,
	0x79, 0x2d, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x2d, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x6b,
	0x0a, 0x16, 0x52, 0x65, 0x74, 0x72, 0x79, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x50, 0x68, 0x6f,
	0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x79, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x18, 0x22, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x74, 0x72, 0x79, 0x2d, 0x76,
	0x65, 0x72, 0x69, 0x66, 0x79, 0x2d, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x09, 0x5a, 0x07, 0x2e,
	0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_verification_proto_rawDescOnce sync.Once
	file_verification_proto_rawDescData = file_verification_proto_rawDesc
)

func file_verification_proto_rawDescGZIP() []byte {
	file_verification_proto_rawDescOnce.Do(func() {
		file_verification_proto_rawDescData = protoimpl.X.CompressGZIP(file_verification_proto_rawDescData)
	})
	return file_verification_proto_rawDescData
}

var file_verification_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_verification_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_verification_proto_goTypes = []interface{}{
	(RetryPhoneRequest_PhoneNumberVerificationType)(0), // 0: wallet.RetryPhoneRequest.PhoneNumberVerificationType
	(*VerifyRequest)(nil),                              // 1: wallet.VerifyRequest
	(*RetryVerifyRequest)(nil),                         // 2: wallet.RetryVerifyRequest
	(*VerifyResponse)(nil),                             // 3: wallet.VerifyResponse
	(*RetryPhoneRequest)(nil),                          // 4: wallet.RetryPhoneRequest
}
var file_verification_proto_depIdxs = []int32{
	0, // 0: wallet.RetryPhoneRequest.verification_type:type_name -> wallet.RetryPhoneRequest.PhoneNumberVerificationType
	1, // 1: wallet.Verification.VerifyEmail:input_type -> wallet.VerifyRequest
	2, // 2: wallet.Verification.RetryVerifyEmail:input_type -> wallet.RetryVerifyRequest
	1, // 3: wallet.Verification.VerifyPhoneNumber:input_type -> wallet.VerifyRequest
	4, // 4: wallet.Verification.RetryVerifyPhoneNumber:input_type -> wallet.RetryPhoneRequest
	3, // 5: wallet.Verification.VerifyEmail:output_type -> wallet.VerifyResponse
	3, // 6: wallet.Verification.RetryVerifyEmail:output_type -> wallet.VerifyResponse
	3, // 7: wallet.Verification.VerifyPhoneNumber:output_type -> wallet.VerifyResponse
	3, // 8: wallet.Verification.RetryVerifyPhoneNumber:output_type -> wallet.VerifyResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_verification_proto_init() }
func file_verification_proto_init() {
	if File_verification_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_verification_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyRequest); i {
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
		file_verification_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetryVerifyRequest); i {
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
		file_verification_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyResponse); i {
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
		file_verification_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetryPhoneRequest); i {
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
			RawDescriptor: file_verification_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_verification_proto_goTypes,
		DependencyIndexes: file_verification_proto_depIdxs,
		EnumInfos:         file_verification_proto_enumTypes,
		MessageInfos:      file_verification_proto_msgTypes,
	}.Build()
	File_verification_proto = out.File
	file_verification_proto_rawDesc = nil
	file_verification_proto_goTypes = nil
	file_verification_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// VerificationClient is the client API for Verification service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VerificationClient interface {
	// Will verify emailaddress with a token
	VerifyEmail(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*VerifyResponse, error)
	// Will retry verification process
	RetryVerifyEmail(ctx context.Context, in *RetryVerifyRequest, opts ...grpc.CallOption) (*VerifyResponse, error)
	// Will retry verification process
	VerifyPhoneNumber(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*VerifyResponse, error)
	// Will retry verification process
	RetryVerifyPhoneNumber(ctx context.Context, in *RetryPhoneRequest, opts ...grpc.CallOption) (*VerifyResponse, error)
}

type verificationClient struct {
	cc grpc.ClientConnInterface
}

func NewVerificationClient(cc grpc.ClientConnInterface) VerificationClient {
	return &verificationClient{cc}
}

func (c *verificationClient) VerifyEmail(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*VerifyResponse, error) {
	out := new(VerifyResponse)
	err := c.cc.Invoke(ctx, "/wallet.Verification/VerifyEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *verificationClient) RetryVerifyEmail(ctx context.Context, in *RetryVerifyRequest, opts ...grpc.CallOption) (*VerifyResponse, error) {
	out := new(VerifyResponse)
	err := c.cc.Invoke(ctx, "/wallet.Verification/RetryVerifyEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *verificationClient) VerifyPhoneNumber(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*VerifyResponse, error) {
	out := new(VerifyResponse)
	err := c.cc.Invoke(ctx, "/wallet.Verification/VerifyPhoneNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *verificationClient) RetryVerifyPhoneNumber(ctx context.Context, in *RetryPhoneRequest, opts ...grpc.CallOption) (*VerifyResponse, error) {
	out := new(VerifyResponse)
	err := c.cc.Invoke(ctx, "/wallet.Verification/RetryVerifyPhoneNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VerificationServer is the server API for Verification service.
type VerificationServer interface {
	// Will verify emailaddress with a token
	VerifyEmail(context.Context, *VerifyRequest) (*VerifyResponse, error)
	// Will retry verification process
	RetryVerifyEmail(context.Context, *RetryVerifyRequest) (*VerifyResponse, error)
	// Will retry verification process
	VerifyPhoneNumber(context.Context, *VerifyRequest) (*VerifyResponse, error)
	// Will retry verification process
	RetryVerifyPhoneNumber(context.Context, *RetryPhoneRequest) (*VerifyResponse, error)
}

// UnimplementedVerificationServer can be embedded to have forward compatible implementations.
type UnimplementedVerificationServer struct {
}

func (*UnimplementedVerificationServer) VerifyEmail(context.Context, *VerifyRequest) (*VerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyEmail not implemented")
}
func (*UnimplementedVerificationServer) RetryVerifyEmail(context.Context, *RetryVerifyRequest) (*VerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetryVerifyEmail not implemented")
}
func (*UnimplementedVerificationServer) VerifyPhoneNumber(context.Context, *VerifyRequest) (*VerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyPhoneNumber not implemented")
}
func (*UnimplementedVerificationServer) RetryVerifyPhoneNumber(context.Context, *RetryPhoneRequest) (*VerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetryVerifyPhoneNumber not implemented")
}

func RegisterVerificationServer(s *grpc.Server, srv VerificationServer) {
	s.RegisterService(&_Verification_serviceDesc, srv)
}

func _Verification_VerifyEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerificationServer).VerifyEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.Verification/VerifyEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerificationServer).VerifyEmail(ctx, req.(*VerifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Verification_RetryVerifyEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetryVerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerificationServer).RetryVerifyEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.Verification/RetryVerifyEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerificationServer).RetryVerifyEmail(ctx, req.(*RetryVerifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Verification_VerifyPhoneNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerificationServer).VerifyPhoneNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.Verification/VerifyPhoneNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerificationServer).VerifyPhoneNumber(ctx, req.(*VerifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Verification_RetryVerifyPhoneNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetryPhoneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerificationServer).RetryVerifyPhoneNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.Verification/RetryVerifyPhoneNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerificationServer).RetryVerifyPhoneNumber(ctx, req.(*RetryPhoneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Verification_serviceDesc = grpc.ServiceDesc{
	ServiceName: "wallet.Verification",
	HandlerType: (*VerificationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VerifyEmail",
			Handler:    _Verification_VerifyEmail_Handler,
		},
		{
			MethodName: "RetryVerifyEmail",
			Handler:    _Verification_RetryVerifyEmail_Handler,
		},
		{
			MethodName: "VerifyPhoneNumber",
			Handler:    _Verification_VerifyPhoneNumber_Handler,
		},
		{
			MethodName: "RetryVerifyPhoneNumber",
			Handler:    _Verification_RetryVerifyPhoneNumber_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "verification.proto",
}
