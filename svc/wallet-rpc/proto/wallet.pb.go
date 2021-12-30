// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: wallet.proto

package proto

import (
	context "context"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type CreateConsentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AccessToken string               `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	ClientId    string               `protobuf:"bytes,3,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Description string               `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Name        string               `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	UserPseudo  string               `protobuf:"bytes,6,opt,name=user_pseudo,json=userPseudo,proto3" json:"user_pseudo,omitempty"`
	GrantedAt   *timestamp.Timestamp `protobuf:"bytes,7,opt,name=granted_at,json=grantedAt,proto3" json:"granted_at,omitempty"`
}

func (x *CreateConsentRequest) Reset() {
	*x = CreateConsentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateConsentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateConsentRequest) ProtoMessage() {}

func (x *CreateConsentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateConsentRequest.ProtoReflect.Descriptor instead.
func (*CreateConsentRequest) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{0}
}

func (x *CreateConsentRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CreateConsentRequest) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *CreateConsentRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *CreateConsentRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateConsentRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateConsentRequest) GetUserPseudo() string {
	if x != nil {
		return x.UserPseudo
	}
	return ""
}

func (x *CreateConsentRequest) GetGrantedAt() *timestamp.Timestamp {
	if x != nil {
		return x.GrantedAt
	}
	return nil
}

type ConsentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AccessToken string               `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	ClientId    string               `protobuf:"bytes,3,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Description string               `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Name        string               `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	UserId      string               `protobuf:"bytes,6,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	GrantedAt   *timestamp.Timestamp `protobuf:"bytes,7,opt,name=granted_at,json=grantedAt,proto3" json:"granted_at,omitempty"`
}

func (x *ConsentResponse) Reset() {
	*x = ConsentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsentResponse) ProtoMessage() {}

func (x *ConsentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsentResponse.ProtoReflect.Descriptor instead.
func (*ConsentResponse) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{1}
}

func (x *ConsentResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ConsentResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *ConsentResponse) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *ConsentResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ConsentResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ConsentResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ConsentResponse) GetGrantedAt() *timestamp.Timestamp {
	if x != nil {
		return x.GrantedAt
	}
	return nil
}

type GetBSNForPseudonymRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pseudonym string `protobuf:"bytes,1,opt,name=pseudonym,proto3" json:"pseudonym,omitempty"`
}

func (x *GetBSNForPseudonymRequest) Reset() {
	*x = GetBSNForPseudonymRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBSNForPseudonymRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBSNForPseudonymRequest) ProtoMessage() {}

func (x *GetBSNForPseudonymRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBSNForPseudonymRequest.ProtoReflect.Descriptor instead.
func (*GetBSNForPseudonymRequest) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{2}
}

func (x *GetBSNForPseudonymRequest) GetPseudonym() string {
	if x != nil {
		return x.Pseudonym
	}
	return ""
}

type GetBSNForPseudonymResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bsn string `protobuf:"bytes,1,opt,name=bsn,proto3" json:"bsn,omitempty"`
}

func (x *GetBSNForPseudonymResponse) Reset() {
	*x = GetBSNForPseudonymResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBSNForPseudonymResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBSNForPseudonymResponse) ProtoMessage() {}

func (x *GetBSNForPseudonymResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBSNForPseudonymResponse.ProtoReflect.Descriptor instead.
func (*GetBSNForPseudonymResponse) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{3}
}

func (x *GetBSNForPseudonymResponse) GetBsn() string {
	if x != nil {
		return x.Bsn
	}
	return ""
}

var File_wallet_proto protoreflect.FileDescriptor

var file_wallet_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8c,
	0x02, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x25, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01,
	0x01, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x70, 0x73, 0x65, 0x75, 0x64, 0x6f,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x50, 0x73, 0x65, 0x75,
	0x64, 0x6f, 0x12, 0x39, 0x0a, 0x0a, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0xeb, 0x01,
	0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x39, 0x0a, 0x0a, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x39, 0x0a, 0x19, 0x47,
	0x65, 0x74, 0x42, 0x53, 0x4e, 0x46, 0x6f, 0x72, 0x50, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x73, 0x65, 0x75,
	0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x73, 0x65,
	0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x22, 0x2e, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x42, 0x53, 0x4e,
	0x46, 0x6f, 0x72, 0x50, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x62, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x62, 0x73, 0x6e, 0x32, 0xc4, 0x01, 0x0a, 0x06, 0x57, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x12, 0x5b, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x73, 0x65,
	0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x17, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x13, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0d, 0x22, 0x0b, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x74, 0x12, 0x5d,
	0x0a, 0x12, 0x47, 0x65, 0x74, 0x42, 0x53, 0x4e, 0x46, 0x6f, 0x72, 0x50, 0x73, 0x65, 0x75, 0x64,
	0x6f, 0x6e, 0x79, 0x6d, 0x12, 0x21, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x47, 0x65,
	0x74, 0x42, 0x53, 0x4e, 0x46, 0x6f, 0x72, 0x50, 0x73, 0x65, 0x75, 0x64, 0x6f, 0x6e, 0x79, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x2e, 0x47, 0x65, 0x74, 0x42, 0x53, 0x4e, 0x46, 0x6f, 0x72, 0x50, 0x73, 0x65, 0x75, 0x64, 0x6f,
	0x6e, 0x79, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x09, 0x5a,
	0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_wallet_proto_rawDescOnce sync.Once
	file_wallet_proto_rawDescData = file_wallet_proto_rawDesc
)

func file_wallet_proto_rawDescGZIP() []byte {
	file_wallet_proto_rawDescOnce.Do(func() {
		file_wallet_proto_rawDescData = protoimpl.X.CompressGZIP(file_wallet_proto_rawDescData)
	})
	return file_wallet_proto_rawDescData
}

var file_wallet_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_wallet_proto_goTypes = []interface{}{
	(*CreateConsentRequest)(nil),       // 0: wallet.CreateConsentRequest
	(*ConsentResponse)(nil),            // 1: wallet.ConsentResponse
	(*GetBSNForPseudonymRequest)(nil),  // 2: wallet.GetBSNForPseudonymRequest
	(*GetBSNForPseudonymResponse)(nil), // 3: wallet.GetBSNForPseudonymResponse
	(*timestamp.Timestamp)(nil),        // 4: google.protobuf.Timestamp
}
var file_wallet_proto_depIdxs = []int32{
	4, // 0: wallet.CreateConsentRequest.granted_at:type_name -> google.protobuf.Timestamp
	4, // 1: wallet.ConsentResponse.granted_at:type_name -> google.protobuf.Timestamp
	0, // 2: wallet.Wallet.CreateConsent:input_type -> wallet.CreateConsentRequest
	2, // 3: wallet.Wallet.GetBSNForPseudonym:input_type -> wallet.GetBSNForPseudonymRequest
	1, // 4: wallet.Wallet.CreateConsent:output_type -> wallet.ConsentResponse
	3, // 5: wallet.Wallet.GetBSNForPseudonym:output_type -> wallet.GetBSNForPseudonymResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_wallet_proto_init() }
func file_wallet_proto_init() {
	if File_wallet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_wallet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateConsentRequest); i {
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
		file_wallet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsentResponse); i {
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
		file_wallet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBSNForPseudonymRequest); i {
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
		file_wallet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBSNForPseudonymResponse); i {
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
			RawDescriptor: file_wallet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_wallet_proto_goTypes,
		DependencyIndexes: file_wallet_proto_depIdxs,
		MessageInfos:      file_wallet_proto_msgTypes,
	}.Build()
	File_wallet_proto = out.File
	file_wallet_proto_rawDesc = nil
	file_wallet_proto_goTypes = nil
	file_wallet_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// WalletClient is the client API for Wallet service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WalletClient interface {
	// CreateConsent
	//
	// Will create a new consent
	CreateConsent(ctx context.Context, in *CreateConsentRequest, opts ...grpc.CallOption) (*ConsentResponse, error)
	GetBSNForPseudonym(ctx context.Context, in *GetBSNForPseudonymRequest, opts ...grpc.CallOption) (*GetBSNForPseudonymResponse, error)
}

type walletClient struct {
	cc grpc.ClientConnInterface
}

func NewWalletClient(cc grpc.ClientConnInterface) WalletClient {
	return &walletClient{cc}
}

func (c *walletClient) CreateConsent(ctx context.Context, in *CreateConsentRequest, opts ...grpc.CallOption) (*ConsentResponse, error) {
	out := new(ConsentResponse)
	err := c.cc.Invoke(ctx, "/wallet.Wallet/CreateConsent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) GetBSNForPseudonym(ctx context.Context, in *GetBSNForPseudonymRequest, opts ...grpc.CallOption) (*GetBSNForPseudonymResponse, error) {
	out := new(GetBSNForPseudonymResponse)
	err := c.cc.Invoke(ctx, "/wallet.Wallet/GetBSNForPseudonym", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WalletServer is the server API for Wallet service.
type WalletServer interface {
	// CreateConsent
	//
	// Will create a new consent
	CreateConsent(context.Context, *CreateConsentRequest) (*ConsentResponse, error)
	GetBSNForPseudonym(context.Context, *GetBSNForPseudonymRequest) (*GetBSNForPseudonymResponse, error)
}

// UnimplementedWalletServer can be embedded to have forward compatible implementations.
type UnimplementedWalletServer struct {
}

func (*UnimplementedWalletServer) CreateConsent(context.Context, *CreateConsentRequest) (*ConsentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateConsent not implemented")
}
func (*UnimplementedWalletServer) GetBSNForPseudonym(context.Context, *GetBSNForPseudonymRequest) (*GetBSNForPseudonymResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBSNForPseudonym not implemented")
}

func RegisterWalletServer(s *grpc.Server, srv WalletServer) {
	s.RegisterService(&_Wallet_serviceDesc, srv)
}

func _Wallet_CreateConsent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateConsentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).CreateConsent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.Wallet/CreateConsent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).CreateConsent(ctx, req.(*CreateConsentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_GetBSNForPseudonym_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBSNForPseudonymRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).GetBSNForPseudonym(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wallet.Wallet/GetBSNForPseudonym",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).GetBSNForPseudonym(ctx, req.(*GetBSNForPseudonymRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Wallet_serviceDesc = grpc.ServiceDesc{
	ServiceName: "wallet.Wallet",
	HandlerType: (*WalletServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateConsent",
			Handler:    _Wallet_CreateConsent_Handler,
		},
		{
			MethodName: "GetBSNForPseudonym",
			Handler:    _Wallet_GetBSNForPseudonym_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wallet.proto",
}
