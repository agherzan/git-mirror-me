// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package main

import (
	"testing"

	"github.com/agherzan/git-mirror-me-action"
	"github.com/google/go-cmp/cmp"
)

// TestParseArgs tests that basic command line parsing works as expected.
func TestParseArgs(t *testing.T) {
	t.Parallel()
	{
		// Test passing -source-repository.
		config, _, err := parseArgs("test", []string{"-source-repository=src"})
		if err != nil {
			t.Fatalf("setting src failed: %s", err)
		}
		if !cmp.Equal(*config, mirror.Config{SrcRepo: "src"}) {
			t.Fatalf("unexpected src value: %s", config.Pretty())
		}
	}
	{
		// Test passing -destination-repository.
		config, _, err := parseArgs("test",
			[]string{"-destination-repository=dst"})
		if err != nil {
			t.Fatalf("setting dst failed: %s", err)
		}
		if !cmp.Equal(*config, mirror.Config{DstRepo: "dst"}) {
			t.Fatalf("unexpected dst value: %s", config.Pretty())
		}
	}
	{
		// Test passing -ssh-known-hosts-path.
		config, _, err := parseArgs("test",
			[]string{"-ssh-known-hosts-path=file"})
		if err != nil {
			t.Fatalf("setting host key failed: %s", err)
		}
		if !cmp.Equal(*config, mirror.Config{
			SSH: mirror.SSHConf{
				KnownHostsPath: "file",
			},
		}) {
			t.Fatalf("unexpected host key value: %s", config.Pretty())
		}
	}
	{
		// Test passing -debug.
		config, _, err := parseArgs("test",
			[]string{"-debug"})
		if err != nil {
			t.Fatalf("setting debug failed: %s", err)
		}
		if !cmp.Equal(*config, mirror.Config{
			Debug: true,
		}) {
			t.Fatalf("unexpected debug value: %s", config.Pretty())
		}
	}
	{
		// Test passing invalid flag.
		_, _, err := parseArgs("test", []string{"-invalid-flag"})
		if err == nil {
			t.Fatal("invalid flag succeeded")
		}
	}
}
