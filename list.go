// Copyright 2021 Mineiros GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package terramate

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ListStacks walks the basedir directory looking for terraform stacks.
// It returns a lexicographic sorted list of stack directories.
func ListStacks(basedir string) ([]Entry, error) {
	entries := []Entry{}

	err := filepath.Walk(
		basedir,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				stackfile := filepath.Join(path, ConfigFilename)
				st, err := os.Stat(stackfile)
				if err != nil || !st.Mode().IsRegular() {
					return nil
				}

				entries = append(entries, Entry{Dir: path})
			}

			return nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("while walking dir: %w", err)
	}

	return entries, nil
}