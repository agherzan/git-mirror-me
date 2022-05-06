// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package main

import (
	"io/ioutil"
	"os"
	"testing"

	mirror "github.com/agherzan/git-mirror-me"
	"github.com/agherzan/git-mirror-me/internal/utils"
)

// TestRun tests the run function.
func TestRun(t *testing.T) {
	t.Parallel()

	// no need for logs
	devnull, _ := os.Open(os.DevNull)
	defer devnull.Close()
	logger := mirror.NewLogger(devnull)

	// Create a source repository.
	srcRepoPath, err := ioutil.TempDir("/tmp", "git-mirror-me-test-src-")
	if err != nil {
		t.Fatalf("failed to create a temporary src repo: %s", err)
	}

	defer os.RemoveAll(srcRepoPath)

	_, srcHead, err := utils.NewTestRepo(srcRepoPath, []string{
		"refs/heads/a",
		"refs/heads/b",
		"refs/pull/1",
		"refs/pull/2",
		"refs/meta/foo",
	})
	if err != nil {
		t.Fatalf("failed to create a test src repo: %s", err)
	}

	// Create a destination repository.
	dstRepoPath, err := ioutil.TempDir("/tmp", "git-mirror-me-test-dst-")
	if err != nil {
		t.Fatalf("failed to create a temporary dst repo: %s", err)
	}

	defer os.RemoveAll(dstRepoPath)

	dstRepo, _, err := utils.NewTestRepo(dstRepoPath, []string{
		"refs/heads/a",
		"refs/heads/b",
		"refs/heads/c",
	})
	if err != nil {
		t.Fatalf("failed to create a test dst repo: %s", err)
	}

	// Test help.
	args := []string{"-help"}
	if err := run(logger, map[string]string{}, "test", args); err != nil {
		t.Fatalf("help failed: %s", err)
	}

	// Test invalid argument.
	args = []string{"-invalidflag"}
	if err := run(logger, map[string]string{}, "test", args); err == nil {
		t.Fatal("invalid argument passed")
	}

	// Fail configuration.
	if err := run(logger, map[string]string{}, "test", []string{}); err == nil {
		t.Fatal("invalid configuration passed")
	}

	// Fail on an invalid destination repository.
	env := map[string]string{"GMM_SRC_REPO": srcRepoPath}
	args = []string{"--destination-repository", "invalid"}

	if err := run(logger, env, "test", args); err == nil {
		t.Fatal("run succeeded with an invalid dst repository")
	}

	// Valid run.
	env = map[string]string{"GMM_SRC_REPO": srcRepoPath}
	args = []string{"--destination-repository", dstRepoPath}

	if err := run(logger, env, "test", args); err != nil {
		t.Fatalf("run failed: %s", err)
	}

	// Verify the destination.
	dstRepoRefs, err := utils.RepoRefsSlice(dstRepo)
	if err != nil {
		t.Fatalf("failed to get the dst repo refs: %s", err)
	}

	if !utils.SlicesAreEqual(dstRepoRefs, []string{
		"HEAD",
		"refs/heads/master",
		"refs/heads/a",
		"refs/heads/b",
		"refs/meta/foo",
	}) {
		t.Fatalf("unexpected refs in the dst repo: %s", dstRepoRefs)
	}

	ok, err := utils.RepoRefsCheckHash(dstRepo, srcHead)
	if err != nil {
		t.Fatalf("dst repo hash check failed: %s", err)
	}

	if !ok {
		t.Fatal("unexpected hash test result for the dst repo")
	}
}
