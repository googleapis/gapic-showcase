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
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Room_lifecycle(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())

	first, err := s.CreateRoom(
		context.Background(),
		&pb.CreateRoomRequest{
			Room: &pb.Room{DisplayName: "Living Room"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	delete, err := s.CreateRoom(
		context.Background(),
		&pb.CreateRoomRequest{
			Room: &pb.Room{DisplayName: "Library"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	_, err = s.DeleteRoom(
		context.Background(),
		&pb.DeleteRoomRequest{Name: delete.Name})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
	}

	created, err := s.CreateRoom(
		context.Background(),
		&pb.CreateRoomRequest{
			Room: &pb.Room{DisplayName: "Weight Room"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	got, err := s.GetRoom(
		context.Background(),
		&pb.GetRoomRequest{Name: created.GetName()})
	if err != nil {
		t.Errorf("Get: unexpected err %+v", err)
	}
	if !proto.Equal(created, got) {
		t.Error("Expected to get created room.")
	}

	got.DisplayName = "Library"
	_, err = s.UpdateRoom(
		context.Background(),
		&pb.UpdateRoomRequest{Room: got, UpdateMask: nil})
	if err != nil {
		t.Errorf("Update: unexpected err %+v", err)
	}

	updated, err := s.GetRoom(
		context.Background(),
		&pb.GetRoomRequest{Name: got.GetName()})
	if err != nil {
		t.Errorf("Get: unexpected err %+v", err)
	}
	// Cannot use proto.Equal here because the update time is changed on updates.
	if updated.GetName() != got.GetName() ||
		updated.GetDisplayName() != got.GetDisplayName() ||
		!proto.Equal(updated.GetCreateTime(), got.GetCreateTime()) ||
		proto.Equal(updated.GetUpdateTime(), got.GetUpdateTime()) {
		t.Errorf("Expected to get updated room. Want: %+v, got %+v", got, updated)
	}

	r, err := s.ListRooms(
		context.Background(),
		&pb.ListRoomsRequest{PageSize: 1, PageToken: ""})
	if len(r.GetRooms()) != 1 {
		t.Errorf("List want: page size %d, got %d", 1, len(r.GetRooms()))
	}
	if !proto.Equal(first, r.GetRooms()[0]) {
		t.Errorf("List want: first room %+v, got %+v", first, r.GetRooms()[0])
	}
	if r.GetNextPageToken() == "" {
		t.Error("List want: non empty next page token")
	}

	r, err = s.ListRooms(
		context.Background(),
		&pb.ListRoomsRequest{PageSize: 10, PageToken: r.GetNextPageToken()})
	if len(r.GetRooms()) != 1 {
		t.Errorf("List want: page size %d, got %d", 1, len(r.GetRooms()))
	}
	if !proto.Equal(updated, r.GetRooms()[0]) {
		t.Errorf("List want: updated room %+v, got %+v", updated, r.GetRooms()[0])
	}
	if r.GetNextPageToken() != "" {
		t.Error("List want: empty next page token")
	}
}

func Test_CreateRoom_invalid(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())
	_, err := s.CreateRoom(
		context.Background(),
		&pb.CreateRoomRequest{Room: &pb.Room{DisplayName: ""}})
	status, _ := status.FromError(err)
	if status.Code() != codes.InvalidArgument {
		t.Errorf(
			"Create: Want error code %d got %d",
			codes.InvalidArgument,
			status.Code())
	}
}

func Test_CreateRoom_alreadyPresent(t *testing.T) {
	room := &pb.Room{DisplayName: "Living Room"}

	s := NewMessagingServer(NewIdentityServer())
	_, err := s.CreateRoom(context.Background(), &pb.CreateRoomRequest{Room: room})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}
	_, err = s.CreateRoom(context.Background(), &pb.CreateRoomRequest{Room: room})
	status, _ := status.FromError(err)
	if status.Code() != codes.AlreadyExists {
		t.Errorf(
			"Create: Want error code %d got %d",
			codes.AlreadyExists,
			status.Code())
	}
}

func Test_GetRoom_notFound(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())
	_, err := s.GetRoom(
		context.Background(),
		&pb.GetRoomRequest{Name: "Weight Room"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Get: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_GetRoom_deleted(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())
	created, err := s.CreateRoom(
		context.Background(),
		&pb.CreateRoomRequest{
			Room: &pb.Room{DisplayName: "Weight Room"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	_, err = s.DeleteRoom(
		context.Background(),
		&pb.DeleteRoomRequest{Name: created.GetName()})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
	}

	_, err = s.GetRoom(
		context.Background(),
		&pb.GetRoomRequest{Name: created.GetName()})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Get deleted: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_UpdateRoom_fieldmask(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())
	_, err := s.UpdateRoom(
		context.Background(),
		&pb.UpdateRoomRequest{
			Room:       nil,
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

func Test_UpdateRoom_notFound(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())
	_, err := s.UpdateRoom(
		context.Background(),
		&pb.UpdateRoomRequest{
			Room: &pb.Room{
				Name:        "rooms/abc",
				DisplayName: "Weight Room",
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

func Test_UpdateRoom_invalid(t *testing.T) {
	first := &pb.Room{DisplayName: "Living Room"}
	second := &pb.Room{DisplayName: ""}
	s := NewMessagingServer(NewIdentityServer())
	created, err := s.CreateRoom(
		context.Background(),
		&pb.CreateRoomRequest{Room: first})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}
	second.Name = created.GetName()
	_, err = s.UpdateRoom(
		context.Background(),
		&pb.UpdateRoomRequest{Room: second, UpdateMask: nil})
	status, _ := status.FromError(err)
	if status.Code() != codes.InvalidArgument {
		t.Errorf(
			"Update: Want error code %d got %d",
			codes.InvalidArgument,
			status.Code())
	}
}

func Test_UpdateRoom_alreadyPresent(t *testing.T) {
	first := []*pb.Room{
		&pb.Room{DisplayName: "Living Room"},
		&pb.Room{DisplayName: "Weight Room"},
	}
	second := &pb.Room{DisplayName: "Weight Room"}

	s := NewMessagingServer(NewIdentityServer())
	for i, r := range first {
		created, err := s.CreateRoom(
			context.Background(),
			&pb.CreateRoomRequest{Room: r})
		if err != nil {
			t.Errorf("Create: unexpected err %+v", err)
		}
		first[i].Name = created.GetName()
	}
	second.Name = first[0].GetName()
	_, err := s.UpdateRoom(
		context.Background(),
		&pb.UpdateRoomRequest{Room: second, UpdateMask: nil})
	status, _ := status.FromError(err)
	if status.Code() != codes.AlreadyExists {
		t.Errorf(
			"Update: Want error code %d got %d",
			codes.AlreadyExists,
			status.Code())
	}
}

func Test_DeleteRoom_notFound(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())
	_, err := s.DeleteRoom(
		context.Background(),
		&pb.DeleteRoomRequest{Name: "Weight Room"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Delete: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_ListRooms_invalidToken(t *testing.T) {
	s := messagingServerImpl{
		token:    &tokenGenerator{salt: "salt"},
		roomKeys: map[string]int{},
	}

	tests := []string{
		"1", // Not base64 encoded
		base64.StdEncoding.EncodeToString([]byte("1")),        // No salt
		base64.StdEncoding.EncodeToString([]byte("saltblah")), // Invalid index
	}

	for _, token := range tests {
		_, err := s.ListRooms(
			context.Background(),
			&pb.ListRoomsRequest{PageSize: 1, PageToken: token})
		status, _ := status.FromError(err)
		if status.Code() != codes.InvalidArgument {
			t.Errorf(
				"List: Want error code %d got %d",
				codes.InvalidArgument,
				status.Code())
		}
	}
}

type mockIdentityServer struct{}

func (m *mockIdentityServer) GetUser(_ context.Context, _ *pb.GetUserRequest) (*pb.User, error) {
	return &pb.User{}, nil
}

func (m *mockIdentityServer) ListUsers(_ context.Context, _ *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	return nil, nil
}

func Test_Blurb_lifecycle(t *testing.T) {
	s := NewMessagingServer(&mockIdentityServer{})

	first, err := s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/rumble",
				Content: &pb.Blurb_Text{Text: "woof"},
			},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	delete, err := s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/ekko",
				Content: &pb.Blurb_Text{Text: "bark"},
			},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	_, err = s.DeleteBlurb(
		context.Background(),
		&pb.DeleteBlurbRequest{Name: delete.Name})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
	}

	created, err := s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/musubi",
				Content: &pb.Blurb_Text{Text: "meow"},
			},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	got, err := s.GetBlurb(
		context.Background(),
		&pb.GetBlurbRequest{Name: created.GetName()})
	if err != nil {
		t.Errorf("Get: unexpected err %+v", err)
	}
	if !proto.Equal(created, got) {
		t.Error("Expected to get created blurb.")
	}

	got.Content = &pb.Blurb_Text{Text: "purrr"}
	_, err = s.UpdateBlurb(
		context.Background(),
		&pb.UpdateBlurbRequest{Blurb: got, UpdateMask: nil})
	if err != nil {
		t.Errorf("Update: unexpected err %+v", err)
	}

	updated, err := s.GetBlurb(
		context.Background(),
		&pb.GetBlurbRequest{Name: got.GetName()})
	if err != nil {
		t.Errorf("Get: unexpected err %+v", err)
	}
	// Cannot use proto.Equal here because the update time is changed on updates.
	if updated.GetName() != got.GetName() ||
		updated.GetUser() != got.GetUser() ||
		updated.GetText() != got.GetText() ||
		!proto.Equal(updated.GetCreateTime(), got.GetCreateTime()) ||
		proto.Equal(updated.GetUpdateTime(), got.GetUpdateTime()) {
		t.Errorf("Expected to get updated blurb. Want: %+v, got %+v", got, updated)
	}

	r, err := s.ListBlurbs(
		context.Background(),
		&pb.ListBlurbsRequest{
			Parent:    "users/rumble/profile",
			PageSize:  1,
			PageToken: "",
		})
	if err != nil {
		t.Errorf("List: unexpected err %+v", err)
	}
	if len(r.GetBlurbs()) != 1 {
		t.Errorf("List want: page size %d, got %d", 1, len(r.GetBlurbs()))
	}
	if !proto.Equal(first, r.GetBlurbs()[0]) {
		t.Errorf("List want: first blurb %+v, got %+v", first, r.GetBlurbs()[0])
	}
	if r.GetNextPageToken() == "" {
		t.Error("List want: non empty next page token")
	}

	r, err = s.ListBlurbs(
		context.Background(),
		&pb.ListBlurbsRequest{
			Parent:    "users/rumble/profile",
			PageSize:  10,
			PageToken: r.GetNextPageToken(),
		})
	if err != nil {
		t.Errorf("List: unexpected err %+v", err)
	}
	if len(r.GetBlurbs()) != 1 {
		t.Errorf("List want: page size %d, got %d", 1, len(r.GetBlurbs()))
	}
	if !proto.Equal(updated, r.GetBlurbs()[0]) {
		t.Errorf("List want: updated blurb %+v, got %+v", updated, r.GetBlurbs()[0])
	}
	if r.GetNextPageToken() != "" {
		t.Error("List want: empty next page token")
	}
}

func Test_CreateBlurb_invalid(t *testing.T) {
	s := NewMessagingServer(&mockIdentityServer{})

	_, err := s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{Blurb: &pb.Blurb{}})
	status, _ := status.FromError(err)
	if status.Code() != codes.InvalidArgument {
		t.Errorf(
			"Create: Want error code %d got %d",
			codes.InvalidArgument,
			status.Code())
	}
}

func Test_CreateBlurb_parentNotFound(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())

	_, err := s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{Parent: "users/rumble/profile", Blurb: &pb.Blurb{}})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Create: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_GetBlurb_notFound(t *testing.T) {
	s := NewMessagingServer(&mockIdentityServer{})
	_, err := s.GetBlurb(
		context.Background(),
		&pb.GetBlurbRequest{Name: "users/rumble/profile/blurbs/1"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Get: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_GetBlurb_deleted(t *testing.T) {
	s := NewMessagingServer(&mockIdentityServer{})
	created, err := s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/rumble",
				Content: &pb.Blurb_Text{Text: "woof"},
			},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	_, err = s.DeleteBlurb(
		context.Background(),
		&pb.DeleteBlurbRequest{Name: created.GetName()})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
	}

	_, err = s.GetBlurb(
		context.Background(),
		&pb.GetBlurbRequest{Name: created.GetName()})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Get deleted: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_UpdateBlurb_fieldmask(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())
	_, err := s.UpdateBlurb(
		context.Background(),
		&pb.UpdateBlurbRequest{
			Blurb:      nil,
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

func Test_UpdateBlurb_notFound(t *testing.T) {
	s := NewMessagingServer(&mockIdentityServer{})
	_, err := s.UpdateBlurb(
		context.Background(),
		&pb.UpdateBlurbRequest{
			Blurb: &pb.Blurb{
				User:    "users/rumble",
				Content: &pb.Blurb_Text{Text: "woof"},
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

func Test_UpdateBlurb_invalid(t *testing.T) {
	first := &pb.Blurb{
		User:    "users/rumble",
		Content: &pb.Blurb_Text{Text: "woof"},
	}
	second := &pb.Blurb{
		User:    "",
		Content: &pb.Blurb_Text{Text: "woof"},
	}
	s := NewMessagingServer(&mockIdentityServer{})
	created, err := s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{Blurb: first})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}
	second.Name = created.GetName()
	_, err = s.UpdateBlurb(
		context.Background(),
		&pb.UpdateBlurbRequest{Blurb: second, UpdateMask: nil})
	status, _ := status.FromError(err)
	if status.Code() != codes.InvalidArgument {
		t.Errorf(
			"Update: Want error code %d got %d",
			codes.InvalidArgument,
			status.Code())
	}
}

func Test_DeleteBlurb_notFound(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())
	_, err := s.DeleteBlurb(
		context.Background(),
		&pb.DeleteBlurbRequest{Name: "user/rumble/profile/blurbs/1"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Delete: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_ListBlurbs_invalidToken(t *testing.T) {
	s := messagingServerImpl{
		identityServer: &mockIdentityServer{},
		token:          &tokenGenerator{salt: "salt"},
		roomKeys:       map[string]int{},
		blurbKeys:      map[string]blurbIndex{},
		blurbs:         map[string][]blurbEntry{},
		parentUids:     map[string]*uniqID{},
	}

	tests := []string{
		"1", // Not base64 encoded
		base64.StdEncoding.EncodeToString([]byte("1")),        // No salt
		base64.StdEncoding.EncodeToString([]byte("saltblah")), // Invalid index
	}

	_, err := s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/rumble",
				Content: &pb.Blurb_Text{Text: "woof"},
			},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	for _, token := range tests {
		_, err := s.ListBlurbs(
			context.Background(),
			&pb.ListBlurbsRequest{
				Parent:    "users/rumble/profile",
				PageSize:  1,
				PageToken: token})
		status, _ := status.FromError(err)
		if status.Code() != codes.InvalidArgument {
			t.Errorf(
				"List: Want error code %d got %d",
				codes.InvalidArgument,
				status.Code())
		}
	}
}

func Test_ListBlurbs_parentNotFound(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())

	_, err := s.ListBlurbs(
		context.Background(),
		&pb.ListBlurbsRequest{Parent: "users/rumble/profile"})
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"List: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_ListBlurbs_noneCreated(t *testing.T) {
	s := NewMessagingServer(&mockIdentityServer{})

	resp, err := s.ListBlurbs(
		context.Background(),
		&pb.ListBlurbsRequest{Parent: "users/rumble/profile"})
	if err != nil {
		t.Errorf("List: unexpected err %+v", err)
	}
	if len(resp.GetBlurbs()) > 0 {
		t.Errorf("List none created: want empty list got %+v", resp)
	}
}

func Test_SearchBlurbs(t *testing.T) {
	s := NewMessagingServer(&mockIdentityServer{})

	req := &pb.SearchBlurbsRequest{
		Query:    "woof bark",
		Parent:   "users/rumble/profile",
		PageSize: 100,
	}

	op, err := s.SearchBlurbs(
		context.Background(),
		req)
	if err != nil {
		t.Errorf("SearchBlurbs: unexpected err %+v", err)
	}

	expName := "operations/google.showcase.v1alpha3.Messaging/SearchBlurbs/"
	if !strings.HasPrefix(op.Name, expName) {
		t.Errorf(
			"SearchBlurbs op.Name prefex want: '%s', got: %s'",
			expName,
			op.Name)
	}

	reqProto := &pb.SearchBlurbsRequest{}
	encodedBytes := strings.TrimPrefix(
		op.Name,
		expName)
	bytes, _ := base64.StdEncoding.DecodeString(encodedBytes)
	proto.Unmarshal(bytes, reqProto)
	if !proto.Equal(reqProto, req) {
		t.Errorf(
			"SearchBlurbs for %q expected unmarshalled %q, got %q",
			req,
			req,
			reqProto)
	}
}

func Test_SearchBlurbs_parentNotFound(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())

	req := &pb.SearchBlurbsRequest{
		Query:  "woof bark",
		Parent: "users/rumble/profile",
	}

	_, err := s.SearchBlurbs(
		context.Background(),
		req)
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"SearchBlurbs: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}
