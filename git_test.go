// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/agherzan/git-mirror-me/internal/utils"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	testSSHKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAlwAAAAdzc2gtcn
NhAAAAAwEAAQAAAIEAxi9PL+/GEMwIsmQeWm50/LNZqSrxk4Oa3D+W2iTDKbmE/2RgroHX
/Lc+V6r4HDZTticNDeHr3mLMSBKR4YSgySp+TQWIflLEM8wO2MnkmJ07oH3BHV/bIz7rIp
h06nddaxq5hixNubVqYqZTjppwYPnu3nJzeV1V2IK/UgCAMpsAAAIIGvv3sxr797MAAAAH
c3NoLXJzYQAAAIEAxi9PL+/GEMwIsmQeWm50/LNZqSrxk4Oa3D+W2iTDKbmE/2RgroHX/L
c+V6r4HDZTticNDeHr3mLMSBKR4YSgySp+TQWIflLEM8wO2MnkmJ07oH3BHV/bIz7rIph0
6nddaxq5hixNubVqYqZTjppwYPnu3nJzeV1V2IK/UgCAMpsAAAADAQABAAAAgEGLpgn5qC
0n/fxaBnvsKj7lZlL/w/QAw7fyRAcTv4ROOkFpRlyQzwli5XiDMBnMkfUdh0C/Jo5faKax
lZPblH0+CZqYU0gDPYDUjkIWbrhDVd8b4j56dJ9Oa5e+exuNaS9oR+1IaFyFtwBvkhs3pk
Rs/AmRK25/vvWLIASAAFp5AAAAQQCi+oPpPJmqkDXzshEgkoaRNIp6s5QNaQsA4Ra7Sk0K
9mi1sY8lngYO+4ln5Rr2lcp8ZsxPleEuA6ISIChoNaeKAAAAQQD5xWVL8NgAdDNA0F4Th/
KEAAVL5xBzHfH3q+OV30mfE5pPItvRRrkdzO6uQTqlaKF+9vQWTt3DJOdnZw+fqhSXAAAA
QQDLIJLCNDXRLbDX+mHaq7PPb+Y+ZAU9TLJw8MXgki3cNm+oSzYM3g4RSL5t/BEobOIhYL
Hau0thh3byP4srEz6dAAAADmFuZHJlaUBxd2lya2xlAQIDBA==
-----END OPENSSH PRIVATE KEY-----`
	testKnownHost = "github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAA" +
		"AIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl"
)

// TestFilterOutRefsMatch tests the filterOutRefs function when the filter matches
// some references.
func TestFilterOutRefsMatch(t *testing.T) {
	t.Parallel()

	path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary repo: %s", err)
	}

	defer os.RemoveAll(path)

	repo, head, err := utils.NewTestRepo(path, []string{
		"refs/heads/a",
		"refs/heads/b",
		"refs/meta/a",
		"refs/meta/b",
	})
	if err != nil {
		t.Fatalf("failed to create a test repo: %s", err)
	}

	err = filterOutRefs(repo, []string{"refs/meta"})
	if err != nil {
		t.Fatalf("failed to filter refs: %s", err)
	}

	refs, err := utils.RepoRefsSlice(repo)
	if err != nil {
		t.Fatalf("failed to get repo's refs: %s", err)
	}

	if !utils.SlicesAreEqual(refs, []string{
		"HEAD",
		"refs/heads/master",
		"refs/heads/a",
		"refs/heads/b",
	}) {
		t.Fatalf("unexpected refs in repo: %s", refs)
	}

	check, err := utils.RepoRefsCheckHash(repo, head)
	if err != nil {
		t.Fatal("failed to check repo refs hash")
	}

	if !check {
		t.Fatal("unexpected ref hash")
	}
}

// TestFilterOutRefsNoMatch tests the filterOutRefs function when the filter doesn't
// match.
func TestFilterOutRefsNoMatch(t *testing.T) {
	t.Parallel()

	path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary repo: %s", err)
	}

	defer os.RemoveAll(path)

	repo, head, err := utils.NewTestRepo(path, []string{
		"refs/heads/a",
		"refs/heads/b",
		"refs/meta/a",
		"refs/meta/b",
	})
	if err != nil {
		t.Fatalf("failed to create a test repo: %s", err)
	}

	err = filterOutRefs(repo, []string{"refs/nonexistent"})
	if err != nil {
		t.Fatalf("failed to filter refs: %s", err)
	}

	refs, err := utils.RepoRefsSlice(repo)
	if err != nil {
		t.Fatalf("failed to get repo's refs: %s", err)
	}

	if !utils.SlicesAreEqual(refs, []string{
		"HEAD",
		"refs/heads/master",
		"refs/heads/a",
		"refs/heads/b",
		"refs/meta/a",
		"refs/meta/b",
	}) {
		t.Fatalf("unexpected refs in repo: %s", refs)
	}

	check, err := utils.RepoRefsCheckHash(repo, head)
	if err != nil {
		t.Fatal("failed to check repo refs hash")
	}

	if !check {
		t.Fatal("unexpected ref hash")
	}
}

// TestFilterOutRefsDeleteAll tests the filterOutRefs function for deleting all
// references.
func TestFilterOutRefsDeleteAll(t *testing.T) {
	t.Parallel()

	path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary repo: %s", err)
	}

	defer os.RemoveAll(path)

	repo, head, err := utils.NewTestRepo(path, []string{
		"refs/heads/a",
		"refs/heads/b",
		"refs/meta/a",
		"refs/meta/b",
	})
	if err != nil {
		t.Fatalf("failed to create a test repo: %s", err)
	}

	err = filterOutRefs(repo, []string{""})
	if err != nil {
		t.Fatalf("failed to filter refs: %s", err)
	}

	refs, err := utils.RepoRefsSlice(repo)
	if err != nil {
		t.Fatalf("failed to get repo's refs: %s", err)
	}

	if !utils.SlicesAreEqual(refs, []string{}) {
		t.Fatalf("unexpected refs in repo: %s", refs)
	}

	check, err := utils.RepoRefsCheckHash(repo, head)
	if err != nil {
		t.Fatal("failed to check repo refs hash")
	}

	if !check {
		t.Fatal("unexpected ref hash")
	}
}

// TestFilterOutRefsNoPrefix tests the filterOutRefs function when there is no prefix
// provided.
func TestFilterOutRefsNoPrefix(t *testing.T) {
	t.Parallel()

	path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary repo: %s", err)
	}

	defer os.RemoveAll(path)

	repo, head, err := utils.NewTestRepo(path, []string{})
	if err != nil {
		t.Fatalf("failed to create a test repo: %s", err)
	}

	err = filterOutRefs(repo, []string{})
	if err != nil {
		t.Fatalf("failed to filter refs: %s", err)
	}

	refs, err := utils.RepoRefsSlice(repo)
	if err != nil {
		t.Fatalf("failed to get repo's refs: %s", err)
	}

	if !utils.SlicesAreEqual(refs, []string{
		"HEAD",
		"refs/heads/master",
	}) {
		t.Fatalf("unexpected refs in repo: %s", refs)
	}

	check, err := utils.RepoRefsCheckHash(repo, head)
	if err != nil {
		t.Fatal("failed to check repo refs hash")
	}

	if !check {
		t.Fatal("unexpected ref hash")
	}
}

// TestRefsToDeleteSpecs tests refsToDeleteSpecs function.
func TestRefsToDeleteSpecs(t *testing.T) {
	t.Parallel()

	{
		specs := refsToDeleteSpecs([]*plumbing.Reference{})
		if !utils.SlicesAreEqual(utils.SpecsToStrings(specs), []string{}) {
			t.Fatal("unexpected delete specs")
		}
	}
	{
		specs := refsToDeleteSpecs([]*plumbing.Reference{
			plumbing.NewReferenceFromStrings("foo", ""),
			plumbing.NewReferenceFromStrings("bar", ""),
		})
		if !utils.SlicesAreEqual(utils.SpecsToStrings(specs), []string{
			":foo",
			":bar",
		}) {
			t.Fatal("unexpected delete specs")
		}
	}
}

// TestExtraRefs tests extraRefs function.
func TestExtraRefs(t *testing.T) {
	t.Parallel()

	{
		path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
		if err != nil {
			t.Fatalf("failed to create a temporary repo: %s", err)
		}
		defer os.RemoveAll(path)
		repo, _, err := utils.NewTestRepo(path, []string{
			"refs/heads/a",
			"refs/heads/b",
		})
		if err != nil {
			t.Fatalf("failed to create a test repo: %s", err)
		}
		refs, err := extraRefs(repo, []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("refs/heads/a", ""),
			plumbing.NewReferenceFromStrings("refs/heads/b", ""),
			plumbing.NewReferenceFromStrings("refs/meta/a", ""),
			plumbing.NewReferenceFromStrings("refs/meta/b", ""),
		})
		if err != nil {
			t.Fatalf("failed to get extra refs: %s", err)
		}
		if !utils.SlicesAreEqual(utils.RefsToStrings(refs), []string{
			"refs/meta/a",
			"refs/meta/b",
		}) {
			t.Fatal("unexpected extra refs")
		}
	}
	{
		path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
		if err != nil {
			t.Fatalf("failed to create a temporary repo: %s", err)
		}
		defer os.RemoveAll(path)
		repo, _, err := utils.NewTestRepo(path, []string{
			"refs/heads/a",
			"refs/heads/b",
			"refs/meta/a",
			"refs/meta/b",
		})
		if err != nil {
			t.Fatalf("failed to create a test repo: %s", err)
		}
		refs, err := extraRefs(repo, []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("refs/heads/a", ""),
			plumbing.NewReferenceFromStrings("refs/heads/b", ""),
			plumbing.NewReferenceFromStrings("refs/meta/a", ""),
			plumbing.NewReferenceFromStrings("refs/meta/b", ""),
		})
		if err != nil {
			t.Fatalf("failed to get extra refs: %s", err)
		}
		if !utils.SlicesAreEqual(utils.RefsToStrings(refs), []string{}) {
			t.Fatal("unexpected extra refs")
		}
	}
	{
		path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
		if err != nil {
			t.Fatalf("failed to create a temporary repo: %s", err)
		}
		defer os.RemoveAll(path)
		repo, _, err := utils.NewTestRepo(path, []string{
			"refs/heads/a",
			"refs/heads/b",
		})
		if err != nil {
			t.Fatalf("failed to create a test repo: %s", err)
		}
		refs, err := extraRefs(repo, []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("refs/meta/a", ""),
			plumbing.NewReferenceFromStrings("refs/meta/b", ""),
		})
		if err != nil {
			t.Fatalf("failed to get extra refs: %s", err)
		}
		if !utils.SlicesAreEqual(utils.RefsToStrings(refs), []string{
			"refs/meta/a",
			"refs/meta/b",
		}) {
			t.Fatal("unexpected extra refs")
		}
	}
}

// TestExtraSpecs tests extraSpecs function.
func TestExtraSpecs(t *testing.T) {
	t.Parallel()

	{
		path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
		if err != nil {
			t.Fatalf("failed to create a temporary repo: %s", err)
		}
		defer os.RemoveAll(path)
		repo, _, err := utils.NewTestRepo(path, []string{
			"refs/heads/a",
			"refs/heads/b",
		})
		if err != nil {
			t.Fatalf("failed to create a test repo: %s", err)
		}
		refs, err := extraSpecs(repo, []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("refs/heads/a", ""),
			plumbing.NewReferenceFromStrings("refs/heads/b", ""),
			plumbing.NewReferenceFromStrings("refs/meta/a", ""),
			plumbing.NewReferenceFromStrings("refs/meta/b", ""),
		})
		if err != nil {
			t.Fatalf("failed to get extra refs: %s", err)
		}
		if !utils.SlicesAreEqual(utils.SpecsToStrings(refs), []string{
			":refs/meta/a",
			":refs/meta/b",
		}) {
			t.Fatal("unexpected extra refs")
		}
	}
	{
		path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
		if err != nil {
			t.Fatalf("failed to create a temporary repo: %s", err)
		}
		defer os.RemoveAll(path)
		repo, _, err := utils.NewTestRepo(path, []string{
			"refs/heads/a",
			"refs/heads/b",
			"refs/meta/a",
			"refs/meta/b",
		})
		if err != nil {
			t.Fatalf("failed to create a test repo: %s", err)
		}
		refs, err := extraSpecs(repo, []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("refs/heads/a", ""),
			plumbing.NewReferenceFromStrings("refs/heads/b", ""),
			plumbing.NewReferenceFromStrings("refs/meta/a", ""),
			plumbing.NewReferenceFromStrings("refs/meta/b", ""),
		})
		if err != nil {
			t.Fatalf("failed to get extra refs: %s", err)
		}
		if !utils.SlicesAreEqual(utils.SpecsToStrings(refs), []string{}) {
			t.Fatal("unexpected extra refs")
		}
	}
	{
		path, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
		if err != nil {
			t.Fatalf("failed to create a temporary repo: %s", err)
		}
		defer os.RemoveAll(path)
		repo, _, err := utils.NewTestRepo(path, []string{
			"refs/heads/a",
			"refs/heads/b",
		})
		if err != nil {
			t.Fatalf("failed to create a test repo: %s", err)
		}
		refs, err := extraSpecs(repo, []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("refs/meta/a", ""),
			plumbing.NewReferenceFromStrings("refs/meta/b", ""),
		})
		if err != nil {
			t.Fatalf("failed to get extra refs: %s", err)
		}
		if !utils.SlicesAreEqual(utils.SpecsToStrings(refs), []string{
			":refs/meta/a",
			":refs/meta/b",
		}) {
			t.Fatal("unexpected extra refs")
		}
	}
}

// TestSetupStagingRepo tests setupStagingRepo function.
func TestSetupStagingRepo(t *testing.T) {
	t.Parallel()

	// no need for logs
	devnull, _ := os.Open(os.DevNull)
	defer devnull.Close()
	logger := NewLogger(devnull)

	// Create a test repository and setup a Staging using the test repository
	// as source.
	srcRepoPath, err := ioutil.TempDir("/tmp", "git-mirror-me-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary repo: %s", err)
	}

	defer os.RemoveAll(srcRepoPath)

	_, srcHead, err := utils.NewTestRepo(srcRepoPath, []string{
		"refs/heads/a",
		"refs/heads/b",
		"refs/meta/a",
		"refs/meta/b",
	})
	if err != nil {
		t.Fatalf("failed to create a test repo: %s", err)
	}

	// First test that it fails with an invalid source.
	_, err = setupStagingRepo(Config{
		SrcRepo: "/invalid",
	}, logger)
	if err == nil {
		t.Fatal("setupStagingRepo with an invalid source")
	}

	stagingRepo, err := setupStagingRepo(Config{
		SrcRepo: srcRepoPath,
	}, logger)
	if err != nil {
		t.Fatalf("failed to setup the staging repo: %s", err)
	}

	// Check that all the refs are in place and they all point to the right
	// hash.
	stagingRepoRefs, err := utils.RepoRefsSlice(stagingRepo)
	if err != nil {
		t.Fatalf("failed to get the refs: %s", err)
	}

	if !utils.SlicesAreEqual(stagingRepoRefs, []string{
		"HEAD",
		"refs/heads/master",
		"refs/heads/a",
		"refs/heads/b",
		"refs/meta/a",
		"refs/meta/b",
	}) {
		t.Fatalf("unexpected refs in staging repo: %s", stagingRepoRefs)
	}

	ok, err := utils.RepoRefsCheckHash(stagingRepo, srcHead)
	if err != nil {
		t.Fatalf("utils.RepoRefsCheckHash failed: %s", err)
	}

	if ok {
		t.Fatal("unexpected hash test result")
	}
}

// TestDoMirror tests DoMirror function.
func TestDoMirror(t *testing.T) {
	t.Parallel()

	// no need for logs
	devnull, _ := os.Open(os.DevNull)
	defer devnull.Close()
	logger := NewLogger(devnull)

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

	// Define the configuration and run the function under test.
	conf := Config{
		SrcRepo: srcRepoPath,
		DstRepo: dstRepoPath,
		SSH: SSHConf{
			PrivateKey: testSSHKey,
			KnownHosts: testKnownHost,
		},
	}

	err = DoMirror(conf, logger)
	if err != nil {
		t.Fatalf("DoMirror failed: %s", err)
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
