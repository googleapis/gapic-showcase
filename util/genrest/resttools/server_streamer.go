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
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const ServerStreamingChunkSize = 30

// ServerStreamer implements REST server streaming functionality that can be called by streaming
// RPCs to stream their responses over HTTP/JSON. The messages are encoded such that once the stream
// is finished, the total payload represents a properly formed JSON array of objects. In order to
// ensure this, users must call the End() method to terminate the stream properly, typically by
// using `defer`.
type ServerStreamer struct {
	buffer    bytes.Buffer // buffers the output to be chunked
	chunkSize int

	output    io.Writer    // receives the actual output, flushed after each chunk
	flusher   http.Flusher // flushes the output to create the chunks
	marshaler *protojson.MarshalOptions
	started   bool
	grpc.ServerStream
}

// NewServerStreamer returns a ServerStreamer instance initialized to write to responseWriter. Users
// must call the End() method to terminate the stream properly, typically by using `defer`. Note
// that responseWriter must also be an http.Flusher. If chunkSize is positive, messages written to
// this ServerStreamer are chunked-encoded into chunks of size chunkSize (except for the final
// chunk, which could be smaller). If chunkSize is zero, each message is written into a single
// chunk.
func NewServerStreamer(responseWriter io.Writer, chunkSize int) (*ServerStreamer, error) {
	if responseWriter == nil {
		return nil, fmt.Errorf("error: responseWriter provided is nil")
	}

	flusher, ok := responseWriter.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("error: responseWriter provided does not implement http.Flusher")
	}

	if chunkSize < 0 {
		return nil, fmt.Errorf("error: chunkSize must be non-negative")
	}

	streamer := &ServerStreamer{
		chunkSize: chunkSize,
		output:    responseWriter,
		flusher:   flusher,
		marshaler: ToJSON(),
	}

	return streamer, nil
}

// Send sends a `message` over the REST stream using chunked-encoding according to the constructor
// chunkSize parameter.
func (streamer *ServerStreamer) Send(message proto.Message) error {
	json, err := streamer.marshaler.Marshal(message)
	if err != nil {
		return fmt.Errorf("error json-encoding message: %s", err.Error())
	}

	return streamer.sendJSONArrayMessage(json)
}

// End terminates the REST stream by sending the trailing bytes (the closing bracket for the array).
func (streamer *ServerStreamer) End() error {
	if !streamer.started {
		return nil
	}

	return streamer.sendAsChunks([]byte("]"), true)
}

// Context is needed to satisfy grpc.ServerStream.
func (streamer *ServerStreamer) Context() context.Context {
	return context.Background()
}

// sendJSONArrayMessage buffers `message` and then flushes as appropriate so that the chunks written
// to streamer.output are of size streamer.chunkSize. Each `message` is assumed to be a JSON object
// delimited by curly braces, inasmuch as `message`s are separated by commas, the first `message` is
// preceded by an opening square bracket, and End() will send the final closing square
// bracket. Empty `messages` do not get written.
func (streamer *ServerStreamer) sendJSONArrayMessage(message []byte) error {
	if len(message) == 0 {
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

	return streamer.sendAsChunks(append(prefix, message...), false)
}

// sendAsChunks writes `content` to streamer.output and flushes streamer.output so that each flushed
// chunk is of size streamer.chunkSize. If streamer.chunkSize is 0, all of `content` (and any
// previous contents of streamer.buffer) is flushed as a single chunk. If forceFlush is true, this
// function ensures all bytes in `content` are flushed, even if that results in a trailing chunk of
// size less than streamer.chunkSize.
func (streamer *ServerStreamer) sendAsChunks(content []byte, forceFlush bool) error {
	// we're writing in two places below, so we define a local function for conciseness.
	writeOneChunk := func(data []byte) error {
		if _, err := streamer.output.Write(data); err != nil {
			return fmt.Errorf("error writing streamed http chunk: %s\n  chunk: %q", err.Error(), string(data))
		}

		// Flush() causes the chunk to be sent immediately, and chunked transfer encoding to
		// be set in the HTTP headers before the first chunk.
		// cf. https://stackoverflow.com/a/30603654
		streamer.flusher.Flush()
		return nil
	}

	if _, err := streamer.buffer.Write(content); err != nil {
		return err
	}

	if streamer.chunkSize > 0 {
		for streamer.buffer.Len() >= streamer.chunkSize {
			if err := writeOneChunk(streamer.buffer.Next(streamer.chunkSize)); err != nil {
				return err
			}
		}
		if streamer.buffer.Len() == 0 {
			// reuse memory when possible
			streamer.buffer.Reset()
		}

	}

	if streamer.chunkSize <= 0 || forceFlush {
		if err := writeOneChunk(streamer.buffer.Bytes()); err != nil {
			return err
		}
		streamer.buffer.Reset()
	}

	return nil
}
