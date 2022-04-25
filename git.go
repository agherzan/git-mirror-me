// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

const (
	refsFilterPrefix       = "refs/pull"
	srcRemoteName          = "src"
	dstRemoteName          = "dst"
	tmpKnownHostPathPrefix = "git-mirror-me-known_hosts-"
)

// FilterOutRefs takes a repository and removes references based on a slice of
// prefixes.
func filterOutRefs(repo *git.Repository, prefixes []string) error {
	if len(prefixes) == 0 {
		return nil
	}
	refs, err := repo.References()
	if err != nil {
		return err
	}
	if err = refs.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().String()
		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) {
				if err := repo.Storer.RemoveReference(ref.Name()); err != nil {
					return err
				}
				break
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// refsToDeleteSpecs returns a slice of delete refspecs for a slice of
// references.
func refsToDeleteSpecs(refs []*plumbing.Reference) []config.RefSpec {
	var specs []config.RefSpec
	for _, ref := range refs {
		specs = append(specs, config.RefSpec(":"+ref.Name().String()))
	}
	return specs
}

// extraRefs returns a slice of references that are in refs but not in the
// repository.
func extraRefs(repo *git.Repository, refs []*plumbing.Reference) ([]*plumbing.Reference, error) {
	var retRefs []*plumbing.Reference
	for _, ref := range refs {
		repoRefs, err := repo.References()
		if err != nil {
			return nil, err
		}
		found := false
		repoRefs.ForEach(func(repoRef *plumbing.Reference) error {
			if repoRef.Name().String() == ref.Name().String() {
				found = true
			}
			return nil
		})
		if !found {
			retRefs = append(retRefs, ref)
		}
	}
	return retRefs, nil
}

// extraSpecs takes a repository and a slice of refs and returns the refs
// that are not in the repository as a slice of delete refspecs.
func extraSpecs(repo *git.Repository, refs []*plumbing.Reference) ([]config.RefSpec, error) {
	diffRefs, err := extraRefs(repo, refs)
	if err != nil {
		return nil, err
	}
	return refsToDeleteSpecs(diffRefs), nil
}

// setupStagingRepo initialises an in-memory git repositry populated with the
// source's references.
func setupStagingRepo(conf Config, logger *Logger) (*git.Repository, error) {
	// Setup a working repository.
	logger.Info("Setting up a staging git repository.")
	repo, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed initialising staging git repository: %w",
			err)
	}

	// Set up the source remote.
	src, err := repo.CreateRemote(&config.RemoteConfig{
		Name: srcRemoteName,
		URLs: []string{conf.SrcRepo},
	})
	if err != nil {
		return nil, fmt.Errorf("failed configuring source remote: %w", err)
	}

	// Fetch the source.
	logger.Info("Fetching all refs from", conf.SrcRepo, "...")
	if err := src.Fetch(&git.FetchOptions{
		RemoteName: srcRemoteName,
		RefSpecs:   []config.RefSpec{"refs/*:refs/*"},
	}); err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, fmt.Errorf("failed to fetch source remote: %w", err)
	}

	return repo, nil
}

// pushWithAuth sets authentication based on configuration and pushes all
// references to the configured destination repository (as a mirror).
func pushWithAuth(conf Config, logger *Logger, stagingRepo *git.Repository) error {
	var auth transport.AuthMethod
	// Set up SSH authentication.
	if len(conf.Ssh.PrivateKey) > 0 {
		logger.Debug(conf.Debug, "Using SSH authentication.")
		sshKeys, err := ssh.NewPublicKeys("git", []byte(conf.Ssh.PrivateKey), "")
		if err != nil {
			return fmt.Errorf("failed to setup the SSH key: %w", err)
		}

		// The host public keys can be provided via both content and path. When
		// it is provided via content, we need to use a temporary known_hosts
		// file.
		knownHostsPath := conf.GetKnownHostsPath()
		if len(conf.Ssh.KnownHosts) != 0 {
			f, err := ioutil.TempFile("/tmp", tmpKnownHostPathPrefix)
			if err != nil {
				return fmt.Errorf("error creating known_hosts tmp file: %w", err)
			}
			defer func() {
				f.Close()
				os.Remove(f.Name())
			}()
			knownHostsPath = f.Name()
			err = os.WriteFile(knownHostsPath, []byte(conf.Ssh.KnownHosts), 0600)
			if err != nil {
				return fmt.Errorf("error writing known_hosts tmp file: %w", err)
			}
		}
		hostKeyCallback, err := ssh.NewKnownHostsCallback(knownHostsPath)
		if err != nil {
			return fmt.Errorf("failed to set up host keys: %w", err)
		}
		hostKeyCallbackHelper := ssh.HostKeyCallbackHelper{
			HostKeyCallback: hostKeyCallback,
		}
		sshKeys.HostKeyCallbackHelper = hostKeyCallbackHelper
		auth = sshKeys
	}

	// Set up the destination remote.
	dst, err := stagingRepo.CreateRemote(&config.RemoteConfig{
		Name: dstRemoteName,
		URLs: []string{conf.DstRepo},
	})
	if err != nil {
		return fmt.Errorf("failed configuring destination remote: %w", err)
	}

	logger.Info("Pushing to destination...")
	err = dst.Push(&git.PushOptions{
		RemoteName: dstRemoteName,
		Auth:       auth,
		RefSpecs:   []config.RefSpec{"refs/*:refs/*"},
		Force:      true,
		Prune:      false, // https://github.com/go-git/go-git/issues/520
	})
	switch err {
	case nil:
		logger.Info("Successfully mirrored pushed to destination repository.")
	case git.NoErrAlreadyUpToDate:
		logger.Info("Destination already up to date.")
	default:
		return fmt.Errorf("failed to push to destination: %w", err)
	}

	// We can not use prune in git.Push due to an existing bug
	// https://github.com/go-git/go-git/issues/520 so we workaround it dealing
	// with the prunning with a separate push.
	dstRefs, err := dst.List(&git.ListOptions{
		Auth: auth,
	})
	if err != nil {
		return fmt.Errorf("failed to list the destination remote: %w", err)
	}
	deleteSpecs, err := extraSpecs(stagingRepo, dstRefs)
	if err != nil {
		return fmt.Errorf("failed to prune destination: %w", err)
	}
	if len(deleteSpecs) > 0 {
		logger.Info("Pruning the destination...")
		err := dst.Push(&git.PushOptions{
			RemoteName: dstRemoteName,
			Auth:       auth,
			RefSpecs:   deleteSpecs,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}
	}

	return nil
}

// DoMirror mirrors the source to the destination git repository based on the
// provided configuration. Special references (for example GitHub's
// refs/pull/*) are ignored.
func DoMirror(conf Config, logger *Logger) error {
	repo, err := setupStagingRepo(conf, logger)
	if err != nil {
		return err
	}
	// Do not push GitHub special references used for dealing with pull
	// requests.
	if err := filterOutRefs(repo, []string{refsFilterPrefix}); err != nil {
		return fmt.Errorf("failed to filter out the refs: %w", err)
	}
	if err := pushWithAuth(conf, logger, repo); err != nil {
		return err
	}

	return nil
}
