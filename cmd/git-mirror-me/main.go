// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	mirror "github.com/agherzan/git-mirror-me"
)

func run(logger *mirror.Logger, env map[string]string, progName string, args []string) error {
	conf, output, err := parseArgs(progName, args)

	switch {
	case errors.Is(err, flag.ErrHelp):
		fmt.Fprintf(logger.GetOutput(), output)

		return nil
	case errors.Is(err, ErrVersion):
		debugInfo, _ := debug.ReadBuildInfo()
		fmt.Fprintf(logger.GetOutput(), debugInfo.String())

		return nil
	case err != nil:
		return fmt.Errorf("%w", err)
	}

	conf.ProcessEnv(logger, env)
	logger.Debug(conf.Debug, conf.Pretty())

	err = conf.Validate(logger)
	if err != nil {
		return fmt.Errorf("configuration failed: %w", err)
	}

	err = mirror.DoMirror(*conf, logger)
	if err != nil {
		return fmt.Errorf("mirror operation failed: %w", err)
	}

	return nil
}

func main() {
	// Keep the main function minimum as it is not covered by testing.
	logger := mirror.NewLogger(os.Stderr)

	env := make(map[string]string)

	envVars := []string{
		"GMM_SRC_REPO",
		"GITHUB_SERVER_URL",
		"GITHUB_REPOSITORY",
		"GMM_DEST_REPO",
		"GMM_SSH_PRIVATE_KEY",
		"GMM_SSH_KNOWN_HOSTS",
		"GMM_DEBUG",
	}

	for _, envVar := range envVars {
		if val, set := os.LookupEnv(envVar); set {
			env[envVar] = val
		}
	}

	if err := run(logger, env, os.Args[0], os.Args[1:]); err != nil {
		logger.Fatal(err)
	}
}
