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
	"encoding/base64"
	"sync"
	"testing"

	"github.com/golang/protobuf/proto"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_User_lifecycle(t *testing.T) {
	s := NewIdentityServer()

	first, err := s.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{
			User: &pb.User{DisplayName: "ekkodog", Email: "ekko@google.com"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	delete, err := s.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{
			User: &pb.User{DisplayName: "mishacat", Email: "misha@google.com"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	_, err = s.DeleteUser(
		context.Background(),
		&pb.DeleteUserRequest{Name: delete.Name})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
	}

	created, err := s.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{
			User: &pb.User{DisplayName: "rumbledog", Email: "rumble@google.com"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	got, err := s.GetUser(
		context.Background(),
		&pb.GetUserRequest{Name: created.GetName()})
	if err != nil {
		t.Errorf("Get: unexpected err %+v", err)
	}
	if !proto.Equal(created, got) {
		t.Error("Expected to get created user.")
	}

	got.DisplayName = "musubi"
	updated, err := s.UpdateUser(
		context.Background(),
		&pb.UpdateUserRequest{User: got, UpdateMask: nil})
	if err != nil {
		t.Errorf("Update: unexpected err %+v", err)
	}

	got, err = s.GetUser(
		context.Background(),
		&pb.GetUserRequest{Name: updated.GetName()})
	if err != nil {
		t.Errorf("Get: unexpected err %+v", err)
	}
	// Cannot use proto.Equal here because the update time is changed on updates.
	if updated.GetName() != got.GetName() ||
		updated.GetDisplayName() != got.GetDisplayName() ||
		updated.GetEmail() != got.GetEmail() ||
		!proto.Equal(updated.GetCreateTime(), got.GetCreateTime()) ||
		proto.Equal(updated.GetUpdateTime(), got.GetUpdateTime()) {
		t.Error("Expected to get updated user.")
	}

	r, err := s.ListUsers(
		context.Background(),
		&pb.ListUsersRequest{PageSize: 1, PageToken: ""})
	if len(r.GetUsers()) != 1 {
		t.Errorf("List want: page size %d, got %d", 1, len(r.GetUsers()))
	}
	if !proto.Equal(first, r.GetUsers()[0]) {
		t.Errorf("List want: first user %+v, got %+v", first, r.GetUsers()[0])
	}
	if r.GetNextPageToken() == "" {
		t.Error("List want: non empty next page token")
	}

	r, err = s.ListUsers(
		context.Background(),
		&pb.ListUsersRequest{PageSize: 10, PageToken: r.GetNextPageToken()})
	if len(r.GetUsers()) != 1 {
		t.Errorf("List want: page size %d, got %d", 1, len(r.GetUsers()))
	}
	if !proto.Equal(got, r.GetUsers()[0]) {
		t.Errorf("List want: updated user %+v, got %+v", first, r.GetUsers()[0])
	}
	if r.GetNextPageToken() != "" {
		t.Error("List want: empty next page token")
	}
}

func Test_Create_invalid(t *testing.T) {
	tests := []*pb.User{
		&pb.User{DisplayName: "", Email: "rumble@google.com"},
		&pb.User{DisplayName: "Rumble", Email: ""},
	}
	s := NewIdentityServer()
	for _, u := range tests {
		_, err := s.CreateUser(
			context.Background(),
			&pb.CreateUserRequest{User: u})
		status, _ := status.FromError(err)
		if status.Code() != codes.InvalidArgument {
			t.Errorf(
				"Create: Want error code %d got %d",
				codes.InvalidArgument,
				status.Code())
		}
	}
}

func Test_Create_alreadyPresent(t *testing.T) {
	first := []*pb.User{
		&pb.User{DisplayName: "Ekko", Email: "ekko@google.com"},
		&pb.User{DisplayName: "Rumble", Email: "rumble@google.com"},
	}
	second := []*pb.User{
		&pb.User{DisplayName: "Ekko", Email: "musubigoogle.com"},
		&pb.User{DisplayName: "Musubi", Email: "rumble@google.com"},
	}

	s := NewIdentityServer()
	for _, u := range first {
		_, err := s.CreateUser(context.Background(), &pb.CreateUserRequest{User: u})
		if err != nil {
			t.Errorf("Create: unexpected err %+v", err)
		}
	}
	for _, u := range second {
		_, err := s.CreateUser(context.Background(), &pb.CreateUserRequest{User: u})
		status, _ := status.FromError(err)
		if status.Code() != codes.AlreadyExists {
			t.Errorf(
				"Create: Want error code %d got %d",
				codes.AlreadyExists,
				status.Code())
		}
	}
}

func Test_Get_notFound(t *testing.T) {
	s := NewIdentityServer()
	_, err := s.GetUser(
		context.Background(),
		&pb.GetUserRequest{Name: "Rumble"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Get: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_Get_deleted(t *testing.T) {
	s := NewIdentityServer()
	created, err := s.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{
			User: &pb.User{DisplayName: "rumbledog", Email: "rumble@google.com"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	_, err = s.DeleteUser(
		context.Background(),
		&pb.DeleteUserRequest{Name: created.Name})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
	}

	_, err = s.GetUser(
		context.Background(),
		&pb.GetUserRequest{Name: created.Name})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Get deleted: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_Update_fieldmask(t *testing.T) {
	s := NewIdentityServer()
	_, err := s.UpdateUser(
		context.Background(),
		&pb.UpdateUserRequest{
			User:       nil,
			UpdateMask: &field_mask.FieldMask{Paths: []string{"email"}},
		})
	status, _ := status.FromError(err)
	if status.Code() != codes.Unimplemented {
		t.Errorf(
			"Update: Want error code %d got %d",
			codes.Unimplemented,
			status.Code())
	}
}

func Test_Update_notFound(t *testing.T) {
	s := NewIdentityServer()
	_, err := s.UpdateUser(
		context.Background(),
		&pb.UpdateUserRequest{
			User: &pb.User{
				Name:        "Rumble",
				DisplayName: "rumbledog",
				Email:       "rumble@google.com",
			},
			UpdateMask: nil,
		})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Update: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_Update_invalid(t *testing.T) {
	first := []*pb.User{
		&pb.User{DisplayName: "Ekko", Email: "ekko@google.com"},
		&pb.User{DisplayName: "Rumble", Email: "rumble@google.com"},
	}
	second := []*pb.User{
		&pb.User{DisplayName: "", Email: "ekko@google.com"},
		&pb.User{DisplayName: "Rumble", Email: ""},
	}
	s := NewIdentityServer()
	for i, u := range first {
		created, err := s.CreateUser(
			context.Background(),
			&pb.CreateUserRequest{User: u})
		if err != nil {
			t.Errorf("Create: unexpected err %+v", err)
		}
		second[i].Name = created.GetName()
	}
	for _, u := range second {
		_, err := s.UpdateUser(
			context.Background(),
			&pb.UpdateUserRequest{User: u, UpdateMask: nil})
		status, _ := status.FromError(err)
		if status.Code() != codes.InvalidArgument {
			t.Errorf(
				"Update: Want error code %d got %d",
				codes.InvalidArgument,
				status.Code())
		}
	}
}

func Test_Update_alreadyPresent(t *testing.T) {
	first := []*pb.User{
		&pb.User{DisplayName: "Ekko", Email: "ekko@google.com"},
		&pb.User{DisplayName: "Rumble", Email: "rumble@google.com"},
	}
	second := []*pb.User{
		&pb.User{DisplayName: "Rumble", Email: "ekko@google.com"},
		&pb.User{DisplayName: "Rumble", Email: "ekko@google.com"},
	}

	s := NewIdentityServer()
	for i, u := range first {
		created, err := s.CreateUser(
			context.Background(),
			&pb.CreateUserRequest{User: u})
		if err != nil {
			t.Errorf("Create: unexpected err %+v", err)
		}
		second[i].Name = created.GetName()

	}
	for _, u := range second {
		_, err := s.UpdateUser(
			context.Background(),
			&pb.UpdateUserRequest{User: u, UpdateMask: nil})
		status, _ := status.FromError(err)
		if status.Code() != codes.AlreadyExists {
			t.Errorf(
				"Update: Want error code %d got %d",
				codes.AlreadyExists,
				status.Code())
		}
	}
}

func Test_Delete_notFound(t *testing.T) {
	s := NewIdentityServer()
	_, err := s.DeleteUser(
		context.Background(),
		&pb.DeleteUserRequest{Name: "Rumble"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Delete: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_List_invalidToken(t *testing.T) {
	db := &userDb{
		uid:   &uniqID{},
		token: &tokenGenerator{salt: "Ekko"},
		mu:    sync.Mutex{},
		keys:  map[string]int{},
		users: []userEntry{},
	}
	s := identityServerImpl{db: db}

	tests := []string{
		"1", // Not base64 encoded
		base64.StdEncoding.EncodeToString([]byte("1")),          // No salt
		base64.StdEncoding.EncodeToString([]byte("EkkoRumble")), // Invalid index
	}

	for _, token := range tests {
		_, err := s.ListUsers(
			context.Background(),
			&pb.ListUsersRequest{PageSize: 1, PageToken: token})
		status, _ := status.FromError(err)
		if status.Code() != codes.InvalidArgument {
			t.Errorf(
				"List: Want error code %d got %d",
				codes.InvalidArgument,
				status.Code())
		}
	}
}
