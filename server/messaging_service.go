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
	"fmt"
	"sync"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewMessagingServer() pb.MessagingServer {
	return &messagingServerImpl{
		token: NewTokenGenerator(),
		keys:  map[string]int{},
	}
}

type messagingServerImpl struct {
	uid   uniqID
	token TokenGenerator

	mu    sync.Mutex
	keys  map[string]int
	rooms []roomEntry
}

type roomEntry struct {
	room    *pb.Room
	deleted bool
}

// Creates a room.
func (s *messagingServerImpl) CreateRoom(ctx context.Context, in *pb.CreateRoomRequest) (*pb.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := in.GetRoom()
	err := validateRoom(r)
	if err != nil {
		return nil, err
	}

	// Validate Unique Fields.
	uniqName := func(x *pb.Room) bool {
		return (r.GetDisplayName() == x.GetDisplayName())
	}
	if s.anyRoom(uniqName) {
		return nil, status.Errorf(
			codes.AlreadyExists,
			"A room with display_name %s already exists.",
			r.GetDisplayName())
	}

	// Assign info.
	id := s.uid.next()
	name := fmt.Sprintf("rooms/%d", id)
	now := ptypes.TimestampNow()

	r.Name = name
	r.CreateTime = now
	r.UpdateTime = now

	// Insert.
	index := len(s.rooms)
	s.rooms = append(s.rooms, roomEntry{room: r})
	s.keys[name] = index

	return r, nil
}

// Retrieves the Room with the given resource name.
func (s *messagingServerImpl) GetRoom(ctx context.Context, in *pb.GetRoomRequest) (*pb.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := in.GetName()
	if i, ok := s.keys[name]; ok {
		entry := s.rooms[i]
		if !entry.deleted {
			return entry.room, nil
		}
	}

	return nil, status.Errorf(
		codes.NotFound, "A room with name %s not found.",
		name)
}

// Updates a room.
func (s *messagingServerImpl) UpdateRoom(ctx context.Context, in *pb.UpdateRoomRequest) (*pb.Room, error) {
	f := in.GetUpdateMask()
	r := in.GetRoom()
	if f != nil && len(f.GetPaths()) > 0 {
		return nil, status.Error(
			codes.Unimplemented,
			"Field masks are currently not supported.")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	err := validateRoom(r)
	if err != nil {
		return nil, err
	}
	i, ok := s.keys[r.GetName()]

	if !ok || s.rooms[i].deleted {
		return nil, status.Errorf(
			codes.NotFound,
			"A room with name %s not found.", r.GetName())
	}

	entry := s.rooms[i]
	// Validate Unique Fields.
	uniqName := func(x *pb.Room) bool {
		return x != entry.room && (r.GetDisplayName() == x.GetDisplayName())
	}
	if s.anyRoom(uniqName) {
		return nil, status.Errorf(
			codes.AlreadyExists,
			"A room with either display_name, %s already exists.",
			r.GetDisplayName())
	}

	// Update store.
	updated := &pb.Room{
		Name:        r.GetName(),
		DisplayName: r.GetDisplayName(),
		Description: r.GetDescription(),
		CreateTime:  entry.room.GetCreateTime(),
		UpdateTime:  ptypes.TimestampNow(),
	}
	s.rooms[i] = roomEntry{room: updated}
	return updated, nil
}

// Deletes a room and all of its blurbs.
func (s *messagingServerImpl) DeleteRoom(ctx context.Context, in *pb.DeleteRoomRequest) (*empty.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	i, ok := s.keys[in.GetName()]

	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"A room with name %s not found.", in.GetName())
	}

	entry := s.rooms[i]
	s.rooms[i] = roomEntry{room: entry.room, deleted: true}

	return &empty.Empty{}, nil
}

// Lists all chat rooms.
func (s *messagingServerImpl) ListRooms(ctx context.Context, in *pb.ListRoomsRequest) (*pb.ListRoomsResponse, error) {
	start, err := s.token.GetIndex(in.GetPageToken())
	if err != nil {
		return nil, err
	}

	offset := 0
	rooms := []*pb.Room{}
	for _, entry := range s.rooms[start:] {
		offset++
		if entry.deleted {
			continue
		}
		rooms = append(rooms, entry.room)
		if len(rooms) >= int(in.GetPageSize()) {
			break
		}
	}

	nextToken := ""
	if start+offset < len(s.rooms) {
		nextToken = s.token.ForIndex(start + offset)
	}

	return &pb.ListRoomsResponse{Rooms: rooms, NextPageToken: nextToken}, nil
}

func (s *messagingServerImpl) anyRoom(f func(*pb.Room) bool) bool {
	for _, entry := range s.rooms {
		if !entry.deleted && f(entry.room) {
			return true
		}
	}
	return false
}

func validateRoom(r *pb.Room) error {
	// Validate Required Fields.
	if r.GetDisplayName() == "" {
		return status.Errorf(
			codes.InvalidArgument,
			"The field `display_name` is required.")
	}
	return nil
}

// Creates a blurb. If the parent is a room, the blurb is understood to be a
// message in that room. If the parent is a profile, the blurb is understood
// to be a post on the profile.
func (s *messagingServerImpl) CreateBlurb(ctx context.Context, in *pb.CreateBlurbRequest) (*pb.Blurb, error) {
	return nil, nil
}

// Retrieves the Blurb with the given resource name.
func (s *messagingServerImpl) GetBlurb(ctx context.Context, in *pb.GetBlurbRequest) (*pb.Blurb, error) {
	return nil, nil
}

// Updates a blurb.
func (s *messagingServerImpl) UpdateBlurb(ctx context.Context, in *pb.UpdateBlurbRequest) (*pb.Blurb, error) {
	return nil, nil
}

// Deletes a blurb.
func (s *messagingServerImpl) DeleteBlurb(ctx context.Context, in *pb.DeleteBlurbRequest) (*empty.Empty, error) {
	return nil, nil
}

// Lists blurbs for a specific chat room or user profile depending on the
// parent resource name.
func (s *messagingServerImpl) ListBlurbs(ctx context.Context, in *pb.ListBlurbsRequest) (*pb.ListBlurbsResponse, error) {
	return nil, nil
}

// This method searches through all blurbs across all rooms and profiles
// for blurbs containing to words found in the query. Only posts that
// contain an exact match of a queried word will be returned.
func (s *messagingServerImpl) SearchBlurbs(ctx context.Context, in *pb.SearchBlurbsRequest) (*longrunning.Operation, error) {
	return nil, nil
}

// This returns a stream that emits the blurbs that are created for a
// particular chat room or user profile.
func (s *messagingServerImpl) StreamBlurbs(in *pb.StreamBlurbsRequest, stream pb.Messaging_StreamBlurbsServer) error {
	return nil
}

// This is a stream to create multiple blurbs. If an invalid blurb is
// requested to be created, the stream will close with an error.
func (s *messagingServerImpl) SendBlurbs(stream pb.Messaging_SendBlurbsServer) error {
	return nil
}

// This method starts a bidirectional stream that receives all blurbs that
// are being created after the stream has started and sends requests to create
// blurbs. If an invalid blurb is requested to be created, the stream will
// close with an error.
func (s *messagingServerImpl) Connect(stream pb.Messaging_ConnectServer) error {
	return nil
}
