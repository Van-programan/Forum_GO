// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: proto/forumservice/forum.proto

package forumservice

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type Topic struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	AuthorId      int64                  `protobuf:"varint,3,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Topic) Reset() {
	*x = Topic{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Topic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Topic) ProtoMessage() {}

func (x *Topic) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Topic.ProtoReflect.Descriptor instead.
func (*Topic) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{0}
}

func (x *Topic) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Topic) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Topic) GetAuthorId() int64 {
	if x != nil {
		return x.AuthorId
	}
	return 0
}

func (x *Topic) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	TopicId       int64                  `protobuf:"varint,2,opt,name=topic_id,json=topicId,proto3" json:"topic_id,omitempty"`
	UserId        int64                  `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Content       string                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{1}
}

func (x *Message) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Message) GetTopicId() int64 {
	if x != nil {
		return x.TopicId
	}
	return 0
}

func (x *Message) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Message) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type CreateTopicRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	AuthorId      int64                  `protobuf:"varint,2,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateTopicRequest) Reset() {
	*x = CreateTopicRequest{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateTopicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTopicRequest) ProtoMessage() {}

func (x *CreateTopicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTopicRequest.ProtoReflect.Descriptor instead.
func (*CreateTopicRequest) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{2}
}

func (x *CreateTopicRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CreateTopicRequest) GetAuthorId() int64 {
	if x != nil {
		return x.AuthorId
	}
	return 0
}

type GetTopicRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTopicRequest) Reset() {
	*x = GetTopicRequest{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTopicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTopicRequest) ProtoMessage() {}

func (x *GetTopicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTopicRequest.ProtoReflect.Descriptor instead.
func (*GetTopicRequest) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{3}
}

func (x *GetTopicRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetTopicsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTopicsRequest) Reset() {
	*x = GetTopicsRequest{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTopicsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTopicsRequest) ProtoMessage() {}

func (x *GetTopicsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTopicsRequest.ProtoReflect.Descriptor instead.
func (*GetTopicsRequest) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{4}
}

type GetTopicsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Topics        []*Topic               `protobuf:"bytes,1,rep,name=topics,proto3" json:"topics,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTopicsResponse) Reset() {
	*x = GetTopicsResponse{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTopicsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTopicsResponse) ProtoMessage() {}

func (x *GetTopicsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTopicsResponse.ProtoReflect.Descriptor instead.
func (*GetTopicsResponse) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{5}
}

func (x *GetTopicsResponse) GetTopics() []*Topic {
	if x != nil {
		return x.Topics
	}
	return nil
}

type DeleteTopicRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteTopicRequest) Reset() {
	*x = DeleteTopicRequest{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteTopicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTopicRequest) ProtoMessage() {}

func (x *DeleteTopicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTopicRequest.ProtoReflect.Descriptor instead.
func (*DeleteTopicRequest) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteTopicRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type CreateMessageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TopicId       int64                  `protobuf:"varint,1,opt,name=topic_id,json=topicId,proto3" json:"topic_id,omitempty"`
	UserId        int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Content       string                 `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateMessageRequest) Reset() {
	*x = CreateMessageRequest{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMessageRequest) ProtoMessage() {}

func (x *CreateMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMessageRequest.ProtoReflect.Descriptor instead.
func (*CreateMessageRequest) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{7}
}

func (x *CreateMessageRequest) GetTopicId() int64 {
	if x != nil {
		return x.TopicId
	}
	return 0
}

func (x *CreateMessageRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CreateMessageRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type GetMessagesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TopicId       int64                  `protobuf:"varint,1,opt,name=topic_id,json=topicId,proto3" json:"topic_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMessagesRequest) Reset() {
	*x = GetMessagesRequest{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessagesRequest) ProtoMessage() {}

func (x *GetMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessagesRequest.ProtoReflect.Descriptor instead.
func (*GetMessagesRequest) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{8}
}

func (x *GetMessagesRequest) GetTopicId() int64 {
	if x != nil {
		return x.TopicId
	}
	return 0
}

type GetMessagesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Messages      []*Message             `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMessagesResponse) Reset() {
	*x = GetMessagesResponse{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessagesResponse) ProtoMessage() {}

func (x *GetMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessagesResponse.ProtoReflect.Descriptor instead.
func (*GetMessagesResponse) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{9}
}

func (x *GetMessagesResponse) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

type DeleteMessageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteMessageRequest) Reset() {
	*x = DeleteMessageRequest{}
	mi := &file_proto_forumservice_forum_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteMessageRequest) ProtoMessage() {}

func (x *DeleteMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_forumservice_forum_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteMessageRequest.ProtoReflect.Descriptor instead.
func (*DeleteMessageRequest) Descriptor() ([]byte, []int) {
	return file_proto_forumservice_forum_proto_rawDescGZIP(), []int{10}
}

func (x *DeleteMessageRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_proto_forumservice_forum_proto protoreflect.FileDescriptor

const file_proto_forumservice_forum_proto_rawDesc = "" +
	"\n" +
	"\x1eproto/forumservice/forum.proto\x12\x05forum\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1bgoogle/protobuf/empty.proto\"\x85\x01\n" +
	"\x05Topic\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12\x1b\n" +
	"\tauthor_id\x18\x03 \x01(\x03R\bauthorId\x129\n" +
	"\n" +
	"created_at\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\"\xa2\x01\n" +
	"\aMessage\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x19\n" +
	"\btopic_id\x18\x02 \x01(\x03R\atopicId\x12\x17\n" +
	"\auser_id\x18\x03 \x01(\x03R\x06userId\x12\x18\n" +
	"\acontent\x18\x04 \x01(\tR\acontent\x129\n" +
	"\n" +
	"created_at\x18\x05 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\"G\n" +
	"\x12CreateTopicRequest\x12\x14\n" +
	"\x05title\x18\x01 \x01(\tR\x05title\x12\x1b\n" +
	"\tauthor_id\x18\x02 \x01(\x03R\bauthorId\"!\n" +
	"\x0fGetTopicRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\"\x12\n" +
	"\x10GetTopicsRequest\"9\n" +
	"\x11GetTopicsResponse\x12$\n" +
	"\x06topics\x18\x01 \x03(\v2\f.forum.TopicR\x06topics\"$\n" +
	"\x12DeleteTopicRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\"d\n" +
	"\x14CreateMessageRequest\x12\x19\n" +
	"\btopic_id\x18\x01 \x01(\x03R\atopicId\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x03R\x06userId\x12\x18\n" +
	"\acontent\x18\x03 \x01(\tR\acontent\"/\n" +
	"\x12GetMessagesRequest\x12\x19\n" +
	"\btopic_id\x18\x01 \x01(\x03R\atopicId\"A\n" +
	"\x13GetMessagesResponse\x12*\n" +
	"\bmessages\x18\x01 \x03(\v2\x0e.forum.MessageR\bmessages\"&\n" +
	"\x14DeleteMessageRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id2\xc4\x03\n" +
	"\fForumService\x126\n" +
	"\vCreateTopic\x12\x19.forum.CreateTopicRequest\x1a\f.forum.Topic\x120\n" +
	"\bGetTopic\x12\x16.forum.GetTopicRequest\x1a\f.forum.Topic\x12>\n" +
	"\tGetTopics\x12\x17.forum.GetTopicsRequest\x1a\x18.forum.GetTopicsResponse\x12@\n" +
	"\vDeleteTopic\x12\x19.forum.DeleteTopicRequest\x1a\x16.google.protobuf.Empty\x12<\n" +
	"\rCreateMessage\x12\x1b.forum.CreateMessageRequest\x1a\x0e.forum.Message\x12D\n" +
	"\vGetMessages\x12\x19.forum.GetMessagesRequest\x1a\x1a.forum.GetMessagesResponse\x12D\n" +
	"\rDeleteMessage\x12\x1b.forum.DeleteMessageRequest\x1a\x16.google.protobuf.EmptyB6Z4github.com/Van-programan/Forum_GO/proto/forumserviceb\x06proto3"

var (
	file_proto_forumservice_forum_proto_rawDescOnce sync.Once
	file_proto_forumservice_forum_proto_rawDescData []byte
)

func file_proto_forumservice_forum_proto_rawDescGZIP() []byte {
	file_proto_forumservice_forum_proto_rawDescOnce.Do(func() {
		file_proto_forumservice_forum_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_forumservice_forum_proto_rawDesc), len(file_proto_forumservice_forum_proto_rawDesc)))
	})
	return file_proto_forumservice_forum_proto_rawDescData
}

var file_proto_forumservice_forum_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_forumservice_forum_proto_goTypes = []any{
	(*Topic)(nil),                 // 0: forum.Topic
	(*Message)(nil),               // 1: forum.Message
	(*CreateTopicRequest)(nil),    // 2: forum.CreateTopicRequest
	(*GetTopicRequest)(nil),       // 3: forum.GetTopicRequest
	(*GetTopicsRequest)(nil),      // 4: forum.GetTopicsRequest
	(*GetTopicsResponse)(nil),     // 5: forum.GetTopicsResponse
	(*DeleteTopicRequest)(nil),    // 6: forum.DeleteTopicRequest
	(*CreateMessageRequest)(nil),  // 7: forum.CreateMessageRequest
	(*GetMessagesRequest)(nil),    // 8: forum.GetMessagesRequest
	(*GetMessagesResponse)(nil),   // 9: forum.GetMessagesResponse
	(*DeleteMessageRequest)(nil),  // 10: forum.DeleteMessageRequest
	(*timestamppb.Timestamp)(nil), // 11: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 12: google.protobuf.Empty
}
var file_proto_forumservice_forum_proto_depIdxs = []int32{
	11, // 0: forum.Topic.created_at:type_name -> google.protobuf.Timestamp
	11, // 1: forum.Message.created_at:type_name -> google.protobuf.Timestamp
	0,  // 2: forum.GetTopicsResponse.topics:type_name -> forum.Topic
	1,  // 3: forum.GetMessagesResponse.messages:type_name -> forum.Message
	2,  // 4: forum.ForumService.CreateTopic:input_type -> forum.CreateTopicRequest
	3,  // 5: forum.ForumService.GetTopic:input_type -> forum.GetTopicRequest
	4,  // 6: forum.ForumService.GetTopics:input_type -> forum.GetTopicsRequest
	6,  // 7: forum.ForumService.DeleteTopic:input_type -> forum.DeleteTopicRequest
	7,  // 8: forum.ForumService.CreateMessage:input_type -> forum.CreateMessageRequest
	8,  // 9: forum.ForumService.GetMessages:input_type -> forum.GetMessagesRequest
	10, // 10: forum.ForumService.DeleteMessage:input_type -> forum.DeleteMessageRequest
	0,  // 11: forum.ForumService.CreateTopic:output_type -> forum.Topic
	0,  // 12: forum.ForumService.GetTopic:output_type -> forum.Topic
	5,  // 13: forum.ForumService.GetTopics:output_type -> forum.GetTopicsResponse
	12, // 14: forum.ForumService.DeleteTopic:output_type -> google.protobuf.Empty
	1,  // 15: forum.ForumService.CreateMessage:output_type -> forum.Message
	9,  // 16: forum.ForumService.GetMessages:output_type -> forum.GetMessagesResponse
	12, // 17: forum.ForumService.DeleteMessage:output_type -> google.protobuf.Empty
	11, // [11:18] is the sub-list for method output_type
	4,  // [4:11] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_proto_forumservice_forum_proto_init() }
func file_proto_forumservice_forum_proto_init() {
	if File_proto_forumservice_forum_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_forumservice_forum_proto_rawDesc), len(file_proto_forumservice_forum_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_forumservice_forum_proto_goTypes,
		DependencyIndexes: file_proto_forumservice_forum_proto_depIdxs,
		MessageInfos:      file_proto_forumservice_forum_proto_msgTypes,
	}.Build()
	File_proto_forumservice_forum_proto = out.File
	file_proto_forumservice_forum_proto_goTypes = nil
	file_proto_forumservice_forum_proto_depIdxs = nil
}
