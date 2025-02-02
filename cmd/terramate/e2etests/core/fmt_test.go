// Copyright 2023 Terramate GmbH
// SPDX-License-Identifier: MPL-2.0

package core_test

import (
	stdfmt "fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/madlambda/spells/assert"
	. "github.com/terramate-io/terramate/cmd/terramate/e2etests/internal/runner"
	"github.com/terramate-io/terramate/hcl/fmt"
	"github.com/terramate-io/terramate/test/sandbox"
)

func TestFormatRecursively(t *testing.T) {
	t.Parallel()

	const unformattedHCL = `
globals {
name = "name"
		description = "desc"
	test = true
	}
	`
	formattedHCL, err := fmt.Format(unformattedHCL, "")
	assert.NoError(t, err)

	s := sandbox.New(t)
	cli := NewCLI(t, s.RootDir())

	t.Run("checking succeeds when there is no Terramate files", func(t *testing.T) {
		AssertRunResult(t, cli.Run("fmt", "--check"), RunExpected{})
	})

	t.Run("formatting succeeds when there is no Terramate files", func(t *testing.T) {
		AssertRunResult(t, cli.Run("fmt"), RunExpected{})
	})

	sprintf := stdfmt.Sprintf
	writeUnformattedFiles := func() {
		s.BuildTree([]string{
			sprintf("f:globals.tm:%s", unformattedHCL),
			sprintf("f:another-stacks/globals.tm.hcl:%s", unformattedHCL),
			sprintf("f:another-stacks/stack-1/globals.tm.hcl:%s", unformattedHCL),
			sprintf("f:another-stacks/stack-2/globals.tm.hcl:%s", unformattedHCL),
			sprintf("f:stacks/globals.tm:%s", unformattedHCL),
			sprintf("f:stacks/stack-1/globals.tm:%s", unformattedHCL),
			sprintf("f:stacks/stack-2/globals.tm:%s", unformattedHCL),
		})
	}

	writeUnformattedFiles()

	wantedFiles := []string{
		"globals.tm",
		"another-stacks/globals.tm.hcl",
		"another-stacks/stack-1/globals.tm.hcl",
		"another-stacks/stack-2/globals.tm.hcl",
		"stacks/globals.tm",
		"stacks/stack-1/globals.tm",
		"stacks/stack-2/globals.tm",
	}
	filesListOutput := func(files []string) string {
		portablewantedFiles := make([]string, len(files))
		for i, f := range files {
			portablewantedFiles[i] = filepath.FromSlash(f)
		}
		return strings.Join(portablewantedFiles, "\n") + "\n"
	}
	wantedFilesStr := filesListOutput(wantedFiles)

	assertFileContents := func(t *testing.T, path string, want string) {
		t.Helper()
		got := s.RootEntry().ReadFile(path)
		assert.EqualStrings(t, want, string(got))
	}

	assertWantedFilesContents := func(t *testing.T, want string) {
		t.Helper()

		for _, file := range wantedFiles {
			assertFileContents(t, file, want)
		}
	}

	t.Run("checking fails with unformatted files", func(t *testing.T) {
		AssertRunResult(t, cli.Run("fmt", "--check"), RunExpected{
			Status: 1,
			Stdout: wantedFilesStr,
		})
		assertWantedFilesContents(t, unformattedHCL)
	})

	t.Run("checking fails with unformatted files on subdirs", func(t *testing.T) {
		subdir := filepath.Join(s.RootDir(), "another-stacks")
		cli := NewCLI(t, subdir)
		AssertRunResult(t, cli.Run("fmt", "--check"), RunExpected{
			Status: 1,
			Stdout: filesListOutput([]string{
				"globals.tm.hcl",
				"stack-1/globals.tm.hcl",
				"stack-2/globals.tm.hcl",
			}),
		})
		assertWantedFilesContents(t, unformattedHCL)
	})

	t.Run("update unformatted files in place", func(t *testing.T) {
		AssertRunResult(t, cli.Run("fmt"), RunExpected{
			Stdout: wantedFilesStr,
		})
		assertWantedFilesContents(t, formattedHCL)
	})

	t.Run("checking succeeds when all files are formatted", func(t *testing.T) {
		AssertRunResult(t, cli.Run("fmt", "--check"), RunExpected{})
		assertWantedFilesContents(t, formattedHCL)
	})

	t.Run("formatting succeeds when all files are formatted", func(t *testing.T) {
		AssertRunResult(t, cli.Run("fmt"), RunExpected{})
		assertWantedFilesContents(t, formattedHCL)
	})

	t.Run("update unformatted files in subdirs", func(t *testing.T) {
		writeUnformattedFiles()

		anotherStacks := filepath.Join(s.RootDir(), "another-stacks")
		cli := NewCLI(t, anotherStacks)
		AssertRunResult(t, cli.Run("fmt"), RunExpected{
			Stdout: filesListOutput([]string{
				"globals.tm.hcl",
				"stack-1/globals.tm.hcl",
				"stack-2/globals.tm.hcl",
			}),
		})

		assertFileContents(t, "another-stacks/globals.tm.hcl", formattedHCL)
		assertFileContents(t, "another-stacks/stack-1/globals.tm.hcl", formattedHCL)
		assertFileContents(t, "another-stacks/stack-2/globals.tm.hcl", formattedHCL)

		assertFileContents(t, "globals.tm", unformattedHCL)
		assertFileContents(t, "stacks/globals.tm", unformattedHCL)
		assertFileContents(t, "stacks/stack-1/globals.tm", unformattedHCL)
		assertFileContents(t, "stacks/stack-2/globals.tm", unformattedHCL)

		stacks := filepath.Join(s.RootDir(), "stacks")
		cli = NewCLI(t, stacks)
		AssertRunResult(t, cli.Run("fmt"), RunExpected{
			Stdout: filesListOutput([]string{
				"globals.tm",
				"stack-1/globals.tm",
				"stack-2/globals.tm",
			}),
		})

		assertFileContents(t, "another-stacks/globals.tm.hcl", formattedHCL)
		assertFileContents(t, "another-stacks/stack-1/globals.tm.hcl", formattedHCL)
		assertFileContents(t, "another-stacks/stack-2/globals.tm.hcl", formattedHCL)
		assertFileContents(t, "stacks/globals.tm", formattedHCL)
		assertFileContents(t, "stacks/stack-1/globals.tm", formattedHCL)
		assertFileContents(t, "stacks/stack-2/globals.tm", formattedHCL)

		assertFileContents(t, "globals.tm", unformattedHCL)

		cli = NewCLI(t, s.RootDir())
		AssertRunResult(t, cli.Run("fmt"), RunExpected{
			Stdout: filesListOutput([]string{"globals.tm"}),
		})

		assertWantedFilesContents(t, formattedHCL)
	})
}
