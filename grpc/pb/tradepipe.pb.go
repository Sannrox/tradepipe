// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: api/proto/tradepipe.proto

package pb

import (
	login "github.com/Sannrox/tradepipe/grpc/pb/login"
	portfolio "github.com/Sannrox/tradepipe/grpc/pb/portfolio"
	timeline "github.com/Sannrox/tradepipe/grpc/pb/timeline"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Alive struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status     string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	ServerTime int64  `protobuf:"varint,2,opt,name=serverTime,proto3" json:"serverTime,omitempty"`
}

func (x *Alive) Reset() {
	*x = Alive{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_tradepipe_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Alive) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Alive) ProtoMessage() {}

func (x *Alive) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_tradepipe_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Alive.ProtoReflect.Descriptor instead.
func (*Alive) Descriptor() ([]byte, []int) {
	return file_api_proto_tradepipe_proto_rawDescGZIP(), []int{0}
}

func (x *Alive) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Alive) GetServerTime() int64 {
	if x != nil {
		return x.ServerTime
	}
	return 0
}

var File_api_proto_tradepipe_proto protoreflect.FileDescriptor

var file_api_proto_tradepipe_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x72, 0x61, 0x64,
	0x65, 0x70, 0x69, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a,
	0x1b, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e,
	0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x61, 0x70,
	0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x23, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6f, 0x72, 0x74, 0x66,
	0x6f, 0x6c, 0x69, 0x6f, 0x2f, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x3f, 0x0a, 0x05, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x54, 0x69,
	0x6d, 0x65, 0x32, 0xf7, 0x02, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x64, 0x65, 0x50, 0x69, 0x70, 0x65,
	0x12, 0x2c, 0x0a, 0x05, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x22, 0x00, 0x12, 0x2f,
	0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x12, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e,
	0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x1a, 0x10, 0x2e, 0x6c, 0x6f,
	0x67, 0x69, 0x6e, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x64, 0x22, 0x00, 0x12,
	0x30, 0x0a, 0x06, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x12, 0x10, 0x2e, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x2e, 0x54, 0x77, 0x6f, 0x46, 0x41, 0x41, 0x73, 0x6b, 0x73, 0x1a, 0x12, 0x2e, 0x6c, 0x6f,
	0x67, 0x69, 0x6e, 0x2e, 0x54, 0x77, 0x6f, 0x46, 0x41, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x22,
	0x00, 0x12, 0x43, 0x0a, 0x08, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x19, 0x2e,
	0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x1a, 0x1a, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x6c,
	0x69, 0x6e, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x6c, 0x69, 0x6e, 0x65, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x0f, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69,
	0x6e, 0x65, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x19, 0x2e, 0x74, 0x69, 0x6d, 0x65,
	0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x6c, 0x69, 0x6e, 0x65, 0x1a, 0x1a, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x22, 0x00, 0x12, 0x48, 0x0a, 0x09, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x1b, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x1c, 0x2e, 0x70,
	0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x00, 0x42, 0x26, 0x5a, 0x24,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x61, 0x6e, 0x6e, 0x72,
	0x6f, 0x78, 0x2f, 0x74, 0x72, 0x61, 0x64, 0x65, 0x70, 0x69, 0x70, 0x65, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_tradepipe_proto_rawDescOnce sync.Once
	file_api_proto_tradepipe_proto_rawDescData = file_api_proto_tradepipe_proto_rawDesc
)

func file_api_proto_tradepipe_proto_rawDescGZIP() []byte {
	file_api_proto_tradepipe_proto_rawDescOnce.Do(func() {
		file_api_proto_tradepipe_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_tradepipe_proto_rawDescData)
	})
	return file_api_proto_tradepipe_proto_rawDescData
}

var file_api_proto_tradepipe_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_proto_tradepipe_proto_goTypes = []interface{}{
	(*Alive)(nil),                       // 0: pb.Alive
	(*emptypb.Empty)(nil),               // 1: google.protobuf.Empty
	(*login.Credentials)(nil),           // 2: login.Credentials
	(*login.TwoFAAsks)(nil),             // 3: login.TwoFAAsks
	(*timeline.RequestTimeline)(nil),    // 4: timeline.RequestTimeline
	(*portfolio.RequestPositions)(nil),  // 5: portfolio.RequestPositions
	(*login.ProcessId)(nil),             // 6: login.ProcessId
	(*login.TwoFAReturn)(nil),           // 7: login.TwoFAReturn
	(*timeline.ResponseTimeline)(nil),   // 8: timeline.ResponseTimeline
	(*portfolio.ResponsePositions)(nil), // 9: portfolio.ResponsePositions
}
var file_api_proto_tradepipe_proto_depIdxs = []int32{
	1, // 0: pb.TradePipe.Alive:input_type -> google.protobuf.Empty
	2, // 1: pb.TradePipe.Login:input_type -> login.Credentials
	3, // 2: pb.TradePipe.Verify:input_type -> login.TwoFAAsks
	4, // 3: pb.TradePipe.Timeline:input_type -> timeline.RequestTimeline
	4, // 4: pb.TradePipe.TimelineDetails:input_type -> timeline.RequestTimeline
	5, // 5: pb.TradePipe.Positions:input_type -> portfolio.RequestPositions
	0, // 6: pb.TradePipe.Alive:output_type -> pb.Alive
	6, // 7: pb.TradePipe.Login:output_type -> login.ProcessId
	7, // 8: pb.TradePipe.Verify:output_type -> login.TwoFAReturn
	8, // 9: pb.TradePipe.Timeline:output_type -> timeline.ResponseTimeline
	8, // 10: pb.TradePipe.TimelineDetails:output_type -> timeline.ResponseTimeline
	9, // 11: pb.TradePipe.Positions:output_type -> portfolio.ResponsePositions
	6, // [6:12] is the sub-list for method output_type
	0, // [0:6] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_tradepipe_proto_init() }
func file_api_proto_tradepipe_proto_init() {
	if File_api_proto_tradepipe_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_tradepipe_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Alive); i {
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
			RawDescriptor: file_api_proto_tradepipe_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_tradepipe_proto_goTypes,
		DependencyIndexes: file_api_proto_tradepipe_proto_depIdxs,
		MessageInfos:      file_api_proto_tradepipe_proto_msgTypes,
	}.Build()
	File_api_proto_tradepipe_proto = out.File
	file_api_proto_tradepipe_proto_rawDesc = nil
	file_api_proto_tradepipe_proto_goTypes = nil
	file_api_proto_tradepipe_proto_depIdxs = nil
}
