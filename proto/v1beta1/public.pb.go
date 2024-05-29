// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: gotocompany/shield/v1beta1/public.proto

package shieldv1beta1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type ServiceDataKeyRequestBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Project     string `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	Key         string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *ServiceDataKeyRequestBody) Reset() {
	*x = ServiceDataKeyRequestBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceDataKeyRequestBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceDataKeyRequestBody) ProtoMessage() {}

func (x *ServiceDataKeyRequestBody) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceDataKeyRequestBody.ProtoReflect.Descriptor instead.
func (*ServiceDataKeyRequestBody) Descriptor() ([]byte, []int) {
	return file_gotocompany_shield_v1beta1_public_proto_rawDescGZIP(), []int{0}
}

func (x *ServiceDataKeyRequestBody) GetProject() string {
	if x != nil {
		return x.Project
	}
	return ""
}

func (x *ServiceDataKeyRequestBody) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ServiceDataKeyRequestBody) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type ServiceDataKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urn string `protobuf:"bytes,1,opt,name=urn,proto3" json:"urn,omitempty"`
	Id  string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ServiceDataKey) Reset() {
	*x = ServiceDataKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceDataKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceDataKey) ProtoMessage() {}

func (x *ServiceDataKey) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceDataKey.ProtoReflect.Descriptor instead.
func (*ServiceDataKey) Descriptor() ([]byte, []int) {
	return file_gotocompany_shield_v1beta1_public_proto_rawDescGZIP(), []int{1}
}

func (x *ServiceDataKey) GetUrn() string {
	if x != nil {
		return x.Urn
	}
	return ""
}

func (x *ServiceDataKey) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CreateServiceDataKeyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body *ServiceDataKeyRequestBody `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *CreateServiceDataKeyRequest) Reset() {
	*x = CreateServiceDataKeyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateServiceDataKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateServiceDataKeyRequest) ProtoMessage() {}

func (x *CreateServiceDataKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateServiceDataKeyRequest.ProtoReflect.Descriptor instead.
func (*CreateServiceDataKeyRequest) Descriptor() ([]byte, []int) {
	return file_gotocompany_shield_v1beta1_public_proto_rawDescGZIP(), []int{2}
}

func (x *CreateServiceDataKeyRequest) GetBody() *ServiceDataKeyRequestBody {
	if x != nil {
		return x.Body
	}
	return nil
}

type CreateServiceDataKeyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceDataKey *ServiceDataKey `protobuf:"bytes,1,opt,name=service_data_key,json=serviceDataKey,proto3" json:"service_data_key,omitempty"`
}

func (x *CreateServiceDataKeyResponse) Reset() {
	*x = CreateServiceDataKeyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateServiceDataKeyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateServiceDataKeyResponse) ProtoMessage() {}

func (x *CreateServiceDataKeyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateServiceDataKeyResponse.ProtoReflect.Descriptor instead.
func (*CreateServiceDataKeyResponse) Descriptor() ([]byte, []int) {
	return file_gotocompany_shield_v1beta1_public_proto_rawDescGZIP(), []int{3}
}

func (x *CreateServiceDataKeyResponse) GetServiceDataKey() *ServiceDataKey {
	if x != nil {
		return x.ServiceDataKey
	}
	return nil
}

type UpsertServiceDataRequestBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Project string            `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	Data    map[string]string `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *UpsertServiceDataRequestBody) Reset() {
	*x = UpsertServiceDataRequestBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertServiceDataRequestBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertServiceDataRequestBody) ProtoMessage() {}

func (x *UpsertServiceDataRequestBody) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertServiceDataRequestBody.ProtoReflect.Descriptor instead.
func (*UpsertServiceDataRequestBody) Descriptor() ([]byte, []int) {
	return file_gotocompany_shield_v1beta1_public_proto_rawDescGZIP(), []int{4}
}

func (x *UpsertServiceDataRequestBody) GetProject() string {
	if x != nil {
		return x.Project
	}
	return ""
}

func (x *UpsertServiceDataRequestBody) GetData() map[string]string {
	if x != nil {
		return x.Data
	}
	return nil
}

type UpsertUserServiceDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string                        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Body *UpsertServiceDataRequestBody `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *UpsertUserServiceDataRequest) Reset() {
	*x = UpsertUserServiceDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertUserServiceDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertUserServiceDataRequest) ProtoMessage() {}

func (x *UpsertUserServiceDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertUserServiceDataRequest.ProtoReflect.Descriptor instead.
func (*UpsertUserServiceDataRequest) Descriptor() ([]byte, []int) {
	return file_gotocompany_shield_v1beta1_public_proto_rawDescGZIP(), []int{5}
}

func (x *UpsertUserServiceDataRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpsertUserServiceDataRequest) GetBody() *UpsertServiceDataRequestBody {
	if x != nil {
		return x.Body
	}
	return nil
}

type UpsertGroupServiceDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string                        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Body *UpsertServiceDataRequestBody `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *UpsertGroupServiceDataRequest) Reset() {
	*x = UpsertGroupServiceDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertGroupServiceDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertGroupServiceDataRequest) ProtoMessage() {}

func (x *UpsertGroupServiceDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertGroupServiceDataRequest.ProtoReflect.Descriptor instead.
func (*UpsertGroupServiceDataRequest) Descriptor() ([]byte, []int) {
	return file_gotocompany_shield_v1beta1_public_proto_rawDescGZIP(), []int{6}
}

func (x *UpsertGroupServiceDataRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpsertGroupServiceDataRequest) GetBody() *UpsertServiceDataRequestBody {
	if x != nil {
		return x.Body
	}
	return nil
}

type UpsertUserServiceDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urn string `protobuf:"bytes,1,opt,name=urn,proto3" json:"urn,omitempty"`
}

func (x *UpsertUserServiceDataResponse) Reset() {
	*x = UpsertUserServiceDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertUserServiceDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertUserServiceDataResponse) ProtoMessage() {}

func (x *UpsertUserServiceDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertUserServiceDataResponse.ProtoReflect.Descriptor instead.
func (*UpsertUserServiceDataResponse) Descriptor() ([]byte, []int) {
	return file_gotocompany_shield_v1beta1_public_proto_rawDescGZIP(), []int{7}
}

func (x *UpsertUserServiceDataResponse) GetUrn() string {
	if x != nil {
		return x.Urn
	}
	return ""
}

type UpsertGroupServiceDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urn string `protobuf:"bytes,1,opt,name=urn,proto3" json:"urn,omitempty"`
}

func (x *UpsertGroupServiceDataResponse) Reset() {
	*x = UpsertGroupServiceDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertGroupServiceDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertGroupServiceDataResponse) ProtoMessage() {}

func (x *UpsertGroupServiceDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gotocompany_shield_v1beta1_public_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertGroupServiceDataResponse.ProtoReflect.Descriptor instead.
func (*UpsertGroupServiceDataResponse) Descriptor() ([]byte, []int) {
	return file_gotocompany_shield_v1beta1_public_proto_rawDescGZIP(), []int{8}
}

func (x *UpsertGroupServiceDataResponse) GetUrn() string {
	if x != nil {
		return x.Urn
	}
	return ""
}

var File_gotocompany_shield_v1beta1_public_proto protoreflect.FileDescriptor

var file_gotocompany_shield_v1beta1_public_proto_rawDesc = []byte{
	0x0a, 0x27, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2f, 0x73, 0x68,
	0x69, 0x65, 0x6c, 0x64, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x67, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d,
	0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x69, 0x0a, 0x19,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x32, 0x0a, 0x0e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x68, 0x0a, 0x1b, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x49, 0x0a, 0x04, 0x62, 0x6f,
	0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x52,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x74, 0x0a, 0x1c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x54, 0x0a, 0x10, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2a, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x73, 0x68,
	0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x52, 0x0e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x22, 0xd5, 0x01, 0x0a, 0x1c,
	0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x18, 0x0a, 0x07,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x62, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x42, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61,
	0x6e, 0x79, 0x2e, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2e, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x9a, 0x01, 0x04,
	0x08, 0x01, 0x10, 0x01, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x37, 0x0a, 0x09, 0x44, 0x61,
	0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0x7c, 0x0a, 0x1c, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x4c, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x38, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e,
	0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x55,
	0x70, 0x73, 0x65, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x04, 0x62, 0x6f, 0x64,
	0x79, 0x22, 0x7d, 0x0a, 0x1d, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x4c, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x38, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x73,
	0x68, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x55, 0x70,
	0x73, 0x65, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x22, 0x31, 0x0a, 0x1d, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6e, 0x22, 0x32, 0x0a, 0x1e, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6e, 0x32, 0xc6, 0x05, 0x0a, 0x13, 0x53, 0x68, 0x69, 0x65,
	0x6c, 0x64, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0xd7, 0x01, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x12, 0x37, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x38, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e,
	0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x4b, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4c, 0x92, 0x41, 0x27,
	0x0a, 0x0c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20, 0x44, 0x61, 0x74, 0x61, 0x12, 0x17,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20, 0x44,
	0x61, 0x74, 0x61, 0x20, 0x4b, 0x65, 0x79, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1c, 0x3a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x22, 0x14, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x64, 0x61, 0x74, 0x61, 0x12, 0xe6, 0x01, 0x0a, 0x15, 0x55, 0x70,
	0x73, 0x65, 0x72, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x38, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e,
	0x79, 0x2e, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2e, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x39, 0x2e,
	0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x73, 0x68, 0x69, 0x65,
	0x6c, 0x64, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x55, 0x70, 0x73, 0x65, 0x72,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x58, 0x92, 0x41, 0x28, 0x0a, 0x0c, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20, 0x44, 0x61, 0x74, 0x61, 0x12, 0x18, 0x55, 0x70, 0x73,
	0x65, 0x72, 0x74, 0x20, 0x55, 0x73, 0x65, 0x72, 0x20, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x20, 0x44, 0x61, 0x74, 0x61, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x27, 0x3a, 0x04, 0x62, 0x6f, 0x64,
	0x79, 0x1a, 0x1f, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x64, 0x61,
	0x74, 0x61, 0x12, 0xeb, 0x01, 0x0a, 0x16, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x39, 0x2e,
	0x67, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x73, 0x68, 0x69, 0x65,
	0x6c, 0x64, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x55, 0x70, 0x73, 0x65, 0x72,
	0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3a, 0x2e, 0x67, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x5a, 0x92, 0x41, 0x29, 0x0a, 0x0c, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x20, 0x44, 0x61, 0x74, 0x61, 0x12, 0x19, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x20,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x20, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20, 0x44, 0x61,
	0x74, 0x61, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x28, 0x3a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x1a, 0x20,
	0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f,
	0x7b, 0x69, 0x64, 0x7d, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x64, 0x61, 0x74, 0x61,
	0x42, 0x76, 0x92, 0x41, 0x14, 0x12, 0x0f, 0x0a, 0x06, 0x53, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x32,
	0x05, 0x30, 0x2e, 0x31, 0x2e, 0x30, 0x2a, 0x01, 0x01, 0x0a, 0x25, 0x63, 0x6f, 0x6d, 0x2e, 0x67,
	0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x6e, 0x2e, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x42, 0x06, 0x53, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6e,
	0x2f, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x68, 0x69, 0x65, 0x6c,
	0x64, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gotocompany_shield_v1beta1_public_proto_rawDescOnce sync.Once
	file_gotocompany_shield_v1beta1_public_proto_rawDescData = file_gotocompany_shield_v1beta1_public_proto_rawDesc
)

func file_gotocompany_shield_v1beta1_public_proto_rawDescGZIP() []byte {
	file_gotocompany_shield_v1beta1_public_proto_rawDescOnce.Do(func() {
		file_gotocompany_shield_v1beta1_public_proto_rawDescData = protoimpl.X.CompressGZIP(file_gotocompany_shield_v1beta1_public_proto_rawDescData)
	})
	return file_gotocompany_shield_v1beta1_public_proto_rawDescData
}

var file_gotocompany_shield_v1beta1_public_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_gotocompany_shield_v1beta1_public_proto_goTypes = []interface{}{
	(*ServiceDataKeyRequestBody)(nil),      // 0: gotocompany.shield.v1beta1.ServiceDataKeyRequestBody
	(*ServiceDataKey)(nil),                 // 1: gotocompany.shield.v1beta1.ServiceDataKey
	(*CreateServiceDataKeyRequest)(nil),    // 2: gotocompany.shield.v1beta1.CreateServiceDataKeyRequest
	(*CreateServiceDataKeyResponse)(nil),   // 3: gotocompany.shield.v1beta1.CreateServiceDataKeyResponse
	(*UpsertServiceDataRequestBody)(nil),   // 4: gotocompany.shield.v1beta1.UpsertServiceDataRequestBody
	(*UpsertUserServiceDataRequest)(nil),   // 5: gotocompany.shield.v1beta1.UpsertUserServiceDataRequest
	(*UpsertGroupServiceDataRequest)(nil),  // 6: gotocompany.shield.v1beta1.UpsertGroupServiceDataRequest
	(*UpsertUserServiceDataResponse)(nil),  // 7: gotocompany.shield.v1beta1.UpsertUserServiceDataResponse
	(*UpsertGroupServiceDataResponse)(nil), // 8: gotocompany.shield.v1beta1.UpsertGroupServiceDataResponse
	nil,                                    // 9: gotocompany.shield.v1beta1.UpsertServiceDataRequestBody.DataEntry
}
var file_gotocompany_shield_v1beta1_public_proto_depIdxs = []int32{
	0, // 0: gotocompany.shield.v1beta1.CreateServiceDataKeyRequest.body:type_name -> gotocompany.shield.v1beta1.ServiceDataKeyRequestBody
	1, // 1: gotocompany.shield.v1beta1.CreateServiceDataKeyResponse.service_data_key:type_name -> gotocompany.shield.v1beta1.ServiceDataKey
	9, // 2: gotocompany.shield.v1beta1.UpsertServiceDataRequestBody.data:type_name -> gotocompany.shield.v1beta1.UpsertServiceDataRequestBody.DataEntry
	4, // 3: gotocompany.shield.v1beta1.UpsertUserServiceDataRequest.body:type_name -> gotocompany.shield.v1beta1.UpsertServiceDataRequestBody
	4, // 4: gotocompany.shield.v1beta1.UpsertGroupServiceDataRequest.body:type_name -> gotocompany.shield.v1beta1.UpsertServiceDataRequestBody
	2, // 5: gotocompany.shield.v1beta1.ShieldPublicService.CreateServiceDataKey:input_type -> gotocompany.shield.v1beta1.CreateServiceDataKeyRequest
	5, // 6: gotocompany.shield.v1beta1.ShieldPublicService.UpsertUserServiceData:input_type -> gotocompany.shield.v1beta1.UpsertUserServiceDataRequest
	6, // 7: gotocompany.shield.v1beta1.ShieldPublicService.UpsertGroupServiceData:input_type -> gotocompany.shield.v1beta1.UpsertGroupServiceDataRequest
	3, // 8: gotocompany.shield.v1beta1.ShieldPublicService.CreateServiceDataKey:output_type -> gotocompany.shield.v1beta1.CreateServiceDataKeyResponse
	7, // 9: gotocompany.shield.v1beta1.ShieldPublicService.UpsertUserServiceData:output_type -> gotocompany.shield.v1beta1.UpsertUserServiceDataResponse
	8, // 10: gotocompany.shield.v1beta1.ShieldPublicService.UpsertGroupServiceData:output_type -> gotocompany.shield.v1beta1.UpsertGroupServiceDataResponse
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_gotocompany_shield_v1beta1_public_proto_init() }
func file_gotocompany_shield_v1beta1_public_proto_init() {
	if File_gotocompany_shield_v1beta1_public_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gotocompany_shield_v1beta1_public_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceDataKeyRequestBody); i {
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
		file_gotocompany_shield_v1beta1_public_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceDataKey); i {
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
		file_gotocompany_shield_v1beta1_public_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateServiceDataKeyRequest); i {
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
		file_gotocompany_shield_v1beta1_public_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateServiceDataKeyResponse); i {
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
		file_gotocompany_shield_v1beta1_public_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertServiceDataRequestBody); i {
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
		file_gotocompany_shield_v1beta1_public_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertUserServiceDataRequest); i {
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
		file_gotocompany_shield_v1beta1_public_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertGroupServiceDataRequest); i {
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
		file_gotocompany_shield_v1beta1_public_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertUserServiceDataResponse); i {
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
		file_gotocompany_shield_v1beta1_public_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertGroupServiceDataResponse); i {
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
			RawDescriptor: file_gotocompany_shield_v1beta1_public_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gotocompany_shield_v1beta1_public_proto_goTypes,
		DependencyIndexes: file_gotocompany_shield_v1beta1_public_proto_depIdxs,
		MessageInfos:      file_gotocompany_shield_v1beta1_public_proto_msgTypes,
	}.Build()
	File_gotocompany_shield_v1beta1_public_proto = out.File
	file_gotocompany_shield_v1beta1_public_proto_rawDesc = nil
	file_gotocompany_shield_v1beta1_public_proto_goTypes = nil
	file_gotocompany_shield_v1beta1_public_proto_depIdxs = nil
}
