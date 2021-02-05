// Copyright 2021 Google LLC
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

package services

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/protobuf/proto"
)

func TestIamPolicySetGetTest(t *testing.T) {
	s := NewIAMPolicyServer()

	resource := "projects/myproject/foos/foo123"
	policy := &iampb.Policy{
		Version: 1,
		Bindings: []*iampb.Binding{
			{Role: "project/fooEditor"},
		},
	}
	permissions := []string{"foo.write", "foo.read"}

	p, err := s.SetIamPolicy(context.Background(), &iampb.SetIamPolicyRequest{Resource: resource, Policy: policy})
	if err != nil {
		t.Errorf("TestIamPolicySetGetTest > SetIamPolicy: unexpected error %q", err)
	}

	if diff := cmp.Diff(p, policy, cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("TestIamPolicySetGetTest > SetIamPolicy: got(-),want(+):\n%s", diff)
	}

	p, err = s.GetIamPolicy(context.Background(), &iampb.GetIamPolicyRequest{Resource: resource})
	if err != nil {
		t.Errorf("TestIamPolicySetGetTest > GetIamPolicy: unexpected error %q", err)
	}

	if diff := cmp.Diff(p, policy, cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("TestIamPolicySetGetTest > GetIamPolicy: got(-),want(+):\n%s", diff)
	}

	res, err := s.TestIamPermissions(context.Background(), &iampb.TestIamPermissionsRequest{
		Resource:    resource,
		Permissions: permissions,
	})
	if err != nil {
		t.Errorf("TestIamPolicySetGetTest > TestIamPermissions: unexpected error %q", err)
	}

	if diff := cmp.Diff(res.Permissions, permissions); diff != "" {
		t.Errorf("TestIamPolicySetGetTest > TestIamPermissions: got(-),want(+):\n%s", diff)
	}
}

func TestGetIamPolicy_errors(t *testing.T) {
	s := NewIAMPolicyServer()

	for _, tst := range []struct {
		name, resource string
		err            error
	}{
		{
			name: "missing resource field",
			err:  missingResource,
		},
		{
			name:     "resource does not exist",
			resource: "foo/bar",
			err:      resourceDNE,
		},
	} {
		res, err := s.GetIamPolicy(context.Background(), &iampb.GetIamPolicyRequest{Resource: tst.resource})
		if err != nil {
			if diff := cmp.Diff(err, tst.err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetIamPolicy_errors(%s): got(-),want(+):\n%s", tst.name, diff)
			}
		} else {
			t.Errorf("GetIamPolicy_errors(%s): expected an error but got response: %v", tst.name, res)
		}
	}
}

func TestSetIamPolicy_errors(t *testing.T) {
	s := NewIAMPolicyServer()

	for _, tst := range []struct {
		name, resource string
		err            error
	}{
		{
			name: "missing resource field",
			err:  missingResource,
		},
		{
			name:     "missing policy field",
			resource: "foo/bar",
			err:      missingPolicy,
		},
	} {
		res, err := s.SetIamPolicy(context.Background(), &iampb.SetIamPolicyRequest{Resource: tst.resource})
		if err != nil {
			if diff := cmp.Diff(err, tst.err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("SetIamPolicy_errors(%s): got(-),want(+):\n%s", tst.name, diff)
			}
		} else {
			t.Errorf("SetIamPolicy_errors(%s): expected an error but got response: %v", tst.name, res)
		}
	}
}

func TestTestIamPermissions_errors(t *testing.T) {
	s := NewIAMPolicyServer()

	for _, tst := range []struct {
		name, resource string
		err            error
	}{
		{
			name: "missing resource field",
			err:  missingResource,
		},
		{
			name:     "resource does not exist",
			resource: "foo/bar",
			err:      resourceDNE,
		},
	} {
		res, err := s.TestIamPermissions(context.Background(), &iampb.TestIamPermissionsRequest{Resource: tst.resource})
		if err != nil {
			if diff := cmp.Diff(err, tst.err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("TestIamPermissions_errors(%s): got(-),want(+):\n%s", tst.name, diff)
			}
		} else {
			t.Errorf("TestIamPermissions_errors(%s): expected an error but got response: %v", tst.name, res)
		}
	}
}
