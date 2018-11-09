// Copyright 2018 Google LLC
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

package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
)

func NewIdentityServer() pb.IdentityServer {
	return &identityServerImpl{db: NewUserDb()}
}

type identityServerImpl struct {
	db UserDb
}

// Creates a user.
func (s *identityServerImpl) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.User, error) {
	return s.db.Create(in.GetUser())
}

// Retrieves the User with the given uri.
func (s *identityServerImpl) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	return s.db.Get(in.GetName())
}

// Updates a user.
func (s *identityServerImpl) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.User, error) {
	return s.db.Update(in.GetUser(), in.GetUpdateMask())
}

// Deletes a user, their profile, and all of their authored messages.
func (s *identityServerImpl) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*empty.Empty, error) {
	return nil, s.db.Delete(in.GetName())
}

// Lists all users.
func (s *identityServerImpl) ListUsers(ctx context.Context, in *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	return s.db.List(in.GetPageSize(), in.GetPageToken())
}
