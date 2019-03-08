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
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewTokenGenerator() TokenGenerator {
	return &tokenGenerator{salt: strconv.FormatInt(time.Now().Unix(), 10)}
}

type TokenGenerator interface {
	ForIndex(int) string
	GetIndex(string) (int, error)
}

var InvalidTokenErr error = status.Errorf(
	codes.InvalidArgument,
	"The field `page_token` is invalid.")

type tokenGenerator struct {
	salt string
}

func (t *tokenGenerator) ForIndex(i int) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s%d", t.salt, i)))
}

func (t *tokenGenerator) GetIndex(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	bs, err := base64.StdEncoding.DecodeString(s)

	if err != nil {
		return -1, InvalidTokenErr
	}

	if !strings.HasPrefix(string(bs), t.salt) {
		return -1, InvalidTokenErr
	}

	i, err := strconv.Atoi(strings.TrimPrefix(string(bs), t.salt))
	if err != nil {
		return -1, InvalidTokenErr
	}
	return i, nil
}
