// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/showcase/v1alpha1/testing.proto

package genproto // import "github.com/googleapis/gapic-showcase/server/genproto"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"
import _ "google.golang.org/genproto/googleapis/api/annotations"

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

// The topline state of the report.
type ReportSessionResponse_State int32

const (
	ReportSessionResponse_STATE_UNSPECIFIED ReportSessionResponse_State = 0
	// The session is complete, and everything passed.
	ReportSessionResponse_PASSED ReportSessionResponse_State = 1
	// The session had an explicit failure.
	ReportSessionResponse_FAILED ReportSessionResponse_State = 2
	// The session is incomplete.
	// This is a failure response.
	ReportSessionResponse_INCOMPLETE ReportSessionResponse_State = 3
)

var ReportSessionResponse_State_name = map[int32]string{
	0: "STATE_UNSPECIFIED",
	1: "PASSED",
	2: "FAILED",
	3: "INCOMPLETE",
}
var ReportSessionResponse_State_value = map[string]int32{
	"STATE_UNSPECIFIED": 0,
	"PASSED":            1,
	"FAILED":            2,
	"INCOMPLETE":        3,
}

func (x ReportSessionResponse_State) String() string {
	return proto.EnumName(ReportSessionResponse_State_name, int32(x))
}
func (ReportSessionResponse_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{1, 0}
}

// The different potential types of issues.
type ReportSessionResponse_Issue_Type int32

const (
	ReportSessionResponse_Issue_TYPE_UNSPECIFIED ReportSessionResponse_Issue_Type = 0
	// The test was never instrumented.
	ReportSessionResponse_Issue_SKIPPED ReportSessionResponse_Issue_Type = 1
	// The test was instrumented, but Showcase got an unexpected
	// value when the generator tried to confirm success.
	ReportSessionResponse_Issue_INCORRECT_CONFIRMATION ReportSessionResponse_Issue_Type = 2
)

var ReportSessionResponse_Issue_Type_name = map[int32]string{
	0: "TYPE_UNSPECIFIED",
	1: "SKIPPED",
	2: "INCORRECT_CONFIRMATION",
}
var ReportSessionResponse_Issue_Type_value = map[string]int32{
	"TYPE_UNSPECIFIED":       0,
	"SKIPPED":                1,
	"INCORRECT_CONFIRMATION": 2,
}

func (x ReportSessionResponse_Issue_Type) String() string {
	return proto.EnumName(ReportSessionResponse_Issue_Type_name, int32(x))
}
func (ReportSessionResponse_Issue_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{1, 0, 0}
}

// Severity levels.
type ReportSessionResponse_Issue_Severity int32

const (
	ReportSessionResponse_Issue_SEVERITY_UNSPECIFIED ReportSessionResponse_Issue_Severity = 0
	// Errors.
	ReportSessionResponse_Issue_ERROR ReportSessionResponse_Issue_Severity = 1
	// Warnings.
	ReportSessionResponse_Issue_WARNING ReportSessionResponse_Issue_Severity = 2
)

var ReportSessionResponse_Issue_Severity_name = map[int32]string{
	0: "SEVERITY_UNSPECIFIED",
	1: "ERROR",
	2: "WARNING",
}
var ReportSessionResponse_Issue_Severity_value = map[string]int32{
	"SEVERITY_UNSPECIFIED": 0,
	"ERROR":                1,
	"WARNING":              2,
}

func (x ReportSessionResponse_Issue_Severity) String() string {
	return proto.EnumName(ReportSessionResponse_Issue_Severity_name, int32(x))
}
func (ReportSessionResponse_Issue_Severity) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{1, 0, 1}
}

// The specification versions understood by Showcase.
type Session_Version int32

const (
	Session_VERSION_UNSPECIFIED Session_Version = 0
	// The latest v1. Currently, this is v1.0.
	Session_V1_LATEST Session_Version = 1
	// v1.0. (Until the spec is "GA", this will be a moving target.)
	Session_V1_0 Session_Version = 2
)

var Session_Version_name = map[int32]string{
	0: "VERSION_UNSPECIFIED",
	1: "V1_LATEST",
	2: "V1_0",
}
var Session_Version_value = map[string]int32{
	"VERSION_UNSPECIFIED": 0,
	"V1_LATEST":           1,
	"V1_0":                2,
}

func (x Session_Version) String() string {
	return proto.EnumName(Session_Version_name, int32(x))
}
func (Session_Version) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{4, 0}
}

// Whether or not a test is required, recommended, or optional.
type Test_ExpectationLevel int32

const (
	Test_EXPECTATION_LEVEL_UNSPECIFIED Test_ExpectationLevel = 0
	// This test is strictly required.
	Test_REQUIRED Test_ExpectationLevel = 1
	// This test is recommended.
	//
	// If a generator explicitly ignores a recommended test (see `DeleteTest`),
	// then the report may still pass, but with a warning.
	//
	// If a generator skips a recommended test and does not explicitly
	// express that intention, the report will fail.
	Test_RECOMMENDED Test_ExpectationLevel = 2
	// This test is optional.
	//
	// If a generator explicitly ignores an optional test (see `DeleteTest`),
	// then the report may still pass, and no warning will be issued.
	//
	// If a generator skips an optional test and does not explicitly
	// express that intention, the report may still pass, but with a
	// warning.
	Test_OPTIONAL Test_ExpectationLevel = 3
)

var Test_ExpectationLevel_name = map[int32]string{
	0: "EXPECTATION_LEVEL_UNSPECIFIED",
	1: "REQUIRED",
	2: "RECOMMENDED",
	3: "OPTIONAL",
}
var Test_ExpectationLevel_value = map[string]int32{
	"EXPECTATION_LEVEL_UNSPECIFIED": 0,
	"REQUIRED":                      1,
	"RECOMMENDED":                   2,
	"OPTIONAL":                      3,
}

func (x Test_ExpectationLevel) String() string {
	return proto.EnumName(Test_ExpectationLevel_name, int32(x))
}
func (Test_ExpectationLevel) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{5, 0}
}

// Request message for reporting on a session.
type ReportSessionRequest struct {
	// The session to be reported on.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportSessionRequest) Reset()         { *m = ReportSessionRequest{} }
func (m *ReportSessionRequest) String() string { return proto.CompactTextString(m) }
func (*ReportSessionRequest) ProtoMessage()    {}
func (*ReportSessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{0}
}
func (m *ReportSessionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportSessionRequest.Unmarshal(m, b)
}
func (m *ReportSessionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportSessionRequest.Marshal(b, m, deterministic)
}
func (dst *ReportSessionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportSessionRequest.Merge(dst, src)
}
func (m *ReportSessionRequest) XXX_Size() int {
	return xxx_messageInfo_ReportSessionRequest.Size(m)
}
func (m *ReportSessionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportSessionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReportSessionRequest proto.InternalMessageInfo

func (m *ReportSessionRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Response message for reporting on a session.
type ReportSessionResponse struct {
	// The state of the report.
	State ReportSessionResponse_State `protobuf:"varint,1,opt,name=state,proto3,enum=google.showcase.v1alpha2.ReportSessionResponse_State" json:"state,omitempty"`
	// Issues found which constitute errors.
	// A non-zero list here requires a failure state for the report overall.
	Errors []*ReportSessionResponse_Issue `protobuf:"bytes,2,rep,name=errors,proto3" json:"errors,omitempty"`
	// Issues found which do not constitute errors.
	Warnings []*ReportSessionResponse_Issue `protobuf:"bytes,3,rep,name=warnings,proto3" json:"warnings,omitempty"`
	// The overall completion in the report.
	Completion           *ReportSessionResponse_Completion `protobuf:"bytes,4,opt,name=completion,proto3" json:"completion,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                          `json:"-"`
	XXX_unrecognized     []byte                            `json:"-"`
	XXX_sizecache        int32                             `json:"-"`
}

func (m *ReportSessionResponse) Reset()         { *m = ReportSessionResponse{} }
func (m *ReportSessionResponse) String() string { return proto.CompactTextString(m) }
func (*ReportSessionResponse) ProtoMessage()    {}
func (*ReportSessionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{1}
}
func (m *ReportSessionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportSessionResponse.Unmarshal(m, b)
}
func (m *ReportSessionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportSessionResponse.Marshal(b, m, deterministic)
}
func (dst *ReportSessionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportSessionResponse.Merge(dst, src)
}
func (m *ReportSessionResponse) XXX_Size() int {
	return xxx_messageInfo_ReportSessionResponse.Size(m)
}
func (m *ReportSessionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportSessionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ReportSessionResponse proto.InternalMessageInfo

func (m *ReportSessionResponse) GetState() ReportSessionResponse_State {
	if m != nil {
		return m.State
	}
	return ReportSessionResponse_STATE_UNSPECIFIED
}

func (m *ReportSessionResponse) GetErrors() []*ReportSessionResponse_Issue {
	if m != nil {
		return m.Errors
	}
	return nil
}

func (m *ReportSessionResponse) GetWarnings() []*ReportSessionResponse_Issue {
	if m != nil {
		return m.Warnings
	}
	return nil
}

func (m *ReportSessionResponse) GetCompletion() *ReportSessionResponse_Completion {
	if m != nil {
		return m.Completion
	}
	return nil
}

// An issue highlighted in the report.
type ReportSessionResponse_Issue struct {
	// The test that is the root of the issue.
	Test string `protobuf:"bytes,1,opt,name=test,proto3" json:"test,omitempty"`
	// The type of the issue.
	Type ReportSessionResponse_Issue_Type `protobuf:"varint,2,opt,name=type,proto3,enum=google.showcase.v1alpha2.ReportSessionResponse_Issue_Type" json:"type,omitempty"`
	// The severity of the issue.
	Severity ReportSessionResponse_Issue_Severity `protobuf:"varint,3,opt,name=severity,proto3,enum=google.showcase.v1alpha2.ReportSessionResponse_Issue_Severity" json:"severity,omitempty"`
	// A description of the issue.
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportSessionResponse_Issue) Reset()         { *m = ReportSessionResponse_Issue{} }
func (m *ReportSessionResponse_Issue) String() string { return proto.CompactTextString(m) }
func (*ReportSessionResponse_Issue) ProtoMessage()    {}
func (*ReportSessionResponse_Issue) Descriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{1, 0}
}
func (m *ReportSessionResponse_Issue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportSessionResponse_Issue.Unmarshal(m, b)
}
func (m *ReportSessionResponse_Issue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportSessionResponse_Issue.Marshal(b, m, deterministic)
}
func (dst *ReportSessionResponse_Issue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportSessionResponse_Issue.Merge(dst, src)
}
func (m *ReportSessionResponse_Issue) XXX_Size() int {
	return xxx_messageInfo_ReportSessionResponse_Issue.Size(m)
}
func (m *ReportSessionResponse_Issue) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportSessionResponse_Issue.DiscardUnknown(m)
}

var xxx_messageInfo_ReportSessionResponse_Issue proto.InternalMessageInfo

func (m *ReportSessionResponse_Issue) GetTest() string {
	if m != nil {
		return m.Test
	}
	return ""
}

func (m *ReportSessionResponse_Issue) GetType() ReportSessionResponse_Issue_Type {
	if m != nil {
		return m.Type
	}
	return ReportSessionResponse_Issue_TYPE_UNSPECIFIED
}

func (m *ReportSessionResponse_Issue) GetSeverity() ReportSessionResponse_Issue_Severity {
	if m != nil {
		return m.Severity
	}
	return ReportSessionResponse_Issue_SEVERITY_UNSPECIFIED
}

func (m *ReportSessionResponse_Issue) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

// A message describing the completion level for tests.
type ReportSessionResponse_Completion struct {
	// The ratio of required and recommended tests completed.
	//
	// The behavior of the denominator for recommended and optional tests
	// matches the description in their individual ratios below.
	Total float32 `protobuf:"fixed32,1,opt,name=total,proto3" json:"total,omitempty"`
	// The ratio of required tests completed.
	Required float32 `protobuf:"fixed32,2,opt,name=required,proto3" json:"required,omitempty"`
	// The ratio of recommended tests completed.
	// Deleting a recommended test *does not* remove it from the denominator.
	Recommended float32 `protobuf:"fixed32,3,opt,name=recommended,proto3" json:"recommended,omitempty"`
	// The ratio of optional tests completed.
	// Deleting an optional test *does* remove it from the denominator.
	Optional             float32  `protobuf:"fixed32,4,opt,name=optional,proto3" json:"optional,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportSessionResponse_Completion) Reset()         { *m = ReportSessionResponse_Completion{} }
func (m *ReportSessionResponse_Completion) String() string { return proto.CompactTextString(m) }
func (*ReportSessionResponse_Completion) ProtoMessage()    {}
func (*ReportSessionResponse_Completion) Descriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{1, 1}
}
func (m *ReportSessionResponse_Completion) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportSessionResponse_Completion.Unmarshal(m, b)
}
func (m *ReportSessionResponse_Completion) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportSessionResponse_Completion.Marshal(b, m, deterministic)
}
func (dst *ReportSessionResponse_Completion) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportSessionResponse_Completion.Merge(dst, src)
}
func (m *ReportSessionResponse_Completion) XXX_Size() int {
	return xxx_messageInfo_ReportSessionResponse_Completion.Size(m)
}
func (m *ReportSessionResponse_Completion) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportSessionResponse_Completion.DiscardUnknown(m)
}

var xxx_messageInfo_ReportSessionResponse_Completion proto.InternalMessageInfo

func (m *ReportSessionResponse_Completion) GetTotal() float32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *ReportSessionResponse_Completion) GetRequired() float32 {
	if m != nil {
		return m.Required
	}
	return 0
}

func (m *ReportSessionResponse_Completion) GetRecommended() float32 {
	if m != nil {
		return m.Recommended
	}
	return 0
}

func (m *ReportSessionResponse_Completion) GetOptional() float32 {
	if m != nil {
		return m.Optional
	}
	return 0
}

// Request message for deleting a test.
type DeleteTestRequest struct {
	// The test to be deleted.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteTestRequest) Reset()         { *m = DeleteTestRequest{} }
func (m *DeleteTestRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteTestRequest) ProtoMessage()    {}
func (*DeleteTestRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{2}
}
func (m *DeleteTestRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteTestRequest.Unmarshal(m, b)
}
func (m *DeleteTestRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteTestRequest.Marshal(b, m, deterministic)
}
func (dst *DeleteTestRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteTestRequest.Merge(dst, src)
}
func (m *DeleteTestRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteTestRequest.Size(m)
}
func (m *DeleteTestRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteTestRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteTestRequest proto.InternalMessageInfo

func (m *DeleteTestRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Register an answer from a test.
type RegisterTestRequest struct {
	// The test to have an answer registered to it.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The answer from the test.
	Answers              []string `protobuf:"bytes,2,rep,name=answers,proto3" json:"answers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterTestRequest) Reset()         { *m = RegisterTestRequest{} }
func (m *RegisterTestRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterTestRequest) ProtoMessage()    {}
func (*RegisterTestRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{3}
}
func (m *RegisterTestRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterTestRequest.Unmarshal(m, b)
}
func (m *RegisterTestRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterTestRequest.Marshal(b, m, deterministic)
}
func (dst *RegisterTestRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterTestRequest.Merge(dst, src)
}
func (m *RegisterTestRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterTestRequest.Size(m)
}
func (m *RegisterTestRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterTestRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterTestRequest proto.InternalMessageInfo

func (m *RegisterTestRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RegisterTestRequest) GetAnswers() []string {
	if m != nil {
		return m.Answers
	}
	return nil
}

// A session is a suite of tests, generally being made in the context
// of testing code generation.
//
// A session defines tests it may expect, based on which version of the
// code generation spec is in use.
type Session struct {
	// The name of the session. The ID must conform to ^[a-z]+$
	// If this is not provided, Showcase chooses one at random.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Required. The version this session is using.
	Version              Session_Version `protobuf:"varint,2,opt,name=version,proto3,enum=google.showcase.v1alpha2.Session_Version" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Session) Reset()         { *m = Session{} }
func (m *Session) String() string { return proto.CompactTextString(m) }
func (*Session) ProtoMessage()    {}
func (*Session) Descriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{4}
}
func (m *Session) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Session.Unmarshal(m, b)
}
func (m *Session) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Session.Marshal(b, m, deterministic)
}
func (dst *Session) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Session.Merge(dst, src)
}
func (m *Session) XXX_Size() int {
	return xxx_messageInfo_Session.Size(m)
}
func (m *Session) XXX_DiscardUnknown() {
	xxx_messageInfo_Session.DiscardUnknown(m)
}

var xxx_messageInfo_Session proto.InternalMessageInfo

func (m *Session) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Session) GetVersion() Session_Version {
	if m != nil {
		return m.Version
	}
	return Session_VERSION_UNSPECIFIED
}

// A test is a single test expected by Showcase.
// The exact way that each test is exercised may vary from test to test.
type Test struct {
	// The name of the test.
	// The tests/* portion of the names are hard-coded, and do not change
	// from session to session.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The expectation level for this test.
	ExpectationLevel Test_ExpectationLevel `protobuf:"varint,2,opt,name=expectation_level,json=expectationLevel,proto3,enum=google.showcase.v1alpha2.Test_ExpectationLevel" json:"expectation_level,omitempty"`
	// A description of the test.
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Test) Reset()         { *m = Test{} }
func (m *Test) String() string { return proto.CompactTextString(m) }
func (*Test) ProtoMessage()    {}
func (*Test) Descriptor() ([]byte, []int) {
	return fileDescriptor_testing_b204758f5e20f31d, []int{5}
}
func (m *Test) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Test.Unmarshal(m, b)
}
func (m *Test) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Test.Marshal(b, m, deterministic)
}
func (dst *Test) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Test.Merge(dst, src)
}
func (m *Test) XXX_Size() int {
	return xxx_messageInfo_Test.Size(m)
}
func (m *Test) XXX_DiscardUnknown() {
	xxx_messageInfo_Test.DiscardUnknown(m)
}

var xxx_messageInfo_Test proto.InternalMessageInfo

func (m *Test) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Test) GetExpectationLevel() Test_ExpectationLevel {
	if m != nil {
		return m.ExpectationLevel
	}
	return Test_EXPECTATION_LEVEL_UNSPECIFIED
}

func (m *Test) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterType((*ReportSessionRequest)(nil), "google.showcase.v1alpha2.ReportSessionRequest")
	proto.RegisterType((*ReportSessionResponse)(nil), "google.showcase.v1alpha2.ReportSessionResponse")
	proto.RegisterType((*ReportSessionResponse_Issue)(nil), "google.showcase.v1alpha2.ReportSessionResponse.Issue")
	proto.RegisterType((*ReportSessionResponse_Completion)(nil), "google.showcase.v1alpha2.ReportSessionResponse.Completion")
	proto.RegisterType((*DeleteTestRequest)(nil), "google.showcase.v1alpha2.DeleteTestRequest")
	proto.RegisterType((*RegisterTestRequest)(nil), "google.showcase.v1alpha2.RegisterTestRequest")
	proto.RegisterType((*Session)(nil), "google.showcase.v1alpha2.Session")
	proto.RegisterType((*Test)(nil), "google.showcase.v1alpha2.Test")
	proto.RegisterEnum("google.showcase.v1alpha2.ReportSessionResponse_State", ReportSessionResponse_State_name, ReportSessionResponse_State_value)
	proto.RegisterEnum("google.showcase.v1alpha2.ReportSessionResponse_Issue_Type", ReportSessionResponse_Issue_Type_name, ReportSessionResponse_Issue_Type_value)
	proto.RegisterEnum("google.showcase.v1alpha2.ReportSessionResponse_Issue_Severity", ReportSessionResponse_Issue_Severity_name, ReportSessionResponse_Issue_Severity_value)
	proto.RegisterEnum("google.showcase.v1alpha2.Session_Version", Session_Version_name, Session_Version_value)
	proto.RegisterEnum("google.showcase.v1alpha2.Test_ExpectationLevel", Test_ExpectationLevel_name, Test_ExpectationLevel_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TestingClient is the client API for Testing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TestingClient interface {
	// Report on the status of a session.
	// This generates a report detailing which tests have been completed,
	// and an overall rollup.
	ReportSession(ctx context.Context, in *ReportSessionRequest, opts ...grpc.CallOption) (*ReportSessionResponse, error)
	// Explicitly decline to implement a test.
	//
	// This removes the test from subsequent `ListTests` calls, and
	// attempting to do the test will error.
	//
	// This method will error if attempting to delete a required test.
	DeleteTest(ctx context.Context, in *DeleteTestRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Register a response to a test.
	//
	// In cases where a test involves registering a final answer at the
	// end of the test, this method provides the means to do so.
	RegisterTest(ctx context.Context, in *RegisterTestRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type testingClient struct {
	cc *grpc.ClientConn
}

func NewTestingClient(cc *grpc.ClientConn) TestingClient {
	return &testingClient{cc}
}

func (c *testingClient) ReportSession(ctx context.Context, in *ReportSessionRequest, opts ...grpc.CallOption) (*ReportSessionResponse, error) {
	out := new(ReportSessionResponse)
	err := c.cc.Invoke(ctx, "/google.showcase.v1alpha2.Testing/ReportSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testingClient) DeleteTest(ctx context.Context, in *DeleteTestRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/google.showcase.v1alpha2.Testing/DeleteTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testingClient) RegisterTest(ctx context.Context, in *RegisterTestRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/google.showcase.v1alpha2.Testing/RegisterTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TestingServer is the server API for Testing service.
type TestingServer interface {
	// Report on the status of a session.
	// This generates a report detailing which tests have been completed,
	// and an overall rollup.
	ReportSession(context.Context, *ReportSessionRequest) (*ReportSessionResponse, error)
	// Explicitly decline to implement a test.
	//
	// This removes the test from subsequent `ListTests` calls, and
	// attempting to do the test will error.
	//
	// This method will error if attempting to delete a required test.
	DeleteTest(context.Context, *DeleteTestRequest) (*empty.Empty, error)
	// Register a response to a test.
	//
	// In cases where a test involves registering a final answer at the
	// end of the test, this method provides the means to do so.
	RegisterTest(context.Context, *RegisterTestRequest) (*empty.Empty, error)
}

func RegisterTestingServer(s *grpc.Server, srv TestingServer) {
	s.RegisterService(&_Testing_serviceDesc, srv)
}

func _Testing_ReportSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestingServer).ReportSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.showcase.v1alpha2.Testing/ReportSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestingServer).ReportSession(ctx, req.(*ReportSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Testing_DeleteTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestingServer).DeleteTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.showcase.v1alpha2.Testing/DeleteTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestingServer).DeleteTest(ctx, req.(*DeleteTestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Testing_RegisterTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterTestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestingServer).RegisterTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.showcase.v1alpha2.Testing/RegisterTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestingServer).RegisterTest(ctx, req.(*RegisterTestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Testing_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.showcase.v1alpha2.Testing",
	HandlerType: (*TestingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReportSession",
			Handler:    _Testing_ReportSession_Handler,
		},
		{
			MethodName: "DeleteTest",
			Handler:    _Testing_DeleteTest_Handler,
		},
		{
			MethodName: "RegisterTest",
			Handler:    _Testing_RegisterTest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/showcase/v1alpha1/testing.proto",
}

func init() {
	proto.RegisterFile("google/showcase/v1alpha1/testing.proto", fileDescriptor_testing_b204758f5e20f31d)
}

var fileDescriptor_testing_b204758f5e20f31d = []byte{
	// 985 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0x5f, 0x6f, 0xdb, 0x54,
	0x14, 0x27, 0x4e, 0xda, 0xa4, 0xa7, 0xdd, 0xe6, 0xde, 0x75, 0x5b, 0x94, 0x0d, 0x54, 0x3c, 0x40,
	0xa5, 0x65, 0xf6, 0x5a, 0xca, 0x10, 0xe1, 0x8f, 0xe4, 0x25, 0xb7, 0xc8, 0x5a, 0xfe, 0xed, 0xda,
	0x0b, 0xac, 0x42, 0x8a, 0xdc, 0xf4, 0x90, 0x5a, 0x72, 0x6c, 0xcf, 0xd7, 0x49, 0xa9, 0x60, 0x3c,
	0xf0, 0xca, 0x23, 0x0f, 0x08, 0x89, 0x0f, 0xc1, 0xd7, 0x80, 0x57, 0xbe, 0x02, 0x4f, 0x7b, 0xe0,
	0x03, 0xec, 0x09, 0xf9, 0xda, 0x4e, 0xda, 0x94, 0xa0, 0x76, 0x3c, 0x25, 0xf7, 0xdc, 0xf3, 0xfb,
	0x9d, 0x73, 0x7e, 0xf7, 0x9c, 0x7b, 0x0d, 0xef, 0x0c, 0x7c, 0x7f, 0xe0, 0xa2, 0xc6, 0x8f, 0xfc,
	0xe3, 0xbe, 0xcd, 0x51, 0x1b, 0x6f, 0xdb, 0x6e, 0x70, 0x64, 0x6f, 0x6b, 0x11, 0xf2, 0xc8, 0xf1,
	0x06, 0x6a, 0x10, 0xfa, 0x91, 0x4f, 0xca, 0x89, 0x9f, 0x9a, 0xf9, 0xa9, 0xa9, 0xdf, 0x4e, 0xe5,
	0x4e, 0xca, 0x60, 0x07, 0x8e, 0x66, 0x7b, 0x9e, 0x1f, 0xd9, 0x91, 0xe3, 0x7b, 0x3c, 0xc1, 0x55,
	0x6e, 0xa7, 0xbb, 0x62, 0x75, 0x30, 0xfa, 0x5a, 0xc3, 0x61, 0x10, 0x9d, 0x24, 0x9b, 0xca, 0x03,
	0x58, 0x63, 0x18, 0xf8, 0x61, 0x64, 0x22, 0xe7, 0x8e, 0xef, 0x31, 0x7c, 0x36, 0x42, 0x1e, 0x91,
	0x37, 0xa0, 0xe0, 0xd9, 0x43, 0x2c, 0xe7, 0xd6, 0x73, 0x1b, 0x4b, 0x0f, 0xe1, 0xa5, 0x5e, 0xcc,
	0x1c, 0x84, 0x5d, 0xf9, 0xa5, 0x08, 0x37, 0x66, 0x80, 0x3c, 0xf0, 0x3d, 0x8e, 0xe4, 0x11, 0x2c,
	0xf0, 0xc8, 0x8e, 0x12, 0xe8, 0xd5, 0x9d, 0x0f, 0xd4, 0x79, 0x69, 0xab, 0xff, 0x8a, 0x57, 0xcd,
	0x18, 0xcc, 0x12, 0x0e, 0xd2, 0x84, 0x45, 0x0c, 0x43, 0x3f, 0xe4, 0x65, 0x69, 0x3d, 0xbf, 0xb1,
	0x7c, 0x79, 0x36, 0x83, 0xf3, 0x11, 0xb2, 0x94, 0x84, 0x3c, 0x86, 0xd2, 0xb1, 0x1d, 0x7a, 0x8e,
	0x37, 0xe0, 0xe5, 0xfc, 0xff, 0x21, 0x9c, 0xd0, 0x90, 0x7d, 0x80, 0xbe, 0x3f, 0x0c, 0x5c, 0x8c,
	0x25, 0x2f, 0x17, 0xd6, 0x73, 0x1b, 0xcb, 0x3b, 0xd5, 0xcb, 0x92, 0xd6, 0x26, 0x0c, 0xec, 0x14,
	0x5b, 0xe5, 0x6f, 0x09, 0x16, 0x44, 0x3c, 0x42, 0xa0, 0x10, 0x37, 0x43, 0x72, 0x1c, 0x4c, 0xfc,
	0x27, 0x2d, 0x28, 0x44, 0x27, 0x01, 0x96, 0x25, 0xa1, 0x73, 0xf5, 0x95, 0x0a, 0x51, 0xad, 0x93,
	0x00, 0x99, 0xe0, 0x21, 0xfb, 0x50, 0xe2, 0x38, 0xc6, 0xd0, 0x89, 0x4e, 0xca, 0x79, 0xc1, 0xf9,
	0xd9, 0xab, 0x71, 0x9a, 0x29, 0x0b, 0x9b, 0xf0, 0x91, 0x75, 0x58, 0x3e, 0x44, 0xde, 0x0f, 0x9d,
	0x60, 0x22, 0xd3, 0x12, 0x3b, 0x6d, 0x52, 0x28, 0x14, 0xe2, 0x5c, 0xc8, 0x1a, 0xc8, 0xd6, 0xd3,
	0x0e, 0xed, 0x3d, 0x69, 0x99, 0x1d, 0x5a, 0x33, 0xf6, 0x0c, 0x5a, 0x97, 0x5f, 0x23, 0xcb, 0x50,
	0x34, 0x1f, 0x19, 0x9d, 0x0e, 0xad, 0xcb, 0x39, 0x52, 0x81, 0x9b, 0x46, 0xab, 0xd6, 0x66, 0x8c,
	0xd6, 0xac, 0x5e, 0xad, 0xdd, 0xda, 0x33, 0x58, 0x53, 0xb7, 0x8c, 0x76, 0x4b, 0x96, 0x94, 0x4f,
	0xa0, 0x94, 0x85, 0x27, 0x65, 0x58, 0x33, 0x69, 0x97, 0x32, 0xc3, 0x7a, 0x3a, 0x43, 0xb7, 0x04,
	0x0b, 0x94, 0xb1, 0x36, 0x93, 0x73, 0x31, 0xf3, 0x17, 0x3a, 0x6b, 0x19, 0xad, 0xcf, 0x65, 0xa9,
	0xf2, 0x1d, 0xc0, 0xf4, 0x28, 0xc8, 0x1a, 0x2c, 0x44, 0x7e, 0x64, 0xbb, 0x42, 0x75, 0x89, 0x25,
	0x0b, 0x52, 0x81, 0x52, 0x88, 0xcf, 0x46, 0x4e, 0x88, 0x87, 0x42, 0x7a, 0x89, 0x4d, 0xd6, 0x71,
	0x99, 0x21, 0xf6, 0xfd, 0xe1, 0x10, 0xbd, 0x43, 0x3c, 0x14, 0x2a, 0x4a, 0xec, 0xb4, 0x29, 0x46,
	0xfb, 0xa2, 0x60, 0xdb, 0x15, 0x2a, 0x48, 0x6c, 0xb2, 0x56, 0xf6, 0x60, 0x41, 0x34, 0x3f, 0xb9,
	0x01, 0xab, 0xa6, 0xa5, 0x5b, 0xb3, 0x22, 0x00, 0x2c, 0x76, 0x74, 0xd3, 0x14, 0x1a, 0x00, 0x2c,
	0xee, 0xe9, 0x46, 0x83, 0xd6, 0x65, 0x89, 0x5c, 0x05, 0x88, 0xf5, 0x68, 0x76, 0x1a, 0xd4, 0xa2,
	0x72, 0x5e, 0xb9, 0x0f, 0xab, 0x75, 0x74, 0x31, 0x42, 0x0b, 0x79, 0x94, 0x0d, 0xf4, 0xed, 0x33,
	0x03, 0x5d, 0x7c, 0xa9, 0x17, 0xc4, 0x6e, 0x32, 0xcd, 0x0d, 0xb8, 0xce, 0x70, 0xe0, 0xf0, 0x08,
	0xc3, 0x8b, 0x62, 0x48, 0x19, 0x8a, 0xb6, 0xc7, 0x8f, 0x31, 0x9d, 0xcd, 0x25, 0x96, 0x2d, 0x95,
	0xdf, 0x72, 0x90, 0xdd, 0x16, 0xe4, 0xee, 0x19, 0x8a, 0x6b, 0x2f, 0xf4, 0x15, 0x00, 0x9e, 0xec,
	0x71, 0x6d, 0x33, 0xa5, 0xaa, 0x41, 0x71, 0x8c, 0x61, 0x6c, 0x4b, 0x9b, 0xf9, 0xdd, 0xf9, 0x8d,
	0x97, 0x12, 0xab, 0xdd, 0x04, 0xc0, 0x32, 0xa4, 0xf2, 0x31, 0x14, 0x53, 0x1b, 0xb9, 0x05, 0xd7,
	0xbb, 0x94, 0x99, 0x46, 0xbb, 0x35, 0xa3, 0xe0, 0x15, 0x58, 0xea, 0x6e, 0xf7, 0x1a, 0xba, 0x45,
	0x4d, 0x4b, 0xce, 0x91, 0x12, 0x14, 0xba, 0xdb, 0xbd, 0xfb, 0xb2, 0xa4, 0xfc, 0x2c, 0x81, 0xa8,
	0x8d, 0x6c, 0x9d, 0xc9, 0xf7, 0xd6, 0x0b, 0x7d, 0x0d, 0xc8, 0x34, 0x5f, 0x71, 0x23, 0x4f, 0xf3,
	0xfe, 0x0a, 0x56, 0xf1, 0x9b, 0x00, 0xfb, 0xc9, 0x7d, 0xdb, 0x73, 0x71, 0x8c, 0x6e, 0x5a, 0x81,
	0x36, 0xbf, 0x82, 0x38, 0x8e, 0x4a, 0xa7, 0xb8, 0x46, 0x0c, 0x63, 0x32, 0xce, 0x58, 0x66, 0x67,
	0x26, 0x7f, 0x7e, 0x66, 0x0e, 0x40, 0x9e, 0xe5, 0x21, 0x6f, 0xc2, 0xeb, 0xf4, 0xcb, 0x0e, 0xad,
	0x59, 0x62, 0x22, 0x7a, 0x0d, 0xda, 0xa5, 0x8d, 0x19, 0x15, 0x56, 0xa0, 0xc4, 0xe8, 0xe3, 0x27,
	0x06, 0x13, 0x9d, 0x74, 0x0d, 0x96, 0x19, 0xad, 0xb5, 0x9b, 0x4d, 0xda, 0xaa, 0x8b, 0x76, 0x5a,
	0x81, 0x52, 0xbb, 0x13, 0x83, 0xf5, 0x86, 0x9c, 0xdf, 0xf9, 0x3d, 0x0f, 0x45, 0x2b, 0x79, 0x87,
	0xc8, 0xaf, 0x39, 0xb8, 0x72, 0x66, 0xf0, 0x89, 0x7a, 0xe1, 0x1b, 0x42, 0x74, 0x54, 0x45, 0xbb,
	0xe4, 0x8d, 0xa2, 0x6c, 0xfe, 0xf0, 0xe7, 0x5f, 0x3f, 0x49, 0x6f, 0x29, 0x4a, 0xf6, 0x2a, 0xee,
	0x68, 0xdf, 0xc6, 0xda, 0x7f, 0x3a, 0x3d, 0x94, 0xe7, 0xd5, 0x50, 0x40, 0xc9, 0xf7, 0x00, 0xd3,
	0xbe, 0x27, 0x5b, 0xf3, 0x43, 0x9d, 0x9b, 0x8e, 0xca, 0xcd, 0xcc, 0x39, 0x7b, 0x24, 0x55, 0x1a,
	0x3f, 0x92, 0xca, 0x96, 0x08, 0xff, 0xf6, 0xe6, 0xdd, 0xf9, 0xe1, 0xb3, 0x9e, 0x78, 0x4e, 0x7e,
	0xcc, 0xc1, 0xca, 0xe9, 0x31, 0x22, 0xf7, 0xfe, 0xab, 0xda, 0x73, 0xe3, 0x36, 0x37, 0x89, 0x5d,
	0x91, 0x84, 0xaa, 0xbc, 0x77, 0x81, 0x24, 0xaa, 0x61, 0x4a, 0x5c, 0x59, 0xfd, 0x43, 0xbf, 0xea,
	0xfa, 0x7d, 0xdb, 0x3d, 0xf2, 0x79, 0x54, 0xfd, 0x70, 0xf7, 0xc1, 0x47, 0x0f, 0x2d, 0xb8, 0xd3,
	0xf7, 0x87, 0xf3, 0x92, 0xda, 0xde, 0xdf, 0x1d, 0x38, 0xd1, 0xd1, 0xe8, 0x40, 0xed, 0xfb, 0x43,
	0x2d, 0x71, 0xb2, 0x03, 0x87, 0x6b, 0x03, 0x3b, 0x70, 0xfa, 0xf7, 0x26, 0x9f, 0x27, 0x1c, 0xc3,
	0x31, 0x86, 0xda, 0x00, 0xbd, 0x24, 0xcf, 0x45, 0xf1, 0xf3, 0xfe, 0x3f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x92, 0x8f, 0xa4, 0xf2, 0xc8, 0x08, 0x00, 0x00,
}
