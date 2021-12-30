// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: documentation.proto

package proto

import (
	context "context"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "lab.weave.nl/devops/proto-istio-auth-generator/proto"
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

type RefType int32

const (
	RefType_UNSPECIFIED RefType = 0
	RefType_TAG         RefType = 1
	RefType_BRANCH      RefType = 2
)

// Enum value maps for RefType.
var (
	RefType_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "TAG",
		2: "BRANCH",
	}
	RefType_value = map[string]int32{
		"UNSPECIFIED": 0,
		"TAG":         1,
		"BRANCH":      2,
	}
)

func (x RefType) Enum() *RefType {
	p := new(RefType)
	*p = x
	return p
}

func (x RefType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RefType) Descriptor() protoreflect.EnumDescriptor {
	return file_documentation_proto_enumTypes[0].Descriptor()
}

func (RefType) Type() protoreflect.EnumType {
	return &file_documentation_proto_enumTypes[0]
}

func (x RefType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RefType.Descriptor instead.
func (RefType) EnumDescriptor() ([]byte, []int) {
	return file_documentation_proto_rawDescGZIP(), []int{0}
}

type GetFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FilePath    string `protobuf:"bytes,1,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"`
	Ref         string `protobuf:"bytes,2,opt,name=ref,proto3" json:"ref,omitempty"`
	ServiceName string `protobuf:"bytes,3,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
}

func (x *GetFileRequest) Reset() {
	*x = GetFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documentation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileRequest) ProtoMessage() {}

func (x *GetFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_documentation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileRequest.ProtoReflect.Descriptor instead.
func (*GetFileRequest) Descriptor() ([]byte, []int) {
	return file_documentation_proto_rawDescGZIP(), []int{0}
}

func (x *GetFileRequest) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *GetFileRequest) GetRef() string {
	if x != nil {
		return x.Ref
	}
	return ""
}

func (x *GetFileRequest) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

type GetFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content      string         `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	SwaggerFiles []*SwaggerFile `protobuf:"bytes,2,rep,name=swagger_files,json=swaggerFiles,proto3" json:"swagger_files,omitempty"`
}

func (x *GetFileResponse) Reset() {
	*x = GetFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documentation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileResponse) ProtoMessage() {}

func (x *GetFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_documentation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileResponse.ProtoReflect.Descriptor instead.
func (*GetFileResponse) Descriptor() ([]byte, []int) {
	return file_documentation_proto_rawDescGZIP(), []int{1}
}

func (x *GetFileResponse) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *GetFileResponse) GetSwaggerFiles() []*SwaggerFile {
	if x != nil {
		return x.SwaggerFiles
	}
	return nil
}

type SwaggerFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	SignedUrl string `protobuf:"bytes,2,opt,name=signed_url,json=signedUrl,proto3" json:"signed_url,omitempty"`
}

func (x *SwaggerFile) Reset() {
	*x = SwaggerFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documentation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SwaggerFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SwaggerFile) ProtoMessage() {}

func (x *SwaggerFile) ProtoReflect() protoreflect.Message {
	mi := &file_documentation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SwaggerFile.ProtoReflect.Descriptor instead.
func (*SwaggerFile) Descriptor() ([]byte, []int) {
	return file_documentation_proto_rawDescGZIP(), []int{2}
}

func (x *SwaggerFile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SwaggerFile) GetSignedUrl() string {
	if x != nil {
		return x.SignedUrl
	}
	return ""
}

type ListDirectoryFilesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FilePath string `protobuf:"bytes,1,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"`
	Ref      string `protobuf:"bytes,2,opt,name=ref,proto3" json:"ref,omitempty"`
}

func (x *ListDirectoryFilesRequest) Reset() {
	*x = ListDirectoryFilesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documentation_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListDirectoryFilesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDirectoryFilesRequest) ProtoMessage() {}

func (x *ListDirectoryFilesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_documentation_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDirectoryFilesRequest.ProtoReflect.Descriptor instead.
func (*ListDirectoryFilesRequest) Descriptor() ([]byte, []int) {
	return file_documentation_proto_rawDescGZIP(), []int{3}
}

func (x *ListDirectoryFilesRequest) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *ListDirectoryFilesRequest) GetRef() string {
	if x != nil {
		return x.Ref
	}
	return ""
}

type ListDirectoryFilesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files []*File `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *ListDirectoryFilesResponse) Reset() {
	*x = ListDirectoryFilesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documentation_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListDirectoryFilesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDirectoryFilesResponse) ProtoMessage() {}

func (x *ListDirectoryFilesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_documentation_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDirectoryFilesResponse.ProtoReflect.Descriptor instead.
func (*ListDirectoryFilesResponse) Descriptor() ([]byte, []int) {
	return file_documentation_proto_rawDescGZIP(), []int{4}
}

func (x *ListDirectoryFilesResponse) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Extension  string `protobuf:"bytes,2,opt,name=extension,proto3" json:"extension,omitempty"`
	PrettyName string `protobuf:"bytes,3,opt,name=pretty_name,json=prettyName,proto3" json:"pretty_name,omitempty"`
	Path       string `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	Type       string `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documentation_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_documentation_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_documentation_proto_rawDescGZIP(), []int{5}
}

func (x *File) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *File) GetExtension() string {
	if x != nil {
		return x.Extension
	}
	return ""
}

func (x *File) GetPrettyName() string {
	if x != nil {
		return x.PrettyName
	}
	return ""
}

func (x *File) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *File) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type ListRepositoryRefsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Refs []*Ref `protobuf:"bytes,1,rep,name=refs,proto3" json:"refs,omitempty"`
}

func (x *ListRepositoryRefsResponse) Reset() {
	*x = ListRepositoryRefsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documentation_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRepositoryRefsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRepositoryRefsResponse) ProtoMessage() {}

func (x *ListRepositoryRefsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_documentation_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRepositoryRefsResponse.ProtoReflect.Descriptor instead.
func (*ListRepositoryRefsResponse) Descriptor() ([]byte, []int) {
	return file_documentation_proto_rawDescGZIP(), []int{6}
}

func (x *ListRepositoryRefsResponse) GetRefs() []*Ref {
	if x != nil {
		return x.Refs
	}
	return nil
}

type Ref struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type RefType `protobuf:"varint,2,opt,name=type,proto3,enum=documentation.RefType" json:"type,omitempty"`
}

func (x *Ref) Reset() {
	*x = Ref{}
	if protoimpl.UnsafeEnabled {
		mi := &file_documentation_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ref) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ref) ProtoMessage() {}

func (x *Ref) ProtoReflect() protoreflect.Message {
	mi := &file_documentation_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ref.ProtoReflect.Descriptor instead.
func (*Ref) Descriptor() ([]byte, []int) {
	return file_documentation_proto_rawDescGZIP(), []int{7}
}

func (x *Ref) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Ref) GetType() RefType {
	if x != nil {
		return x.Type
	}
	return RefType_UNSPECIFIED
}

var File_documentation_proto protoreflect.FileDescriptor

var file_documentation_proto_rawDesc = []byte{
	0x0a, 0x13, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x40, 0x6c, 0x61, 0x62, 0x2e, 0x77, 0x65, 0x61, 0x76, 0x65, 0x2e, 0x6e, 0x6c, 0x2f,
	0x64, 0x65, 0x76, 0x6f, 0x70, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2d, 0x69, 0x73, 0x74,
	0x69, 0x6f, 0x2d, 0x61, 0x75, 0x74, 0x68, 0x2d, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f,
	0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x77, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27,
	0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x42, 0x03, 0x2e, 0x6d, 0x64, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x19, 0x0a, 0x03, 0x72, 0x65, 0x66, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x03, 0x72,
	0x65, 0x66, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x6c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x12, 0x3f, 0x0a, 0x0d, 0x73, 0x77, 0x61, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x77, 0x61, 0x67, 0x67, 0x65,
	0x72, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x0c, 0x73, 0x77, 0x61, 0x67, 0x67, 0x65, 0x72, 0x46, 0x69,
	0x6c, 0x65, 0x73, 0x22, 0x40, 0x0a, 0x0b, 0x53, 0x77, 0x61, 0x67, 0x67, 0x65, 0x72, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x55, 0x72, 0x6c, 0x22, 0x53, 0x0a, 0x19, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12,
	0x19, 0x0a, 0x03, 0x72, 0x65, 0x66, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x03, 0x72, 0x65, 0x66, 0x22, 0x47, 0x0a, 0x1a, 0x4c, 0x69,
	0x73, 0x74, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x05, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x22, 0x81, 0x01, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1f,
	0x0a, 0x0b, 0x70, 0x72, 0x65, 0x74, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x72, 0x65, 0x74, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x44, 0x0a, 0x1a, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x66, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x04, 0x72, 0x65, 0x66, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x66, 0x52, 0x04, 0x72, 0x65, 0x66, 0x73, 0x22, 0x45, 0x0a,
	0x03, 0x52, 0x65, 0x66, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x66, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x2a, 0x2f, 0x0a, 0x07, 0x52, 0x65, 0x66, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x07, 0x0a, 0x03, 0x54, 0x41, 0x47, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x42, 0x52, 0x41,
	0x4e, 0x43, 0x48, 0x10, 0x02, 0x32, 0x81, 0x03, 0x0a, 0x0d, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x63, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x1d, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1e, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x19, 0xca, 0x3e, 0x16, 0x67, 0x65, 0x74, 0x5f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x8a, 0x01, 0x0a,
	0x12, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x46, 0x69,
	0x6c, 0x65, 0x73, 0x12, 0x28, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x79, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e,
	0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1f, 0xca, 0x3e, 0x1c, 0x6c, 0x69, 0x73,
	0x74, 0x5f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x7e, 0x0a, 0x12, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x66, 0x73, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x29, 0x2e, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x66, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x25, 0xca, 0x3e, 0x22, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x6f, 0x72, 0x79, 0x5f, 0x72, 0x65, 0x66, 0x73, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_documentation_proto_rawDescOnce sync.Once
	file_documentation_proto_rawDescData = file_documentation_proto_rawDesc
)

func file_documentation_proto_rawDescGZIP() []byte {
	file_documentation_proto_rawDescOnce.Do(func() {
		file_documentation_proto_rawDescData = protoimpl.X.CompressGZIP(file_documentation_proto_rawDescData)
	})
	return file_documentation_proto_rawDescData
}

var file_documentation_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_documentation_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_documentation_proto_goTypes = []interface{}{
	(RefType)(0),                       // 0: documentation.RefType
	(*GetFileRequest)(nil),             // 1: documentation.GetFileRequest
	(*GetFileResponse)(nil),            // 2: documentation.GetFileResponse
	(*SwaggerFile)(nil),                // 3: documentation.SwaggerFile
	(*ListDirectoryFilesRequest)(nil),  // 4: documentation.ListDirectoryFilesRequest
	(*ListDirectoryFilesResponse)(nil), // 5: documentation.ListDirectoryFilesResponse
	(*File)(nil),                       // 6: documentation.File
	(*ListRepositoryRefsResponse)(nil), // 7: documentation.ListRepositoryRefsResponse
	(*Ref)(nil),                        // 8: documentation.Ref
	(*empty.Empty)(nil),                // 9: google.protobuf.Empty
}
var file_documentation_proto_depIdxs = []int32{
	3, // 0: documentation.GetFileResponse.swagger_files:type_name -> documentation.SwaggerFile
	6, // 1: documentation.ListDirectoryFilesResponse.files:type_name -> documentation.File
	8, // 2: documentation.ListRepositoryRefsResponse.refs:type_name -> documentation.Ref
	0, // 3: documentation.Ref.type:type_name -> documentation.RefType
	1, // 4: documentation.Documentation.GetFile:input_type -> documentation.GetFileRequest
	4, // 5: documentation.Documentation.ListDirectoryFiles:input_type -> documentation.ListDirectoryFilesRequest
	9, // 6: documentation.Documentation.ListRepositoryRefs:input_type -> google.protobuf.Empty
	2, // 7: documentation.Documentation.GetFile:output_type -> documentation.GetFileResponse
	5, // 8: documentation.Documentation.ListDirectoryFiles:output_type -> documentation.ListDirectoryFilesResponse
	7, // 9: documentation.Documentation.ListRepositoryRefs:output_type -> documentation.ListRepositoryRefsResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_documentation_proto_init() }
func file_documentation_proto_init() {
	if File_documentation_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_documentation_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileRequest); i {
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
		file_documentation_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileResponse); i {
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
		file_documentation_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SwaggerFile); i {
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
		file_documentation_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListDirectoryFilesRequest); i {
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
		file_documentation_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListDirectoryFilesResponse); i {
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
		file_documentation_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*File); i {
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
		file_documentation_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRepositoryRefsResponse); i {
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
		file_documentation_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ref); i {
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
			RawDescriptor: file_documentation_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_documentation_proto_goTypes,
		DependencyIndexes: file_documentation_proto_depIdxs,
		EnumInfos:         file_documentation_proto_enumTypes,
		MessageInfos:      file_documentation_proto_msgTypes,
	}.Build()
	File_documentation_proto = out.File
	file_documentation_proto_rawDesc = nil
	file_documentation_proto_goTypes = nil
	file_documentation_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DocumentationClient is the client API for Documentation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DocumentationClient interface {
	// GetFile
	//
	// Get markdown file on given file path
	GetFile(ctx context.Context, in *GetFileRequest, opts ...grpc.CallOption) (*GetFileResponse, error)
	// ListDirectoryFiles
	//
	// List markdown files in given directory
	ListDirectoryFiles(ctx context.Context, in *ListDirectoryFilesRequest, opts ...grpc.CallOption) (*ListDirectoryFilesResponse, error)
	// ListRepositoryRefs
	//
	// List refs of the given repository
	ListRepositoryRefs(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ListRepositoryRefsResponse, error)
}

type documentationClient struct {
	cc grpc.ClientConnInterface
}

func NewDocumentationClient(cc grpc.ClientConnInterface) DocumentationClient {
	return &documentationClient{cc}
}

func (c *documentationClient) GetFile(ctx context.Context, in *GetFileRequest, opts ...grpc.CallOption) (*GetFileResponse, error) {
	out := new(GetFileResponse)
	err := c.cc.Invoke(ctx, "/documentation.Documentation/GetFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *documentationClient) ListDirectoryFiles(ctx context.Context, in *ListDirectoryFilesRequest, opts ...grpc.CallOption) (*ListDirectoryFilesResponse, error) {
	out := new(ListDirectoryFilesResponse)
	err := c.cc.Invoke(ctx, "/documentation.Documentation/ListDirectoryFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *documentationClient) ListRepositoryRefs(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ListRepositoryRefsResponse, error) {
	out := new(ListRepositoryRefsResponse)
	err := c.cc.Invoke(ctx, "/documentation.Documentation/ListRepositoryRefs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DocumentationServer is the server API for Documentation service.
type DocumentationServer interface {
	// GetFile
	//
	// Get markdown file on given file path
	GetFile(context.Context, *GetFileRequest) (*GetFileResponse, error)
	// ListDirectoryFiles
	//
	// List markdown files in given directory
	ListDirectoryFiles(context.Context, *ListDirectoryFilesRequest) (*ListDirectoryFilesResponse, error)
	// ListRepositoryRefs
	//
	// List refs of the given repository
	ListRepositoryRefs(context.Context, *empty.Empty) (*ListRepositoryRefsResponse, error)
}

// UnimplementedDocumentationServer can be embedded to have forward compatible implementations.
type UnimplementedDocumentationServer struct {
}

func (*UnimplementedDocumentationServer) GetFile(context.Context, *GetFileRequest) (*GetFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFile not implemented")
}
func (*UnimplementedDocumentationServer) ListDirectoryFiles(context.Context, *ListDirectoryFilesRequest) (*ListDirectoryFilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDirectoryFiles not implemented")
}
func (*UnimplementedDocumentationServer) ListRepositoryRefs(context.Context, *empty.Empty) (*ListRepositoryRefsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRepositoryRefs not implemented")
}

func RegisterDocumentationServer(s *grpc.Server, srv DocumentationServer) {
	s.RegisterService(&_Documentation_serviceDesc, srv)
}

func _Documentation_GetFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DocumentationServer).GetFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/documentation.Documentation/GetFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DocumentationServer).GetFile(ctx, req.(*GetFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Documentation_ListDirectoryFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDirectoryFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DocumentationServer).ListDirectoryFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/documentation.Documentation/ListDirectoryFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DocumentationServer).ListDirectoryFiles(ctx, req.(*ListDirectoryFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Documentation_ListRepositoryRefs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DocumentationServer).ListRepositoryRefs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/documentation.Documentation/ListRepositoryRefs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DocumentationServer).ListRepositoryRefs(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Documentation_serviceDesc = grpc.ServiceDesc{
	ServiceName: "documentation.Documentation",
	HandlerType: (*DocumentationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFile",
			Handler:    _Documentation_GetFile_Handler,
		},
		{
			MethodName: "ListDirectoryFiles",
			Handler:    _Documentation_ListDirectoryFiles_Handler,
		},
		{
			MethodName: "ListRepositoryRefs",
			Handler:    _Documentation_ListRepositoryRefs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "documentation.proto",
}