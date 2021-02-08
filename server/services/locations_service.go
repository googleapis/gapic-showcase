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
	"strings"

	locpb "google.golang.org/genproto/googleapis/cloud/location"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var missingName error = status.Error(codes.InvalidArgument, "Missing required argument: name")

// NewLocationsServer returns a new LocationsServer for the Showcase API.
func NewLocationsServer() locpb.LocationsServer {
	return &locationsServerImpl{}
}

type locationsServerImpl struct{}

// GetLocation echoes back a Location with the same name as
// in.name and has the last segment of name as the display_name.
func (l *locationsServerImpl) GetLocation(ctx context.Context, in *locpb.GetLocationRequest) (*locpb.Location, error) {
	if in.GetName() == "" {
		return nil, missingName
	}

	name := in.GetName()
	display := name[strings.LastIndex(name, "/")+1:]

	return &locpb.Location{
		Name:        name,
		DisplayName: display,
	}, nil
}

// ListLocations returns a fixed list of Locations prefixed with the
// in.name value.
func (l *locationsServerImpl) ListLocations(ctx context.Context, in *locpb.ListLocationsRequest) (*locpb.ListLocationsResponse, error) {
	if in.GetName() == "" {
		return nil, missingName
	}

	return &locpb.ListLocationsResponse{
		Locations: []*locpb.Location{
			{
				Name:        in.GetName() + "/locations/us-north",
				DisplayName: "us-north",
			},
			{
				Name:        in.GetName() + "/locations/us-south",
				DisplayName: "us-south",
			},
			{
				Name:        in.GetName() + "/locations/us-east",
				DisplayName: "us-east",
			},
			{
				Name:        in.GetName() + "/locations/us-west",
				DisplayName: "us-west",
			},
		},
	}, nil
}
