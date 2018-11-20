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
	"fmt"
	"sync"

	"github.com/golang/protobuf/ptypes"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRoomDb() RoomDb {
	return &roomDb{
		uid:   &uniqID{},
		token: NewTokenGenerator(),

		mu:    sync.Mutex{},
		keys:  map[string]int{},
		rooms: []roomEntry{},
	}
}

type RoomDb interface {
	Create(u *pb.Room) (*pb.Room, error)
	Get(name string) (*pb.Room, error)
	Update(u *pb.Room, f *field_mask.FieldMask) (*pb.Room, error)
	Delete(name string) error
	List(page_size int32, token string) (*pb.ListRoomsResponse, error)
}

type ReadOnlyRoomDb interface {
	Get(name string) (*pb.Room, error)
}

type roomEntry struct {
	room    *pb.Room
	deleted bool
}

type roomDb struct {
	uid   *uniqID
	token TokenGenerator

	mu    sync.Mutex
	keys  map[string]int
	rooms []roomEntry
}

func (db *roomDb) Create(r *pb.Room) (*pb.Room, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	err := validateRoom(r)
	if err != nil {
		return nil, err
	}

	// Validate Unique Fields.
	uniqName := func(x *pb.Room) bool {
		return (r.GetDisplayName() == x.GetDisplayName())
	}
	if db.anyRoom(uniqName) {
		return nil, status.Errorf(
			codes.AlreadyExists,
			"A room with display_name %s already exists.",
			r.GetDisplayName())
	}

	// Assign info.
	id := db.uid.next()
	name := fmt.Sprintf("rooms/%d", id)
	now := ptypes.TimestampNow()

	r.Name = name
	r.CreateTime = now
	r.UpdateTime = now

	// Insert.
	index := len(db.rooms)
	db.rooms = append(db.rooms, roomEntry{room: r, deleted: false})
	db.keys[name] = index

	return r, nil
}

func (db *roomDb) Get(s string) (*pb.Room, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if i, ok := db.keys[s]; ok {
		entry := db.rooms[i]
		if !entry.deleted {
			return entry.room, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "A room with name %s not found.", s)
}

func (db *roomDb) Update(r *pb.Room, f *field_mask.FieldMask) (*pb.Room, error) {
	if f != nil && len(f.GetPaths()) > 0 {
		return nil, status.Error(
			codes.Unimplemented,
			"Field masks are currently not supported.")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	err := validateRoom(r)
	if err != nil {
		return nil, err
	}
	i, ok := db.keys[r.GetName()]

	if !ok || db.rooms[i].deleted {
		return nil, status.Errorf(
			codes.NotFound,
			"A room with name %s not found.", r.GetName())
	}

	entry := db.rooms[i]
	// Validate Unique Fields.
	uniqName := func(x *pb.Room) bool {
		return x != entry.room && (r.GetDisplayName() == x.GetDisplayName())
	}
	if db.anyRoom(uniqName) {
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
	db.rooms[i] = roomEntry{room: updated, deleted: false}
	return updated, nil
}

func (db *roomDb) Delete(s string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i, ok := db.keys[s]

	if !ok {
		return status.Errorf(
			codes.NotFound,
			"A room with name %s not found.", s)
	}

	entry := db.rooms[i]
	db.rooms[i] = roomEntry{room: entry.room, deleted: true}

	return nil
}

func (db *roomDb) List(pageSize int32, pageToken string) (*pb.ListRoomsResponse, error) {
	start, err := db.token.GetIndex(pageToken)
	if err != nil {
		return nil, err
	}

	offset := 0
	rooms := []*pb.Room{}
	for _, entry := range db.rooms[start:] {
		offset++
		if entry.deleted {
			continue
		}
		rooms = append(rooms, entry.room)
		if len(rooms) >= int(pageSize) {
			break
		}
	}

	nextToken := ""
	if start+offset < len(db.rooms) {
		nextToken = db.token.ForIndex(start + offset)
	}

	return &pb.ListRoomsResponse{Rooms: rooms, NextPageToken: nextToken}, nil
}

func (db *roomDb) anyRoom(f func(*pb.Room) bool) bool {
	for _, entry := range db.rooms {
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
