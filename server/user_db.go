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

func NewUserDb() UserDb {
	return &userDb{
		uid:   &uniqID{},
		token: NewTokenGenerator(),
		mu:    sync.Mutex{},
		keys:  map[string]dbIndex{},
		users: []userEntry{},
	}
}

type UserDb interface {
	Create(u *pb.User) (*pb.User, error)
	Get(name string) (*pb.User, error)
	Update(u *pb.User, f *field_mask.FieldMask) (*pb.User, error)
	Delete(name string) error
	List(page_size int32, token string) (*pb.ListUsersResponse, error)
}

type userEntry struct {
	user    *pb.User
	deleted bool
}

type dbIndex int

type userDb struct {
	uid   *uniqID
	token TokenGenerator

	mu    sync.Mutex
	keys  map[string]dbIndex
	users []userEntry
}

func (db *userDb) Create(u *pb.User) (*pb.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	err := validate(u)
	if err != nil {
		return nil, err
	}

	// Validate Unique Fields.
	uniqNameEmail := func(x *pb.User) bool {
		return (u.GetDisplayName() == x.GetDisplayName() || u.GetEmail() == x.GetEmail())
	}
	if db.anyUser(uniqNameEmail) {
		return nil, status.Errorf(
			codes.AlreadyExists,
			"A user with display_name %s or email %s already exists.",
			u.GetDisplayName(),
			u.GetEmail())
	}

	// Assign info.
	id := db.uid.next()
	name := fmt.Sprintf("users/%d", id)
	now := ptypes.TimestampNow()

	u.Name = name
	u.CreateTime = now
	u.UpdateTime = now

	// Insert.
	index := dbIndex(len(db.users))
	db.users = append(db.users, userEntry{user: u, deleted: false})
	db.keys[name] = index

	return u, nil
}

func (db *userDb) Get(s string) (*pb.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if v, ok := db.keys[s]; ok {
		entry := db.users[int(v)]
		if !entry.deleted {
			return entry.user, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "A user with name %s not found.", s)
}

func (db *userDb) Update(u *pb.User, f *field_mask.FieldMask) (*pb.User, error) {
	if f != nil && len(f.GetPaths()) > 0 {
		return nil, status.Error(
			codes.Unimplemented,
			"Field masks are currently not supported.")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	err := validate(u)
	if err != nil {
		return nil, err
	}
	i, ok := db.keys[u.GetName()]

	if !ok || db.users[int(i)].deleted {
		return nil, status.Errorf(
			codes.NotFound,
			"A user with name %s not found.", u.GetName())
	}

	entry := db.users[int(i)]
	// Validate Unique Fields.
	uniqNameEmail := func(x *pb.User) bool {
		return x != entry.user && (u.GetDisplayName() == x.GetDisplayName() || u.GetEmail() == x.GetEmail())
	}
	if db.anyUser(uniqNameEmail) {
		return nil, status.Errorf(
			codes.AlreadyExists,
			"A user with either display_name, %s, or email %s already exists.",
			u.GetDisplayName(),
			u.GetEmail())
	}

	// Update store.
	updated := &pb.User{
		Name:        u.GetName(),
		DisplayName: u.GetDisplayName(),
		Email:       u.GetEmail(),
		CreateTime:  entry.user.GetCreateTime(),
		UpdateTime:  ptypes.TimestampNow(),
	}
	db.users[int(i)] = userEntry{user: updated, deleted: false}
	return u, nil
}

func (db *userDb) Delete(s string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	i, ok := db.keys[s]

	if !ok {
		return status.Errorf(
			codes.NotFound,
			"A user with name %s not found.", s)
	}

	entry := db.users[int(i)]
	db.users[int(i)] = userEntry{user: entry.user, deleted: true}

	return nil
}

func (db *userDb) List(pageSize int32, pageToken string) (*pb.ListUsersResponse, error) {
	start, err := db.token.GetIndex(pageToken)
	if err != nil {
		return nil, err
	}

	offset := 0
	users := []*pb.User{}
	for _, entry := range db.users[start:] {
		offset++
		if entry.deleted {
			continue
		}
		users = append(users, entry.user)
		if len(users) >= int(pageSize) {
			break
		}
	}

	nextToken := ""
	if start+offset < len(db.users) {
		nextToken = db.token.ForIndex(start + offset)
	}

	return &pb.ListUsersResponse{Users: users, NextPageToken: nextToken}, nil
}

func (db *userDb) anyUser(f func(*pb.User) bool) bool {
	for _, entry := range db.users {
		if !entry.deleted && f(entry.user) {
			return true
		}
	}
	return false
}

func validate(u *pb.User) error {
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
	return nil
}
