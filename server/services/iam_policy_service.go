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
	"sync"

	iampb "cloud.google.com/go/iam/apiv1/iampb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var missingResource error = status.Error(codes.InvalidArgument, "Missing required argument: resource")
var missingPolicy error = status.Error(codes.InvalidArgument, "Missing required argument: policy")
var resourceDoesNotExist error = status.Error(codes.NotFound, "Requested resource has no Policy")

// NewIAMPolicyServer returns a new LocationsServer for the Showcase API.
func NewIAMPolicyServer() iampb.IAMPolicyServer {
	return &iamPolicyServerImpl{
		policies: map[string]*iampb.Policy{},
	}
}

type iamPolicyServerImpl struct {
	mu sync.Mutex
	// The key is the resource name.
	// The value is the IAM Policy assigned to the resource.
	policies map[string]*iampb.Policy
}

// GetIamPolicy retrieves the Policy for the named resource, if it has one.
func (i *iamPolicyServerImpl) GetIamPolicy(ctx context.Context, in *iampb.GetIamPolicyRequest) (*iampb.Policy, error) {
	if in.GetResource() == "" {
		return nil, missingResource
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	p, ok := i.policies[in.GetResource()]
	if !ok {
		return nil, resourceDoesNotExist
	}
	res := proto.Clone(p)

	return res.(*iampb.Policy), nil
}

// SetIamPolicy assigns the given Policy to the resource.
func (i *iamPolicyServerImpl) SetIamPolicy(ctx context.Context, in *iampb.SetIamPolicyRequest) (*iampb.Policy, error) {
	if in.GetResource() == "" {
		return nil, missingResource
	}
	if in.GetPolicy() == nil {
		return nil, missingPolicy
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	clone := proto.Clone(in.GetPolicy())
	p := clone.(*iampb.Policy)
	i.policies[in.GetResource()] = p

	return in.GetPolicy(), nil
}

// TestIamPermissions verifies that the requested resource has been Set at one point, and echoes the back permissions to be tested.
func (i *iamPolicyServerImpl) TestIamPermissions(ctx context.Context, in *iampb.TestIamPermissionsRequest) (*iampb.TestIamPermissionsResponse, error) {
	if in.GetResource() == "" {
		return nil, missingResource
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.policies[in.GetResource()]; !ok {
		return nil, resourceDoesNotExist
	}

	return &iampb.TestIamPermissionsResponse{
		Permissions: in.GetPermissions(),
	}, nil
}
