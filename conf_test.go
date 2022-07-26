// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"os"
	"testing"
)

const (
	testPrivateKey     = "setkey"
	testKnownHostsPath = "setkeypath"
)

// TestSetGetSSHKey tests the getter and setter for the SSH private key.
func TestSetGetSSHKey(t *testing.T) {
	t.Parallel()

	config := Config{
		SSH: SSHConf{
			PrivateKey: "key",
		},
	}

	config.SetSSHKey(testPrivateKey)

	if config.SSH.PrivateKey != testPrivateKey {
		t.Fatal("ssh key setter failed")
	}

	if config.GetSSHKey() != testPrivateKey {
		t.Fatal("ssh key getter failed")
	}
}

// TestSetGetKnownHosts tests the getter and setter for the host public key.
func TestSetGetKnownHosts(t *testing.T) {
	t.Parallel()

	config := Config{
		SSH: SSHConf{
			KnownHosts: "key",
		},
	}

	config.SetKnownHosts(testPrivateKey)

	if config.SSH.KnownHosts != testPrivateKey {
		t.Fatal("host key (by value) setter failed")
	}

	if config.GetKnownHosts() != testPrivateKey {
		t.Fatal("host key (by value) getter failed")
	}
}

// TestSetGetKnownHostsPath tests the getter and setter for the host public
// key provided by file path.
func TestSetGetKnownHostsPath(t *testing.T) {
	t.Parallel()

	config := Config{
		SSH: SSHConf{
			KnownHostsPath: "keypath",
		},
	}

	config.SetKnownHostsPath(testKnownHostsPath)

	if config.SSH.KnownHostsPath != testKnownHostsPath {
		t.Fatal("host key (by file path) setter failed")
	}

	if config.GetKnownHostsPath() != testKnownHostsPath {
		t.Fatal("host key (by file path) getter failed")
	}
}

// TestPretty tests the pretty output of a configuration structure.
func TestPretty(t *testing.T) {
	t.Parallel()

	// This also verifies that the sensitive fields are masked.
	out := Config{
		SrcRepo: "src",
		DstRepo: "dst",
		SSH: SSHConf{
			PrivateKey:     "key",
			KnownHosts:     "khkey",
			KnownHostsPath: "khpath",
		},
		Debug: true,
	}.Pretty()
	expectedOut := `{
	"SrcRepo": "src",
	"DstRepo": "dst",
	"SSH": {
		"PrivateKey": "2c70e12b7a0646f92279f427c7b38e7334d8e5389cff167a1dc30e73f826b683",
		"KnownHosts": "b3f1ba1ea27e621a8cab09c9e601097fd84c3c438dee43d9ee7b0efe8cfd0ecd",
		"KnownHostsPath": "khpath"
	},
	"Debug": true
}`

	if out != expectedOut {
		t.Fatalf("unexpected Pretty(): %s", out)
	}
}

// TestProcessEnv tests that the environment variables are used as expected.
func TestProcessEnv(t *testing.T) {
	t.Parallel()

	// No need for logs.
	devnull, _ := os.Open(os.DevNull)
	defer devnull.Close()
	logger := NewLogger(devnull)

	{
		// Source repository can be set from an environment variable.
		conf := Config{}
		env := map[string]string{
			"GMM_SRC_REPO": "srcenv",
		}
		conf.ProcessEnv(logger, env)
		if conf.SrcRepo != "srcenv" {
			t.Fatal("failed setting source repository from an env variable")
		}
	}
	{
		// Source can be set by GitHub CI environment variables.
		conf := Config{}
		env := map[string]string{
			"GITHUB_SERVER_URL": "http://github.com",
			"GITHUB_REPOSITORY": "user/repo",
		}
		conf.ProcessEnv(logger, env)
		if conf.SrcRepo != "http://github.com/user/repo" {
			t.Fatal("failed setting source repository from GitHub env variables")
		}
	}
	{
		// Environment variables don't override existing source configuration.
		conf := Config{SrcRepo: "src"}
		env := map[string]string{
			"GMM_SRC_REPO":      "srcenv",
			"GITHUB_SERVER_URL": "foo",
			"GITHUB_REPOSITORY": "bar",
		}
		conf.ProcessEnv(logger, env)
		if conf.SrcRepo != "src" {
			t.Fatal("env variables override existing configuration for the src " +
				"repo")
		}
	}
	{
		// GMM_SRC_REPO env variable has higher priority than GitHub CI
		// environment variables.
		conf := Config{}
		env := map[string]string{
			"GMM_SRC_REPO":      "srcenv",
			"GITHUB_SERVER_URL": "foo",
			"GITHUB_REPOSITORY": "bar",
		}
		conf.ProcessEnv(logger, env)
		if conf.SrcRepo != "srcenv" {
			t.Fatal("source repository priority for env variables failed")
		}
	}
	{
		// Destination repository can be set from an environment variable.
		conf := Config{}
		env := map[string]string{
			"GMM_DST_REPO": "dstenv",
		}
		conf.ProcessEnv(logger, env)
		if conf.DstRepo != "dstenv" {
			t.Fatal("failed setting destination repository from an env " +
				"variable")
		}
	}
	{
		// Environment variables don't override existing source configuration.
		conf := Config{DstRepo: "dst"}
		env := map[string]string{
			"GMM_DST_REPO": "dstenv",
		}
		conf.ProcessEnv(logger, env)
		if conf.DstRepo != "dst" {
			t.Fatal("env variables override existing configuration for the " +
				"destination repository")
		}
	}
	{
		// Populating the SSH private key from an environment variable.
		conf := Config{}
		env := map[string]string{
			"GMM_SSH_PRIVATE_KEY": "keyenv",
		}
		conf.ProcessEnv(logger, env)
		if conf.SSH.PrivateKey != "keyenv" {
			t.Fatal("failed setting SSH private key from an env variable")
		}
	}
	{
		// Populating the host public key from an environment variable.
		conf := Config{}
		env := map[string]string{
			"GMM_SSH_KNOWN_HOSTS": "khkeyenv",
		}
		conf.ProcessEnv(logger, env)
		if conf.SSH.KnownHosts != "khkeyenv" {
			t.Fatal("failed setting host key from an env variable")
		}
	}

	// Debug mode via an environment variable tests.
	{
		conf := Config{Debug: false}
		env := map[string]string{
			"GMM_DEBUG": "1",
		}
		conf.ProcessEnv(logger, env)
		if conf.Debug != true {
			t.Fatal("failed debug mode env variable test")
		}
	}
	{
		conf := Config{Debug: false}
		env := map[string]string{
			"GMM_DEBUG": "",
		}
		conf.ProcessEnv(logger, env)
		if conf.Debug != false {
			t.Fatal("failed debug mode env variable test")
		}
	}
	{
		conf := Config{Debug: false}
		env := map[string]string{}
		conf.ProcessEnv(logger, env)
		if conf.Debug != false {
			t.Fatal("failed debug mode env variable test")
		}
	}
	{
		conf := Config{Debug: true}
		env := map[string]string{
			"GMM_DEBUG": "1",
		}
		conf.ProcessEnv(logger, env)
		if conf.Debug != true {
			t.Fatal("failed debug mode env variable test")
		}
	}
	{
		conf := Config{Debug: true}
		env := map[string]string{
			"GMM_DEBUG": "",
		}
		conf.ProcessEnv(logger, env)
		if conf.Debug != true {
			t.Fatal("failed debug mode env variable test")
		}
	}
	{
		conf := Config{Debug: true}
		env := map[string]string{}
		conf.ProcessEnv(logger, env)
		if conf.Debug != true {
			t.Fatal("failed debug mode env variable test")
		}
	}
}

// TestValidate tests various valid/invalid configurations.
func TestValidate(t *testing.T) {
	t.Parallel()

	// No need for logs.
	devnull, _ := os.Open(os.DevNull)
	defer devnull.Close()
	logger := NewLogger(devnull)

	{
		// Source and destination repositories are required.
		conf := Config{}
		if err := conf.Validate(logger); err == nil {
			t.Fatal("source repository was not required")
		}
		conf.SrcRepo = "src"
		if err := conf.Validate(logger); err == nil {
			t.Fatal("destination repository was not required")
		}
		conf = Config{
			SrcRepo: "src",
			DstRepo: "dst",
		}
		if err := conf.Validate(logger); err != nil {
			// This also tests that no authentication is allowed.
			t.Fatal("src and dst defined but function failed")
		}
	}
	{
		// SSH private key configration requires host key configuration.
		conf := Config{
			SrcRepo: "src",
			DstRepo: "dst",
			SSH: SSHConf{
				PrivateKey: "key",
			},
		}
		if err := conf.Validate(logger); err == nil {
			t.Fatal("SSH key configuration didn't require host key " +
				"configuration")
		}
	}
	{
		// Test that Validate works when both SSH key and host public key are
		// provided.
		conf := Config{
			SrcRepo: "src",
			DstRepo: "dst",
			SSH: SSHConf{
				PrivateKey: "key",
				KnownHosts: "khkey",
			},
		}
		if err := conf.Validate(logger); err != nil {
			t.Fatal("ssh key and host key defined but function failed")
		}
	}
	{
		// Host key configurations as value and file path are mutually
		// exclusive.
		conf := Config{
			SrcRepo: "src",
			DstRepo: "dst",
			SSH: SSHConf{
				PrivateKey:     "key",
				KnownHosts:     "khkey",
				KnownHostsPath: "khpath",
			},
		}
		if err := conf.Validate(logger); err == nil {
			t.Fatal("host key provided as value and file path was allowed")
		}
	}
	{
		// Allow host key provided by value.
		conf := Config{
			SrcRepo: "src",
			DstRepo: "dst",
			SSH: SSHConf{
				PrivateKey: "key",
				KnownHosts: "khkey",
			},
		}
		if err := conf.Validate(logger); err != nil {
			t.Fatal("host key provided by value was not allowed")
		}
	}
	{
		// Allow host key provided by file path.
		conf := Config{
			SrcRepo: "src",
			DstRepo: "dst",
			SSH: SSHConf{
				PrivateKey:     "key",
				KnownHostsPath: "khpath",
			},
		}
		if err := conf.Validate(logger); err != nil {
			t.Fatal("host key provided by file path was not allowed")
		}
	}
}
