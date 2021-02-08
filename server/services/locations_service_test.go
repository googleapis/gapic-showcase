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
	locpb "google.golang.org/genproto/googleapis/cloud/location"
	"google.golang.org/protobuf/proto"
)

func TestListLocations(t *testing.T) {
	s := NewLocationsServer()
	name := "projects/foo"

	got, err := s.ListLocations(context.Background(), &locpb.ListLocationsRequest{Name: name})
	if err != nil {
		t.Error(err)
	}

	want := &locpb.ListLocationsResponse{
		Locations: []*locpb.Location{
			{
				Name:        name + "/locations/us-north",
				DisplayName: "us-north",
			},
			{
				Name:        name + "/locations/us-south",
				DisplayName: "us-south",
			},
			{
				Name:        name + "/locations/us-east",
				DisplayName: "us-east",
			},
			{
				Name:        name + "/locations/us-west",
				DisplayName: "us-west",
			},
		},
	}

	if diff := cmp.Diff(got, want, cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("ListLocations: got(-),want(+):\n%s", diff)
	}
}

func TestListLocations_missingName(t *testing.T) {
	s := NewLocationsServer()

	res, err := s.ListLocations(context.Background(), &locpb.ListLocationsRequest{})
	if err == nil {
		t.Errorf("ListLocations_missingName: expected an error but got response: %v", res)
	}

	if diff := cmp.Diff(err, missingName, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("ListLocations_missingName: got(-),want(+):\n%s", diff)
	}
}

func TestGetLocation(t *testing.T) {
	s := NewLocationsServer()
	name := "projects/foo/locations/us-west"
	want := &locpb.Location{
		Name:        name,
		DisplayName: "us-west",
	}

	got, err := s.GetLocation(context.Background(), &locpb.GetLocationRequest{Name: name})
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want, cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("GetLocation: got(-),want(+):\n%s", diff)
	}
}

func TestGetLocation_missingName(t *testing.T) {
	s := NewLocationsServer()

	res, err := s.GetLocation(context.Background(), &locpb.GetLocationRequest{})
	if err == nil {
		t.Errorf("GetLocation_missingName: expected an error but got response: %v", res)
	}

	if diff := cmp.Diff(err, missingName, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("GetLocation_missingName: got(-),want(+):\n%s", diff)
	}
}
