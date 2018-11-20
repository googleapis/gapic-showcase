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

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewBlurbDb() BlurbDb {
	return &blurbDb{
		token:      NewTokenGenerator(),
		blurbMu:    sync.Mutex{},
		keys:       map[string]blurbIndex{},
		blurbs:     map[string][]blurbEntry{},
		parentUids: map[string]*uniqID{},
	}
}

type BlurbDb interface {
	Create(parent string, b *pb.Blurb) (*pb.Blurb, error)
	Get(name string) (*pb.Blurb, error)
	Update(b *pb.Blurb, f *field_mask.FieldMask) (*pb.Blurb, error)
	Delete(name string) error
	List(r *ListBlurbsDbRequest) (*pb.ListBlurbsResponse, error)
}

type ReadOnlyBlurbDb interface {
	Get(name string) (*pb.Blurb, error)
	List(r *ListBlurbsDbRequest) (*pb.ListBlurbsResponse, error)
}

type ListBlurbsDbRequest struct {
	Parent    string
	PageSize  int32
	PageToken string
	Filter    func(*pb.Blurb) bool
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

type blurbDb struct {
	// Generates Page Tokens
	token TokenGenerator

	blurbMu sync.Mutex
	// Mapping from blurb name to index.
	keys map[string]blurbIndex
	// Mapping from parent name to list of blurbs.
	blurbs map[string][]blurbEntry
	// Mapping from parent name to its context.
	parentUids map[string]*uniqID
}

func (db *blurbDb) Create(parent string, b *pb.Blurb) (*pb.Blurb, error) {
	db.blurbMu.Lock()
	defer db.blurbMu.Unlock()

	if err := validateBlurb(b); err != nil {
		return nil, err
	}

	// Assign info.
	parentBs, ok := db.blurbs[parent]
	if !ok {
		parentBs = []blurbEntry{}
	}
	puid, ok := db.parentUids[parent]
	if !ok {
		puid = &uniqID{}
		db.parentUids[parent] = puid
	}

	id := puid.next()
	name := fmt.Sprintf("%s/blurbs/%d", parent, id)
	now := ptypes.TimestampNow()

	b.Name = name
	b.CreateTime = now
	b.UpdateTime = now

	// Insert.
	index := len(parentBs)
	db.blurbs[parent] = append(parentBs, blurbEntry{blurb: b, deleted: false})
	db.keys[name] = blurbIndex{row: parent, col: index}

	return b, nil
}

func (db *blurbDb) Get(s string) (*pb.Blurb, error) {
	db.blurbMu.Lock()
	defer db.blurbMu.Unlock()

	if i, ok := db.keys[s]; ok {
		entry := db.blurbs[i.row][i.col]
		if !entry.deleted {
			return entry.blurb, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "A blurb with name %s not found.", s)
}

func (db *blurbDb) Update(b *pb.Blurb, f *field_mask.FieldMask) (*pb.Blurb, error) {
	if f != nil && len(f.GetPaths()) > 0 {
		return nil, status.Error(
			codes.Unimplemented,
			"Field masks are currently not supported.")
	}

	db.blurbMu.Lock()
	defer db.blurbMu.Unlock()

	i, ok := db.keys[b.GetName()]
	if !ok || db.blurbs[i.row][i.col].deleted {
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
	db.blurbs[i.row][i.col] = blurbEntry{blurb: updated, deleted: false}

	return updated, nil
}

func (db *blurbDb) Delete(s string) error {
	db.blurbMu.Lock()
	defer db.blurbMu.Unlock()

	i, ok := db.keys[s]

	if !ok {
		return status.Errorf(
			codes.NotFound,
			"A blurb with name %s not found.", s)
	}

	entry := db.blurbs[i.row][i.col]
	db.blurbs[i.row][i.col] = blurbEntry{blurb: entry.blurb, deleted: true}

	return nil
}

func (db *blurbDb) List(r *ListBlurbsDbRequest) (*pb.ListBlurbsResponse, error) {
	bs, ok := db.blurbs[r.Parent]
	if !ok {
		return &pb.ListBlurbsResponse{}, nil
	}

	start, err := db.token.GetIndex(r.PageToken)
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
		if r.Filter == nil || r.Filter(entry.blurb) {
			blurbs = append(blurbs, entry.blurb)
		}
		if len(blurbs) >= int(r.PageSize) {
			break
		}
	}

	nextToken := ""
	if start+offset < len(db.blurbs[r.Parent]) {
		nextToken = db.token.ForIndex(start + offset)
	}

	return &pb.ListBlurbsResponse{Blurbs: blurbs, NextPageToken: nextToken}, nil
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
