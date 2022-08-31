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

package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewIdentityServer returns a new instance of showcase identity server.
func NewIdentityServer() pb.IdentityServer {
	return &identityServerImpl{
		token: server.NewTokenGenerator(),
		keys:  map[string]int{},
	}
}

type userEntry struct {
	user    *pb.User
	deleted bool
}

// ReadOnlyIdentityServer provides a read-only interface of an identity server.
type ReadOnlyIdentityServer interface {
	GetUser(context.Context, *pb.GetUserRequest) (*pb.User, error)
	ListUsers(context.Context, *pb.ListUsersRequest) (*pb.ListUsersResponse, error)
}

type identityServerImpl struct {
	uid   server.UniqID
	token server.TokenGenerator

	mu    sync.Mutex
	keys  map[string]int
	users []userEntry
}

// Creates a user.
func (s *identityServerImpl) CreateUser(_ context.Context, in *pb.CreateUserRequest) (*pb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u := in.GetUser()

	// Ignore passed in name.
	u.Name = ""

	err := s.validate(u)
	if err != nil {
		return nil, err
	}

	// Assign info.
	id := s.uid.Next()
	name := fmt.Sprintf("users/%d", id)
	now := ptypes.TimestampNow()

	u.Name = name
	u.CreateTime = now
	u.UpdateTime = now

	// Insert.
	index := len(s.users)
	s.users = append(s.users, userEntry{user: u})
	s.keys[name] = index

	return u, nil
}

// Retrieves the User with the given uri.
func (s *identityServerImpl) GetUser(_ context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := in.GetName()
	if i, ok := s.keys[name]; ok {
		entry := s.users[i]
		if !entry.deleted {
			return entry.user, nil
		}
	}

	return nil, status.Errorf(
		codes.NotFound, "A user with name %s not found.",
		name)
}

// Updates a user.
func (s *identityServerImpl) UpdateUser(_ context.Context, in *pb.UpdateUserRequest) (*pb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	mask := in.GetUpdateMask()
	u := in.GetUser()
	i, ok := s.keys[u.GetName()]
	entry := s.users[i]
	if !ok || entry.deleted {
		return nil, status.Errorf(
			codes.NotFound,
			"A user with name %s not found.", u.GetName())
	}

	err := s.validate(u)
	if err != nil {
		return nil, err
	}

	// Update store.
	existing := entry.user
	updated := proto.Clone(existing).(*pb.User)
	applyFieldMask(u.ProtoReflect(), updated.ProtoReflect(), mask.GetPaths())
	updated.CreateTime = existing.GetCreateTime()
	updated.UpdateTime = ptypes.TimestampNow()

	s.users[i] = userEntry{user: updated}
	return updated, nil
}

// Deletes a user, their profile, and all of their authored messages.
func (s *identityServerImpl) DeleteUser(_ context.Context, in *pb.DeleteUserRequest) (*empty.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	i, ok := s.keys[in.GetName()]

	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"A user with name %s not found.", in.GetName())
	}

	entry := s.users[i]
	s.users[i] = userEntry{user: entry.user, deleted: true}

	return &empty.Empty{}, nil
}

// Lists all users.
func (s *identityServerImpl) ListUsers(_ context.Context, in *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	start, err := s.token.GetIndex(in.GetPageToken())
	if err != nil {
		return nil, err
	}

	offset := 0
	users := []*pb.User{}
	for _, entry := range s.users[start:] {
		offset++
		if entry.deleted {
			continue
		}
		users = append(users, entry.user)
		if len(users) >= int(in.GetPageSize()) {
			break
		}
	}

	nextToken := ""
	if start+offset < len(s.users) {
		nextToken = s.token.ForIndex(start + offset)
	}

	return &pb.ListUsersResponse{Users: users, NextPageToken: nextToken}, nil
}

func (s *identityServerImpl) validate(u *pb.User) error {
	// Validate Required Fields.
	if u.GetDisplayName() == "" {
		return status.Errorf(
			codes.InvalidArgument,
			"The field `display_name` is required.")
	}
	if u.GetEmail() == "" {
		return status.Errorf(
			codes.InvalidArgument,
			"The field `email` is required.")
	}
	// Validate Unique Fields.
	for _, x := range s.users {
		if x.deleted {
			continue
		}
		if (u.GetEmail() == x.user.GetEmail()) &&
			(u.GetName() != x.user.GetName()) {
			return status.Errorf(
				codes.AlreadyExists,
				"A user with email %s already exists.",
				u.GetEmail())
		}
	}
	return nil
}
