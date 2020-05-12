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
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/googleapis/gapic-showcase/server"
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
		token:    server.TokenGeneratorWithSalt("salt"),
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
	// create a Blurb with legacy_user_id
	second, err := s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/rumble",
				Content: &pb.Blurb_Text{Text: "non-slash resource test."},
				LegacyUserId: "legacy_rumble",
			},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}
	// get the second Blurb from database and verify the legacy_user_id, then delete it.
	gotSecond, err := s.GetBlurb(
		context.Background(),
		&pb.GetBlurbRequest{Name: second.GetName()})
	if err != nil {
		t.Errorf("Get: unexpected err %+v", err)
	}
	if !proto.Equal(second, gotSecond) {
		t.Error("Expected to get created blurb.")
	}
	_, err = s.DeleteBlurb(
		context.Background(),
		&pb.DeleteBlurbRequest{Name: second.GetName()})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
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
		&pb.DeleteBlurbRequest{Name: delete.GetName()})
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
		token:          server.TokenGeneratorWithSalt("salt"),
		roomKeys:       map[string]int{},
		blurbKeys:      map[string]blurbIndex{},
		blurbs:         map[string][]blurbEntry{},
		parentUids:     map[string]*server.UniqID{},
		observers:      map[string]map[string]blurbObserver{},
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

	expName := "operations/google.showcase.v1beta1.Messaging/SearchBlurbs/"
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

type mockStreamBlurbsStream struct {
	mu    sync.Mutex
	resps []*pb.StreamBlurbsResponse
	pb.Messaging_StreamBlurbsServer
}

func (m *mockStreamBlurbsStream) Send(resp *pb.StreamBlurbsResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resps = append(m.resps, resp)
	return nil
}

func TestStreamBlurbs_lifecycle(t *testing.T) {
	// We specify the now function so we can control when the stream ends.
	s := &messagingServerImpl{
		identityServer: &mockIdentityServer{},
		token:          server.NewTokenGenerator(),
		roomKeys:       map[string]int{},
		blurbKeys:      map[string]blurbIndex{},
		blurbs:         map[string][]blurbEntry{},
		parentUids:     map[string]*server.UniqID{},
		observers:      map[string]map[string]blurbObserver{},
		nowF: func() time.Time {
			return time.Unix(int64(0), int64(0))
		},
	}

	// Set up the mock stream.
	m := &mockStreamBlurbsStream{resps: []*pb.StreamBlurbsResponse{}}

	// Make the end time some time after the time the now function will return.
	endTime, err := ptypes.TimestampProto(time.Unix(int64(1), int64(0)))
	if err != nil {
		t.Errorf("TimestampProto: unexpected err %+v", err)
	}

	// The parent to stream
	p := "users/rumble/profile"

	// Start the stream.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go (func() {
		err := s.StreamBlurbs(
			&pb.StreamBlurbsRequest{
				Name:       p,
				ExpireTime: endTime,
			},
			m,
		)
		if err != nil {
			t.Errorf("StreamBlurbs: unexpected err %+v", err)
		}
		wg.Done()
	})()

	// Wait for the stream request to propogate the observer to the database.
	for {
		if s.hasObservers(p) {
			break
		}
	}

	// Create a blurb.
	created, err := s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{
			Parent: p,
			Blurb: &pb.Blurb{
				User:    "users/musubi",
				Content: &pb.Blurb_Text{Text: "meow"},
			},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	// Check that the stream sent the blurb info.
	m.mu.Lock()
	streamResp := m.resps[len(m.resps)-1]
	m.mu.Unlock()
	if streamResp.GetAction() != pb.StreamBlurbsResponse_CREATE {
		t.Errorf(
			"StreamBlurbs: want blurb with action %s, got %s",
			pb.StreamBlurbsResponse_Action_name[int32(pb.StreamBlurbsResponse_CREATE)],
			pb.StreamBlurbsResponse_Action_name[int32(streamResp.GetAction())])
	}
	if !proto.Equal(streamResp.GetBlurb(), created) {
		t.Errorf(
			"StreamBlurbs: want created blurb %+v, got %+v",
			created,
			streamResp.GetBlurb())
	}

	// Update the blurb.
	created.Content = &pb.Blurb_Text{Text: "purrr"}
	updated, err := s.UpdateBlurb(
		context.Background(),
		&pb.UpdateBlurbRequest{Blurb: created, UpdateMask: nil})
	if err != nil {
		t.Errorf("Update: unexpected err %+v", err)
	}

	// Check that the stream sent the blurb info.
	m.mu.Lock()
	streamResp = m.resps[len(m.resps)-1]
	m.mu.Unlock()
	if streamResp.GetAction() != pb.StreamBlurbsResponse_UPDATE {
		t.Errorf(
			"StreamBlurbs: want blurb with action %s, got %s",
			pb.StreamBlurbsResponse_Action_name[int32(pb.StreamBlurbsResponse_UPDATE)],
			pb.StreamBlurbsResponse_Action_name[int32(streamResp.GetAction())])
	}
	if !proto.Equal(streamResp.GetBlurb(), updated) {
		t.Errorf(
			"StreamBlurbs: want updated blurb %+v, got %+v",
			updated,
			streamResp.GetBlurb())
	}

	// Delete the blurb.
	_, err = s.DeleteBlurb(
		context.Background(),
		&pb.DeleteBlurbRequest{Name: updated.Name})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
	}

	// Check that the stream sent the blurb info.
	m.mu.Lock()
	streamResp = m.resps[len(m.resps)-1]
	m.mu.Unlock()
	if streamResp.GetAction() != pb.StreamBlurbsResponse_DELETE {
		t.Errorf(
			"StreamBlurbs: want blurb with action %s, got %s",
			pb.StreamBlurbsResponse_Action_name[int32(pb.StreamBlurbsResponse_DELETE)],
			pb.StreamBlurbsResponse_Action_name[int32(streamResp.GetAction())])
	}
	if !proto.Equal(streamResp.GetBlurb(), updated) {
		t.Errorf(
			"StreamBlurbs: want deleted blurb %+v, got %+v",
			updated,
			streamResp.GetBlurb())
	}

	// Set the now function to return a time after the expire time to close the
	// stream.
	s.nowF = func() time.Time {
		return time.Unix(int64(2), int64(0))
	}

	// Wait til the stream is closed.
	wg.Wait()
}

func Test_StreamBlurbs_parentNotFound(t *testing.T) {
	s := NewMessagingServer(NewIdentityServer())

	err := s.StreamBlurbs(
		&pb.StreamBlurbsRequest{Name: "users/rumble/profile"},
		nil)
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Create: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

func Test_StreamBlurbs_invalidTimestamp(t *testing.T) {
	is := NewIdentityServer()
	first, err := is.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{
			User: &pb.User{DisplayName: "rumbledog", Email: "rumble@google.com"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}
	s := NewMessagingServer(is)

	maxValidSeconds := int64(253402300800)
	err = s.StreamBlurbs(
		&pb.StreamBlurbsRequest{Name: first.GetName(), ExpireTime: &timestamp.Timestamp{Seconds: maxValidSeconds + 1}},
		nil)
	status, _ := status.FromError(err)
	if status.Code() != codes.InvalidArgument {
		t.Errorf(
			"Create: Want error code %d got %d",
			codes.InvalidArgument,
			status.Code())
	}
}

type errorStreamBlurbsStream struct {
	pb.Messaging_StreamBlurbsServer
}

func (m *errorStreamBlurbsStream) Send(_ *pb.StreamBlurbsResponse) error {
	return status.Error(codes.Unknown, "Error")
}

func Test_StreamBlurbs_sendError(t *testing.T) {
	// We specify the now function so we can control when the stream ends.
	s := &messagingServerImpl{

		identityServer: &mockIdentityServer{},
		token:          server.NewTokenGenerator(),
		roomKeys:       map[string]int{},
		blurbKeys:      map[string]blurbIndex{},
		blurbs:         map[string][]blurbEntry{},
		parentUids:     map[string]*server.UniqID{},
		observers:      map[string]map[string]blurbObserver{},
		nowF: func() time.Time {
			return time.Unix(int64(0), int64(0))
		},
	}

	// Make the end time some time after the time the now function will return.
	endTime, err := ptypes.TimestampProto(time.Unix(int64(1), int64(0)))
	if err != nil {
		t.Errorf("TimestampProto: unexpected err %+v", err)
	}

	// The parent to stream.
	p := "users/rumble/profile"

	// Start Stream.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go (func() {
		err = s.StreamBlurbs(
			&pb.StreamBlurbsRequest{Name: p, ExpireTime: endTime},
			&errorStreamBlurbsStream{})
		status, _ := status.FromError(err)
		if status.Code() != codes.Unknown {
			t.Errorf(
				"Create: Want error code %d got %d",
				codes.Unknown,
				status.Code())
		}
		wg.Done()
	})()

	// Wait for the stream request to propogate the observer to the database.
	for {
		if s.hasObservers(p) {
			break
		}
	}

	// Create a blurb to trigger observer.
	_, err = s.CreateBlurb(
		context.Background(),
		&pb.CreateBlurbRequest{
			Parent: p,
			Blurb: &pb.Blurb{
				User:    "users/musubi",
				Content: &pb.Blurb_Text{Text: "meow"},
			},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	wg.Wait()
}

type nilStreamBlurbsStream struct {
	pb.Messaging_StreamBlurbsServer
}

func (m *nilStreamBlurbsStream) Send(resp *pb.StreamBlurbsResponse) error {
	return nil
}

func Test_StreamBlurbs_parentNotFoundLater(t *testing.T) {
	// Setup Identity server to validate parent against.
	is := NewIdentityServer()
	first, err := is.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{
			User: &pb.User{DisplayName: "rumbledog", Email: "rumble@google.com"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	// We specify the now function so we can control when the stream ends.
	s := &messagingServerImpl{
		identityServer: is,

		token:      server.NewTokenGenerator(),
		roomKeys:   map[string]int{},
		blurbKeys:  map[string]blurbIndex{},
		blurbs:     map[string][]blurbEntry{},
		parentUids: map[string]*server.UniqID{},
		observers:  map[string]map[string]blurbObserver{},
		nowF: func() time.Time {
			return time.Unix(int64(0), int64(0))
		},
	}

	// Make the end time some time after the time the now function will return.
	endTime, err := ptypes.TimestampProto(time.Unix(int64(1), int64(0)))
	if err != nil {
		t.Errorf("TimestampProto: unexpected err %+v", err)
	}

	parent := fmt.Sprintf("%s/profile", first.GetName())

	// Start Stream.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go (func() {
		err := s.StreamBlurbs(
			&pb.StreamBlurbsRequest{
				Name:       parent,
				ExpireTime: endTime,
			},
			&nilStreamBlurbsStream{})
		status, _ := status.FromError(err)
		if status.Code() != codes.NotFound {
			t.Errorf(
				"StreamBlurbs: Want error code %d got %d",
				codes.NotFound,
				status.Code())
		}
		wg.Done()
	})()

	for {
		if s.hasObservers(parent) {
			break
		}
	}

	// Delete the user so that the parent is invalid.
	is.DeleteUser(
		context.Background(),
		&pb.DeleteUserRequest{Name: first.GetName()})

	// Wait til the stream closes.
	wg.Wait()
}

type mockSendBlurbsStream struct {
	reqs []*pb.CreateBlurbRequest
	resp *pb.SendBlurbsResponse
	t    *testing.T

	mu   sync.Mutex
	next int
	pb.Messaging_SendBlurbsServer
}

func (m *mockSendBlurbsStream) SendAndClose(r *pb.SendBlurbsResponse) error {
	m.resp = r
	return nil
}

func (m *mockSendBlurbsStream) Recv() (*pb.CreateBlurbRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.next < len(m.reqs) {
		cur := m.next
		m.next++
		return m.reqs[cur], nil
	}
	return nil, io.EOF
}

func TestSendBlurbs(t *testing.T) {
	reqs := []*pb.CreateBlurbRequest{
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/rumble",
				Content: &pb.Blurb_Text{Text: "woof"},
			},
		},
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/musubi",
				Content: &pb.Blurb_Text{Text: "meow"},
			},
		},
		&pb.CreateBlurbRequest{
			Parent: "users/musubi/profile",
			Blurb: &pb.Blurb{
				User:    "users/ekko",
				Content: &pb.Blurb_Text{Text: "bark"},
			},
		},
	}
	m := &mockSendBlurbsStream{
		reqs: reqs,
		t:    t,
	}
	s := NewMessagingServer(&mockIdentityServer{})

	err := s.SendBlurbs(m)
	if err != nil {
		t.Errorf("SendBlurbs: unexpected err %+v", err)
	}
	for i, name := range m.resp.GetNames() {
		got, err := s.GetBlurb(
			context.Background(),
			&pb.GetBlurbRequest{Name: name})
		if err != nil {
			t.Errorf("Get: unexpected err %+v", err)
		}
		if reqs[i].GetBlurb().GetUser() != got.GetUser() ||
			reqs[i].GetBlurb().GetText() != got.GetText() {
			t.Errorf(
				"Expected to get created blurb. Want: %+v, got %+v",
				reqs[i].GetBlurb(),
				got)
		}
	}
}

type errorSendBlurbsStream struct {
	reqs []*pb.CreateBlurbRequest
	resp *pb.SendBlurbsResponse
	t    *testing.T

	mu   sync.Mutex
	next int
	pb.Messaging_SendBlurbsServer
}

func (m *errorSendBlurbsStream) SendAndClose(r *pb.SendBlurbsResponse) error {
	m.resp = r
	return nil
}

func (m *errorSendBlurbsStream) Recv() (*pb.CreateBlurbRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.next < len(m.reqs) {
		cur := m.next
		m.next++
		return m.reqs[cur], nil
	}
	return nil, status.Error(codes.Unknown, "Error")
}

func TestSendBlurbs_error(t *testing.T) {
	reqs := []*pb.CreateBlurbRequest{
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/rumble",
				Content: &pb.Blurb_Text{Text: "woof"},
			},
		},
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/musubi",
				Content: &pb.Blurb_Text{Text: "meow"},
			},
		},
		&pb.CreateBlurbRequest{
			Parent: "users/musubi/profile",
			Blurb: &pb.Blurb{
				User:    "users/ekko",
				Content: &pb.Blurb_Text{Text: "bark"},
			},
		},
	}
	m := &errorSendBlurbsStream{
		reqs: reqs,
		t:    t,
	}
	s := NewMessagingServer(&mockIdentityServer{})

	err := s.SendBlurbs(m)
	if err == nil {
		t.Error("SendBlurbs: expected err")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Errorf("SendBlurbs: expected err to be status %+v", err)
	}
	details := st.Proto().GetDetails()
	if len(details) > 1 {
		t.Errorf("SendBlurbs: expected err details to be of length 1")
	}
	resp := &pb.SendBlurbsResponse{}
	ptypes.UnmarshalAny(details[0], resp)

	for i, name := range resp.GetNames() {
		got, err := s.GetBlurb(
			context.Background(),
			&pb.GetBlurbRequest{Name: name})
		if err != nil {
			t.Errorf("Get: unexpected err %+v", err)
		}
		if reqs[i].GetBlurb().GetName() != got.GetName() ||
			reqs[i].GetBlurb().GetUser() != got.GetUser() ||
			reqs[i].GetBlurb().GetText() != got.GetText() {
			t.Errorf("Expected to get updated blurb. Want: %+v, got %+v", reqs, got)
		}
	}
}

func TestSendBlurbs_invalidParent(t *testing.T) {
	// Setup Identity server to validate parent against.
	is := NewIdentityServer()
	first, err := is.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{
			User: &pb.User{DisplayName: "rumbledog", Email: "rumble@google.com"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}

	parent := fmt.Sprintf("%s/profile", first.GetName())
	reqs := []*pb.CreateBlurbRequest{
		&pb.CreateBlurbRequest{
			Parent: parent,
			Blurb: &pb.Blurb{
				User:    "users/rumble",
				Content: &pb.Blurb_Text{Text: "woof"},
			},
		},
		&pb.CreateBlurbRequest{
			Parent: parent,
			Blurb: &pb.Blurb{
				User:    "users/musubi",
				Content: &pb.Blurb_Text{Text: "meow"},
			},
		},
		&pb.CreateBlurbRequest{
			Parent: "does/not/exist",
			Blurb: &pb.Blurb{
				User:    "users/ekko",
				Content: &pb.Blurb_Text{Text: "bark"},
			},
		},
	}
	m := &mockSendBlurbsStream{
		reqs: reqs,
		t:    t,
	}
	s := NewMessagingServer(is)

	err = s.SendBlurbs(m)
	if err == nil {
		t.Error("SendBlurbs: expected err")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Errorf("SendBlurbs: expected err to be status %+v", err)
	}
	details := st.Proto().GetDetails()
	if len(details) > 1 {
		t.Errorf("SendBlurbs: expected err details to be of length 1")
	}
	resp := &pb.SendBlurbsResponse{}
	ptypes.UnmarshalAny(details[0], resp)

	for i, name := range resp.GetNames() {
		got, err := s.GetBlurb(
			context.Background(),
			&pb.GetBlurbRequest{Name: name})
		if err != nil {
			t.Errorf("Get: unexpected err %+v", err)
		}
		if reqs[i].GetBlurb().GetName() != got.GetName() ||
			reqs[i].GetBlurb().GetUser() != got.GetUser() ||
			reqs[i].GetBlurb().GetText() != got.GetText() {
			t.Errorf("Expected to get updated blurb. Want: %+v, got %+v", reqs, got)
		}
	}
}

func TestSendBlurbs_invalidBlurb(t *testing.T) {
	reqs := []*pb.CreateBlurbRequest{
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/rumble",
				Content: &pb.Blurb_Text{Text: "woof"},
			},
		},
		&pb.CreateBlurbRequest{
			Parent: "users/rumble/profile",
			Blurb: &pb.Blurb{
				User:    "users/musubi",
				Content: &pb.Blurb_Text{Text: "meow"},
			},
		},
		&pb.CreateBlurbRequest{
			Parent: "users/musubi/profile",
			Blurb: &pb.Blurb{
				Content: &pb.Blurb_Text{Text: "bark"},
			},
		},
	}
	m := &mockSendBlurbsStream{
		reqs: reqs,
		t:    t,
	}
	s := NewMessagingServer(&mockIdentityServer{})

	err := s.SendBlurbs(m)
	if err == nil {
		t.Error("SendBlurbs: expected err")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Errorf("SendBlurbs: expected err to be status %+v", err)
	}
	details := st.Proto().GetDetails()
	if len(details) > 1 {
		t.Errorf("SendBlurbs: expected err details to be of length 1")
	}
	resp := &pb.SendBlurbsResponse{}
	ptypes.UnmarshalAny(details[0], resp)

	for i, name := range resp.GetNames() {
		got, err := s.GetBlurb(
			context.Background(),
			&pb.GetBlurbRequest{Name: name})
		if err != nil {
			t.Errorf("Get: unexpected err %+v", err)
		}
		if reqs[i].GetBlurb().GetName() != got.GetName() ||
			reqs[i].GetBlurb().GetUser() != got.GetUser() ||
			reqs[i].GetBlurb().GetText() != got.GetText() {
			t.Errorf("Expected to get updated blurb. Want: %+v, got %+v", reqs, got)
		}
	}
}

type mockConnectStream struct {
	reqs []*pb.ConnectRequest
	t    *testing.T
	stop bool

	respMu sync.Mutex
	resps  []*pb.StreamBlurbsResponse

	nextMu sync.Mutex
	next   int

	pb.Messaging_ConnectServer
}

func (m *mockConnectStream) Recv() (*pb.ConnectRequest, error) {
	if m.next < len(m.reqs) {
		req := m.reqs[m.next]
		m.next++
		return req, nil
	}
	if m.stop {
		return nil, io.EOF
	}
	return nil, nil
}

func (m *mockConnectStream) Send(r *pb.StreamBlurbsResponse) error {
	m.respMu.Lock()
	defer m.respMu.Unlock()
	m.resps = append(m.resps, r)
	return nil
}

func TestConnect(t *testing.T) {
	reqs := []*pb.ConnectRequest{
		&pb.ConnectRequest{
			Request: &pb.ConnectRequest_Config{
				Config: &pb.ConnectRequest_ConnectConfig{
					Parent: "users/rumble/profile",
				},
			},
		},
		&pb.ConnectRequest{
			Request: &pb.ConnectRequest_Blurb{
				Blurb: &pb.Blurb{
					User:    "users/rumble",
					Content: &pb.Blurb_Text{Text: "woof"},
				},
			},
		},
	}
	m := &mockConnectStream{
		reqs:  reqs,
		t:     t,
		resps: []*pb.StreamBlurbsResponse{},
	}
	s := NewMessagingServer(&mockIdentityServer{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go (func() {
		err := s.Connect(m)
		if err != nil {
			t.Fatalf("Connect: unexpected err %+v", err)
		}
		wg.Done()
	})()

	for {
		if len(m.resps) > 0 {
			break
		}
	}

	// Check that the stream sent the blurb info.
	m.respMu.Lock()
	streamResp := m.resps[len(m.resps)-1]
	m.respMu.Unlock()
	if streamResp.GetAction() != pb.StreamBlurbsResponse_CREATE {
		t.Errorf(
			"StreamBlurbs: want blurb with action %s, got %s",
			pb.StreamBlurbsResponse_Action_name[int32(pb.StreamBlurbsResponse_CREATE)],
			pb.StreamBlurbsResponse_Action_name[int32(streamResp.GetAction())])
	}
	if reqs[1].GetBlurb().GetName() != streamResp.GetBlurb().GetName() ||
		reqs[1].GetBlurb().GetUser() != streamResp.GetBlurb().GetUser() ||
		reqs[1].GetBlurb().GetText() != streamResp.GetBlurb().GetText() {
		t.Errorf(
			"Expected to get created blurb. Want: %+v, got %+v",
			reqs[1].GetBlurb(),
			streamResp.GetBlurb())
	}

	// Create a blurb.
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

	// Check that the stream sent the blurb info.
	m.respMu.Lock()
	streamResp = m.resps[len(m.resps)-1]
	m.respMu.Unlock()
	if streamResp.GetAction() != pb.StreamBlurbsResponse_CREATE {
		t.Errorf(
			"StreamBlurbs: want blurb with action %s, got %s",
			pb.StreamBlurbsResponse_Action_name[int32(pb.StreamBlurbsResponse_CREATE)],
			pb.StreamBlurbsResponse_Action_name[int32(streamResp.GetAction())])
	}
	if !proto.Equal(streamResp.GetBlurb(), created) {
		t.Errorf(
			"StreamBlurbs: want created blurb %+v, got %+v",
			created,
			streamResp.GetBlurb())
	}

	// Update the blurb.
	created.Content = &pb.Blurb_Text{Text: "purrr"}
	updated, err := s.UpdateBlurb(
		context.Background(),
		&pb.UpdateBlurbRequest{Blurb: created, UpdateMask: nil})
	if err != nil {
		t.Errorf("Update: unexpected err %+v", err)
	}

	// Check that the stream sent the blurb info.
	m.respMu.Lock()
	streamResp = m.resps[len(m.resps)-1]
	m.respMu.Unlock()
	if streamResp.GetAction() != pb.StreamBlurbsResponse_UPDATE {
		t.Errorf(
			"StreamBlurbs: want blurb with action %s, got %s",
			pb.StreamBlurbsResponse_Action_name[int32(pb.StreamBlurbsResponse_UPDATE)],
			pb.StreamBlurbsResponse_Action_name[int32(streamResp.GetAction())])
	}
	if !proto.Equal(streamResp.GetBlurb(), updated) {
		t.Errorf(
			"StreamBlurbs: want updated blurb %+v, got %+v",
			updated,
			streamResp.GetBlurb())
	}

	// Delete the blurb.
	_, err = s.DeleteBlurb(
		context.Background(),
		&pb.DeleteBlurbRequest{Name: updated.Name})
	if err != nil {
		t.Errorf("Delete: unexpected err %+v", err)
	}

	// Check that the stream sent the blurb info.
	m.respMu.Lock()
	streamResp = m.resps[len(m.resps)-1]
	m.respMu.Unlock()
	if streamResp.GetAction() != pb.StreamBlurbsResponse_DELETE {
		t.Errorf(
			"StreamBlurbs: want blurb with action %s, got %s",
			pb.StreamBlurbsResponse_Action_name[int32(pb.StreamBlurbsResponse_DELETE)],
			pb.StreamBlurbsResponse_Action_name[int32(streamResp.GetAction())])
	}
	if !proto.Equal(streamResp.GetBlurb(), updated) {
		t.Errorf(
			"StreamBlurbs: want deleted blurb %+v, got %+v",
			updated,
			streamResp.GetBlurb())
	}

	m.stop = true
	wg.Wait()
}

type errorConnectStream struct {
	pb.Messaging_ConnectServer
}

func (m *errorConnectStream) Recv() (*pb.ConnectRequest, error) {
	return nil, status.Error(codes.Unknown, "Error")
}

func TestConnect_error(t *testing.T) {
	m := &errorConnectStream{}
	s := NewMessagingServer(&mockIdentityServer{})

	err := s.Connect(m)
	status, _ := status.FromError(err)
	if status.Code() != codes.Unknown {
		t.Errorf(
			"Connect: Want error code %d got %d",
			codes.Unknown,
			status.Code())
	}
}

func TestConnect_notConfigured(t *testing.T) {
	reqs := []*pb.ConnectRequest{
		&pb.ConnectRequest{
			Request: &pb.ConnectRequest_Blurb{
				Blurb: &pb.Blurb{
					User:    "users/rumble",
					Content: &pb.Blurb_Text{Text: "woof"},
				},
			},
		},
	}
	m := &mockConnectStream{
		reqs:  reqs,
		t:     t,
		resps: []*pb.StreamBlurbsResponse{},
	}
	s := NewMessagingServer(&mockIdentityServer{})

	err := s.Connect(m)
	if err == nil {
		t.Fatalf("Connect: expected err")
	}
	status, _ := status.FromError(err)
	if status.Code() != codes.InvalidArgument {
		t.Errorf(
			"Connect: Want error code %d got %d",
			codes.InvalidArgument,
			status.Code())
	}
}

func TestConnect_parentNotFound(t *testing.T) {
	reqs := []*pb.ConnectRequest{
		&pb.ConnectRequest{
			Request: &pb.ConnectRequest_Config{
				Config: &pb.ConnectRequest_ConnectConfig{
					Parent: "users/rumble/profile",
				},
			},
		},
	}
	m := &mockConnectStream{
		reqs:  reqs,
		t:     t,
		resps: []*pb.StreamBlurbsResponse{},
	}
	s := NewMessagingServer(NewIdentityServer())

	err := s.Connect(m)
	if err == nil {
		t.Fatalf("Connect: expected err")
	}
	status, _ := status.FromError(err)
	if status.Code() != codes.NotFound {
		t.Errorf(
			"Connect: Want error code %d got %d",
			codes.NotFound,
			status.Code())
	}
}

type sendErrorConnectStream struct {
	reqs []*pb.ConnectRequest
	t    *testing.T
	stop bool

	respMu sync.Mutex
	resps  []*pb.StreamBlurbsResponse

	nextMu sync.Mutex
	next   int

	pb.Messaging_ConnectServer
}

func (m *sendErrorConnectStream) Recv() (*pb.ConnectRequest, error) {
	if m.next < len(m.reqs) {
		req := m.reqs[m.next]
		m.next++
		return req, nil
	}
	if m.stop {
		return nil, io.EOF
	}
	return nil, nil
}

func (m *sendErrorConnectStream) Send(r *pb.StreamBlurbsResponse) error {
	return status.Error(codes.Unknown, "Error")
}

func TestConnect_sendError(t *testing.T) {
	reqs := []*pb.ConnectRequest{
		&pb.ConnectRequest{
			Request: &pb.ConnectRequest_Config{
				Config: &pb.ConnectRequest_ConnectConfig{
					Parent: "users/rumble/profile",
				},
			},
		},
		&pb.ConnectRequest{
			Request: &pb.ConnectRequest_Blurb{
				Blurb: &pb.Blurb{
					User:    "users/rumble",
					Content: &pb.Blurb_Text{Text: "woof"},
				},
			},
		},
	}
	m := &sendErrorConnectStream{
		reqs:  reqs,
		t:     t,
		resps: []*pb.StreamBlurbsResponse{},
	}
	s := NewMessagingServer(&mockIdentityServer{})

	err := s.Connect(m)
	if err == nil {
		t.Fatalf("Connect: expected err")
	}
	status, _ := status.FromError(err)
	if status.Code() != codes.Unknown {
		t.Errorf(
			"Connect: Want error code %d got %d",
			codes.Unknown,
			status.Code())
	}
}

func TestConnect_parentNotFoundLater(t *testing.T) {
	// Setup Identity server to validate parent against.
	is := NewIdentityServer()
	first, err := is.CreateUser(
		context.Background(),
		&pb.CreateUserRequest{
			User: &pb.User{DisplayName: "rumbledog", Email: "rumble@google.com"},
		})
	if err != nil {
		t.Errorf("Create: unexpected err %+v", err)
	}
	parent := fmt.Sprintf("%s/profile", first.GetName())
	reqs := []*pb.ConnectRequest{
		&pb.ConnectRequest{
			Request: &pb.ConnectRequest_Config{
				Config: &pb.ConnectRequest_ConnectConfig{
					Parent: parent,
				},
			},
		},
		&pb.ConnectRequest{
			Request: &pb.ConnectRequest_Blurb{
				Blurb: &pb.Blurb{
					User:    "users/rumble",
					Content: &pb.Blurb_Text{Text: "woof"},
				},
			},
		},
	}
	m := &mockConnectStream{
		reqs:  reqs,
		t:     t,
		resps: []*pb.StreamBlurbsResponse{},
	}

	// We specify the now function so we can control when the stream ends.
	s := &messagingServerImpl{
		identityServer: is,

		token:      server.NewTokenGenerator(),
		roomKeys:   map[string]int{},
		blurbKeys:  map[string]blurbIndex{},
		blurbs:     map[string][]blurbEntry{},
		parentUids: map[string]*server.UniqID{},
		observers:  map[string]map[string]blurbObserver{},
		nowF:       time.Now,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go (func() {
		err := s.Connect(m)
		if err == nil {
			t.Fatalf("Connect: expected err")
		}
		status, _ := status.FromError(err)
		if status.Code() != codes.NotFound {
			t.Errorf(
				"Connect: Want error code %d got %d",
				codes.NotFound,
				status.Code())
		}
		wg.Done()
	})()

	for {
		if s.hasObservers(parent) {
			break
		}
	}

	// Delete the user so that the parent is invalid.
	is.DeleteUser(
		context.Background(),
		&pb.DeleteUserRequest{Name: first.GetName()})

	// Wait til the stream closes.
	wg.Wait()
}

func Test_Connect_creationFailure(t *testing.T) {
	reqs := []*pb.ConnectRequest{
		&pb.ConnectRequest{
			Request: &pb.ConnectRequest_Config{
				Config: &pb.ConnectRequest_ConnectConfig{
					Parent: "users/rumble/profile",
				},
			},
		},
		&pb.ConnectRequest{
			Request: &pb.ConnectRequest_Blurb{
				Blurb: &pb.Blurb{
					Content: &pb.Blurb_Text{Text: "woof"},
				},
			},
		},
	}
	m := &mockConnectStream{
		reqs:  reqs,
		t:     t,
		resps: []*pb.StreamBlurbsResponse{},
	}
	s := NewMessagingServer(&mockIdentityServer{})

	err := s.Connect(m)
	if err == nil {
		t.Fatalf("Connect: expected err")
	}
	status, _ := status.FromError(err)
	if status.Code() != codes.InvalidArgument {
		t.Errorf(
			"Connect: Want error code %d got %d",
			codes.InvalidArgument,
			status.Code())
	}
}
