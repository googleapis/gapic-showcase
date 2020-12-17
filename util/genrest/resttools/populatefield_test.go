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

	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/protobuf/encoding/prototext"
)

func TestPopulateOneFieldError(t *testing.T) {
	for idx, testCase := range []struct {
		field string
		value string
	}{
		// field path errors

		{".f_string", "hi"},
		{"subpack..f_string", "hi"},
		{"subpack.", "hi"},
		{"subpack.x", "hi"},
		{"subpack.x.f_string", "hi"},
		{"subpack.f_string.subpack", "hi"},
		{"subpack.subpack", "hi"},

		// parsing errors

		{"f_int32", "hello"},
		{"f_sint32", "hello"},
		{"f_sfixed32", "hello"},
		{"f_int32", "2147483648"},    // max int32 + 1
		{"f_sint32", "2147483648"},   // max int32 + 1
		{"f_sfixed32", "2147483648"}, // max int32 + 1
		{"f_int32", "1.1"},
		{"f_sint32", "1.1"},
		{"f_sfixed32", "1.1"},

		{"f_uint32", "hello"},
		{"f_fixed32", "hello"},
		{"f_uint32", "4294967296"},  // max uint32 + 1
		{"f_fixed32", "4294967296"}, // max uint32 + 1
		{"f_uint32", "1.1"},
		{"f_fixed32", "1.1"},
		{"f_uint32", "-1"},
		{"f_fixed32", "-1"},

		{"f_int64", "hello"},
		{"f_sint64", "hello"},
		{"f_sfixed64", "hello"},
		{"f_int64", "9223372036854775808"},    // max int64 + 1
		{"f_sint64", "9223372036854775808"},   // max int64 + 1
		{"f_sfixed64", "9223372036854775808"}, // max int64 + 1
		{"f_int64", "1.1"},
		{"f_sint64", "1.1"},
		{"f_sfixed64", "1.1"},

		{"f_uint64", "hello"},
		{"f_fixed64", "hello"},
		{"f_uint64", "18446744073709551616"},  // max uint64 + 1
		{"f_fixed64", "18446744073709551616"}, // max uint64 + 1
		{"f_uint64", "1.1"},
		{"f_fixed64", "1.1"},
		{"f_uint64", "-1"},
		{"f_fixed64", "-1"},

		{"f_float", "hello"},
		{"f_double", "hello"},
		{"f_float", "1e+39"},   // exponent too large
		{"f_double", "1e+309"}, // exponent too large

		{"f_bool", "hello"},
		{"f_bool", "13"},
	} {
		dataPack := &genprotopb.DataPack{}
		err := PopulateOneField(dataPack, testCase.field, []string{testCase.value})
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
				"f_string": "alphabet",

				"f_int32":    "2147483647", // max int32
				"f_sint32":   "2147483647", // max int32
				"f_sfixed32": "2147483647", // max int32

				"f_uint32":  "4294967295", // max uint32
				"f_fixed32": "4294967295", // max uint32

				"f_int64":    "9223372036854775807", // max int64
				"f_sint64":   "9223372036854775807", // max int64
				"f_sfixed64": "9223372036854775807", // max int64

				"f_uint64":  "18446744073709551615", // max uint64
				"f_fixed64": "18446744073709551615", // max uint64

				"f_float":  "3.40282346638528859811704183484516925440e+38",   // max float32 (https://golang.org/pkg/math/#pkg-constants)
				"f_double": "1.797693134862315708145274237317043567981e+308", // max float64

				"f_bool": "true",

				"f_bytes": "greetings",
			},
			expectProtoText: `f_string:"alphabet" f_int32:2147483647 f_sint32:2147483647 f_sfixed32:2147483647 f_uint32:4294967295 f_fixed32:4294967295 f_int64:9223372036854775807 f_sint64:9223372036854775807 f_sfixed64:9223372036854775807 f_uint64:18446744073709551615 f_fixed64:18446744073709551615 f_double:1.7976931348623157e+308 f_float:3.4028235e+38 f_bool:true f_bytes:"greetings"`,
		},
		{
			label: "nested messages",
			fields: map[string]string{
				"f_string":                 "alphabet",
				"subpack.subpack.f_string": "lexicon",
				"f_int32":                  "5",
				"subpack.subpack.f_double": "53.47",
				"subpack.subpack.f_int32":  "-6",
				"subpack.f_bool":           "1", // NOTE: this gets parsed as "true"
			},
			expectProtoText: `subpack:{subpack:{f_string:"lexicon" f_int32:-6 f_double:53.47} f_bool:true} f_string:"alphabet" f_int32:5`,
		},
		{
			label: "presence/zero values",
			fields: map[string]string{
				"f_string": "",

				"f_int32":    "0",
				"f_sint32":   "0",
				"f_sfixed32": "0",

				"f_uint32":  "0",
				"f_fixed32": "0",

				"f_int64":    "0",
				"f_sint64":   "0",
				"f_sfixed64": "0",

				"f_uint64":  "0",
				"f_fixed64": "0",

				"f_float":  "0",
				"f_double": "0",

				"f_bool": "false",

				"f_bytes": "",

				"p_string": "",
				"p_int32":  "0",
				"p_double": "0",
				"p_bool":   "0", // NOTE: this gets parsed as "false"
			},
			expectProtoText: `p_string:""  p_int32:0  p_double:0  p_bool:false`,
		},
	} {
		dataPack := &genprotopb.DataPack{}
		err := PopulateSingularFields(dataPack, testCase.fields)
		if got, want := (err != nil), testCase.expectError; got != want {
			t.Errorf("test case %d[%q] error: got %v, want %v", idx, testCase.label, err, want)
			continue
		}
		if testCase.expectError {
			continue
		}

		var expectProto genprotopb.DataPack
		err = prototext.Unmarshal([]byte(testCase.expectProtoText), &expectProto)
		if err != nil {
			t.Errorf("test case %d[%q] unexpected error unmarshaling expected proto: %s", idx, testCase.label, err)
			continue
		}

		if got, want := dataPack, &expectProto; !reflect.DeepEqual(got, want) {
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
				"f_string": []string{"alphabet"},
			},
			expectProtoText: `f_string:"alphabet"`,
		},
		{
			label: "repeated fields",
			fields: map[string][]string{
				"f_string": []string{"alphabet", "lexicon"},
			},
			// TODO: Make Populate*Field*() work with repeated fields.
			expectError: true,
		},
	} {
		dataPack := &genprotopb.DataPack{}
		err := PopulateFields(dataPack, testCase.fields)
		if got, want := (err != nil), testCase.expectError; got != want {
			t.Errorf("test case %d[%q] error: got %v, want %v", idx, testCase.label, err, want)
			continue
		}
		if testCase.expectError {
			continue
		}

		var expectProto genprotopb.DataPack
		err = prototext.Unmarshal([]byte(testCase.expectProtoText), &expectProto)
		if err != nil {
			t.Errorf("test case %d[%q] unexpected error unmarshaling expected proto: %s", idx, testCase.label, err)
			continue
		}

		if got, want := dataPack, &expectProto; !reflect.DeepEqual(got, want) {
			gotText, err := prototext.Marshal(got)
			if err != nil {
				gotText = []byte("<error marshalling in test>")
			}
			t.Errorf("test case %d[%q] proto:\n    got: %s\n   want: %s", idx, testCase.label, gotText, testCase.expectProtoText)
		}

	}
}
