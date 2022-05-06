// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"flag"
	"fmt"
	"path"

	"github.com/agherzan/git-mirror-me-action"
)

// parseArgs returns a configuration structure initialised from parsing the
// 'arguments' string slice argument.
func parseArgs(progName string, arguments []string) (*mirror.Config, string, error) {
	var srcRepo, dstRepo, knownHostsPath string

	var debug bool

	var flagsOutput bytes.Buffer

	flags := flag.NewFlagSet(progName, flag.ContinueOnError)
	flags.SetOutput(&flagsOutput)
	flags.Usage = func() {
		output := flags.Output()
		fmt.Fprintf(output,
			`%s is a CLI tool that facilitates mirroring git repository.

CLI arguments/flags
`, path.Base(progName))

		flag.PrintDefaults()
		fmt.Fprintf(output,
			`
Environment variables
  GMM_SRC_REPO
  GITHUB_SERVER_URL
  GITHUB_REPOSITORY
    The source repository can be provided in three ways, listed below in the
    descending order of their precedence:
      * the '-source-repository' CLI flag
      * the 'GMM_SRC_REPO' environment variable
      * using the 'GITHUB_SERVER_URL' and 'GITHUB_REPOSITORY' environment
        variables as 'GITHUB_SERVER_URL/GITHUB_REPOSITORY'
  GMM_DEST_REPO
    Same as '-destination-repository' but overridden by the CLI argument.
  GMM_SSH_PRIVATE_KEY
    The SSH private key used for SSH authentication during git operations. When
    defined, a host public key configuration is required. See
    'GMM_SSH_KNOWN_HOSTS' and '-ssh-known-hosts-path'.
  GMM_SSH_KNOWN_HOSTS
    The host public keys used for host validation. The format needs to be based
    on the 'known_hosts' file. See
    http://man.openbsd.org/sshd#SSH_KNOWN_HOSTS_FILE_FORMAT
    for more information.
    This can't be used in conjunction with '-ssh-known-hosts-path'.
`)
	}
	flags.StringVar(&srcRepo, "source-repository", "",
		"The source repository for the mirroring operation.\nCan also be "+
			"set via environment variables.")
	flags.StringVar(&dstRepo, "destination-repository", "",
		"The destination repository for the mirroring operation.\nCan also "+
			"be set via environment variables.")
	flags.StringVar(&knownHostsPath, "ssh-known-hosts-path", "",
		"Defines the path to the 'known_hosts' file.\nThis is an alternative to "+
			"providing the host public keys via the\n'GMM_SSH_KNOWN_HOSTS' "+
			"environment variable.")
	flags.BoolVar(&debug, "debug", false, "Run this tool in debug mode.")

	if err := flags.Parse(arguments); err != nil {
		return nil, flagsOutput.String(), err
	}

	return &mirror.Config{
		SrcRepo: srcRepo,
		DstRepo: dstRepo,
		SSH: mirror.SSHConf{
			KnownHostsPath: knownHostsPath,
		},
		Debug: debug,
	}, flagsOutput.String(), nil
}
