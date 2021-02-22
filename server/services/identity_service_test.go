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
	"encoding/base64"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/googleapis/gapic-showcase/server"
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
			User: &pb.User{DisplayName: "ekkodog", Email: "ekko@example.com", Nickname: proto.String("ekko"), Age: proto.Int32(26)},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	delete, err := s.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{
			User: &pb.User{DisplayName: "mishacat", Email: "misha@example.com"},
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
			User: &pb.User{
				DisplayName:         "rumbledog",
				Email:               "rumble@example.com",
				Age:                 proto.Int32(42),
				EnableNotifications: proto.Bool(false),
				HeightFeet:          proto.Float64(3.5),
				Nickname:            proto.String("rumble"),
			},
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

	// Make a copy of the User value, then unset proto3_optional fields
	// to scope the update to DisplayName and Nickname.
	clone := proto.Clone(got)
	u := clone.(*pb.User)
	u.DisplayName = "musubi"
	u.Nickname = proto.String("musu")
	u.Age = nil
	u.HeightFeet = nil
	u.EnableNotifications = nil
	updated, err := s.UpdateUser(
		context.Background(),
		&pb.UpdateUserRequest{User: u, UpdateMask: nil})
	if err != nil {
		t.Errorf("Update: unexpected err %+v", err)
	}
	// Ensure the proto3_optional fields that were unset did not get updated.
	if updated.Age == nil {
		t.Errorf("UpdateUser().Age was unexpectedly set to nil")
	}
	if updated.HeightFeet == nil {
		t.Errorf("UpdateUser().HeightFeet was unexpectedly set to nil")
	}
	if updated.EnableNotifications == nil {
		t.Errorf("UpdateUser().EnableNotifications was unexpectedly set to nil")
	}

	got, err = s.GetUser(
		context.Background(),
		&pb.GetUserRequest{Name: updated.GetName()})
	if err != nil {
		t.Errorf("Get: unexpected err %+v", err)
	}
	if !proto.Equal(updated, got) {
		t.Errorf("UpdateUser() = %+v, want %+v", got, updated)
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
		t.Errorf("List want: updated user %+v, got %+v", got, r.GetUsers()[0])
	}
	if r.GetNextPageToken() != "" {
		t.Error("List want: empty next page token")
	}
}

func Test_Create_invalid(t *testing.T) {
	tests := []*pb.User{
		{DisplayName: "", Email: "rumble@example.com"},
		{DisplayName: "Rumble", Email: ""},
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
	first := &pb.User{DisplayName: "Rumble", Email: "rumble@example.com"}
	second := &pb.User{DisplayName: "Musubi", Email: "rumble@example.com"}
	third := &pb.User{DisplayName: "Rumble", Email: "rumble@example.com"}

	s := NewIdentityServer()
	created, err := s.CreateUser(context.Background(), &pb.CreateUserRequest{User: first})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	_, err = s.CreateUser(context.Background(), &pb.CreateUserRequest{User: second})
	stat, _ := status.FromError(err)
	if stat.Code() != codes.AlreadyExists {
		t.Errorf(
			"Create: Want error code %d got %d",
			codes.AlreadyExists,
			stat.Code())
	}

	third.Name = created.GetName()
	_, err = s.CreateUser(context.Background(), &pb.CreateUserRequest{User: third})
	stat, _ = status.FromError(err)
	if stat.Code() != codes.AlreadyExists {
		t.Errorf(
			"Create: Want error code %d got %d",
			codes.AlreadyExists,
			stat.Code())
	}
}

func Test_Create_deleted(t *testing.T) {
	s := NewIdentityServer()
	u := &pb.User{DisplayName: "Rumble", Email: "rumble@example.com"}
	u, err := s.CreateUser(context.Background(), &pb.CreateUserRequest{User: u})
	if err != nil {
		t.Errorf("Create: unexpected error: %v", err)
	}
	_, err = s.DeleteUser(context.Background(), &pb.DeleteUserRequest{Name: u.GetName()})
	if err != nil {
		t.Errorf("Delete: unexpected error: %v", err)
	}
	u, err = s.CreateUser(context.Background(), &pb.CreateUserRequest{User: u})
	if err != nil {
		t.Errorf("Create: a using a deleted email should be allowed, got: %v", err)
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
			User: &pb.User{DisplayName: "rumbledog", Email: "rumble@example.com"},
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
				Email:       "rumble@example.com",
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
		{DisplayName: "Ekko", Email: "ekko@example.com"},
		{DisplayName: "Rumble", Email: "rumble@example.com"},
	}
	second := []*pb.User{
		{DisplayName: "", Email: "ekko@example.com"},
		{DisplayName: "Rumble", Email: ""},
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
	first := &pb.User{DisplayName: "Rumble", Email: "rumble@example.com"}
	second := &pb.User{DisplayName: "Ekko", Email: "ekko@example.com"}

	s := NewIdentityServer()

	first, err := s.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{User: first})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	second, err = s.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{User: second})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	second.Email = first.GetEmail()

	_, err = s.UpdateUser(
		context.Background(),
		&pb.UpdateUserRequest{User: second, UpdateMask: nil})
	status, _ := status.FromError(err)
	if status.Code() != codes.AlreadyExists {
		t.Errorf(
			"Update: Want error code %d got %d",
			codes.AlreadyExists,
			status.Code())
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
	s := &identityServerImpl{
		token: server.NewTokenGenerator(),
		keys:  map[string]int{},
	}

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
