// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resttools

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestParseWellKnownType(t *testing.T) {
	for _, tst := range []struct {
		name  string
		msg   protoreflect.Message
		field protoreflect.Name
		want  proto.Message
	}{
		{
			"google.protobuf.FieldMask",
			(&genprotopb.UpdateUserRequest{}).ProtoReflect(),
			"update_mask",
			&fieldmaskpb.FieldMask{Paths: []string{"foo", "bar", "baz"}},
		},
		{
			"google.protobuf.Timestamp",
			(&genprotopb.User{}).ProtoReflect(),
			"create_time",
			timestamppb.Now(),
		},
		{
			"google.protobuf.Duration",
			(&genprotopb.Sequence_Response{}).ProtoReflect(),
			"delay",
			durationpb.New(5 * time.Second),
		},
	} {
		data, _ := protojson.Marshal(tst.want)
		value := string(data)
		fd := tst.msg.Descriptor().Fields().ByName(tst.field)

		gotp, err := parseWellKnownType(tst.msg, fd, value)
		if err != nil {
			t.Fatal(err)
		}
		if gotp == nil {
			t.Fatal("expected non-nil value from parsing")
		}
		got := gotp.Message().Interface()
		if diff := cmp.Diff(got, tst.want, cmp.Comparer(proto.Equal)); diff != "" {
			t.Fatalf("%s: got(-),want(+):\n%s", "FieldMask", diff)
		}
	}
}

func TestPopulateOneFieldError(t *testing.T) {
	for idx, testCase := range []struct {
		field string
		value string
	}{
		// field path errors

		{".fString", "hi"},
		{"fChild..fString", "hi"},
		{"fChild.", "hi"},
		{"fChild.x", "hi"},
		{"fChild.x.fString", "hi"},
		{"fChild.fString.fChild", "hi"},
		{"fChild.fChild", "hi"},

		// parsing errors

		{"fInt32", "hello"},
		{"fSint32", "hello"},
		{"fSfixed32", "hello"},
		{"fInt32", "2147483648"},    // max int32 + 1
		{"fSint32", "2147483648"},   // max int32 + 1
		{"fSfixed32", "2147483648"}, // max int32 + 1
		{"fInt32", "1.1"},
		{"fSint32", "1.1"},
		{"fSfixed32", "1.1"},

		{"fUint32", "hello"},
		{"fFixed32", "hello"},
		{"fUint32", "4294967296"},  // max uint32 + 1
		{"fFixed32", "4294967296"}, // max uint32 + 1
		{"fUint32", "1.1"},
		{"fFixed32", "1.1"},
		{"fUint32", "-1"},
		{"fFixed32", "-1"},

		{"fInt64", "hello"},
		{"fSint64", "hello"},
		{"fSfixed64", "hello"},
		{"fInt64", "9223372036854775808"},    // max int64 + 1
		{"fSint64", "9223372036854775808"},   // max int64 + 1
		{"fSfixed64", "9223372036854775808"}, // max int64 + 1
		{"fInt64", "1.1"},
		{"fSint64", "1.1"},
		{"fSfixed64", "1.1"},

		{"fUint64", "hello"},
		{"fFixed64", "hello"},
		{"fUint64", "18446744073709551616"},  // max uint64 + 1
		{"fFixed64", "18446744073709551616"}, // max uint64 + 1
		{"fUint64", "1.1"},
		{"fFixed64", "1.1"},
		{"fUint64", "-1"},
		{"fFixed64", "-1"},

		{"fFloat", "hello"},
		{"fDouble", "hello"},
		{"fFloat", "1e+39"},   // exponent too large
		{"fDouble", "1e+309"}, // exponent too large

		{"fBool", "hello"},
		{"fBool", "13"},

		// wrong casing
		{"f_string", "alphabet"},
		{"f_int32", "1"},
		{"f_sint32", "2"},
		{"f_sfixed32", "3"},
		{"f_uint32", "4"},
		{"f_fixed32", "5"},
		{"f_int64", "6"},
		{"f_sint64", "7"},
		{"f_sfixed64", "8"},
		{"f_uint64", "9"},
		{"f_fixed64", "10"},
		{"f_float", "11"},
		{"f_double", "12"},
		{"f_bool", "true"},
		{"f_bytes", "greetings"},
	} {
		complianceData := &genprotopb.ComplianceData{}
		err := PopulateOneField(complianceData, testCase.field, []string{testCase.value})
		if err == nil {
			t.Errorf("test case %d: did not get expected error for %q: %q", idx, testCase.field, testCase.value)
		}
	}
}

func TestPopulateSingularFields(t *testing.T) {

	for idx, testCase := range []struct {
		label           string
		fields          map[string]string
		expectError     bool
		expectProtoText string
	}{
		{
			label: "scalar datatypes",
			fields: map[string]string{
				"fString": "alphabet",

				"fInt32":    "2147483647", // max int32
				"fSint32":   "2147483647", // max int32
				"fSfixed32": "2147483647", // max int32

				"fUint32":  "4294967295", // max uint32
				"fFixed32": "4294967295", // max uint32

				"fInt64":    "9223372036854775807", // max int64
				"fSint64":   "9223372036854775807", // max int64
				"fSfixed64": "9223372036854775807", // max int64

				"fUint64":  "18446744073709551615", // max uint64
				"fFixed64": "18446744073709551615", // max uint64

				"fFloat":  "3.40282346638528859811704183484516925440e+38",   // max float32 (https://golang.org/pkg/math/#pkg-constants)
				"fDouble": "1.797693134862315708145274237317043567981e+308", // max float64

				"fBool": "true",

				"fBytes": "greetings",
			},
			expectProtoText: `f_string:"alphabet" f_int32:2147483647 f_sint32:2147483647 f_sfixed32:2147483647 f_uint32:4294967295 f_fixed32:4294967295 f_int64:9223372036854775807 f_sint64:9223372036854775807 f_sfixed64:9223372036854775807 f_uint64:18446744073709551615 f_fixed64:18446744073709551615 f_double:1.7976931348623157e+308 f_float:3.4028235e+38 f_bool:true f_bytes:"greetings"`,
		},
		{
			label: "nested messages",
			fields: map[string]string{
				"fString":               "alphabet",
				"fChild.fChild.fString": "lexicon",
				"fInt32":                "5",
				"fChild.fChild.fDouble": "53.47",
				"fChild.fChild.fBool":   "true",
				"fChild.fBool":          "true",
			},
			expectProtoText: `f_child:{f_child:{f_string:"lexicon" f_bool:true f_double:53.47} f_bool:true} f_string:"alphabet" f_int32:5`,
		},
		{
			label: "presence/zero values",
			fields: map[string]string{
				"fString": "",

				"fInt32":    "0",
				"fSint32":   "0",
				"fSfixed32": "0",

				"fUint32":  "0",
				"fFixed32": "0",

				"fInt64":    "0",
				"fSint64":   "0",
				"fSfixed64": "0",

				"fUint64":  "0",
				"fFixed64": "0",

				"fFloat":  "0",
				"fDouble": "0",

				"fBool": "false",

				"fBytes": "",

				"pString": "",
				"pInt32":  "0",
				"pDouble": "0",
				"pBool":   "false",
			},
			expectProtoText: `p_string:""  p_int32:0  p_double:0  p_bool:false`,
		},
		{
			label: "presence/zero values in nested messages",
			fields: map[string]string{
				"fChild.fString": "",
				"fChild.fFloat":  "0",
				"fChild.fDouble": "0",
				"fChild.fBool":   "false",

				"fChild.pString": "",
				"fChild.pFloat":  "0",
				"fChild.pDouble": "0",
				"fChild.pBool":   "false",

				"pChild.fString": "",
				"pChild.fFloat":  "0",
				"pChild.fDouble": "0",
				"pChild.fBool":   "false",

				"pChild.pString": "",
				"pChild.pFloat":  "0",
				"pChild.pDouble": "0",
				"pChild.pBool":   "false",

				"fChild.fChild.fString": "",
				"fChild.fChild.fDouble": "0",
				"fChild.fChild.fBool":   "false",

				"fChild.pChild.fString": "",
				"fChild.pChild.fDouble": "0",
				"fChild.pChild.fBool":   "false",

				"pChild.fChild.fString": "",
				"pChild.fChild.fDouble": "0",
				"pChild.fChild.fBool":   "false",

				"pChild.pChild.fString": "",
				"pChild.pChild.fDouble": "0",
				"pChild.pChild.fBool":   "false",
			},
			expectProtoText: `f_child:{f_child:{}  p_string:""  p_float:0  p_double:0  p_bool:false  p_child:{}}  p_child:{f_child:{}  p_string:""  p_float:0  p_double:0  p_bool:false  p_child:{}}`,
		},
		{
			label: "default value only in nested message",
			fields: map[string]string{
				"fChild.fString": "",
			},
			expectProtoText: `f_child:{}`,
		},
		{
			label: "default value only in optional nested message",
			fields: map[string]string{
				"pChild.fString": "",
			},
			expectProtoText: `p_child:{}`,
		},
	} {
		complianceData := &genprotopb.ComplianceData{}
		err := PopulateSingularFields(complianceData, testCase.fields)
		if got, want := (err != nil), testCase.expectError; got != want {
			t.Errorf("test case %d[%q] error: got %v, want %v", idx, testCase.label, err, want)
			continue
		}
		if testCase.expectError {
			continue
		}

		var expectProto genprotopb.ComplianceData
		err = prototext.Unmarshal([]byte(testCase.expectProtoText), &expectProto)
		if err != nil {
			t.Errorf("test case %d[%q] unexpected error unmarshaling expected proto: %s", idx, testCase.label, err)
			continue
		}

		if got, want := complianceData, &expectProto; !reflect.DeepEqual(got, want) {
			gotText, err := prototext.Marshal(got)
			if err != nil {
				gotText = []byte("<error marshalling in test>")
			}
			t.Errorf("test case %d[%q] proto:\n    got: %s\n   want: %s", idx, testCase.label, gotText, testCase.expectProtoText)
		}

	}
}

func TestPopulateFields(t *testing.T) {
	for idx, testCase := range []struct {
		label           string
		fields          map[string][]string
		expectError     bool
		expectProtoText string
	}{
		{
			label: "non-repeated fields",
			fields: map[string][]string{
				"fString": []string{"alphabet"},
			},
			expectProtoText: `f_string:"alphabet"`,
		},
		{
			label: "repeated fields",
			fields: map[string][]string{
				"fString": []string{"alphabet", "lexicon"},
			},
			// TODO: Make Populate*Field*() work with repeated fields.
			expectError: true,
		},
	} {
		complianceData := &genprotopb.ComplianceData{}
		err := PopulateFields(complianceData, testCase.fields)
		if got, want := (err != nil), testCase.expectError; got != want {
			t.Errorf("test case %d[%q] error: got %v, want %v", idx, testCase.label, err, want)
			continue
		}
		if testCase.expectError {
			continue
		}

		var expectProto genprotopb.ComplianceData
		err = prototext.Unmarshal([]byte(testCase.expectProtoText), &expectProto)
		if err != nil {
			t.Errorf("test case %d[%q] unexpected error unmarshaling expected proto: %s", idx, testCase.label, err)
			continue
		}

		if got, want := complianceData, &expectProto; !reflect.DeepEqual(got, want) {
			gotText, err := prototext.Marshal(got)
			if err != nil {
				gotText = []byte("<error marshalling in test>")
			}
			t.Errorf("test case %d[%q] proto:\n    got: %s\n   want: %s", idx, testCase.label, gotText, testCase.expectProtoText)
		}

	}
}

func TestParseBool(t *testing.T) {
	for idx, testCase := range []struct {
		asString    string
		expectValue bool
		expectError bool
	}{
		{"true", true, false},
		{"false", false, false},
		{"True", false, true},
		{"False", false, true},
		{"0", false, true},
		{"1", false, true},
	} {
		val, err := parseBool(testCase.asString)
		if got, want := (err != nil), testCase.expectError; got != want {
			t.Errorf("test case %d[%q] error: got %v, want %v", idx, testCase.asString, err, want)
			continue
		}
		if got, want := val, testCase.expectValue; got != want {
			t.Errorf("test case %d[%q] got: %v,   want: %v", idx, testCase.asString, got, want)
		}
	}
}
