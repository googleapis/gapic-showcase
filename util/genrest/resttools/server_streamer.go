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
	"io"
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
	output    io.Writer
	marshaler *protojson.MarshalOptions
	flusher   http.Flusher
	grpc.ServerStream
	started bool
}

// NewServerStreamer returns a ServerStreamer instance initialized to write to responseWriter. Users
// must call the End() method to terminate the stream properly, typically by using `defer`.
func NewServerStreamer(responseWriter io.Writer) (*ServerStreamer, error) {
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

// Send sends `message` over the REST stream as a chunk to be immediately be sent over the wire.
func (streamer *ServerStreamer) Send(message proto.Message) error {
	json, err := streamer.marshaler.Marshal(message)
	if err != nil {
		return fmt.Errorf("error json-encoding message: %s", err.Error())
	}

	return streamer.sendJSONArrayChunk(json)
}

// End terminates the REST stream by sending the trailing bytes (the closing bracket for the array).
func (streamer *ServerStreamer) End() error {
	if !streamer.started {
		return nil
	}

	if _, err := streamer.output.Write([]byte("]")); err != nil {
		return fmt.Errorf("error terminating json stream: %s", err.Error())
	}

	streamer.flusher.Flush()

	return nil
}

// Context is needed to satisfy grpc.ServerStream.
func (streamer *ServerStreamer) Context() context.Context {
	return context.Background()
}

// sendJSONArrayChunk sends chunk over the REST stream by writing chunk and then flushing the
// writer. Each chunk is assumed to be a JSON object delimited by curly braces, inasmuch as chunks
// are separated by commas, the first chunk is preceded by an opening square bracket, and End() will
// send the final closing square bracket. Empty chunks do not get written.
func (streamer *ServerStreamer) sendJSONArrayChunk(chunk []byte) error {
	if len(chunk) == 0 {
		return nil
	}

	var prefix []byte
	switch streamer.started {
	case false:
		prefix = []byte("[")
		streamer.started = true
	case true:
		prefix = []byte(",")
	}

	if _, err := streamer.output.Write(append(prefix, chunk...)); err != nil {
		return fmt.Errorf("error writing streamed json message: %s", err.Error())
	}

	// Flush() causes the chunk to be sent immediately, and chunked transfer encoding to be set
	// in the HTTP headers before the first chunk.
	// cf. https://stackoverflow.com/a/30603654
	streamer.flusher.Flush()

	return nil
}
