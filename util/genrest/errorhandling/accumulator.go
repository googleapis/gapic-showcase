// Copyright 2020 Google LLC
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

package errorhandling

import (
	"fmt"
	"strings"
)

// Accumulator allows storing a series of errors and then concatenating their string representations
// into one big error.
type Accumulator struct {
	errors []error
}

// AccumulateError stores an error to be reported later.
func (ea *Accumulator) AccumulateError(err error) {
	if err == nil {
		return
	}
	if ea.errors == nil {
		ea.errors = []error{}
	}
	ea.errors = append(ea.errors, err)
}

// Error concatenates the string representations of all stored errors and returns it as a single
// error.
func (ea *Accumulator) Error() error {
	if len(ea.errors) == 0 {
		return nil
	}
	errorStrings := make([]string, len(ea.errors))
	for idx, err := range ea.errors {
		errorStrings[idx] = err.Error()
	}
	return fmt.Errorf(strings.Join(errorStrings, "\n"))
}
