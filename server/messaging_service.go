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
	"google.golang.org/genproto/googleapis/longrunning"
)

func NewMessagingServer() pb.MessagingServer {
	return &messagingServerImpl{
		roomDb: NewRoomDb(),
	}
}

type messagingServerImpl struct {
	roomDb RoomDb
}

// Creates a room.
func (s *messagingServerImpl) CreateRoom(ctx context.Context, in *pb.CreateRoomRequest) (*pb.Room, error) {
	return s.roomDb.Create(in.GetRoom())
}

// Retrieves the Room with the given resource name.
func (s *messagingServerImpl) GetRoom(ctx context.Context, in *pb.GetRoomRequest) (*pb.Room, error) {
	return s.roomDb.Get(in.GetName())
}

// Updates a room.
func (s *messagingServerImpl) UpdateRoom(ctx context.Context, in *pb.UpdateRoomRequest) (*pb.Room, error) {
	return s.roomDb.Update(in.GetRoom(), in.GetUpdateMask())
}

// Deletes a room and all of its blurbs.
func (s *messagingServerImpl) DeleteRoom(ctx context.Context, in *pb.DeleteRoomRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.roomDb.Delete(in.GetName())
}

// Lists all chat rooms.
func (s *messagingServerImpl) ListRooms(ctx context.Context, in *pb.ListRoomsRequest) (*pb.ListRoomsResponse, error) {
	return s.roomDb.List(in.GetPageSize(), in.GetPageToken())
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
