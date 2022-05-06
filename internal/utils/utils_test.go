// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package utils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-cmp/cmp"
)

// TestSortSlice tests the SortSlice function.
func TestSortSlice(t *testing.T) {
	t.Parallel()
	{
		slice := []string{"b", "a"}
		sortedSlice := SortSlice(slice)
		if !cmp.Equal(slice, []string{"b", "a"}) {
			t.Fatalf("SortSlice did in place sort: %s", slice)
		}
		if !cmp.Equal(sortedSlice, []string{"a", "b"}) {
			t.Fatalf("unexpected sort return: %s", slice)
		}
	}
	{
		slice := SortSlice([]string{""})
		if !cmp.Equal(slice, []string{""}) {
			t.Fatalf("unexpected sort return: %s", slice)
		}
	}
	{
		slice := SortSlice([]string{"a", "b", "c", "d"})
		if !cmp.Equal(slice, []string{"a", "b", "c", "d"}) {
			t.Fatalf("unexpected sort return: %s", slice)
		}
	}
	{
		slice := SortSlice([]string{"c", "b", "a"})
		if !cmp.Equal(slice, []string{"a", "b", "c"}) {
			t.Fatalf("unexpected sort return: %s", slice)
		}
	}
	{
		slice := SortSlice([]string{"b", "a", "a"})
		if !cmp.Equal(slice, []string{"a", "a", "b"}) {
			t.Fatalf("unexpected sort return: %s", slice)
		}
	}
}

// TestSlicesAreEqual tests the SlicesAreEqual function.
func TestSlicesAreEqual(t *testing.T) {
	t.Parallel()

	if !SlicesAreEqual([]string{"a", "b", "c"}, []string{"a", "b", "c"}) {
		t.Fatalf("equal slice test failed when comparing two sorted copies")
	}

	if !SlicesAreEqual([]string{"a", "b", "c"}, []string{"b", "c", "a"}) {
		t.Fatalf("equal slice test failed when comparing two unsorted copies")
	}

	if !SlicesAreEqual([]string{}, []string{}) {
		t.Fatalf("equal slice test failed when comparing two empty slices")
	}

	if SlicesAreEqual([]string{"a", "b"}, []string{"b", "c"}) {
		t.Fatalf("equal slice test failed when comparing two different slices " +
			"of same size")
	}

	if SlicesAreEqual([]string{"a", "b"}, []string{"b", "c", "a"}) {
		t.Fatalf("equal slice test failed when comparing two different slices " +
			"of different size [1]")
	}

	if SlicesAreEqual([]string{"a", "b"}, []string{"a", "b", "a"}) {
		t.Fatalf("equal slice test failed when comparing two different slices " +
			"of different size [2]")
	}

	if SlicesAreEqual([]string{}, []string{"a", "b", "a"}) {
		t.Fatalf("equal slice test failed when comparing two different slices " +
			"of different size [3]")
	}

	if SlicesAreEqual([]string{"a"}, []string{"a", ""}) {
		t.Fatalf("equal slice test failed when comparing two different slices " +
			"of different size [4]")
	}

	if SlicesAreEqual([]string{"a"}, []string{"a", "a"}) {
		t.Fatalf("equal slice test failed when comparing two different slices " +
			"of different size [5]")
	}

	if SlicesAreEqual([]string{"a"}, []string{"", "a"}) {
		t.Fatalf("equal slice test failed when comparing two different slices " +
			"of different size [6]")
	}

	if SlicesAreEqual([]string{""}, []string{"", ""}) {
		t.Fatalf("equal slice test failed when comparing two different slices " +
			"of different size [6]")
	}
}

// TestNewBareRepo tests the NewBareRepo function.
func TestNewBareRepo(t *testing.T) {
	t.Parallel()

	path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary repo: %s", err)
	}

	defer os.RemoveAll(path)

	repo, err := NewBareRepo(path)
	if err != nil {
		t.Fatalf("failed to create a bare repo: %s", err)
	}

	refs, err := RepoRefsSlice(repo)
	if err != nil {
		t.Fatalf("failed to get repo's refs: %s", err)
	}

	if !SlicesAreEqual(refs, []string{"HEAD"}) {
		t.Fatalf("unexpected refs in repo: %s", refs)
	}
}

// TestNewTestRepo tests the testNewTestRepo function.
func TestNewTestRepo(t *testing.T) {
	t.Parallel()

	path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary repo: %s", err)
	}

	defer os.RemoveAll(path)

	repo, hash, err := NewTestRepo(path, []string{
		"refs/heads/foo",
		"refs/meta/bar",
	})
	if err != nil {
		t.Fatalf("failed to create a test repo: %s", err)
	}

	refs, err := RepoRefsSlice(repo)
	if err != nil {
		t.Fatalf("failed to get repo's refs: %s", err)
	}

	if !SlicesAreEqual(refs, []string{
		"HEAD",
		"refs/heads/master",
		"refs/heads/foo",
		"refs/meta/bar",
	}) {
		t.Fatalf("unexpected refs in repo: %s", refs)
	}

	check, err := RepoRefsCheckHash(repo, hash, "")
	if err != nil {
		t.Fatal("failed to check repo refs hash")
	}

	if !check {
		t.Fatal("unexpected ref hash")
	}
}

// TestRepoRefsSlice tests the RepoRefsSlice function.
func TestRepoRefsSlice(t *testing.T) {
	t.Parallel()

	path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary repo: %s", err)
	}

	defer os.RemoveAll(path)

	repo, _, err := NewTestRepo(path, []string{
		"refs/heads/a",
		"refs/heads/b",
	})
	if err != nil {
		t.Fatalf("error creating a test repo %s", err)
	}

	refs, err := RepoRefsSlice(repo)
	if err != nil {
		t.Fatalf("error getting the refs as slice: %s", err)
	}

	if !SlicesAreEqual(refs, []string{
		"HEAD",
		"refs/heads/master",
		"refs/heads/a",
		"refs/heads/b",
	}) {
		t.Fatalf("unexpected refs slice: %s", refs)
	}
}

// TestSpecsToStrings tests the SpecsToStrings function.
func TestSpecsToStrings(t *testing.T) {
	t.Parallel()
	{
		specs := SpecsToStrings([]config.RefSpec{})
		if !SlicesAreEqual(specs, []string{}) {
			t.Fatalf("unexpected specs return: %s", specs)
		}
	}
	{
		specs := SpecsToStrings([]config.RefSpec{
			"foo:bar",
			":foo",
		})
		if !SlicesAreEqual(specs, []string{
			"foo:bar",
			":foo",
		}) {
			t.Fatalf("unexpected specs return: %s", specs)
		}
	}
}

// TestRefsToStrings tests the RefsToStrings function.
func TestRefsToStrings(t *testing.T) {
	t.Parallel()
	{
		refs := RefsToStrings([]*plumbing.Reference{})
		if !SlicesAreEqual(refs, []string{}) {
			t.Fatalf("unexpected refs return: %s", refs)
		}
	}
	{
		refs := RefsToStrings([]*plumbing.Reference{
			plumbing.NewReferenceFromStrings("foo", ""),
			plumbing.NewReferenceFromStrings("bar", ""),
		})
		if !SlicesAreEqual(refs, []string{
			"foo",
			"bar",
		}) {
			t.Fatalf("unexpected refs return: %s", refs)
		}
	}
}

// TestRepoRefsCheckHash tests the RepoRefsCheckHash function.
func TestRepoRefsCheckHash(t *testing.T) {
	t.Parallel()

	path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary repo: %s", err)
	}

	defer os.RemoveAll(path)

	repo, hash, err := NewTestRepo(path, []string{
		"refs/heads/foo",
		"refs/meta/bar",
	})
	if err != nil {
		t.Fatalf("error creating a test repo %s", err)
	}

	pass, err := RepoRefsCheckHash(repo, hash, "")
	if err != nil {
		t.Fatalf("RepoRefsCheckHash failed: %s", err)
	}

	if !pass {
		t.Fatal("unexpected hash test result")
	}

	pass, err = RepoRefsCheckHash(repo, hash, "refs/")
	if err != nil {
		t.Fatalf("RepoRefsCheckHash failed: %s", err)
	}

	if !pass {
		t.Fatal("unexpected hash test result")
	}

	pass, err = RepoRefsCheckHash(repo, hash, "nomatch/")
	if err != nil {
		t.Fatalf("RepoRefsCheckHash failed: %s", err)
	}

	if !pass {
		t.Fatal("unexpected hash test result")
	}

	pass, err = RepoRefsCheckHash(repo, plumbing.NewHash(""), "")
	if err != nil {
		t.Fatalf("RepoRefsCheckHash failed: %s", err)
	}

	if pass {
		t.Fatal("unexpected hash test result")
	}
}
