// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: presence.proto

package presence

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UpdatePresenceRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Token         string                 `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdatePresenceRequest) Reset() {
	*x = UpdatePresenceRequest{}
	mi := &file_presence_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdatePresenceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePresenceRequest) ProtoMessage() {}

func (x *UpdatePresenceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_presence_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePresenceRequest.ProtoReflect.Descriptor instead.
func (*UpdatePresenceRequest) Descriptor() ([]byte, []int) {
	return file_presence_proto_rawDescGZIP(), []int{0}
}

func (x *UpdatePresenceRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UpdatePresenceRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type UpdatePresenceResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdatePresenceResponse) Reset() {
	*x = UpdatePresenceResponse{}
	mi := &file_presence_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdatePresenceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePresenceResponse) ProtoMessage() {}

func (x *UpdatePresenceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_presence_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePresenceResponse.ProtoReflect.Descriptor instead.
func (*UpdatePresenceResponse) Descriptor() ([]byte, []int) {
	return file_presence_proto_rawDescGZIP(), []int{1}
}

func (x *UpdatePresenceResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetOnlineUsersRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	TimeoutSeconds int32                  `protobuf:"varint,1,opt,name=timeout_seconds,json=timeoutSeconds,proto3" json:"timeout_seconds,omitempty"`
	Token          string                 `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *GetOnlineUsersRequest) Reset() {
	*x = GetOnlineUsersRequest{}
	mi := &file_presence_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOnlineUsersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOnlineUsersRequest) ProtoMessage() {}

func (x *GetOnlineUsersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_presence_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOnlineUsersRequest.ProtoReflect.Descriptor instead.
func (*GetOnlineUsersRequest) Descriptor() ([]byte, []int) {
	return file_presence_proto_rawDescGZIP(), []int{2}
}

func (x *GetOnlineUsersRequest) GetTimeoutSeconds() int32 {
	if x != nil {
		return x.TimeoutSeconds
	}
	return 0
}

func (x *GetOnlineUsersRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type GetOnlineUsersResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserIds       []string               `protobuf:"bytes,1,rep,name=user_ids,json=userIds,proto3" json:"user_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOnlineUsersResponse) Reset() {
	*x = GetOnlineUsersResponse{}
	mi := &file_presence_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOnlineUsersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOnlineUsersResponse) ProtoMessage() {}

func (x *GetOnlineUsersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_presence_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOnlineUsersResponse.ProtoReflect.Descriptor instead.
func (*GetOnlineUsersResponse) Descriptor() ([]byte, []int) {
	return file_presence_proto_rawDescGZIP(), []int{3}
}

func (x *GetOnlineUsersResponse) GetUserIds() []string {
	if x != nil {
		return x.UserIds
	}
	return nil
}

var File_presence_proto protoreflect.FileDescriptor

const file_presence_proto_rawDesc = "" +
	"\n" +
	"\x0epresence.proto\x12\bpresence\"F\n" +
	"\x15UpdatePresenceRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x14\n" +
	"\x05token\x18\x02 \x01(\tR\x05token\"2\n" +
	"\x16UpdatePresenceResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"V\n" +
	"\x15GetOnlineUsersRequest\x12'\n" +
	"\x0ftimeout_seconds\x18\x01 \x01(\x05R\x0etimeoutSeconds\x12\x14\n" +
	"\x05token\x18\x02 \x01(\tR\x05token\"3\n" +
	"\x16GetOnlineUsersResponse\x12\x19\n" +
	"\buser_ids\x18\x01 \x03(\tR\auserIds2\xbb\x01\n" +
	"\x0fPresenceService\x12S\n" +
	"\x0eUpdatePresence\x12\x1f.presence.UpdatePresenceRequest\x1a .presence.UpdatePresenceResponse\x12S\n" +
	"\x0eGetOnlineUsers\x12\x1f.presence.GetOnlineUsersRequest\x1a .presence.GetOnlineUsersResponseB<Z:websocket-backend/services/presence-service/proto;presenceb\x06proto3"

var (
	file_presence_proto_rawDescOnce sync.Once
	file_presence_proto_rawDescData []byte
)

func file_presence_proto_rawDescGZIP() []byte {
	file_presence_proto_rawDescOnce.Do(func() {
		file_presence_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_presence_proto_rawDesc), len(file_presence_proto_rawDesc)))
	})
	return file_presence_proto_rawDescData
}

var file_presence_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_presence_proto_goTypes = []any{
	(*UpdatePresenceRequest)(nil),  // 0: presence.UpdatePresenceRequest
	(*UpdatePresenceResponse)(nil), // 1: presence.UpdatePresenceResponse
	(*GetOnlineUsersRequest)(nil),  // 2: presence.GetOnlineUsersRequest
	(*GetOnlineUsersResponse)(nil), // 3: presence.GetOnlineUsersResponse
}
var file_presence_proto_depIdxs = []int32{
	0, // 0: presence.PresenceService.UpdatePresence:input_type -> presence.UpdatePresenceRequest
	2, // 1: presence.PresenceService.GetOnlineUsers:input_type -> presence.GetOnlineUsersRequest
	1, // 2: presence.PresenceService.UpdatePresence:output_type -> presence.UpdatePresenceResponse
	3, // 3: presence.PresenceService.GetOnlineUsers:output_type -> presence.GetOnlineUsersResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_presence_proto_init() }
func file_presence_proto_init() {
	if File_presence_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_presence_proto_rawDesc), len(file_presence_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_presence_proto_goTypes,
		DependencyIndexes: file_presence_proto_depIdxs,
		MessageInfos:      file_presence_proto_msgTypes,
	}.Build()
	File_presence_proto = out.File
	file_presence_proto_goTypes = nil
	file_presence_proto_depIdxs = nil
}
