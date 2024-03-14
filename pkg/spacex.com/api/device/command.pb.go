// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: github.com/jgulick48/starlink-statsd/pkg/spacex.com/api/device/command.proto

package device

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

type Capability int32

const (
	Capability_READ             Capability = 0
	Capability_READ_INTERNAL    Capability = 13
	Capability_READ_PRIVATE     Capability = 7
	Capability_LOCAL            Capability = 14
	Capability_WRITE            Capability = 1
	Capability_WRITE_PERSISTENT Capability = 11
	Capability_DEBUG            Capability = 2
	Capability_ADMIN            Capability = 3
	Capability_SETUP            Capability = 4
	Capability_SET_SKU          Capability = 5
	Capability_REFRESH          Capability = 6
	Capability_FUSE             Capability = 8
	Capability_RESET            Capability = 9
	Capability_TEST             Capability = 10
	Capability_SSH              Capability = 12
)

// Enum value maps for Capability.
var (
	Capability_name = map[int32]string{
		0:  "READ",
		13: "READ_INTERNAL",
		7:  "READ_PRIVATE",
		14: "LOCAL",
		1:  "WRITE",
		11: "WRITE_PERSISTENT",
		2:  "DEBUG",
		3:  "ADMIN",
		4:  "SETUP",
		5:  "SET_SKU",
		6:  "REFRESH",
		8:  "FUSE",
		9:  "RESET",
		10: "TEST",
		12: "SSH",
	}
	Capability_value = map[string]int32{
		"READ":             0,
		"READ_INTERNAL":    13,
		"READ_PRIVATE":     7,
		"LOCAL":            14,
		"WRITE":            1,
		"WRITE_PERSISTENT": 11,
		"DEBUG":            2,
		"ADMIN":            3,
		"SETUP":            4,
		"SET_SKU":          5,
		"REFRESH":          6,
		"FUSE":             8,
		"RESET":            9,
		"TEST":             10,
		"SSH":              12,
	}
)

func (x Capability) Enum() *Capability {
	p := new(Capability)
	*p = x
	return p
}

func (x Capability) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Capability) Descriptor() protoreflect.EnumDescriptor {
	return file_spacex_api_device_command_proto_enumTypes[0].Descriptor()
}

func (Capability) Type() protoreflect.EnumType {
	return &file_spacex_api_device_command_proto_enumTypes[0]
}

func (x Capability) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Capability.Descriptor instead.
func (Capability) EnumDescriptor() ([]byte, []int) {
	return file_spacex_api_device_command_proto_rawDescGZIP(), []int{0}
}

type PublicKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key          string       `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Capabilities []Capability `protobuf:"varint,2,rep,packed,name=capabilities,proto3,enum=SpaceX.API.Device.Capability" json:"capabilities,omitempty"`
}

func (x *PublicKey) Reset() {
	*x = PublicKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spacex_api_device_command_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublicKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicKey) ProtoMessage() {}

func (x *PublicKey) ProtoReflect() protoreflect.Message {
	mi := &file_spacex_api_device_command_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublicKey.ProtoReflect.Descriptor instead.
func (*PublicKey) Descriptor() ([]byte, []int) {
	return file_spacex_api_device_command_proto_rawDescGZIP(), []int{0}
}

func (x *PublicKey) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *PublicKey) GetCapabilities() []Capability {
	if x != nil {
		return x.Capabilities
	}
	return nil
}

var File_spacex_api_device_command_proto protoreflect.FileDescriptor

var file_spacex_api_device_command_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x73, 0x70, 0x61, 0x63, 0x65, 0x78, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x11, 0x53, 0x70, 0x61, 0x63, 0x65, 0x58, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x22, 0x60, 0x0a, 0x09, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x41, 0x0a, 0x0c, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74,
	0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x53, 0x70, 0x61, 0x63,
	0x65, 0x58, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x61,
	0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x52, 0x0c, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2a, 0xca, 0x01, 0x0a, 0x0a, 0x43, 0x61, 0x70, 0x61, 0x62,
	0x69, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x08, 0x0a, 0x04, 0x52, 0x45, 0x41, 0x44, 0x10, 0x00, 0x12,
	0x11, 0x0a, 0x0d, 0x52, 0x45, 0x41, 0x44, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c,
	0x10, 0x0d, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x45, 0x41, 0x44, 0x5f, 0x50, 0x52, 0x49, 0x56, 0x41,
	0x54, 0x45, 0x10, 0x07, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x10, 0x0e, 0x12,
	0x09, 0x0a, 0x05, 0x57, 0x52, 0x49, 0x54, 0x45, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x57, 0x52,
	0x49, 0x54, 0x45, 0x5f, 0x50, 0x45, 0x52, 0x53, 0x49, 0x53, 0x54, 0x45, 0x4e, 0x54, 0x10, 0x0b,
	0x12, 0x09, 0x0a, 0x05, 0x44, 0x45, 0x42, 0x55, 0x47, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x41,
	0x44, 0x4d, 0x49, 0x4e, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x45, 0x54, 0x55, 0x50, 0x10,
	0x04, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x45, 0x54, 0x5f, 0x53, 0x4b, 0x55, 0x10, 0x05, 0x12, 0x0b,
	0x0a, 0x07, 0x52, 0x45, 0x46, 0x52, 0x45, 0x53, 0x48, 0x10, 0x06, 0x12, 0x08, 0x0a, 0x04, 0x46,
	0x55, 0x53, 0x45, 0x10, 0x08, 0x12, 0x09, 0x0a, 0x05, 0x52, 0x45, 0x53, 0x45, 0x54, 0x10, 0x09,
	0x12, 0x08, 0x0a, 0x04, 0x54, 0x45, 0x53, 0x54, 0x10, 0x0a, 0x12, 0x07, 0x0a, 0x03, 0x53, 0x53,
	0x48, 0x10, 0x0c, 0x42, 0x17, 0x5a, 0x15, 0x73, 0x70, 0x61, 0x63, 0x65, 0x78, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_spacex_api_device_command_proto_rawDescOnce sync.Once
	file_spacex_api_device_command_proto_rawDescData = file_spacex_api_device_command_proto_rawDesc
)

func file_spacex_api_device_command_proto_rawDescGZIP() []byte {
	file_spacex_api_device_command_proto_rawDescOnce.Do(func() {
		file_spacex_api_device_command_proto_rawDescData = protoimpl.X.CompressGZIP(file_spacex_api_device_command_proto_rawDescData)
	})
	return file_spacex_api_device_command_proto_rawDescData
}

var file_spacex_api_device_command_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_spacex_api_device_command_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_spacex_api_device_command_proto_goTypes = []interface{}{
	(Capability)(0),   // 0: SpaceX.API.Device.Capability
	(*PublicKey)(nil), // 1: SpaceX.API.Device.PublicKey
}
var file_spacex_api_device_command_proto_depIdxs = []int32{
	0, // 0: SpaceX.API.Device.PublicKey.capabilities:type_name -> SpaceX.API.Device.Capability
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_spacex_api_device_command_proto_init() }
func file_spacex_api_device_command_proto_init() {
	if File_spacex_api_device_command_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_spacex_api_device_command_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublicKey); i {
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
			RawDescriptor: file_spacex_api_device_command_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_spacex_api_device_command_proto_goTypes,
		DependencyIndexes: file_spacex_api_device_command_proto_depIdxs,
		EnumInfos:         file_spacex_api_device_command_proto_enumTypes,
		MessageInfos:      file_spacex_api_device_command_proto_msgTypes,
	}.Build()
	File_spacex_api_device_command_proto = out.File
	file_spacex_api_device_command_proto_rawDesc = nil
	file_spacex_api_device_command_proto_goTypes = nil
	file_spacex_api_device_command_proto_depIdxs = nil
}
