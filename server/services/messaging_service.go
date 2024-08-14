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
	"strconv"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	errdetails "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewMessagingServer returns an instance of a messaging server.
func NewMessagingServer(identityServer ReadOnlyIdentityServer) MessagingServer {
	return &messagingServerImpl{
		identityServer: identityServer,
		nowF:           time.Now,
		token:          server.NewTokenGenerator(),
		roomKeys:       map[string]int{},
		blurbKeys:      map[string]blurbIndex{},
		blurbs:         map[string][]blurbEntry{},
		parentUids:     map[string]*server.UniqID{},
		observers:      map[string]map[string]blurbObserver{},
	}
}

// MessagingServer provides an interface which is the implementation of the
// MessagingServer proto and as well as a method for filtering the blurbs.
type MessagingServer interface {
	FilteredListBlurbs(context.Context, *pb.ListBlurbsRequest, func(*pb.Blurb) bool) (*pb.ListBlurbsResponse, error)

	pb.MessagingServer
}

type messagingServerImpl struct {
	nowF           func() time.Time
	token          server.TokenGenerator
	identityServer ReadOnlyIdentityServer

	roomUID  server.UniqID
	roomMu   sync.Mutex
	roomKeys map[string]int
	rooms    []roomEntry

	blurbMu    sync.Mutex
	blurbKeys  map[string]blurbIndex
	blurbs     map[string][]blurbEntry
	parentUids map[string]*server.UniqID

	obsMu     sync.Mutex
	obsUID    server.UniqID
	observers map[string]map[string]blurbObserver
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

type blurbObserver interface {
	OnCreate(b *pb.Blurb)
	OnUpdate(b *pb.Blurb)
	OnDelete(b *pb.Blurb)
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
	id := s.roomUID.Next()
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
	mask := in.GetUpdateMask()
	r := in.GetRoom()

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
	existing := s.rooms[i].room

	// Validate Unique Fields.
	uniqName := func(x *pb.Room) bool {
		return x != existing && (r.GetDisplayName() == x.GetDisplayName())
	}
	if s.anyRoom(uniqName) {
		return nil, status.Errorf(
			codes.AlreadyExists,
			"A room with either display_name, %s already exists.",
			r.GetDisplayName())
	}

	// Update store.
	updated := proto.Clone(existing).(*pb.Room)
	applyFieldMask(r.ProtoReflect(), updated.ProtoReflect(), mask.GetPaths())
	updated.CreateTime = existing.GetCreateTime()
	updated.UpdateTime = ptypes.TimestampNow()

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
		puid = &server.UniqID{}
		s.parentUids[parent] = puid
	}

	id := puid.Next()
	name := fmt.Sprintf("%s/blurbs/%d", parent, id)
	now := ptypes.TimestampNow()
	switch legacyID := b.LegacyId.(type) {
	case *pb.Blurb_LegacyRoomId:
		name = fmt.Sprintf("%s/blurbs/legacy/%s.%d", parent, legacyID.LegacyRoomId, id)
	case *pb.Blurb_LegacyUserId:
		name = fmt.Sprintf("%s/blurbs/legacy/%s~%d", parent, legacyID.LegacyUserId, id)
	}
	b.Name = name
	b.CreateTime = now
	b.UpdateTime = now

	// Insert.
	index := len(parentBs)
	s.blurbs[parent] = append(parentBs, blurbEntry{blurb: b})
	s.blurbKeys[name] = blurbIndex{row: parent, col: index}

	// Call observers.
	s.obsMu.Lock()
	defer s.obsMu.Unlock()
	if parentObservers, ok := s.observers[parent]; ok {
		for _, o := range parentObservers {
			o.OnCreate(b)
		}
	}

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
	s.blurbMu.Lock()
	defer s.blurbMu.Unlock()

	mask := in.GetUpdateMask()
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
	existing := s.blurbs[i.row][i.col].blurb
	updated := proto.Clone(existing).(*pb.Blurb)
	applyFieldMask(b.ProtoReflect(), updated.ProtoReflect(), mask.GetPaths())
	updated.CreateTime = existing.GetCreateTime()
	updated.UpdateTime = ptypes.TimestampNow()
	s.blurbs[i.row][i.col] = blurbEntry{blurb: updated}

	// Call observers.
	s.obsMu.Lock()
	defer s.obsMu.Unlock()
	if parentObservers, ok := s.observers[i.row]; ok {
		for _, o := range parentObservers {
			o.OnUpdate(updated)
		}
	}

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

	// Call observers.
	s.obsMu.Lock()
	defer s.obsMu.Unlock()
	if parentObservers, ok := s.observers[i.row]; ok {
		for _, o := range parentObservers {
			o.OnDelete(entry.blurb)
		}
	}

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
func (s *messagingServerImpl) SearchBlurbs(ctx context.Context, in *pb.SearchBlurbsRequest) (*longrunningpb.Operation, error) {
	if err := s.validateParent(in.GetParent()); err != nil {
		return nil, err
	}
	reqBytes, _ := proto.Marshal(in)

	name := fmt.Sprintf(
		"operations/google.showcase.v1beta1.Messaging/SearchBlurbs/%s",
		base64.StdEncoding.EncodeToString(reqBytes))
	// TODO(landrito) Add randomization to the retry delay.
	meta, _ := ptypes.MarshalAny(
		&pb.SearchBlurbsMetadata{
			RetryInfo: &errdetails.RetryInfo{
				RetryDelay: ptypes.DurationProto(time.Duration(1) * time.Second),
			},
		},
	)
	return &longrunningpb.Operation{Name: name, Done: false, Metadata: meta}, nil
}

// This returns a stream that emits the blurbs that are created for a
// particular chat room or user profile.
func (s *messagingServerImpl) StreamBlurbs(in *pb.StreamBlurbsRequest, stream pb.Messaging_StreamBlurbsServer) error {
	parent := in.GetName()
	if err := s.validateParent(parent); err != nil {
		return err
	}

	expireTime, err := ptypes.Timestamp(in.GetExpireTime())
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	observer := &streamBlurbsObserver{
		stream: stream.(BlurbsOutStream),
		mu:     sync.Mutex{},
	}
	name := s.registerObserver(parent, observer)
	defer s.removeObserver(parent, name)
	for {
		if s.nowF().After(expireTime) {
			break
		}
		if observer.getError() != nil {
			return observer.getError()
		}
		if err := s.validateParent(parent); err != nil {
			return err
		}
	}
	return nil
}

// This is a stream to create multiple blurbs. If an invalid blurb is
// requested to be created, the stream will close with an error.
func (s *messagingServerImpl) SendBlurbs(stream pb.Messaging_SendBlurbsServer) error {
	names := []string{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(
				&pb.SendBlurbsResponse{Names: names})
		}
		if err != nil {
			return withCreatedNames(err, names)
		}
		parent := req.GetParent()
		if err := s.validateParent(parent); err != nil {
			return withCreatedNames(err, names)
		}
		blurb, err := s.CreateBlurb(
			context.Background(),
			&pb.CreateBlurbRequest{Parent: parent, Blurb: req.GetBlurb()})
		if err != nil {
			return withCreatedNames(err, names)
		}
		names = append(names, blurb.GetName())
	}
}

func withCreatedNames(err error, names []string) error {
	s, _ := status.FromError(err)
	spb := s.Proto()

	details, err := ptypes.MarshalAny(&pb.SendBlurbsResponse{Names: names})
	if err == nil {
		spb.Details = append(spb.Details, details)
	}

	return status.ErrorProto(spb)
}

// This method starts a bidirectional stream that receives all blurbs that
// are being created after the stream has started and sends requests to create
// blurbs. If an invalid blurb is requested to be created, the stream will
// close with an error.
func (s *messagingServerImpl) Connect(stream pb.Messaging_ConnectServer) error {
	configured := false
	parent := ""
	var observer *streamBlurbsObserver

	for {
		req, err := stream.Recv()

		// Handle stream errors.
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Setup Configuration
		if !configured && req != nil {
			if req.GetConfig() == nil {
				return status.Error(
					codes.InvalidArgument,
					"The first request to Connect, must contain a config field")
			}

			parent = req.GetConfig().GetParent()
			if err := s.validateParent(parent); err != nil {
				return err
			}
			observer = &streamBlurbsObserver{
				stream: stream.(BlurbsOutStream),
				mu:     sync.Mutex{},
			}
			configured = true
			name := s.registerObserver(parent, observer)
			defer s.removeObserver(parent, name)
			continue
		}
		// Check if there was a send error.
		if err := observer.getError(); err != nil {
			return err
		}

		// Check if the parent still exists.
		if err := s.validateParent(parent); err != nil {
			return err
		}

		// Create the blurb
		if req == nil || req.GetBlurb() == nil {
			continue
		}
		_, err = s.CreateBlurb(
			context.Background(),
			&pb.CreateBlurbRequest{Parent: parent, Blurb: req.GetBlurb()})
		if err != nil {
			return err
		}
	}
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

// Observer
type streamBlurbsObserver struct {
	// The stream to send the blurbs.
	stream BlurbsOutStream

	// Blurb parent to observe.
	parent string

	// Holds an error that occurred when sending a blurb.
	mu  sync.Mutex
	err error
}

// BlurbsOutStream is the common interface of the connect and streamblurbs streams.
// This interface allows the observer to handle both streams.
type BlurbsOutStream interface {
	Send(*pb.StreamBlurbsResponse) error
}

func (o *streamBlurbsObserver) getError() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.err
}

func (o *streamBlurbsObserver) updateError(err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.err == nil {
		o.err = err
	}
}

func (o *streamBlurbsObserver) OnCreate(b *pb.Blurb) {
	o.updateError(
		o.stream.Send(
			&pb.StreamBlurbsResponse{
				Blurb:  b,
				Action: pb.StreamBlurbsResponse_CREATE,
			}))
}

func (o *streamBlurbsObserver) OnUpdate(b *pb.Blurb) {
	o.updateError(
		o.stream.Send(
			&pb.StreamBlurbsResponse{
				Blurb:  b,
				Action: pb.StreamBlurbsResponse_UPDATE,
			}))
}

// Do nothing for now. Need to change the
func (o *streamBlurbsObserver) OnDelete(b *pb.Blurb) {
	o.updateError(
		o.stream.Send(
			&pb.StreamBlurbsResponse{
				Blurb:  b,
				Action: pb.StreamBlurbsResponse_DELETE,
			}))
}

func (s *messagingServerImpl) registerObserver(parent string, o blurbObserver) string {
	s.obsMu.Lock()
	defer s.obsMu.Unlock()
	name := strconv.FormatInt(s.obsUID.Next(), 10)
	if _, ok := s.observers[parent]; !ok {
		s.observers[parent] = map[string]blurbObserver{}
	}
	s.observers[parent][name] = o
	return name
}

func (s *messagingServerImpl) hasObservers(parent string) bool {
	s.obsMu.Lock()
	defer s.obsMu.Unlock()
	if os, ok := s.observers[parent]; ok && len(os) > 0 {
		return true
	}
	return false
}

func (s *messagingServerImpl) removeObserver(parent string, name string) {
	s.obsMu.Lock()
	defer s.obsMu.Unlock()
	delete(s.observers[parent], name)
	if len(s.observers[parent]) <= 0 {
		delete(s.observers, parent)
	}
}
