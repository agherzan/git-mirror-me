// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"encoding/json"
	"errors"
)

var (
	ErrNoSrc     = errors.New("no source repository provided")
	ErrNoDst     = errors.New("no destination repository provided")
	ErrNoHostKey = errors.New("SSH authentication requires host public keys")
	ErrHostKey   = errors.New("host public keys provided via both file path " +
		"and content")
)

// SSHConf structure defines SSH configuration used for git authentication over
// SSH.
type SSHConf struct {
	PrivateKey     string
	KnownHosts     string
	KnownHostsPath string
}

// Config structure provides all the configuration need for the tool to perform
// its operations. It can be populated via a CLI component.
type Config struct {
	SrcRepo string
	DstRepo string
	SSH     SSHConf
	Debug   bool
}

// GetSSHKey is the getter function for the private SSH key from a
// configuration struct.
func (conf Config) GetSSHKey() string {
	return conf.SSH.PrivateKey
}

// SetSSHKey is the setter function for the private SSH key from a
// configuration struct.
func (conf *Config) SetSSHKey(key string) {
	conf.SSH.PrivateKey = key
}

// GetKnownHosts is the getter function for the public host key from a
// configuration struct.
func (conf Config) GetKnownHosts() string {
	return conf.SSH.KnownHosts
}

// GetKnownHosts is the setter function for the public host key from a
// configuration struct.
func (conf *Config) SetKnownHosts(key string) {
	conf.SSH.KnownHosts = key
}

// GetKnownHostsPath is the getter function for the public host key by file
// path from a configuration struct.
func (conf Config) GetKnownHostsPath() string {
	return conf.SSH.KnownHostsPath
}

// GetKnownHostsPath is the setter function for the public host key by file
// path from a configuration struct.
func (conf *Config) SetKnownHostsPath(file string) {
	conf.SSH.KnownHostsPath = file
}

// Pretty provides a string representation of the configuration structure. It
// does that by making sure sensitive information is masked using a hash
// function - e.g. the SSH private key.
func (conf Config) Pretty() string {
	// It's important to have a passed by value conf struct to not have the
	// masking affect the struct's actual values.
	conf.SetSSHKey(mask(conf.SSH.PrivateKey))
	conf.SetKnownHosts(mask(conf.SSH.KnownHosts))

	out, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return ""
	}

	return string(out)
}

// ProcessEnv deals with environment configuration. It populates or overrides
// configuration based on a map that models environment variables.
func (conf *Config) ProcessEnv(logger *Logger, env map[string]string) {
	// Fallback to environment variables for the source repository value.
	if len(conf.SrcRepo) == 0 {
		if src, srcSet := env["GMM_SRC_REPO"]; srcSet {
			conf.SrcRepo = src
		} else {
			url := env["GITHUB_SERVER_URL"]
			repo := env["GITHUB_REPOSITORY"]
			conf.SrcRepo = url + "/" + repo
		}
	}

	// Fallback to environment variables for the destination repository value.
	if len(conf.DstRepo) == 0 {
		conf.DstRepo = env["GMM_DST_REPO"]
	}

	conf.SSH.PrivateKey = env["GMM_SSH_PRIVATE_KEY"]
	conf.SSH.KnownHosts = env["GMM_SSH_KNOWN_HOSTS"]
}

// Validate provides the logic of validating a configuration.
func (conf Config) Validate(logger *Logger) error {
	if len(conf.SrcRepo) == 0 {
		return ErrNoSrc
	}

	logger.Info("Source repository:", conf.SrcRepo, ".")

	if len(conf.DstRepo) == 0 {
		return ErrNoDst
	}

	logger.Info("Destination repository:", conf.DstRepo, ".")

	if len(conf.GetSSHKey()) == 0 {
		logger.Warn("Tool configured with no authentication.")
	} else {
		if len(conf.GetKnownHosts()) != 0 &&
			len(conf.GetKnownHostsPath()) != 0 {
			return ErrHostKey
		} else if len(conf.GetKnownHosts()) == 0 &&
			len(conf.GetKnownHostsPath()) == 0 {
			return ErrNoHostKey
		}
	}

	return nil
}
