// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/showcase/v1beta1/echo.proto

package genproto // import "github.com/googleapis/gapic-showcase/server/genproto"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import duration "github.com/golang/protobuf/ptypes/duration"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import longrunning "google.golang.org/genproto/googleapis/longrunning"
import status "google.golang.org/genproto/googleapis/rpc/status"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message used for the Echo, Collect and Chat methods. If content
// is set in this message then the request will succeed. If a status is
type EchoRequest struct {
	// Types that are valid to be assigned to Response:
	//	*EchoRequest_Content
	//	*EchoRequest_Error
	Response             isEchoRequest_Response `protobuf_oneof:"response"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *EchoRequest) Reset()         { *m = EchoRequest{} }
func (m *EchoRequest) String() string { return proto.CompactTextString(m) }
func (*EchoRequest) ProtoMessage()    {}
func (*EchoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_a6cf1584ae086d5e, []int{0}
}
func (m *EchoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoRequest.Unmarshal(m, b)
}
func (m *EchoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoRequest.Marshal(b, m, deterministic)
}
func (dst *EchoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoRequest.Merge(dst, src)
}
func (m *EchoRequest) XXX_Size() int {
	return xxx_messageInfo_EchoRequest.Size(m)
}
func (m *EchoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EchoRequest proto.InternalMessageInfo

type isEchoRequest_Response interface {
	isEchoRequest_Response()
}

type EchoRequest_Content struct {
	Content string `protobuf:"bytes,1,opt,name=content,proto3,oneof"`
}

type EchoRequest_Error struct {
	Error *status.Status `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*EchoRequest_Content) isEchoRequest_Response() {}

func (*EchoRequest_Error) isEchoRequest_Response() {}

func (m *EchoRequest) GetResponse() isEchoRequest_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *EchoRequest) GetContent() string {
	if x, ok := m.GetResponse().(*EchoRequest_Content); ok {
		return x.Content
	}
	return ""
}

func (m *EchoRequest) GetError() *status.Status {
	if x, ok := m.GetResponse().(*EchoRequest_Error); ok {
		return x.Error
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*EchoRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _EchoRequest_OneofMarshaler, _EchoRequest_OneofUnmarshaler, _EchoRequest_OneofSizer, []interface{}{
		(*EchoRequest_Content)(nil),
		(*EchoRequest_Error)(nil),
	}
}

func _EchoRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*EchoRequest)
	// response
	switch x := m.Response.(type) {
	case *EchoRequest_Content:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Content)
	case *EchoRequest_Error:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Error); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("EchoRequest.Response has unexpected type %T", x)
	}
	return nil
}

func _EchoRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*EchoRequest)
	switch tag {
	case 1: // response.content
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Response = &EchoRequest_Content{x}
		return true, err
	case 2: // response.error
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(status.Status)
		err := b.DecodeMessage(msg)
		m.Response = &EchoRequest_Error{msg}
		return true, err
	default:
		return false, nil
	}
}

func _EchoRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*EchoRequest)
	// response
	switch x := m.Response.(type) {
	case *EchoRequest_Content:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Content)))
		n += len(x.Content)
	case *EchoRequest_Error:
		s := proto.Size(x.Error)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// The response message for the Echo methods.
type EchoResponse struct {
	// The content specified in the request.
	Content              string   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EchoResponse) Reset()         { *m = EchoResponse{} }
func (m *EchoResponse) String() string { return proto.CompactTextString(m) }
func (*EchoResponse) ProtoMessage()    {}
func (*EchoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_a6cf1584ae086d5e, []int{1}
}
func (m *EchoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoResponse.Unmarshal(m, b)
}
func (m *EchoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoResponse.Marshal(b, m, deterministic)
}
func (dst *EchoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoResponse.Merge(dst, src)
}
func (m *EchoResponse) XXX_Size() int {
	return xxx_messageInfo_EchoResponse.Size(m)
}
func (m *EchoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EchoResponse proto.InternalMessageInfo

func (m *EchoResponse) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

// The request message for the Expand method.
type ExpandRequest struct {
	// The content that will be split into words and returned on the stream.
	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	// The error that is thrown after all words are sent on the stream.
	Error                *status.Status `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ExpandRequest) Reset()         { *m = ExpandRequest{} }
func (m *ExpandRequest) String() string { return proto.CompactTextString(m) }
func (*ExpandRequest) ProtoMessage()    {}
func (*ExpandRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_a6cf1584ae086d5e, []int{2}
}
func (m *ExpandRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExpandRequest.Unmarshal(m, b)
}
func (m *ExpandRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExpandRequest.Marshal(b, m, deterministic)
}
func (dst *ExpandRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExpandRequest.Merge(dst, src)
}
func (m *ExpandRequest) XXX_Size() int {
	return xxx_messageInfo_ExpandRequest.Size(m)
}
func (m *ExpandRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExpandRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExpandRequest proto.InternalMessageInfo

func (m *ExpandRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *ExpandRequest) GetError() *status.Status {
	if m != nil {
		return m.Error
	}
	return nil
}

// The request for the PagedExpand method.
type PagedExpandRequest struct {
	// The string to expand.
	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	// The amount of words to returned in each page.
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// The position of the page to be returned.
	PageToken            string   `protobuf:"bytes,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PagedExpandRequest) Reset()         { *m = PagedExpandRequest{} }
func (m *PagedExpandRequest) String() string { return proto.CompactTextString(m) }
func (*PagedExpandRequest) ProtoMessage()    {}
func (*PagedExpandRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_a6cf1584ae086d5e, []int{3}
}
func (m *PagedExpandRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PagedExpandRequest.Unmarshal(m, b)
}
func (m *PagedExpandRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PagedExpandRequest.Marshal(b, m, deterministic)
}
func (dst *PagedExpandRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PagedExpandRequest.Merge(dst, src)
}
func (m *PagedExpandRequest) XXX_Size() int {
	return xxx_messageInfo_PagedExpandRequest.Size(m)
}
func (m *PagedExpandRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PagedExpandRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PagedExpandRequest proto.InternalMessageInfo

func (m *PagedExpandRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *PagedExpandRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *PagedExpandRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

// The response for the PagedExpand method.
type PagedExpandResponse struct {
	// The words that were expanded.
	Responses []*EchoResponse `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
	// The next page token.
	NextPageToken        string   `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PagedExpandResponse) Reset()         { *m = PagedExpandResponse{} }
func (m *PagedExpandResponse) String() string { return proto.CompactTextString(m) }
func (*PagedExpandResponse) ProtoMessage()    {}
func (*PagedExpandResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_a6cf1584ae086d5e, []int{4}
}
func (m *PagedExpandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PagedExpandResponse.Unmarshal(m, b)
}
func (m *PagedExpandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PagedExpandResponse.Marshal(b, m, deterministic)
}
func (dst *PagedExpandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PagedExpandResponse.Merge(dst, src)
}
func (m *PagedExpandResponse) XXX_Size() int {
	return xxx_messageInfo_PagedExpandResponse.Size(m)
}
func (m *PagedExpandResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PagedExpandResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PagedExpandResponse proto.InternalMessageInfo

func (m *PagedExpandResponse) GetResponses() []*EchoResponse {
	if m != nil {
		return m.Responses
	}
	return nil
}

func (m *PagedExpandResponse) GetNextPageToken() string {
	if m != nil {
		return m.NextPageToken
	}
	return ""
}

// The request for Wait method.
type WaitRequest struct {
	// Types that are valid to be assigned to End:
	//	*WaitRequest_EndTime
	//	*WaitRequest_Ttl
	End isWaitRequest_End `protobuf_oneof:"end"`
	// Types that are valid to be assigned to Response:
	//	*WaitRequest_Error
	//	*WaitRequest_Success
	Response             isWaitRequest_Response `protobuf_oneof:"response"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *WaitRequest) Reset()         { *m = WaitRequest{} }
func (m *WaitRequest) String() string { return proto.CompactTextString(m) }
func (*WaitRequest) ProtoMessage()    {}
func (*WaitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_a6cf1584ae086d5e, []int{5}
}
func (m *WaitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WaitRequest.Unmarshal(m, b)
}
func (m *WaitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WaitRequest.Marshal(b, m, deterministic)
}
func (dst *WaitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WaitRequest.Merge(dst, src)
}
func (m *WaitRequest) XXX_Size() int {
	return xxx_messageInfo_WaitRequest.Size(m)
}
func (m *WaitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WaitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WaitRequest proto.InternalMessageInfo

type isWaitRequest_End interface {
	isWaitRequest_End()
}

type WaitRequest_EndTime struct {
	EndTime *timestamp.Timestamp `protobuf:"bytes,1,opt,name=end_time,json=endTime,proto3,oneof"`
}

type WaitRequest_Ttl struct {
	Ttl *duration.Duration `protobuf:"bytes,4,opt,name=ttl,proto3,oneof"`
}

func (*WaitRequest_EndTime) isWaitRequest_End() {}

func (*WaitRequest_Ttl) isWaitRequest_End() {}

func (m *WaitRequest) GetEnd() isWaitRequest_End {
	if m != nil {
		return m.End
	}
	return nil
}

func (m *WaitRequest) GetEndTime() *timestamp.Timestamp {
	if x, ok := m.GetEnd().(*WaitRequest_EndTime); ok {
		return x.EndTime
	}
	return nil
}

func (m *WaitRequest) GetTtl() *duration.Duration {
	if x, ok := m.GetEnd().(*WaitRequest_Ttl); ok {
		return x.Ttl
	}
	return nil
}

type isWaitRequest_Response interface {
	isWaitRequest_Response()
}

type WaitRequest_Error struct {
	Error *status.Status `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

type WaitRequest_Success struct {
	Success *WaitResponse `protobuf:"bytes,3,opt,name=success,proto3,oneof"`
}

func (*WaitRequest_Error) isWaitRequest_Response() {}

func (*WaitRequest_Success) isWaitRequest_Response() {}

func (m *WaitRequest) GetResponse() isWaitRequest_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *WaitRequest) GetError() *status.Status {
	if x, ok := m.GetResponse().(*WaitRequest_Error); ok {
		return x.Error
	}
	return nil
}

func (m *WaitRequest) GetSuccess() *WaitResponse {
	if x, ok := m.GetResponse().(*WaitRequest_Success); ok {
		return x.Success
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*WaitRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _WaitRequest_OneofMarshaler, _WaitRequest_OneofUnmarshaler, _WaitRequest_OneofSizer, []interface{}{
		(*WaitRequest_EndTime)(nil),
		(*WaitRequest_Ttl)(nil),
		(*WaitRequest_Error)(nil),
		(*WaitRequest_Success)(nil),
	}
}

func _WaitRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*WaitRequest)
	// end
	switch x := m.End.(type) {
	case *WaitRequest_EndTime:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.EndTime); err != nil {
			return err
		}
	case *WaitRequest_Ttl:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Ttl); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("WaitRequest.End has unexpected type %T", x)
	}
	// response
	switch x := m.Response.(type) {
	case *WaitRequest_Error:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Error); err != nil {
			return err
		}
	case *WaitRequest_Success:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Success); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("WaitRequest.Response has unexpected type %T", x)
	}
	return nil
}

func _WaitRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*WaitRequest)
	switch tag {
	case 1: // end.end_time
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(timestamp.Timestamp)
		err := b.DecodeMessage(msg)
		m.End = &WaitRequest_EndTime{msg}
		return true, err
	case 4: // end.ttl
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(duration.Duration)
		err := b.DecodeMessage(msg)
		m.End = &WaitRequest_Ttl{msg}
		return true, err
	case 2: // response.error
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(status.Status)
		err := b.DecodeMessage(msg)
		m.Response = &WaitRequest_Error{msg}
		return true, err
	case 3: // response.success
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(WaitResponse)
		err := b.DecodeMessage(msg)
		m.Response = &WaitRequest_Success{msg}
		return true, err
	default:
		return false, nil
	}
}

func _WaitRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*WaitRequest)
	// end
	switch x := m.End.(type) {
	case *WaitRequest_EndTime:
		s := proto.Size(x.EndTime)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *WaitRequest_Ttl:
		s := proto.Size(x.Ttl)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	// response
	switch x := m.Response.(type) {
	case *WaitRequest_Error:
		s := proto.Size(x.Error)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *WaitRequest_Success:
		s := proto.Size(x.Success)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// The result of the Wait operation.
type WaitResponse struct {
	// This content of the result.
	Content              string   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WaitResponse) Reset()         { *m = WaitResponse{} }
func (m *WaitResponse) String() string { return proto.CompactTextString(m) }
func (*WaitResponse) ProtoMessage()    {}
func (*WaitResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_a6cf1584ae086d5e, []int{6}
}
func (m *WaitResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WaitResponse.Unmarshal(m, b)
}
func (m *WaitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WaitResponse.Marshal(b, m, deterministic)
}
func (dst *WaitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WaitResponse.Merge(dst, src)
}
func (m *WaitResponse) XXX_Size() int {
	return xxx_messageInfo_WaitResponse.Size(m)
}
func (m *WaitResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_WaitResponse.DiscardUnknown(m)
}

var xxx_messageInfo_WaitResponse proto.InternalMessageInfo

func (m *WaitResponse) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

// The metadata for Wait operation.
type WaitMetadata struct {
	// The time that this operation will complete.
	EndTime              *timestamp.Timestamp `protobuf:"bytes,1,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *WaitMetadata) Reset()         { *m = WaitMetadata{} }
func (m *WaitMetadata) String() string { return proto.CompactTextString(m) }
func (*WaitMetadata) ProtoMessage()    {}
func (*WaitMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_a6cf1584ae086d5e, []int{7}
}
func (m *WaitMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WaitMetadata.Unmarshal(m, b)
}
func (m *WaitMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WaitMetadata.Marshal(b, m, deterministic)
}
func (dst *WaitMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WaitMetadata.Merge(dst, src)
}
func (m *WaitMetadata) XXX_Size() int {
	return xxx_messageInfo_WaitMetadata.Size(m)
}
func (m *WaitMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_WaitMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_WaitMetadata proto.InternalMessageInfo

func (m *WaitMetadata) GetEndTime() *timestamp.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

func init() {
	proto.RegisterType((*EchoRequest)(nil), "google.showcase.v1beta1.EchoRequest")
	proto.RegisterType((*EchoResponse)(nil), "google.showcase.v1beta1.EchoResponse")
	proto.RegisterType((*ExpandRequest)(nil), "google.showcase.v1beta1.ExpandRequest")
	proto.RegisterType((*PagedExpandRequest)(nil), "google.showcase.v1beta1.PagedExpandRequest")
	proto.RegisterType((*PagedExpandResponse)(nil), "google.showcase.v1beta1.PagedExpandResponse")
	proto.RegisterType((*WaitRequest)(nil), "google.showcase.v1beta1.WaitRequest")
	proto.RegisterType((*WaitResponse)(nil), "google.showcase.v1beta1.WaitResponse")
	proto.RegisterType((*WaitMetadata)(nil), "google.showcase.v1beta1.WaitMetadata")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EchoClient is the client API for Echo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EchoClient interface {
	// This method simply echos the request. This method is showcases unary rpcs.
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error)
	// This method split the given content into words and will pass each word back
	// through the stream. This method showcases server-side streaming rpcs.
	Expand(ctx context.Context, in *ExpandRequest, opts ...grpc.CallOption) (Echo_ExpandClient, error)
	// This method will collect the words given to it. When the stream is closed
	// by the client, this method will return the a concatenation of the strings
	// passed to it. This method showcases client-side streaming rpcs.
	Collect(ctx context.Context, opts ...grpc.CallOption) (Echo_CollectClient, error)
	// This method, upon receiving a request on the stream, the same content will
	// be passed  back on the stream. This method showcases bidirectional
	// streaming rpcs.
	Chat(ctx context.Context, opts ...grpc.CallOption) (Echo_ChatClient, error)
	// This is similar to the Expand method but instead of returning a stream of
	// expanded words, this method returns a paged list of expanded words.
	PagedExpand(ctx context.Context, in *PagedExpandRequest, opts ...grpc.CallOption) (*PagedExpandResponse, error)
	// This method will wait the requested amount of and then return.
	// This method showcases how a client handles a request timing out.
	Wait(ctx context.Context, in *WaitRequest, opts ...grpc.CallOption) (*longrunning.Operation, error)
}

type echoClient struct {
	cc *grpc.ClientConn
}

func NewEchoClient(cc *grpc.ClientConn) EchoClient {
	return &echoClient{cc}
}

func (c *echoClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, "/google.showcase.v1beta1.Echo/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoClient) Expand(ctx context.Context, in *ExpandRequest, opts ...grpc.CallOption) (Echo_ExpandClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Echo_serviceDesc.Streams[0], "/google.showcase.v1beta1.Echo/Expand", opts...)
	if err != nil {
		return nil, err
	}
	x := &echoExpandClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Echo_ExpandClient interface {
	Recv() (*EchoResponse, error)
	grpc.ClientStream
}

type echoExpandClient struct {
	grpc.ClientStream
}

func (x *echoExpandClient) Recv() (*EchoResponse, error) {
	m := new(EchoResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *echoClient) Collect(ctx context.Context, opts ...grpc.CallOption) (Echo_CollectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Echo_serviceDesc.Streams[1], "/google.showcase.v1beta1.Echo/Collect", opts...)
	if err != nil {
		return nil, err
	}
	x := &echoCollectClient{stream}
	return x, nil
}

type Echo_CollectClient interface {
	Send(*EchoRequest) error
	CloseAndRecv() (*EchoResponse, error)
	grpc.ClientStream
}

type echoCollectClient struct {
	grpc.ClientStream
}

func (x *echoCollectClient) Send(m *EchoRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *echoCollectClient) CloseAndRecv() (*EchoResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(EchoResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *echoClient) Chat(ctx context.Context, opts ...grpc.CallOption) (Echo_ChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Echo_serviceDesc.Streams[2], "/google.showcase.v1beta1.Echo/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &echoChatClient{stream}
	return x, nil
}

type Echo_ChatClient interface {
	Send(*EchoRequest) error
	Recv() (*EchoResponse, error)
	grpc.ClientStream
}

type echoChatClient struct {
	grpc.ClientStream
}

func (x *echoChatClient) Send(m *EchoRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *echoChatClient) Recv() (*EchoResponse, error) {
	m := new(EchoResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *echoClient) PagedExpand(ctx context.Context, in *PagedExpandRequest, opts ...grpc.CallOption) (*PagedExpandResponse, error) {
	out := new(PagedExpandResponse)
	err := c.cc.Invoke(ctx, "/google.showcase.v1beta1.Echo/PagedExpand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoClient) Wait(ctx context.Context, in *WaitRequest, opts ...grpc.CallOption) (*longrunning.Operation, error) {
	out := new(longrunning.Operation)
	err := c.cc.Invoke(ctx, "/google.showcase.v1beta1.Echo/Wait", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EchoServer is the server API for Echo service.
type EchoServer interface {
	// This method simply echos the request. This method is showcases unary rpcs.
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
	// This method split the given content into words and will pass each word back
	// through the stream. This method showcases server-side streaming rpcs.
	Expand(*ExpandRequest, Echo_ExpandServer) error
	// This method will collect the words given to it. When the stream is closed
	// by the client, this method will return the a concatenation of the strings
	// passed to it. This method showcases client-side streaming rpcs.
	Collect(Echo_CollectServer) error
	// This method, upon receiving a request on the stream, the same content will
	// be passed  back on the stream. This method showcases bidirectional
	// streaming rpcs.
	Chat(Echo_ChatServer) error
	// This is similar to the Expand method but instead of returning a stream of
	// expanded words, this method returns a paged list of expanded words.
	PagedExpand(context.Context, *PagedExpandRequest) (*PagedExpandResponse, error)
	// This method will wait the requested amount of and then return.
	// This method showcases how a client handles a request timing out.
	Wait(context.Context, *WaitRequest) (*longrunning.Operation, error)
}

func RegisterEchoServer(s *grpc.Server, srv EchoServer) {
	s.RegisterService(&_Echo_serviceDesc, srv)
}

func _Echo_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.showcase.v1beta1.Echo/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Echo_Expand_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExpandRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EchoServer).Expand(m, &echoExpandServer{stream})
}

type Echo_ExpandServer interface {
	Send(*EchoResponse) error
	grpc.ServerStream
}

type echoExpandServer struct {
	grpc.ServerStream
}

func (x *echoExpandServer) Send(m *EchoResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Echo_Collect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EchoServer).Collect(&echoCollectServer{stream})
}

type Echo_CollectServer interface {
	SendAndClose(*EchoResponse) error
	Recv() (*EchoRequest, error)
	grpc.ServerStream
}

type echoCollectServer struct {
	grpc.ServerStream
}

func (x *echoCollectServer) SendAndClose(m *EchoResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *echoCollectServer) Recv() (*EchoRequest, error) {
	m := new(EchoRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Echo_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EchoServer).Chat(&echoChatServer{stream})
}

type Echo_ChatServer interface {
	Send(*EchoResponse) error
	Recv() (*EchoRequest, error)
	grpc.ServerStream
}

type echoChatServer struct {
	grpc.ServerStream
}

func (x *echoChatServer) Send(m *EchoResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *echoChatServer) Recv() (*EchoRequest, error) {
	m := new(EchoRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Echo_PagedExpand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PagedExpandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServer).PagedExpand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.showcase.v1beta1.Echo/PagedExpand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServer).PagedExpand(ctx, req.(*PagedExpandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Echo_Wait_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WaitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServer).Wait(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.showcase.v1beta1.Echo/Wait",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServer).Wait(ctx, req.(*WaitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Echo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.showcase.v1beta1.Echo",
	HandlerType: (*EchoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _Echo_Echo_Handler,
		},
		{
			MethodName: "PagedExpand",
			Handler:    _Echo_PagedExpand_Handler,
		},
		{
			MethodName: "Wait",
			Handler:    _Echo_Wait_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Expand",
			Handler:       _Echo_Expand_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Collect",
			Handler:       _Echo_Collect_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Chat",
			Handler:       _Echo_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "google/showcase/v1beta1/echo.proto",
}

func init() {
	proto.RegisterFile("google/showcase/v1beta1/echo.proto", fileDescriptor_echo_a6cf1584ae086d5e)
}

var fileDescriptor_echo_a6cf1584ae086d5e = []byte{
	// 765 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xcf, 0x4f, 0xdb, 0x48,
	0x14, 0x66, 0x92, 0xf0, 0x6b, 0xb2, 0xec, 0x6a, 0x67, 0x77, 0x45, 0x30, 0x84, 0x8d, 0x5c, 0x8a,
	0xa2, 0x14, 0x6c, 0x48, 0x69, 0x51, 0xa3, 0x5e, 0x12, 0x8a, 0x94, 0x4b, 0x55, 0x14, 0xa8, 0x90,
	0x7a, 0x89, 0x26, 0xf6, 0xc3, 0xb1, 0xea, 0xcc, 0xb8, 0x9e, 0x09, 0x20, 0x8e, 0xdc, 0x5a, 0x55,
	0xbd, 0xf4, 0xbf, 0xe2, 0xda, 0x5b, 0x4f, 0x3d, 0xf4, 0x7f, 0xe8, 0xb5, 0xf2, 0xd8, 0x0e, 0x4e,
	0x68, 0x20, 0x95, 0xb8, 0x24, 0x9a, 0x99, 0xef, 0xbd, 0xef, 0xbd, 0xef, 0x7d, 0xe3, 0xc1, 0xba,
	0xc3, 0xb9, 0xe3, 0x81, 0x29, 0xba, 0xfc, 0xcc, 0xa2, 0x02, 0xcc, 0xd3, 0xed, 0x0e, 0x48, 0xba,
	0x6d, 0x82, 0xd5, 0xe5, 0x86, 0x1f, 0x70, 0xc9, 0xc9, 0x62, 0x84, 0x31, 0x12, 0x8c, 0x11, 0x63,
	0xb4, 0x95, 0x38, 0x98, 0xfa, 0xae, 0x49, 0x19, 0xe3, 0x92, 0x4a, 0x97, 0x33, 0x11, 0x85, 0x69,
	0x8b, 0xa9, 0x53, 0xcb, 0x73, 0x81, 0xc9, 0xf8, 0xe0, 0xff, 0xd4, 0xc1, 0x89, 0x0b, 0x9e, 0xdd,
	0xee, 0x40, 0x97, 0x9e, 0xba, 0x3c, 0x88, 0x01, 0x0f, 0x62, 0x80, 0xc7, 0x99, 0x13, 0xf4, 0x19,
	0x73, 0x99, 0x63, 0x72, 0x1f, 0x82, 0xa1, 0xf4, 0xab, 0x31, 0x48, 0xad, 0x3a, 0xfd, 0x13, 0xd3,
	0xee, 0x47, 0x80, 0x11, 0x96, 0xc1, 0xb9, 0x74, 0x7b, 0x20, 0x24, 0xed, 0xf9, 0x23, 0xf5, 0x05,
	0xbe, 0x65, 0x0a, 0x49, 0x65, 0x3f, 0xce, 0xac, 0x53, 0x9c, 0xdf, 0xb7, 0xba, 0xbc, 0x05, 0xef,
	0xfa, 0x20, 0x24, 0xd1, 0xf0, 0xac, 0xc5, 0x99, 0x04, 0x26, 0x0b, 0xa8, 0x84, 0xca, 0xf3, 0xcd,
	0xa9, 0x56, 0xb2, 0x41, 0x2a, 0x78, 0x1a, 0x82, 0x80, 0x07, 0x85, 0x4c, 0x09, 0x95, 0xf3, 0x55,
	0x62, 0xc4, 0x52, 0x05, 0xbe, 0x65, 0x1c, 0xaa, 0x9c, 0xcd, 0xa9, 0x56, 0x04, 0x69, 0x60, 0x3c,
	0x17, 0x80, 0xf0, 0x39, 0x13, 0xa0, 0x97, 0xf1, 0x1f, 0x11, 0x45, 0xb4, 0x26, 0x85, 0x11, 0x8e,
	0x01, 0x83, 0x7e, 0x88, 0x17, 0xf6, 0xcf, 0x7d, 0xca, 0xec, 0xa4, 0x9c, 0xb1, 0x50, 0x52, 0xbe,
	0xb3, 0x98, 0xb8, 0x14, 0x9d, 0x63, 0x72, 0x40, 0x1d, 0xb0, 0x87, 0x33, 0x17, 0x47, 0x32, 0x37,
	0xb2, 0xdf, 0xea, 0x99, 0xeb, 0xf4, 0xcb, 0x78, 0xde, 0xa7, 0x0e, 0xb4, 0x85, 0x7b, 0x01, 0x8a,
	0x62, 0xba, 0x35, 0x17, 0x6e, 0x1c, 0xba, 0x17, 0x40, 0x8a, 0x18, 0xab, 0x43, 0xc9, 0xdf, 0x02,
	0x2b, 0x64, 0x55, 0x61, 0x0a, 0x7e, 0x14, 0x6e, 0xe8, 0x97, 0x08, 0xff, 0x33, 0xc4, 0x18, 0xf7,
	0xbd, 0x87, 0xe7, 0x13, 0x4d, 0x44, 0x01, 0x95, 0xb2, 0xe5, 0x7c, 0xf5, 0xa1, 0x31, 0xc6, 0x6e,
	0x46, 0x5a, 0xb1, 0xd6, 0x75, 0x1c, 0x59, 0xc7, 0x7f, 0x31, 0x38, 0x97, 0xed, 0x54, 0x01, 0x19,
	0x55, 0xc0, 0x42, 0xb8, 0x7d, 0x30, 0x28, 0xe2, 0x07, 0xc2, 0xf9, 0x63, 0xea, 0xca, 0xa4, 0xdf,
	0x5d, 0x3c, 0x07, 0xcc, 0x6e, 0x87, 0xbe, 0x50, 0x0d, 0xe7, 0xab, 0x5a, 0xc2, 0x9d, 0x98, 0xc6,
	0x38, 0x4a, 0x4c, 0x13, 0x4e, 0x1d, 0x98, 0x1d, 0xae, 0xc9, 0x26, 0xce, 0x4a, 0xe9, 0x15, 0x72,
	0x2a, 0x66, 0xe9, 0x46, 0xcc, 0x8b, 0xd8, 0x88, 0xcd, 0xa9, 0x56, 0x88, 0x9b, 0xc4, 0x24, 0x28,
	0x9e, 0x0c, 0xa9, 0xe3, 0x59, 0xd1, 0xb7, 0x2c, 0x10, 0x42, 0x89, 0x78, 0x9b, 0x1c, 0x51, 0x2b,
	0x91, 0x08, 0x4d, 0xd4, 0x4a, 0xe2, 0x1a, 0xd3, 0x38, 0x0b, 0xcc, 0x1e, 0xb5, 0x5b, 0x1a, 0x7d,
	0x8b, 0xdd, 0xf6, 0x23, 0xe4, 0x4b, 0x90, 0xd4, 0xa6, 0x92, 0x92, 0x27, 0xbf, 0xa3, 0xd1, 0x40,
	0xa1, 0xea, 0xc7, 0x19, 0x9c, 0x0b, 0xc7, 0x45, 0x82, 0xf8, 0x7f, 0xed, 0x8e, 0xa9, 0xaa, 0x89,
	0x68, 0x93, 0xcd, 0x5e, 0x2f, 0x5e, 0x7e, 0xf9, 0xfe, 0x39, 0xb3, 0xa8, 0x93, 0xa1, 0xaf, 0x55,
	0x4d, 0xfd, 0xa0, 0x0a, 0xf9, 0x80, 0xf0, 0x4c, 0xe4, 0x33, 0xb2, 0x3e, 0x3e, 0x61, 0xda, 0xfa,
	0x93, 0x12, 0x9b, 0x5f, 0xeb, 0x0b, 0xb1, 0x52, 0x1b, 0x6a, 0x5e, 0xaa, 0x90, 0x25, 0xfd, 0xdf,
	0x91, 0x42, 0x54, 0xee, 0x1a, 0xaa, 0x6c, 0x21, 0x72, 0x81, 0x67, 0xf7, 0xb8, 0xe7, 0x81, 0x25,
	0xef, 0x57, 0x83, 0x92, 0xa2, 0xd6, 0xf4, 0xff, 0x86, 0xa9, 0xad, 0x88, 0xab, 0x86, 0x2a, 0x65,
	0x44, 0x8e, 0x71, 0x6e, 0xaf, 0x4b, 0xef, 0x97, 0xb8, 0x8c, 0xb6, 0x10, 0xf9, 0x84, 0x70, 0x3e,
	0x75, 0x9d, 0xc9, 0xa3, 0xb1, 0xa1, 0x37, 0x3f, 0x33, 0xda, 0xc6, 0x64, 0xe0, 0xb8, 0xcf, 0x35,
	0xd5, 0xe7, 0xaa, 0xbe, 0x34, 0xdc, 0xa7, 0x7f, 0x0d, 0x0d, 0x47, 0xfe, 0x1e, 0xe1, 0x5c, 0xe8,
	0xdb, 0x5b, 0x5a, 0x4d, 0xdd, 0x7c, 0xad, 0x98, 0xa0, 0x52, 0x2f, 0x8c, 0xf1, 0x2a, 0x79, 0x61,
	0xf4, 0xe7, 0x57, 0xf5, 0x95, 0x91, 0x1b, 0x33, 0x74, 0x2b, 0x7e, 0x6d, 0xbf, 0x33, 0xea, 0x86,
	0xba, 0x6b, 0x7f, 0x5f, 0xd5, 0xff, 0xf4, 0xb8, 0x45, 0xbd, 0x2e, 0x17, 0xb2, 0xb6, 0xbb, 0xf3,
	0xf4, 0x59, 0xe3, 0x35, 0x5e, 0xb6, 0x78, 0x6f, 0x5c, 0x69, 0x07, 0xe8, 0xcd, 0x8e, 0xe3, 0xca,
	0x6e, 0xbf, 0x63, 0x58, 0xbc, 0x67, 0x46, 0x28, 0xea, 0xbb, 0xc2, 0x74, 0xa8, 0xef, 0x5a, 0x9b,
	0x83, 0xb7, 0x59, 0x40, 0x70, 0x0a, 0x81, 0xe9, 0x00, 0x8b, 0xee, 0xde, 0x8c, 0xfa, 0x7b, 0xfc,
	0x33, 0x00, 0x00, 0xff, 0xff, 0x07, 0xd3, 0x4c, 0x6a, 0xc5, 0x07, 0x00, 0x00,
}
