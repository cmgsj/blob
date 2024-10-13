// Copyright 2024 Buf Technologies, Inc.
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
// source: buf/plugin/info/v1/license.proto

package infov1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

// A plugin license.
type License struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The SPDX license ID.
	//
	// See https://spdx.org/licenses.
	SpdxLicenseId string `protobuf:"bytes,1,opt,name=spdx_license_id,json=spdxLicenseId,proto3" json:"spdx_license_id,omitempty"`
	// The source of a license is either raw text, or a URL that contains the license.
	//
	// Types that are assignable to Source:
	//
	//	*License_Text
	//	*License_Url
	Source isLicense_Source `protobuf_oneof:"source"`
}

func (x *License) Reset() {
	*x = License{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_plugin_info_v1_license_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *License) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*License) ProtoMessage() {}

func (x *License) ProtoReflect() protoreflect.Message {
	mi := &file_buf_plugin_info_v1_license_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use License.ProtoReflect.Descriptor instead.
func (*License) Descriptor() ([]byte, []int) {
	return file_buf_plugin_info_v1_license_proto_rawDescGZIP(), []int{0}
}

func (x *License) GetSpdxLicenseId() string {
	if x != nil {
		return x.SpdxLicenseId
	}
	return ""
}

func (m *License) GetSource() isLicense_Source {
	if m != nil {
		return m.Source
	}
	return nil
}

func (x *License) GetText() string {
	if x, ok := x.GetSource().(*License_Text); ok {
		return x.Text
	}
	return ""
}

func (x *License) GetUrl() string {
	if x, ok := x.GetSource().(*License_Url); ok {
		return x.Url
	}
	return ""
}

type isLicense_Source interface {
	isLicense_Source()
}

type License_Text struct {
	// The raw text of the license.
	Text string `protobuf:"bytes,2,opt,name=text,proto3,oneof"`
}

type License_Url struct {
	// The url that contains the license
	Url string `protobuf:"bytes,3,opt,name=url,proto3,oneof"`
}

func (*License_Text) isLicense_Source() {}

func (*License_Url) isLicense_Source() {}

var File_buf_plugin_info_v1_license_proto protoreflect.FileDescriptor

var file_buf_plugin_info_v1_license_proto_rawDesc = []byte{
	0x0a, 0x20, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x69, 0x6e, 0x66,
	0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x12, 0x62, 0x75, 0x66, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x69,
	0x6e, 0x66, 0x6f, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x79, 0x0a, 0x07, 0x4c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x26,
	0x0a, 0x0f, 0x73, 0x70, 0x64, 0x78, 0x5f, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x70, 0x64, 0x78, 0x4c, 0x69, 0x63,
	0x65, 0x6e, 0x73, 0x65, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x1f, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xba, 0x48, 0x08, 0xd8, 0x01,
	0x01, 0x72, 0x03, 0x88, 0x01, 0x01, 0x48, 0x00, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x42, 0x0f, 0x0a,
	0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x05, 0xba, 0x48, 0x02, 0x08, 0x01, 0x42, 0x52,
	0x5a, 0x50, 0x62, 0x75, 0x66, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x67, 0x6f, 0x2f, 0x62, 0x75, 0x66, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x62, 0x75, 0x66, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75,
	0x66, 0x66, 0x65, 0x72, 0x73, 0x2f, 0x67, 0x6f, 0x2f, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x76, 0x31, 0x3b, 0x69, 0x6e, 0x66, 0x6f,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_buf_plugin_info_v1_license_proto_rawDescOnce sync.Once
	file_buf_plugin_info_v1_license_proto_rawDescData = file_buf_plugin_info_v1_license_proto_rawDesc
)

func file_buf_plugin_info_v1_license_proto_rawDescGZIP() []byte {
	file_buf_plugin_info_v1_license_proto_rawDescOnce.Do(func() {
		file_buf_plugin_info_v1_license_proto_rawDescData = protoimpl.X.CompressGZIP(file_buf_plugin_info_v1_license_proto_rawDescData)
	})
	return file_buf_plugin_info_v1_license_proto_rawDescData
}

var file_buf_plugin_info_v1_license_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_buf_plugin_info_v1_license_proto_goTypes = []any{
	(*License)(nil), // 0: buf.plugin.info.v1.License
}
var file_buf_plugin_info_v1_license_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_buf_plugin_info_v1_license_proto_init() }
func file_buf_plugin_info_v1_license_proto_init() {
	if File_buf_plugin_info_v1_license_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_buf_plugin_info_v1_license_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*License); i {
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
	file_buf_plugin_info_v1_license_proto_msgTypes[0].OneofWrappers = []any{
		(*License_Text)(nil),
		(*License_Url)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_buf_plugin_info_v1_license_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_buf_plugin_info_v1_license_proto_goTypes,
		DependencyIndexes: file_buf_plugin_info_v1_license_proto_depIdxs,
		MessageInfos:      file_buf_plugin_info_v1_license_proto_msgTypes,
	}.Build()
	File_buf_plugin_info_v1_license_proto = out.File
	file_buf_plugin_info_v1_license_proto_rawDesc = nil
	file_buf_plugin_info_v1_license_proto_goTypes = nil
	file_buf_plugin_info_v1_license_proto_depIdxs = nil
}
