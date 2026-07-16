// Copyright 2026 Google LLC
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

	pb "github.com/googleapis/gapic-showcase/server/genproto"
)

type resumableUploadServerImpl struct {
	pb.UnimplementedResumableUploadServiceServer
}

func NewResumableUploadServer() pb.ResumableUploadServiceServer {
	return &resumableUploadServerImpl{}
}

func (s *resumableUploadServerImpl) UploadMedia(ctx context.Context, req *pb.UploadMediaRequest) (*pb.UploadMediaResponse, error) {
	return &pb.UploadMediaResponse{}, nil
}
