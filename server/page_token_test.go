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
	"encoding/base64"
	"testing"
)

func Test_tokenGenerator_ForIndex(t *testing.T) {
	salt := "salt"
	index := 1
	want := base64.StdEncoding.EncodeToString(
		[]byte("salt1"))
	tok := TokenGeneratorWithSalt(salt)
	if got := tok.ForIndex(index); got != want {
		t.Errorf("tokenGenerator.ForIndex() = %v, want %v", got, want)
	}
}

func Test_tokenGenerator_GetIndex_notParseable(t *testing.T) {
	tok := NewTokenGenerator()
	_, err := tok.GetIndex("invalid")
	if err == nil {
		t.Error("GetIndex: want error for invalid token.")
	}
}
func Test_tokenGenerator_GetIndex_noSalt(t *testing.T) {
	tok := NewTokenGenerator()
	_, err := tok.GetIndex(base64.StdEncoding.EncodeToString([]byte("invalid")))
	if err == nil {
		t.Error("GetIndex: want error for invalid token.")
	}
}

func Test_tokenGenerator_GetIndex_invalidIndex(t *testing.T) {
	tok := TokenGeneratorWithSalt("salt")
	_, err := tok.GetIndex(base64.StdEncoding.EncodeToString([]byte("saltinvalid")))
	if err == nil {
		t.Error("GetIndex: want error for invalid token.")
	}
}

func Test_tokenGenerator_GetIndex(t *testing.T) {
	tok := TokenGeneratorWithSalt("salt")
	i, err := tok.GetIndex(base64.StdEncoding.EncodeToString([]byte("salt1")))
	if err != nil {
		t.Error("GetIndex: unexpected err")
	}
	if i != 1 {
		t.Errorf("GetIndex: want 1, got %d", i)
	}
}
