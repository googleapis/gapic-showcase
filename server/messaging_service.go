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
	"fmt"
	"strings"
  "sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/genproto/googleapis/longrunning"
	errdetails "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewMessagingServer(identityServer ReadOnlyIdentityServer) MessagingServer {
	return &messagingServerImpl{
    identityServer: identityServer,
		token: NewTokenGenerator(),
		roomKeys:  map[string]int{},
  	blurbKeys: map[string]blurbIndex{},
  	blurbs: map[string][]blurbEntry{},
  	parentUids: map[string]*uniqID{},
	}
}

type MessagingServer interface {
  FilteredListBlurbs(context.Context, *pb.ListBlurbsRequest, func(*pb.Blurb) bool) (*pb.ListBlurbsResponse, error)

  pb.MessagingServer
}

type messagingServerImpl struct {
	uid   uniqID
	token TokenGenerator
  identityServer ReadOnlyIdentityServer

	roomMu    sync.Mutex
	roomKeys  map[string]int
	rooms []roomEntry

	blurbMu sync.Mutex
	blurbKeys map[string]blurbIndex
	blurbs map[string][]blurbEntry
	parentUids map[string]*uniqID
}

type roomEntry struct {
	room    *pb.Room
	deleted bool
}

type blurbIndex struct {
	// The parent of the blurb.
	row string
	// The index within the list of blurbs of a parent.
	col int
}

type blurbEntry struct {
	blurb   *pb.Blurb
	deleted bool
}

// Creates a room.
func (s *messagingServerImpl) CreateRoom(ctx context.Context, in *pb.CreateRoomRequest) (*pb.Room, error) {
	s.roomMu.Lock()
	defer s.roomMu.Unlock()

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
	s.roomKeys[name] = index

	return r, nil
}

// Retrieves the Room with the given resource name.
func (s *messagingServerImpl) GetRoom(ctx context.Context, in *pb.GetRoomRequest) (*pb.Room, error) {
	s.roomMu.Lock()
	defer s.roomMu.Unlock()

	name := in.GetName()
	if i, ok := s.roomKeys[name]; ok {
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

	s.roomMu.Lock()
	defer s.roomMu.Unlock()

	err := validateRoom(r)
	if err != nil {
		return nil, err
	}
	i, ok := s.roomKeys[r.GetName()]

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
	s.roomMu.Lock()
	defer s.roomMu.Unlock()

	i, ok := s.roomKeys[in.GetName()]

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
  parent := in.GetParent()
	if err := s.validateParent(parent); err != nil {
		return nil, err
	}

	s.blurbMu.Lock()
	defer s.blurbMu.Unlock()

  b := in.GetBlurb()
	if err := validateBlurb(b); err != nil {
		return nil, err
	}

	// Assign info.
	parentBs, ok := s.blurbs[parent]
	if !ok {
		parentBs = []blurbEntry{}
	}
	puid, ok := s.parentUids[parent]
	if !ok {
		puid = &uniqID{}
		s.parentUids[parent] = puid
	}

	id := puid.next()
	name := fmt.Sprintf("%s/blurbs/%d", parent, id)
	now := ptypes.TimestampNow()

	b.Name = name
	b.CreateTime = now
	b.UpdateTime = now

	// Insert.
	index := len(parentBs)
	s.blurbs[parent] = append(parentBs, blurbEntry{blurb: b})
	s.blurbKeys[name] = blurbIndex{row: parent, col: index}

	return b, nil
}

// Retrieves the Blurb with the given resource name.
func (s *messagingServerImpl) GetBlurb(ctx context.Context, in *pb.GetBlurbRequest) (*pb.Blurb, error) {
	s.blurbMu.Lock()
	defer s.blurbMu.Unlock()

	if i, ok := s.blurbKeys[in.GetName()]; ok {
		entry := s.blurbs[i.row][i.col]
		if !entry.deleted {
			return entry.blurb, nil
		}
	}

	return nil, status.Errorf(
    codes.NotFound,
    "A blurb with name %s not found.",
    in.GetName())
}

// Updates a blurb.
func (s *messagingServerImpl) UpdateBlurb(ctx context.Context, in *pb.UpdateBlurbRequest) (*pb.Blurb, error) {
	if in.GetUpdateMask() != nil && len(in.GetUpdateMask().GetPaths()) > 0 {
		return nil, status.Error(
			codes.Unimplemented,
			"Field masks are currently not supported.")
	}

	s.blurbMu.Lock()
	defer s.blurbMu.Unlock()

  b := in.GetBlurb()
	i, ok := s.blurbKeys[b.GetName()]
	if !ok || s.blurbs[i.row][i.col].deleted {
		return nil, status.Errorf(
			codes.NotFound,
			"A blurb with name %s not found.", b.GetName())
	}

	if err := validateBlurb(b); err != nil {
		return nil, err
	}
	// Update store.
	updated := proto.Clone(b).(*pb.Blurb)
	updated.UpdateTime = ptypes.TimestampNow()
	s.blurbs[i.row][i.col] = blurbEntry{blurb: updated}

	return updated, nil
}

// Deletes a blurb.
func (s *messagingServerImpl) DeleteBlurb(ctx context.Context, in *pb.DeleteBlurbRequest) (*empty.Empty, error) {
	s.blurbMu.Lock()
	defer s.blurbMu.Unlock()

	i, ok := s.blurbKeys[in.GetName()]

	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"A blurb with name %s not found.", in.GetName())
	}

	entry := s.blurbs[i.row][i.col]
	s.blurbs[i.row][i.col] = blurbEntry{blurb: entry.blurb, deleted: true}

	return &empty.Empty{}, nil
}

// Lists blurbs for a specific chat room or user profile depending on the
// parent resource name.
func (s *messagingServerImpl) ListBlurbs(ctx context.Context, in *pb.ListBlurbsRequest) (*pb.ListBlurbsResponse, error) {
  passFilter := func(_ *pb.Blurb) bool { return true }
  return s.FilteredListBlurbs(ctx, in, passFilter)
}

func (s *messagingServerImpl) FilteredListBlurbs(ctx context.Context, in *pb.ListBlurbsRequest, f func(*pb.Blurb) bool) (*pb.ListBlurbsResponse, error) {
	if err := s.validateParent(in.GetParent()); err != nil {
		return nil, err
	}

	bs, ok := s.blurbs[in.GetParent()]
	if !ok {
		return &pb.ListBlurbsResponse{}, nil
	}

	start, err := s.token.GetIndex(in.GetPageToken())
	if err != nil {
		return nil, err
	}

	offset := 0
	blurbs := []*pb.Blurb{}
	for _, entry := range bs[start:] {
		offset++
		if entry.deleted {
			continue
		}
    if f(entry.blurb) {
      blurbs = append(blurbs, entry.blurb)
    }
		if len(blurbs) >= int(in.GetPageSize()) {
			break
		}
	}

	nextToken := ""
	if start+offset < len(s.blurbs[in.GetParent()]) {
		nextToken = s.token.ForIndex(start + offset)
	}

	return &pb.ListBlurbsResponse{Blurbs: blurbs, NextPageToken: nextToken}, nil
}

// This method searches through all blurbs across all rooms and profiles
// for blurbs containing to words found in the query. Only posts that
// contain an exact match of a queried word will be returned.
func (s *messagingServerImpl) SearchBlurbs(ctx context.Context, in *pb.SearchBlurbsRequest) (*longrunning.Operation, error) {
	if err := s.validateParent(in.GetParent()); err != nil {
		return nil, err
	}
	reqBytes, _ := proto.Marshal(in)

	name := fmt.Sprintf(
		"operations/google.showcase.v1alpha3.Messaging/SearchBlurbs/%s",
		base64.StdEncoding.EncodeToString(reqBytes))
	// TODO(landrito) Add randomization to the retry delay.
	meta, _ := ptypes.MarshalAny(
		&pb.SearchBlurbsMetadata{
			RetryInfo: &errdetails.RetryInfo{
				RetryDelay: ptypes.DurationProto(time.Duration(1) * time.Second),
			},
		},
	)
	return &longrunning.Operation{Name: name, Done: false, Metadata: meta}, nil
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

func (s *messagingServerImpl) validateParent(p string) error {
	_, uErr := s.identityServer.GetUser(
		context.Background(),
		&pb.GetUserRequest{
			Name: strings.TrimSuffix(p, "/profile"),
		},
	)
	_, rErr := s.GetRoom(context.Background(), &pb.GetRoomRequest{Name: p})
	if uErr != nil && rErr != nil {
		return status.Errorf(codes.NotFound, "Parent %s not found.", p)
	}
	return nil
}

func validateBlurb(b *pb.Blurb) error {
	// Validate Required Fields.
	if b.GetUser() == "" {
		return status.Errorf(
			codes.InvalidArgument,
			"The field `user` is required.")
	}
	return nil
}
