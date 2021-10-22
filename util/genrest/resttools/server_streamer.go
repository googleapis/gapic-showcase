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
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ServerStreamer implements REST server streaming functionality that can be called by streaming
// RPCs to stream their responses over HTTP/JSON. The messages are encoded such that once th stream
// is finished, the total payload represents a properly formed JSON array of objects. In order to
// ensure this, users must call the End() method to terminate the stream properly, typically by
// using `defer`.
type ServerStreamer struct {
	output    http.ResponseWriter
	marshaler *protojson.MarshalOptions
	flusher   http.Flusher
	grpc.ServerStream
	started bool
}

// NewServerStreamer returns a ServerStreamer instance initialized to write to responseWriter. Users
// must call the End() method to terminate the stream properly, typically by using `defer`.
func NewServerStreamer(responseWriter http.ResponseWriter) (*ServerStreamer, error) {
	if responseWriter == nil {
		return nil, fmt.Errorf("error: responseWriter provided is nil")
	}

	flusher, ok := responseWriter.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("error: responseWriter provided does not implement http.Flusher")
	}

	streamer := &ServerStreamer{
		output:    responseWriter,
		flusher:   flusher,
		marshaler: ToJSON(),
	}

	return streamer, nil
}

// Send sends `response` over the REST stream by writing the message and then flushing the writer.
func (streamer *ServerStreamer) Send(response proto.Message) error {
	json, err := streamer.marshaler.Marshal(response)
	if err != nil {
		return fmt.Errorf("error json-encoding response: %s", err.Error())
	}

	var prefix []byte
	switch streamer.started {
	case false:
		prefix = []byte("[")
		streamer.started = true
	case true:
		prefix = []byte(",")
	}

	if _, err := streamer.output.Write(append(prefix, json...)); err != nil {
		return fmt.Errorf("error writing streamed json message: %s", err.Error())
	}

	streamer.flusher.Flush()

	return nil
}

// End terminates the REST stream by sending the trailing bytes (the closing bracket for the array).
func (streamer *ServerStreamer) End() error {
	if !streamer.started {
		return nil
	}

	if _, err := streamer.output.Write([]byte("]")); err != nil {
		return fmt.Errorf("error terminating json stream: %s", err.Error())
	}

	return nil
}

// Context is needed to satisfy grpc.ServerStream.
func (streamer *ServerStreamer) Context() context.Context {
	return context.Background()
}
