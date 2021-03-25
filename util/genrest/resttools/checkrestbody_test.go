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
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestComputeEnumFields(t *testing.T) {
	complianceData := &genprotopb.ComplianceData{}
	want := [][]protoreflect.Name{{"f_kingdom"}, {"p_kingdom"}}
	got := ComputeEnumFields(complianceData.ProtoReflect())
	if !reflect.DeepEqual(got, want) {
		t.Errorf("enum fields: got %v, want %v", got, want)
	}

}
