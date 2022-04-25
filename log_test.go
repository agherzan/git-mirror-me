// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

// osExitMock is used as a mock function for is.Exit(int).
func osExitMock(exitCode int) {
	return
}

// TestNewLogger tests the NewLogger function.
func TestNewLogger(t *testing.T) {
	if logger := NewLogger(os.Stderr); logger == nil {
		t.Fatal("failed to create a logger on stderr")
	}
}

// TestDebug checks debug logging.
func TestDebug(t *testing.T) {
	{
		// Test debug mode on.
		var b bytes.Buffer
		logger := NewLogger(&b)

		logger.Debug(true, "msg")
		if output := b.String(); output != "[DEBUG]: msg\n" {
			t.Fatalf("unexpected debug on output: %s", output)
		}
	}
	{
		// Test debug mode off.
		var b bytes.Buffer
		logger := NewLogger(&b)

		logger.Debug(false, "msg")
		if output := b.String(); output != "" {
			t.Fatalf("unexpected debug off output: %s", output)
		}
	}
}

// TestInfo checks info logging.
func TestInfo(t *testing.T) {
	var b bytes.Buffer
	logger := NewLogger(&b)

	logger.Info("msg")
	if output := b.String(); output != "[INFO ]: msg\n" {
		t.Fatalf("unexpected info output: %s", output)
	}
}

// TestWarn checks warn logging.
func TestWarn(t *testing.T) {
	var b bytes.Buffer
	logger := NewLogger(&b)

	logger.Warn("msg")
	if output := b.String(); output != "[WARN ]: msg\n" {
		t.Fatalf("unexpected warn output: %s", output)
	}
}

// TestError checks error logging.
func TestError(t *testing.T) {
	var b bytes.Buffer
	logger := NewLogger(&b)

	logger.Error("msg")
	if output := b.String(); output != "[ERROR]: msg\n" {
		t.Fatalf("unexpected error output: %s", output)
	}
}

// TestFatalMock checks fatal logging using a mocked exit function.
func TestFatalMock(t *testing.T) {
	var b bytes.Buffer
	logger := NewLogger(&b)

	origOsExit := osExit
	osExit = osExitMock
	defer func() { osExit = origOsExit }()

	logger.Fatal("msg")
	if output := b.String(); output != "[FATAL]: msg\n" {
		t.Fatalf("unexpected fatal output: %s", output)
	}
}

// TestFatal checks fatal logging.
func TestFatal(t *testing.T) {
	if os.Getenv("RUNTESTFATAL") == "1" {
		var b bytes.Buffer
		logger := NewLogger(&b)
		logger.Fatal("msg")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "RUNTESTFATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("fatal log didn't fail the running process")
}
