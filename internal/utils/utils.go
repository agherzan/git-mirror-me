// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package utils

import (
	"errors"
	"sort"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

// SortSlice gets a slice and returns a sorted copy of it.
func SortSlice(slice []string) []string {
	if sort.StringsAreSorted(slice) {
		return slice
	}

	sortedSlice := make([]string, len(slice))
	copy(sortedSlice, slice)
	sort.Strings(sortedSlice)

	return sortedSlice
}

// SlicesAreEqual checks if two slices are equal. Order is ignored but
// duplicates are not.
func SlicesAreEqual(sliceA, sliceB []string) bool {
	sliceA = SortSlice(sliceA)
	sliceB = SortSlice(sliceB)

	if len(sliceA) != len(sliceB) {
		return false
	}

	for i := 0; i < len(sliceA); i++ {
		if sliceA[i] != sliceB[i] {
			return false
		}
	}

	return true
}

// NewBareRepo creates a new bare repo at a specific path.
func NewBareRepo(path string) (*git.Repository, error) {
	if len(path) == 0 {
		return nil, errors.New("invalid path for repo")
	}
	bareRepo, err := git.PlainInit(path, true)
	if err != nil {
		return nil, err
	}

	return bareRepo, nil
}

// NewTestRepo creates an new bare repo at a specific path initialised with a
// test commit and a set of refs pointing to the HEAD's reference.
func NewTestRepo(path string, refs []string) (*git.Repository, plumbing.Hash, error) {
	var headHash plumbing.Hash

	// Create a bare repository.
	bareRepo, err := NewBareRepo(path)
	if err != nil {
		return nil, headHash, err
	}

	// Create an in-memory repository so we can initialise the references.
	memoryFs := memfs.New()

	repo, err := git.Init(memory.NewStorage(), memoryFs)
	if err != nil {
		return nil, headHash, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, headHash, err
	}

	testFile, err := memoryFs.Create("testfile.txt")
	if err != nil {
		return nil, headHash, err
	}

	defer testFile.Close()

	_, err = testFile.Write([]byte("test"))
	if err != nil {
		return nil, headHash, err
	}

	_, err = worktree.Add("testfile.txt")
	if err != nil {
		return nil, headHash, err
	}

	_, err = worktree.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Example",
			Email: "ex@ample.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, headHash, err
	}

	// Set the references to the HEAD's hash.
	head, err := repo.Head()
	if err != nil {
		return nil, headHash, err
	}

	headHash = head.Hash()

	for _, ref := range refs {
		r := plumbing.NewHashReference(plumbing.ReferenceName(ref), headHash)

		err = repo.Storer.SetReference(r)
		if err != nil {
			return nil, headHash, err
		}
	}

	// Finally push the in-memory repository to the bare one.
	bareRepoRemote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{path},
	})
	if err != nil {
		return nil, headHash, err
	}

	err = bareRepoRemote.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"*:*"},
		Prune:      true,
		Force:      true,
	})
	if err != nil {
		return nil, headHash, err
	}

	// Also, set the HEAD of the bare repo so that all refs point to the same
	// test revision.
	bareHead := plumbing.NewHashReference(plumbing.HEAD, headHash)

	err = bareRepo.Storer.SetReference(bareHead)
	if err != nil {
		return nil, headHash, err
	}

	return bareRepo, headHash, nil
}

// RepoRefsSlice takes a repository and returns the name of all the references
// as a string slice.
func RepoRefsSlice(repo *git.Repository) ([]string, error) {
	var refsSlice []string
	refs, err := repo.References()
	if err != nil {
		return nil, err
	}
	refs.ForEach(func(ref *plumbing.Reference) error {
		refsSlice = append(refsSlice, ref.Name().String())
		return nil
	})
	return refsSlice, nil
}

// SpecsToStrings takes a slice of refspecs and returns them as a slice of
// strings.
func SpecsToStrings(specs []config.RefSpec) []string {
	str := make([]string, 0, len(specs))
	for _, spec := range specs {
		str = append(str, spec.String())
	}
	return str
}

// RefsToStrings takes a slice of references and returns their names as a slice
// of strings.
func RefsToStrings(refs []*plumbing.Reference) []string {
	str := make([]string, 0, len(refs))
	for _, ref := range refs {
		str = append(str, ref.Name().String())
	}
	return str
}

// RepoRefsCheckHash checks if all the references in a repository point to the
// same hash.
func RepoRefsCheckHash(repo *git.Repository, hash plumbing.Hash) (bool, error) {
	var result bool

	refs, err := repo.References()
	if err != nil {
		return result, err
	}

	result = true
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Hash() != hash {
			result = false
		}

		return nil
	})

	return result, nil
}
