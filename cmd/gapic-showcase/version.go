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

package main

func init() {
	// Make roots version option only emit the version. This is used in circleci.
	// The template looks weird on purpose. Leaving as a single line causes the
	// output to append an extra character.
	rootCmd.Version = "0.12.0"
	rootCmd.SetVersionTemplate(
		`{{printf "%s" .Version}}`)
}
