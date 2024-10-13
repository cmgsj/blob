// Copyright 2023-2024 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: buf/registry/owner/v1/user.proto

package ownerv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "buf.build/gen/go/bufbuild/registry/protocolbuffers/go/buf/registry/priv/extension/v1beta1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The state of the a User, either active or inactive.
type UserState int32

const (
	UserState_USER_STATE_UNSPECIFIED UserState = 0
	UserState_USER_STATE_ACTIVE      UserState = 1
	UserState_USER_STATE_INACTIVE    UserState = 2
)

// Enum value maps for UserState.
var (
	UserState_name = map[int32]string{
		0: "USER_STATE_UNSPECIFIED",
		1: "USER_STATE_ACTIVE",
		2: "USER_STATE_INACTIVE",
	}
	UserState_value = map[string]int32{
		"USER_STATE_UNSPECIFIED": 0,
		"USER_STATE_ACTIVE":      1,
		"USER_STATE_INACTIVE":    2,
	}
)

func (x UserState) Enum() *UserState {
	p := new(UserState)
	*p = x
	return p
}

func (x UserState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserState) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_registry_owner_v1_user_proto_enumTypes[0].Descriptor()
}

func (UserState) Type() protoreflect.EnumType {
	return &file_buf_registry_owner_v1_user_proto_enumTypes[0]
}

func (x UserState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserState.Descriptor instead.
func (UserState) EnumDescriptor() ([]byte, []int) {
	return file_buf_registry_owner_v1_user_proto_rawDescGZIP(), []int{0}
}

// The type of a User.
type UserType int32

const (
	UserType_USER_TYPE_UNSPECIFIED UserType = 0
	// Users that are standard users.
	UserType_USER_TYPE_STANDARD UserType = 1
	// Users that are bots.
	UserType_USER_TYPE_BOT UserType = 2
	// Users that are internal system users.
	UserType_USER_TYPE_SYSTEM UserType = 3
)

// Enum value maps for UserType.
var (
	UserType_name = map[int32]string{
		0: "USER_TYPE_UNSPECIFIED",
		1: "USER_TYPE_STANDARD",
		2: "USER_TYPE_BOT",
		3: "USER_TYPE_SYSTEM",
	}
	UserType_value = map[string]int32{
		"USER_TYPE_UNSPECIFIED": 0,
		"USER_TYPE_STANDARD":    1,
		"USER_TYPE_BOT":         2,
		"USER_TYPE_SYSTEM":      3,
	}
)

func (x UserType) Enum() *UserType {
	p := new(UserType)
	*p = x
	return p
}

func (x UserType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserType) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_registry_owner_v1_user_proto_enumTypes[1].Descriptor()
}

func (UserType) Type() protoreflect.EnumType {
	return &file_buf_registry_owner_v1_user_proto_enumTypes[1]
}

func (x UserType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserType.Descriptor instead.
func (UserType) EnumDescriptor() ([]byte, []int) {
	return file_buf_registry_owner_v1_user_proto_rawDescGZIP(), []int{1}
}

// The verification status of an User.
type UserVerificationStatus int32

const (
	UserVerificationStatus_USER_VERIFICATION_STATUS_UNSPECIFIED UserVerificationStatus = 0
	// The User is unverified.
	UserVerificationStatus_USER_VERIFICATION_STATUS_UNVERIFIED UserVerificationStatus = 1
	// The User is verified.
	UserVerificationStatus_USER_VERIFICATION_STATUS_VERIFIED UserVerificationStatus = 2
	// The User is an official user of the BSR owner.
	UserVerificationStatus_USER_VERIFICATION_STATUS_OFFICIAL UserVerificationStatus = 3
)

// Enum value maps for UserVerificationStatus.
var (
	UserVerificationStatus_name = map[int32]string{
		0: "USER_VERIFICATION_STATUS_UNSPECIFIED",
		1: "USER_VERIFICATION_STATUS_UNVERIFIED",
		2: "USER_VERIFICATION_STATUS_VERIFIED",
		3: "USER_VERIFICATION_STATUS_OFFICIAL",
	}
	UserVerificationStatus_value = map[string]int32{
		"USER_VERIFICATION_STATUS_UNSPECIFIED": 0,
		"USER_VERIFICATION_STATUS_UNVERIFIED":  1,
		"USER_VERIFICATION_STATUS_VERIFIED":    2,
		"USER_VERIFICATION_STATUS_OFFICIAL":    3,
	}
)

func (x UserVerificationStatus) Enum() *UserVerificationStatus {
	p := new(UserVerificationStatus)
	*p = x
	return p
}

func (x UserVerificationStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserVerificationStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_registry_owner_v1_user_proto_enumTypes[2].Descriptor()
}

func (UserVerificationStatus) Type() protoreflect.EnumType {
	return &file_buf_registry_owner_v1_user_proto_enumTypes[2]
}

func (x UserVerificationStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserVerificationStatus.Descriptor instead.
func (UserVerificationStatus) EnumDescriptor() ([]byte, []int) {
	return file_buf_registry_owner_v1_user_proto_rawDescGZIP(), []int{2}
}

// A user on the BSR.
//
// A name uniquely identifies a User, however name is mutable.
type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The id for the User.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The time the User was created.
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// The last time the User was updated.
	UpdateTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	// The name of the User.
	//
	// A name uniquely identifies a User, however name is mutable.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// The type of the User.
	Type UserType `protobuf:"varint,5,opt,name=type,proto3,enum=buf.registry.owner.v1.UserType" json:"type,omitempty"`
	// The state of the User.
	State UserState `protobuf:"varint,6,opt,name=state,proto3,enum=buf.registry.owner.v1.UserState" json:"state,omitempty"`
	// The configurable description of the User.
	Description string `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
	// The configurable URL that represents the homepage for a User.
	Url string `protobuf:"bytes,8,opt,name=url,proto3" json:"url,omitempty"`
	// The verification status of the User.
	VerificationStatus UserVerificationStatus `protobuf:"varint,9,opt,name=verification_status,json=verificationStatus,proto3,enum=buf.registry.owner.v1.UserVerificationStatus" json:"verification_status,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_registry_owner_v1_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_buf_registry_owner_v1_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_buf_registry_owner_v1_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *User) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetType() UserType {
	if x != nil {
		return x.Type
	}
	return UserType_USER_TYPE_UNSPECIFIED
}

func (x *User) GetState() UserState {
	if x != nil {
		return x.State
	}
	return UserState_USER_STATE_UNSPECIFIED
}

func (x *User) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *User) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *User) GetVerificationStatus() UserVerificationStatus {
	if x != nil {
		return x.VerificationStatus
	}
	return UserVerificationStatus_USER_VERIFICATION_STATUS_UNSPECIFIED
}

// UserRef is a reference to a User, either an id or a name.
//
// This is used in requests.
type UserRef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//
	//	*UserRef_Id
	//	*UserRef_Name
	Value isUserRef_Value `protobuf_oneof:"value"`
}

func (x *UserRef) Reset() {
	*x = UserRef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_registry_owner_v1_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserRef) ProtoMessage() {}

func (x *UserRef) ProtoReflect() protoreflect.Message {
	mi := &file_buf_registry_owner_v1_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserRef.ProtoReflect.Descriptor instead.
func (*UserRef) Descriptor() ([]byte, []int) {
	return file_buf_registry_owner_v1_user_proto_rawDescGZIP(), []int{1}
}

func (m *UserRef) GetValue() isUserRef_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *UserRef) GetId() string {
	if x, ok := x.GetValue().(*UserRef_Id); ok {
		return x.Id
	}
	return ""
}

func (x *UserRef) GetName() string {
	if x, ok := x.GetValue().(*UserRef_Name); ok {
		return x.Name
	}
	return ""
}

type isUserRef_Value interface {
	isUserRef_Value()
}

type UserRef_Id struct {
	// The id of the User.
	Id string `protobuf:"bytes,1,opt,name=id,proto3,oneof"`
}

type UserRef_Name struct {
	// The name of the User.
	Name string `protobuf:"bytes,2,opt,name=name,proto3,oneof"`
}

func (*UserRef_Id) isUserRef_Value() {}

func (*UserRef_Name) isUserRef_Value() {}

var File_buf_registry_owner_v1_user_proto protoreflect.FileDescriptor

var file_buf_registry_owner_v1_user_proto_rawDesc = []byte{
	0x0a, 0x20, 0x62, 0x75, 0x66, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x15, 0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x2e, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x33, 0x62, 0x75, 0x66, 0x2f, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x69, 0x76, 0x2f, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x65,
	0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb4, 0x04, 0x0a,
	0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x0b, 0xba, 0x48, 0x08, 0xc8, 0x01, 0x01, 0x72, 0x03, 0x88, 0x02, 0x01, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x43, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x43, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01,
	0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x27, 0xba, 0x48, 0x24, 0xc8,
	0x01, 0x01, 0x72, 0x1f, 0x10, 0x02, 0x18, 0x20, 0x32, 0x19, 0x5e, 0x5b, 0x61, 0x2d, 0x7a, 0x5d,
	0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d, 0x39, 0x2d, 0x5d, 0x2a, 0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d,
	0x39, 0x5d, 0x24, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x40, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x42, 0x0b, 0xba, 0x48, 0x08, 0xc8, 0x01, 0x01,
	0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x43, 0x0a, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x62, 0x75, 0x66,
	0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x42, 0x0b, 0xba, 0x48,
	0x08, 0xc8, 0x01, 0x01, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x2a, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x18, 0xde, 0x02, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0e, 0xba, 0x48, 0x0b, 0xd8, 0x01,
	0x01, 0x72, 0x06, 0x18, 0xff, 0x01, 0x88, 0x01, 0x01, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x6b,
	0x0a, 0x13, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2d, 0x2e, 0x62, 0x75,
	0x66, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x6f, 0x77, 0x6e, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x0b, 0xba, 0x48, 0x08, 0xc8,
	0x01, 0x01, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x12, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x3a, 0x06, 0xea, 0xc5, 0x2b,
	0x02, 0x10, 0x01, 0x22, 0x79, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x66, 0x12, 0x1a,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72,
	0x03, 0x88, 0x02, 0x01, 0x48, 0x00, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3a, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x24, 0xba, 0x48, 0x21, 0x72, 0x1f, 0x10,
	0x02, 0x18, 0x20, 0x32, 0x19, 0x5e, 0x5b, 0x61, 0x2d, 0x7a, 0x5d, 0x5b, 0x61, 0x2d, 0x7a, 0x30,
	0x2d, 0x39, 0x2d, 0x5d, 0x2a, 0x5b, 0x61, 0x2d, 0x7a, 0x30, 0x2d, 0x39, 0x5d, 0x24, 0x48, 0x00,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x3a, 0x06, 0xea, 0xc5, 0x2b, 0x02, 0x08, 0x01, 0x42, 0x0e,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x05, 0xba, 0x48, 0x02, 0x08, 0x01, 0x2a, 0x57,
	0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x55,
	0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x55, 0x53, 0x45, 0x52, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12, 0x17,
	0x0a, 0x13, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x49, 0x4e, 0x41,
	0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x02, 0x2a, 0x66, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x15, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16,
	0x0a, 0x12, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x4e,
	0x44, 0x41, 0x52, 0x44, 0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x42, 0x4f, 0x54, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x55, 0x53, 0x45,
	0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x10, 0x03, 0x2a,
	0xb9, 0x01, 0x0a, 0x16, 0x55, 0x73, 0x65, 0x72, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x28, 0x0a, 0x24, 0x55, 0x53,
	0x45, 0x52, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x27, 0x0a, 0x23, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x56, 0x45, 0x52,
	0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53,
	0x5f, 0x55, 0x4e, 0x56, 0x45, 0x52, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x01, 0x12, 0x25, 0x0a,
	0x21, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49,
	0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x02, 0x12, 0x25, 0x0a, 0x21, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x56, 0x45, 0x52,
	0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53,
	0x5f, 0x4f, 0x46, 0x46, 0x49, 0x43, 0x49, 0x41, 0x4c, 0x10, 0x03, 0x42, 0x55, 0x5a, 0x53, 0x62,
	0x75, 0x66, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f,
	0x62, 0x75, 0x66, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72,
	0x73, 0x2f, 0x67, 0x6f, 0x2f, 0x62, 0x75, 0x66, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x2f, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x6f, 0x77, 0x6e, 0x65, 0x72,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_buf_registry_owner_v1_user_proto_rawDescOnce sync.Once
	file_buf_registry_owner_v1_user_proto_rawDescData = file_buf_registry_owner_v1_user_proto_rawDesc
)

func file_buf_registry_owner_v1_user_proto_rawDescGZIP() []byte {
	file_buf_registry_owner_v1_user_proto_rawDescOnce.Do(func() {
		file_buf_registry_owner_v1_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_buf_registry_owner_v1_user_proto_rawDescData)
	})
	return file_buf_registry_owner_v1_user_proto_rawDescData
}

var file_buf_registry_owner_v1_user_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_buf_registry_owner_v1_user_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_buf_registry_owner_v1_user_proto_goTypes = []any{
	(UserState)(0),                // 0: buf.registry.owner.v1.UserState
	(UserType)(0),                 // 1: buf.registry.owner.v1.UserType
	(UserVerificationStatus)(0),   // 2: buf.registry.owner.v1.UserVerificationStatus
	(*User)(nil),                  // 3: buf.registry.owner.v1.User
	(*UserRef)(nil),               // 4: buf.registry.owner.v1.UserRef
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_buf_registry_owner_v1_user_proto_depIdxs = []int32{
	5, // 0: buf.registry.owner.v1.User.create_time:type_name -> google.protobuf.Timestamp
	5, // 1: buf.registry.owner.v1.User.update_time:type_name -> google.protobuf.Timestamp
	1, // 2: buf.registry.owner.v1.User.type:type_name -> buf.registry.owner.v1.UserType
	0, // 3: buf.registry.owner.v1.User.state:type_name -> buf.registry.owner.v1.UserState
	2, // 4: buf.registry.owner.v1.User.verification_status:type_name -> buf.registry.owner.v1.UserVerificationStatus
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_buf_registry_owner_v1_user_proto_init() }
func file_buf_registry_owner_v1_user_proto_init() {
	if File_buf_registry_owner_v1_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_buf_registry_owner_v1_user_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*User); i {
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
		file_buf_registry_owner_v1_user_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*UserRef); i {
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
	file_buf_registry_owner_v1_user_proto_msgTypes[1].OneofWrappers = []any{
		(*UserRef_Id)(nil),
		(*UserRef_Name)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_buf_registry_owner_v1_user_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_buf_registry_owner_v1_user_proto_goTypes,
		DependencyIndexes: file_buf_registry_owner_v1_user_proto_depIdxs,
		EnumInfos:         file_buf_registry_owner_v1_user_proto_enumTypes,
		MessageInfos:      file_buf_registry_owner_v1_user_proto_msgTypes,
	}.Build()
	File_buf_registry_owner_v1_user_proto = out.File
	file_buf_registry_owner_v1_user_proto_rawDesc = nil
	file_buf_registry_owner_v1_user_proto_goTypes = nil
	file_buf_registry_owner_v1_user_proto_depIdxs = nil
}