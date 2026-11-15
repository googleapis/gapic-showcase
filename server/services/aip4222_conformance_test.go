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
	"bytes"
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	showcaseclient "github.com/googleapis/gapic-showcase/client"
	pb "github.com/googleapis/gapic-showcase/server/genproto"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

type aip4222ConformanceSuite struct {
	SchemaVersion          string                        `json:"schema_version"`
	Spec                   string                        `json:"spec"`
	Generator              string                        `json:"generator"`
	ExtractVectors         []aip4222ExtractVector        `json:"extract_vectors"`
	InvalidTemplateVectors []aip4222InvalidTemplateEntry `json:"invalid_template_vectors"`
}

type aip4222ExtractVector struct {
	ID                  string            `json:"id"`
	Category            string            `json:"category"`
	Description         string            `json:"description"`
	TheoremRef          string            `json:"theorem_ref"`
	Request             map[string]string `json:"request"`
	ExpectedAlphaHeader string            `json:"expected_alpha_header"`
	ExpectHeaderSent    bool              `json:"expect_header_sent"`
}

type aip4222InvalidTemplateEntry struct {
	ID           string `json:"id"`
	BugRef       string `json:"bug_ref"`
	Description  string `json:"description"`
	Field        string `json:"field"`
	PathTemplate string `json:"path_template"`
}

func loadAip4222Suite(t *testing.T) ([]byte, aip4222ConformanceSuite) {
	t.Helper()
	suitePath := filepath.Clean("../../schema/google/showcase/v1beta1/aip4222_routing_suite.json")
	raw, err := os.ReadFile(suitePath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", suitePath, err)
	}
	var suite aip4222ConformanceSuite
	if err := json.Unmarshal(raw, &suite); err != nil {
		t.Fatalf("failed to decode %s: %v", suitePath, err)
	}
	return raw, suite
}

func TestAip4222ConformanceSuiteSync(t *testing.T) {
	raw, suite := loadAip4222Suite(t)
	if suite.SchemaVersion != "1.0" {
		t.Errorf("schema_version = %q, want %q", suite.SchemaVersion, "1.0")
	}
	if suite.Spec != "AIP-4222" {
		t.Errorf("spec = %q, want %q", suite.Spec, "AIP-4222")
	}
	if len(suite.ExtractVectors) != 57 {
		t.Fatalf("expected 57 extract_vectors, got %d", len(suite.ExtractVectors))
	}
	if len(suite.InvalidTemplateVectors) != 20 {
		t.Fatalf("expected 20 invalid_template_vectors, got %d", len(suite.InvalidTemplateVectors))
	}

	seenIDs := make(map[string]bool, 77)
	for _, vec := range suite.ExtractVectors {
		if vec.ID == "" {
			t.Error("found extract_vector with empty id")
		}
		if seenIDs[vec.ID] {
			t.Errorf("duplicate vector id %q", vec.ID)
		}
		seenIDs[vec.ID] = true
	}
	for _, vec := range suite.InvalidTemplateVectors {
		if vec.ID == "" {
			t.Error("found invalid_template_vector with empty id")
		}
		if seenIDs[vec.ID] {
			t.Errorf("duplicate vector id %q", vec.ID)
		}
		seenIDs[vec.ID] = true
	}

	oracleBin := filepath.Clean("../../formal/aip4222/.lake/build/bin/aip4222-oracle")
	if _, err := os.Stat(oracleBin); err == nil {
		cmd := exec.CommandContext(t.Context(), oracleBin, "--emit-conformance")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("%s --emit-conformance failed: %v\nstderr: %s", oracleBin, err, stderr.String())
		}
		if !bytes.Equal(raw, stdout.Bytes()) {
			t.Errorf("checked-in aip4222_routing_suite.json drifted from aip4222-oracle --emit-conformance output")
		}
	}
}

// canonicalizeRoutingHeader sorts key=value pairs by routing key because
// client/echo_client.go builds x-goog-request-params by ranging over a Go map.
func canonicalizeRoutingHeader(raw string) string {
	if raw == "" {
		return ""
	}
	parts := strings.Split(raw, "&")
	slices.SortFunc(parts, func(a, b string) int {
		ka, _, _ := strings.Cut(a, "=")
		kb, _, _ := strings.Cut(b, "=")
		return strings.Compare(ka, kb)
	})
	return strings.Join(parts, "&")
}

func TestShowcaseEchoGrpcWireConformance(t *testing.T) {
	_, suite := loadAip4222Suite(t)

	var showcaseVectors []aip4222ExtractVector
	for _, vec := range suite.ExtractVectors {
		if vec.Category == "showcase_echo" {
			showcaseVectors = append(showcaseVectors, vec)
		}
	}
	if len(showcaseVectors) != 7 {
		t.Fatalf("expected 7 showcase_echo vectors, got %d", len(showcaseVectors))
	}

	var capturedHeaders []string
	interceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			capturedHeaders = nil
		} else {
			capturedHeaders = append([]string(nil), md.Get("x-goog-request-params")...)
		}
		return handler(ctx, req)
	}

	listener := bufconn.Listen(1 << 20)
	srv := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	pb.RegisterEchoServer(srv, NewEchoServer())
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- srv.Serve(listener)
	}()
	t.Cleanup(func() {
		srv.Stop()
		if err := <-serveErr; err != nil && !strings.Contains(err.Error(), "closed") {
			t.Errorf("grpc server exited with unexpected error: %v", err)
		}
	})

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient failed: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	echoClient, err := showcaseclient.NewEchoClient(t.Context(), option.WithGRPCConn(conn), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("NewEchoClient failed: %v", err)
	}
	t.Cleanup(func() { _ = echoClient.Close() })

	// Known AIP-4222 divergences in gapic-generator-go (client/echo_client.go:659-693).
	knownGeneratorDivergences := map[string]struct {
		wantBuggySent   bool
		wantBuggyHeader string
		bugSummary      string
	}{
		"showcase_echo_both_empty": {
			wantBuggySent:   true,
			wantBuggyHeader: "",
			bugSummary:      "echo_client.go:690 unconditionally inserts x-goog-request-params=\"\" even when empty",
		},
		"showcase_echo_unanchored_garbage_prefix": {
			wantBuggySent:   true,
			wantBuggyHeader: "header=GARBAGE%2Fprojects%2Fp1%2Finstances%2Fi1&instance_id=instances%2Fi1&routing_id=GARBAGE%2Fprojects%2Fp1%2Finstances%2Fi1&super_id=projects%2Fp1&table_name=projects%2Fp1%2Finstances%2Fi1",
			bugSummary:      "echo_client.go:671-678 uses unanchored regexes without ^...$, matching inside GARBAGE/ prefix",
		},
		"showcase_echo_rfc6570_space_and_tilde": {
			wantBuggySent:   true,
			wantBuggyHeader: "header=hello+world~1&routing_id=hello+world~1",
			bugSummary:      "echo_client.go uses url.QueryEscape (+) instead of RFC 6570 (%20) for spaces",
		},
	}

	for _, vec := range showcaseVectors {
		t.Run(vec.ID, func(t *testing.T) {
			capturedHeaders = nil
			var respHeaders, respTrailers metadata.MD
			resp, err := echoClient.Echo(
				t.Context(),
				&pb.EchoRequest{
					Header:      vec.Request["header"],
					OtherHeader: vec.Request["other_header"],
					Response:    &pb.EchoRequest_Content{Content: vec.ID},
				},
				gax.WithGRPCOptions(grpc.Header(&respHeaders), grpc.Trailer(&respTrailers)),
			)
			if err != nil {
				t.Fatalf("Echo RPC failed: %v", err)
			}
			if resp.GetContent() != vec.ID {
				t.Fatalf("Echo response content = %q, want %q", resp.GetContent(), vec.ID)
			}

			headerSent := len(capturedHeaders) > 0
			var rawWire string
			if headerSent {
				rawWire = capturedHeaders[0]
			}
			canonicalWire := canonicalizeRoutingHeader(rawWire)

			// Verify NewEchoServer echoes x-goog-request-params faithfully in trailers.
			trailerVals := respTrailers.Get("x-goog-request-params")
			if headerSent {
				if len(trailerVals) == 0 {
					t.Errorf("expected echoed x-goog-request-params trailer when header was sent (%q)", rawWire)
				} else if canonicalizeRoutingHeader(trailerVals[0]) != canonicalWire {
					t.Errorf("echoed trailer = %q, want canonical %q", trailerVals[0], canonicalWire)
				}
			}

			if knownBug, isDivergent := knownGeneratorDivergences[vec.ID]; isDivergent {
				if headerSent != knownBug.wantBuggySent || canonicalWire != knownBug.wantBuggyHeader {
					t.Fatalf("observed wire behavior changed for known divergence (%s): got (sent=%v, %q), want (sent=%v, %q)",
						knownBug.bugSummary, headerSent, canonicalWire, knownBug.wantBuggySent, knownBug.wantBuggyHeader)
				}
				if headerSent == vec.ExpectHeaderSent && canonicalWire == vec.ExpectedAlphaHeader {
					t.Errorf("expected divergence from Lean oracle %q, but wire output matched!", vec.ExpectedAlphaHeader)
				}
				t.Logf("pinned known gapic-generator-go divergence (%s): wire=%q vs lean=%q",
					knownBug.bugSummary, canonicalWire, vec.ExpectedAlphaHeader)
				return
			}

			if headerSent != vec.ExpectHeaderSent {
				t.Errorf("header sent = %v, want %v (wire=%q)", headerSent, vec.ExpectHeaderSent, canonicalWire)
			}
			if canonicalWire != vec.ExpectedAlphaHeader {
				t.Errorf("wire header mismatch:\n  got:  %q\n  want: %q", canonicalWire, vec.ExpectedAlphaHeader)
			}
		})
	}
}
