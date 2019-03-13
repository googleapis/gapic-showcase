// Copyright 2019 Google LLC
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

package spec

import (
	"github.com/googleapis/gapic-showcase/server"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	v1 "github.com/googleapis/gapic-showcase/server/spec/v1"
)

// ShowcaseTests returns the set of tests for a specific version of a session.
func ShowcaseTests(sessionName string, version pb.Session_Version) []server.Test {
	switch version {
	case pb.Session_V1_0:
		return v1_0Tests(sessionName)
	case pb.Session_V1_LATEST:
		return v1LatestTests(sessionName)
	default:
		return []server.Test{}
	}
}

func v1_0Tests(sessionName string) []server.Test {
	return []server.Test{
		v1.NewUnaryTest(sessionName),
	}
}

func v1LatestTests(sessionName string) []server.Test {
	return []server.Test{
		v1.NewUnaryTest(sessionName),
	}
}
