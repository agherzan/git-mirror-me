// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"encoding/json"
	"errors"
	"path"
)

// SshConf structure defines SSH configuration used for git authentication over
// SSH.
type SshConf struct {
	PrivateKey     string
	KnownHosts     string
	KnownHostsPath string
}

// Config structure provides all the configuration need for the tool to perform
// its operations. It can be populated via a CLI component.
type Config struct {
	SrcRepo string
	DstRepo string
	Ssh     SshConf
	Debug   bool
}

// GetSshKey is the getter function for the private SSH key from a
// configuration struct.
func (conf Config) GetSshKey() string {
	return conf.Ssh.PrivateKey
}

// SetSshKey is the setter function for the private SSH key from a
// configuration struct.
func (conf *Config) SetSshKey(key string) {
	conf.Ssh.PrivateKey = key
}

// GetKnownHosts is the getter function for the public host key from a
// configuration struct.
func (conf Config) GetKnownHosts() string {
	return conf.Ssh.KnownHosts
}

// GetKnownHosts is the setter function for the public host key from a
// configuration struct.
func (conf *Config) SetKnownHosts(key string) {
	conf.Ssh.KnownHosts = key
}

// GetKnownHostsPath is the getter function for the public host key by file
// path from a configuration struct.
func (conf Config) GetKnownHostsPath() string {
	return conf.Ssh.KnownHostsPath
}

// GetKnownHostsPath is the setter function for the public host key by file
// path from a configuration struct.
func (conf *Config) SetKnownHostsPath(file string) {
	conf.Ssh.KnownHostsPath = file
}

// Pretty provides a string representation of the configuration structure. It
// does that by making sure sensitive information is masked using a hash
// function - e.g. the SSH private key.
func (conf Config) Pretty() string {
	// It's important to have a passed by value conf struct to not have the
	// masking affect the struct's actual values.
	conf.SetSshKey(mask(conf.Ssh.PrivateKey))
	conf.SetKnownHosts(mask(conf.Ssh.KnownHosts))
	out, _ := json.MarshalIndent(conf, "", "\t")

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
			conf.SrcRepo = path.Join(url, repo)
		}
	}

	// Fallback to environment variables for the destination repository value.
	if len(conf.DstRepo) == 0 {
		conf.DstRepo = env["GMM_DST_REPO"]
	}

	conf.Ssh.PrivateKey = env["GMM_SSH_PRIVATE_KEY"]
	conf.Ssh.KnownHosts = env["GMM_SSH_KNOWN_HOSTS"]
}

// Validate provides the logic of validating a configuration.
func (conf Config) Validate(logger *Logger) error {
	if len(conf.SrcRepo) == 0 {
		return errors.New("no source repository provided")
	}
	logger.Info("Source repository:", conf.SrcRepo, ".")

	if len(conf.DstRepo) == 0 {
		return errors.New("no destination repository provided")
	}
	logger.Info("Destination repository:", conf.DstRepo, ".")

	if len(conf.GetSshKey()) == 0 {
		logger.Warn("Tool configured with no authentication.")
	} else {
		if len(conf.GetKnownHosts()) != 0 &&
			len(conf.GetKnownHostsPath()) != 0 {
			return errors.New("host public keys provided via both file path " +
				"and content")
		} else if len(conf.GetKnownHosts()) == 0 &&
			len(conf.GetKnownHostsPath()) == 0 {
			return errors.New("SSH authentication requires host public keys")
		}
	}

	return nil
}
