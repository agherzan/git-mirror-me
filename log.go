// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"io"
	"log"
	"os"
)

// Used for moking exit function when running tests.
var osExit = os.Exit

// Logger structure provides per log level log.Logger.
type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
	fatal   *log.Logger
}

// NewLogger returns a new Logger that will use the passed io.Writer.
func NewLogger(o io.Writer) *Logger {
	var l Logger
	l.debug = log.New(o, "[DEBUG]: ", 0)
	l.info = log.New(o, "[INFO ]: ", 0)
	l.warning = log.New(o, "[WARN ]: ", 0)
	l.err = log.New(o, "[ERROR]: ", 0)
	l.fatal = log.New(o, "[FATAL]: ", 0)
	return &l
}

// Debug is printing a log message using the debug logger when debug mode is
// enabled.
func (l Logger) Debug(debugMode bool, v ...any) {
	if debugMode {
		l.debug.Println(v...)
	}
}

// Info is printing a log message using the info logger.
func (l Logger) Info(v ...any) {
	l.info.Println(v...)
}

// Warn is printing a log message using the warning logger.
func (l Logger) Warn(v ...any) {
	l.warning.Println(v...)
}

// Error is printing a log message using the err logger.
func (l Logger) Error(v ...any) {
	l.err.Println(v...)
}

// Fatal is printing a log message using the fatal logger followed by an
// os.Exit(1).
func (l Logger) Fatal(v ...any) {
	// We avoid Fatalln because we want to have the ability of mocking the exit
	// function for testing purposes.
	l.fatal.Println(v...)
	osExit(1)
}
