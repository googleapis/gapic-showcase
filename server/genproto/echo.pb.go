// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/showcase/v1alpha3/echo.proto

package genproto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	longrunning "google.golang.org/genproto/googleapis/longrunning"
	status "google.golang.org/genproto/googleapis/rpc/status"
	math "math"
)

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
	return fileDescriptor_1c4be78b9cf935ae, []int{0}
}

func (m *EchoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoRequest.Unmarshal(m, b)
}
func (m *EchoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoRequest.Marshal(b, m, deterministic)
}
func (m *EchoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoRequest.Merge(m, src)
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
	return fileDescriptor_1c4be78b9cf935ae, []int{1}
}

func (m *EchoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoResponse.Unmarshal(m, b)
}
func (m *EchoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoResponse.Marshal(b, m, deterministic)
}
func (m *EchoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoResponse.Merge(m, src)
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
	return fileDescriptor_1c4be78b9cf935ae, []int{2}
}

func (m *ExpandRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExpandRequest.Unmarshal(m, b)
}
func (m *ExpandRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExpandRequest.Marshal(b, m, deterministic)
}
func (m *ExpandRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExpandRequest.Merge(m, src)
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
	return fileDescriptor_1c4be78b9cf935ae, []int{3}
}

func (m *PagedExpandRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PagedExpandRequest.Unmarshal(m, b)
}
func (m *PagedExpandRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PagedExpandRequest.Marshal(b, m, deterministic)
}
func (m *PagedExpandRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PagedExpandRequest.Merge(m, src)
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
	return fileDescriptor_1c4be78b9cf935ae, []int{4}
}

func (m *PagedExpandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PagedExpandResponse.Unmarshal(m, b)
}
func (m *PagedExpandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PagedExpandResponse.Marshal(b, m, deterministic)
}
func (m *PagedExpandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PagedExpandResponse.Merge(m, src)
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
	// The time that this operation will complete.
	EndTime *timestamp.Timestamp `protobuf:"bytes,1,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
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
	return fileDescriptor_1c4be78b9cf935ae, []int{5}
}

func (m *WaitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WaitRequest.Unmarshal(m, b)
}
func (m *WaitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WaitRequest.Marshal(b, m, deterministic)
}
func (m *WaitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WaitRequest.Merge(m, src)
}
func (m *WaitRequest) XXX_Size() int {
	return xxx_messageInfo_WaitRequest.Size(m)
}
func (m *WaitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WaitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WaitRequest proto.InternalMessageInfo

func (m *WaitRequest) GetEndTime() *timestamp.Timestamp {
	if m != nil {
		return m.EndTime
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
		(*WaitRequest_Error)(nil),
		(*WaitRequest_Success)(nil),
	}
}

func _WaitRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*WaitRequest)
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
	return fileDescriptor_1c4be78b9cf935ae, []int{6}
}

func (m *WaitResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WaitResponse.Unmarshal(m, b)
}
func (m *WaitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WaitResponse.Marshal(b, m, deterministic)
}
func (m *WaitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WaitResponse.Merge(m, src)
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
	return fileDescriptor_1c4be78b9cf935ae, []int{7}
}

func (m *WaitMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WaitMetadata.Unmarshal(m, b)
}
func (m *WaitMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WaitMetadata.Marshal(b, m, deterministic)
}
func (m *WaitMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WaitMetadata.Merge(m, src)
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
	proto.RegisterType((*EchoRequest)(nil), "google.showcase.v1alpha3.EchoRequest")
	proto.RegisterType((*EchoResponse)(nil), "google.showcase.v1alpha3.EchoResponse")
	proto.RegisterType((*ExpandRequest)(nil), "google.showcase.v1alpha3.ExpandRequest")
	proto.RegisterType((*PagedExpandRequest)(nil), "google.showcase.v1alpha3.PagedExpandRequest")
	proto.RegisterType((*PagedExpandResponse)(nil), "google.showcase.v1alpha3.PagedExpandResponse")
	proto.RegisterType((*WaitRequest)(nil), "google.showcase.v1alpha3.WaitRequest")
	proto.RegisterType((*WaitResponse)(nil), "google.showcase.v1alpha3.WaitResponse")
	proto.RegisterType((*WaitMetadata)(nil), "google.showcase.v1alpha3.WaitMetadata")
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
	err := c.cc.Invoke(ctx, "/google.showcase.v1alpha3.Echo/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoClient) Expand(ctx context.Context, in *ExpandRequest, opts ...grpc.CallOption) (Echo_ExpandClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Echo_serviceDesc.Streams[0], "/google.showcase.v1alpha3.Echo/Expand", opts...)
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
	stream, err := c.cc.NewStream(ctx, &_Echo_serviceDesc.Streams[1], "/google.showcase.v1alpha3.Echo/Collect", opts...)
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
	stream, err := c.cc.NewStream(ctx, &_Echo_serviceDesc.Streams[2], "/google.showcase.v1alpha3.Echo/Chat", opts...)
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
	err := c.cc.Invoke(ctx, "/google.showcase.v1alpha3.Echo/PagedExpand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoClient) Wait(ctx context.Context, in *WaitRequest, opts ...grpc.CallOption) (*longrunning.Operation, error) {
	out := new(longrunning.Operation)
	err := c.cc.Invoke(ctx, "/google.showcase.v1alpha3.Echo/Wait", in, out, opts...)
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
		FullMethod: "/google.showcase.v1alpha3.Echo/Echo",
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
		FullMethod: "/google.showcase.v1alpha3.Echo/PagedExpand",
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
		FullMethod: "/google.showcase.v1alpha3.Echo/Wait",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServer).Wait(ctx, req.(*WaitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Echo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.showcase.v1alpha3.Echo",
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
	Metadata: "google/showcase/v1alpha3/echo.proto",
}

func init() {
	proto.RegisterFile("google/showcase/v1alpha3/echo.proto", fileDescriptor_1c4be78b9cf935ae)
}

var fileDescriptor_1c4be78b9cf935ae = []byte{
	// 703 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xcf, 0x6f, 0xd3, 0x3e,
	0x1c, 0x9d, 0xbb, 0x1f, 0xdd, 0xdc, 0xef, 0xbe, 0x80, 0x27, 0x20, 0xca, 0x36, 0x36, 0x82, 0x36,
	0xa2, 0x6a, 0x4b, 0x46, 0x37, 0x40, 0xf4, 0xd6, 0x8e, 0x49, 0xbb, 0x20, 0xa6, 0x6e, 0x12, 0x82,
	0xcb, 0xe4, 0xb9, 0x26, 0x89, 0x48, 0x6d, 0x13, 0x3b, 0xdb, 0x34, 0x89, 0x0b, 0x1c, 0xb8, 0x82,
	0xf8, 0x87, 0x10, 0x47, 0xae, 0xdc, 0x38, 0x71, 0xe0, 0xff, 0x00, 0xc5, 0x49, 0xd6, 0xb4, 0xa8,
	0x5d, 0x41, 0xbb, 0x34, 0xb2, 0xfb, 0x3e, 0x9f, 0xf7, 0xde, 0xc7, 0xcf, 0x09, 0xbc, 0xe3, 0x71,
	0xee, 0x85, 0xd4, 0x95, 0x3e, 0x3f, 0x21, 0x58, 0x52, 0xf7, 0xf8, 0x1e, 0x0e, 0x85, 0x8f, 0x37,
	0x5d, 0x4a, 0x7c, 0xee, 0x88, 0x88, 0x2b, 0x8e, 0x8c, 0x14, 0xe4, 0xe4, 0x20, 0x27, 0x07, 0x99,
	0x0b, 0x59, 0x39, 0x16, 0x81, 0x8b, 0x19, 0xe3, 0x0a, 0xab, 0x80, 0x33, 0x99, 0xd6, 0x99, 0x79,
	0xf3, 0x90, 0x33, 0x2f, 0x8a, 0x19, 0x0b, 0x98, 0xe7, 0x72, 0x41, 0xa3, 0x1e, 0xd0, 0x52, 0x06,
	0xd2, 0xab, 0xa3, 0xf8, 0xa5, 0xab, 0x82, 0x0e, 0x95, 0x0a, 0x77, 0x44, 0x06, 0xb8, 0x99, 0x01,
	0x22, 0x41, 0x5c, 0xa9, 0xb0, 0x8a, 0xb3, 0x4a, 0x0b, 0xc3, 0xca, 0x0e, 0xf1, 0x79, 0x8b, 0xbe,
	0x8e, 0xa9, 0x54, 0xc8, 0x84, 0x65, 0xc2, 0x99, 0xa2, 0x4c, 0x19, 0x60, 0x19, 0xd8, 0x33, 0xbb,
	0x63, 0xad, 0x7c, 0x03, 0x55, 0xe1, 0x24, 0x8d, 0x22, 0x1e, 0x19, 0xa5, 0x65, 0x60, 0x57, 0x6a,
	0xc8, 0xc9, 0x1c, 0x45, 0x82, 0x38, 0xfb, 0xba, 0xe7, 0xee, 0x58, 0x2b, 0x85, 0x34, 0x21, 0x9c,
	0x8e, 0xa8, 0x14, 0x9c, 0x49, 0x6a, 0xd9, 0xf0, 0xbf, 0x94, 0x22, 0x5d, 0x23, 0xa3, 0x8f, 0xe3,
	0x9c, 0xc1, 0xda, 0x87, 0xb3, 0x3b, 0xa7, 0x02, 0xb3, 0x76, 0x2e, 0x67, 0x20, 0x14, 0xd9, 0x17,
	0x8a, 0xc9, 0xa4, 0x58, 0x1c, 0xa2, 0x3d, 0xec, 0xd1, 0x76, 0x6f, 0xe7, 0xc5, 0xbe, 0xce, 0xcd,
	0xf1, 0x1f, 0x8d, 0x52, 0xb7, 0xfd, 0x3c, 0x9c, 0x11, 0xd8, 0xa3, 0x87, 0x32, 0x38, 0xa3, 0x9a,
	0x62, 0xb2, 0x35, 0x9d, 0x6c, 0xec, 0x07, 0x67, 0x14, 0x2d, 0x42, 0xa8, 0xff, 0x54, 0xfc, 0x15,
	0x65, 0xc6, 0xb8, 0x16, 0xa6, 0xe1, 0x07, 0xc9, 0x86, 0xf5, 0x0e, 0xc0, 0xb9, 0x1e, 0xc6, 0xcc,
	0xf7, 0x63, 0x38, 0x93, 0xcf, 0x44, 0x1a, 0x60, 0x79, 0xdc, 0xae, 0xd4, 0x56, 0x9d, 0x41, 0xa9,
	0x70, 0x8a, 0x23, 0x6b, 0x75, 0x0b, 0xd1, 0x2a, 0xbc, 0xc2, 0xe8, 0xa9, 0x3a, 0x2c, 0x28, 0x28,
	0x69, 0x05, 0xb3, 0xc9, 0xf6, 0xde, 0xb9, 0x8a, 0xcf, 0x00, 0x56, 0x9e, 0xe1, 0x40, 0xe5, 0x86,
	0xef, 0xc3, 0x69, 0xca, 0xda, 0x87, 0x49, 0x30, 0xb4, 0xe3, 0x4a, 0xcd, 0xcc, 0xc9, 0xf3, 0xd4,
	0x38, 0x07, 0x79, 0x6a, 0x5a, 0x65, 0xca, 0xda, 0xc9, 0xea, 0x6f, 0x0e, 0x1d, 0x35, 0x61, 0x59,
	0xc6, 0x84, 0x50, 0x29, 0xf5, 0x50, 0x86, 0xda, 0x4b, 0xa5, 0xa5, 0xa6, 0x92, 0x90, 0x65, 0x85,
	0xfd, 0xc1, 0x29, 0xc2, 0x86, 0x04, 0x67, 0x27, 0x45, 0x3e, 0xa1, 0x0a, 0xb7, 0xb1, 0xc2, 0xff,
	0x68, 0xb6, 0xf6, 0x6b, 0x12, 0x4e, 0x24, 0x73, 0x47, 0x71, 0xf6, 0x5c, 0xb9, 0xe8, 0x7c, 0xf4,
	0x6c, 0xcd, 0x11, 0x8f, 0xd1, 0xba, 0xf5, 0xf6, 0xdb, 0xcf, 0x4f, 0x25, 0xc3, 0x9a, 0xeb, 0x7d,
	0x43, 0xd4, 0xf5, 0x0f, 0xa8, 0xa2, 0x0f, 0x00, 0x4e, 0xa5, 0xa1, 0x41, 0x77, 0x87, 0xb4, 0x2c,
	0x06, 0x79, 0x64, 0xee, 0xcd, 0xef, 0x8d, 0xab, 0xe7, 0xf3, 0xcb, 0x0e, 0x55, 0xcb, 0x31, 0xad,
	0xeb, 0xfd, 0x72, 0x34, 0x41, 0x1d, 0x54, 0x37, 0x00, 0x7a, 0x03, 0xcb, 0xdb, 0x3c, 0x0c, 0x29,
	0x51, 0x97, 0x3d, 0x8c, 0xdb, 0x9a, 0x7d, 0xde, 0xba, 0xd1, 0xc7, 0x4e, 0x52, 0xba, 0x3a, 0xa8,
	0xda, 0x00, 0x3d, 0x87, 0x13, 0xdb, 0x3e, 0xbe, 0x6c, 0x6e, 0x1b, 0x6c, 0x00, 0xf4, 0x11, 0xc0,
	0x4a, 0xe1, 0x9a, 0xa2, 0xb5, 0xc1, 0xb5, 0x7f, 0xbe, 0x3f, 0xcc, 0xf5, 0x11, 0xd1, 0x99, 0xd9,
	0x15, 0x6d, 0x76, 0xc9, 0x32, 0xfb, 0xcc, 0x8a, 0x2e, 0x36, 0x09, 0x40, 0x08, 0x27, 0x92, 0x1c,
	0x0f, 0xb3, 0x5b, 0xb8, 0xd3, 0xe6, 0x62, 0x0e, 0x2b, 0x7c, 0x1c, 0x9c, 0xa7, 0xf9, 0xc7, 0x61,
	0x60, 0xdc, 0x4e, 0x70, 0x90, 0x8c, 0xd7, 0xbc, 0xf6, 0xb5, 0xf1, 0x7f, 0xc8, 0x09, 0x0e, 0x7d,
	0x2e, 0x55, 0xfd, 0xe1, 0xd6, 0x83, 0x47, 0xcd, 0xf7, 0xe0, 0x4b, 0xc3, 0x41, 0x6b, 0xbe, 0x52,
	0x42, 0xd6, 0x5d, 0xd7, 0x0b, 0x94, 0x1f, 0x1f, 0x39, 0x84, 0x77, 0xdc, 0x94, 0x09, 0x8b, 0x40,
	0xba, 0x1e, 0x16, 0x01, 0x59, 0xcf, 0xa5, 0xc1, 0x05, 0xc2, 0x3b, 0x03, 0xf5, 0xee, 0x81, 0x17,
	0x5b, 0xa3, 0x74, 0x71, 0x25, 0x8d, 0x8e, 0x69, 0xe4, 0x7a, 0x94, 0xa5, 0x37, 0x74, 0x4a, 0x3f,
	0x36, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0xa4, 0x67, 0x25, 0x09, 0x5d, 0x07, 0x00, 0x00,
}
