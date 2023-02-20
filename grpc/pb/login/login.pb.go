// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: api/proto/login/login.proto

package login

import (
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

type Credentials struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number string `protobuf:"bytes,1,opt,name=number,proto3" json:"number,omitempty"`
	Pin    string `protobuf:"bytes,2,opt,name=pin,proto3" json:"pin,omitempty"`
}

func (x *Credentials) Reset() {
	*x = Credentials{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_login_login_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Credentials) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Credentials) ProtoMessage() {}

func (x *Credentials) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_login_login_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Credentials.ProtoReflect.Descriptor instead.
func (*Credentials) Descriptor() ([]byte, []int) {
	return file_api_proto_login_login_proto_rawDescGZIP(), []int{0}
}

func (x *Credentials) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *Credentials) GetPin() string {
	if x != nil {
		return x.Pin
	}
	return ""
}

type ProcessId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcessId string `protobuf:"bytes,1,opt,name=processId,proto3" json:"processId,omitempty"`
	Error     string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *ProcessId) Reset() {
	*x = ProcessId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_login_login_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessId) ProtoMessage() {}

func (x *ProcessId) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_login_login_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessId.ProtoReflect.Descriptor instead.
func (*ProcessId) Descriptor() ([]byte, []int) {
	return file_api_proto_login_login_proto_rawDescGZIP(), []int{1}
}

func (x *ProcessId) GetProcessId() string {
	if x != nil {
		return x.ProcessId
	}
	return ""
}

func (x *ProcessId) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type TwoFAAsks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcessId  string `protobuf:"bytes,1,opt,name=processId,proto3" json:"processId,omitempty"`
	VerifyCode int32  `protobuf:"varint,2,opt,name=verifyCode,proto3" json:"verifyCode,omitempty"`
}

func (x *TwoFAAsks) Reset() {
	*x = TwoFAAsks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_login_login_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TwoFAAsks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TwoFAAsks) ProtoMessage() {}

func (x *TwoFAAsks) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_login_login_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TwoFAAsks.ProtoReflect.Descriptor instead.
func (*TwoFAAsks) Descriptor() ([]byte, []int) {
	return file_api_proto_login_login_proto_rawDescGZIP(), []int{2}
}

func (x *TwoFAAsks) GetProcessId() string {
	if x != nil {
		return x.ProcessId
	}
	return ""
}

func (x *TwoFAAsks) GetVerifyCode() int32 {
	if x != nil {
		return x.VerifyCode
	}
	return 0
}

type TwoFAReturn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *TwoFAReturn) Reset() {
	*x = TwoFAReturn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_login_login_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TwoFAReturn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TwoFAReturn) ProtoMessage() {}

func (x *TwoFAReturn) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_login_login_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TwoFAReturn.ProtoReflect.Descriptor instead.
func (*TwoFAReturn) Descriptor() ([]byte, []int) {
	return file_api_proto_login_login_proto_rawDescGZIP(), []int{3}
}

func (x *TwoFAReturn) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_api_proto_login_login_proto protoreflect.FileDescriptor

var file_api_proto_login_login_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6c,
	0x6f, 0x67, 0x69, 0x6e, 0x22, 0x37, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x61, 0x6c, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x70,
	0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x70, 0x69, 0x6e, 0x22, 0x3f, 0x0a,
	0x09, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x49,
	0x0a, 0x09, 0x54, 0x77, 0x6f, 0x46, 0x41, 0x41, 0x73, 0x6b, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x70,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x76, 0x65, 0x72,
	0x69, 0x66, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x76,
	0x65, 0x72, 0x69, 0x66, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x54, 0x77, 0x6f,
	0x46, 0x41, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x2c,
	0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x61, 0x6e,
	0x6e, 0x72, 0x6f, 0x78, 0x2f, 0x74, 0x72, 0x61, 0x64, 0x65, 0x70, 0x69, 0x70, 0x65, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_login_login_proto_rawDescOnce sync.Once
	file_api_proto_login_login_proto_rawDescData = file_api_proto_login_login_proto_rawDesc
)

func file_api_proto_login_login_proto_rawDescGZIP() []byte {
	file_api_proto_login_login_proto_rawDescOnce.Do(func() {
		file_api_proto_login_login_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_login_login_proto_rawDescData)
	})
	return file_api_proto_login_login_proto_rawDescData
}

var file_api_proto_login_login_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_proto_login_login_proto_goTypes = []interface{}{
	(*Credentials)(nil), // 0: login.Credentials
	(*ProcessId)(nil),   // 1: login.ProcessId
	(*TwoFAAsks)(nil),   // 2: login.TwoFAAsks
	(*TwoFAReturn)(nil), // 3: login.TwoFAReturn
}
var file_api_proto_login_login_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_login_login_proto_init() }
func file_api_proto_login_login_proto_init() {
	if File_api_proto_login_login_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_login_login_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Credentials); i {
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
		file_api_proto_login_login_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessId); i {
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
		file_api_proto_login_login_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TwoFAAsks); i {
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
		file_api_proto_login_login_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TwoFAReturn); i {
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
			RawDescriptor: file_api_proto_login_login_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_login_login_proto_goTypes,
		DependencyIndexes: file_api_proto_login_login_proto_depIdxs,
		MessageInfos:      file_api_proto_login_login_proto_msgTypes,
	}.Build()
	File_api_proto_login_login_proto = out.File
	file_api_proto_login_login_proto_rawDesc = nil
	file_api_proto_login_login_proto_goTypes = nil
	file_api_proto_login_login_proto_depIdxs = nil
}