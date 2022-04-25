<!--
SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>

SPDX-License-Identifier: MIT
-->

# Git Mirror-Me (`GMm`)

This CLI tool provides the ability to mirror a repository to any other git
repository with optional SSH authentication. For example, it can be used with
repositories on GitHub, GitLab, Bitbucket, etc.

Why "Me"? The name derives from the tool's "ability" to default the source
repository to the value computed from a GitHub action environment.

## Build

Use the provided `make` script.

## Tool configuration

The tool can be configured via CLI arguments and/or environment variables.
Run the tool in `help` mode (`git-mirror-me -h`) to check its full description.

### Arguments/Flags

#### `-source-repository`

* Sets the source repository for the mirror operation.
* Can also be set via environment variables.

#### `-destination-repository`

* Sets the destination repository for the mirror operation.
* Can also be set via environment variables.

#### `-ssh-known-hosts-path`

* Defines the path to the `known_hosts` file.
* This is an alternative to providing the host public keys via the
  `GMM_SSH_KNOWN_HOSTS` environment variable (see below).

#### `-debug`

* Runs the tool in debug mode.

### Environment variables

This tool uses `GMM_` as prefix for all the environment variables defined in
its scope. That doesn't include the ones prefixed by `GITHUB_` as they are
expected to be provided directly by the GitHub CI environment.

#### `GMM_SRC_REPO`, `GITHUB_SERVER_URL` and `GITHUB_REPOSITORY`

* The source repository can be provided in three ways, listed below in the
descending order of their precedence:
  * the `-source-repository` CLI argument
  * the `GMM_SRC_REPO` environment variable
  * using the `GITHUB_SERVER_URL` and `GITHUB_REPOSITORY` environment variables
    as `GITHUB_SERVER_URL/GITHUB_REPOSITORY`

#### `GMM_DEST_REPO`

* Sets the destination repository for the mirror operation.

#### `GMM_SSH_PRIVATE_KEY`

* The SSH private key used for SSH authentication during git push operation.
* Password protected SSH keys are not supported.
* When not defined, `git` operations will be executed without authentication.
* When defined, a host public key configuration is required.

#### `GMM_SSH_KNOWN_HOSTS`

* The hosts public keys used for host validation.
* The format needs to be based on the`known_hosts` file.

## Tests and Checks

Use the provided `make` script. For tests, a `tests` target is provided: `make
tests`.

Extra checks can be executed using with `make checks`. This make target has
additional but optional host dependencies on:

* [ineffassign](https://github.com/gordonklaus/ineffassign) -  Detect
  ineffectual assignments in Go code.
* [gocyclo](https://github.com/fzipp/gocyclo) -  Calculate cyclomatic
  complexities of functions in Go source code.
* [misspell](https://github.com/client9/misspell) - Correct commonly misspelled
  English words in source files.

## Contributing

Contributions are more than welcome. You can send patches using [GitHub pull
requests](https://github.com/agherzan/git-mirror-me/pulls).

## Maintainers

* Andrei Gherzan `<andrei at gherzan.com>`

## LICENSE

This repository is [reuse](https://reuse.software/) compliant and it is
released under the [MIT](COPYING.MIT) license.

For convenience, a `make` target is provided to validate this compliance: `make
reuse`. This `make` target has an additional but optional host dependency
on [reuse](https://github.com/fsfe/reuse-tool).
