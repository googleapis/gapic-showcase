// Copyright 2021 Google LLC
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

package resttools

import (
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ServerStreamer struct { // implements genprotopb.Echo_ExpandServer
	responses []string
	grpc.ServerStream

	initialization sync.Once
	marshaler      *protojson.MarshalOptions
}

func (streamer *ServerStreamer) Send(response proto.Message) error { // this could just take the interface for proto,
	// same as Marshal, and so we only have one instance in genrest
	streamer.initialization.Do(streamer.initialize)
	json, err := streamer.marshaler.Marshal(response)
	if err != nil {
		return fmt.Errorf("error json-encoding response: %s", err.Error())
	}
	streamer.responses = append(streamer.responses, string(json))
	return nil
}

func (streamer *ServerStreamer) ListJSON() string {
	return fmt.Sprintf("{\n%s\n}", strings.Join(streamer.responses, ",\n"))
}

func (streamer *ServerStreamer) initialize() {
	streamer.marshaler = ToJSON()
}
