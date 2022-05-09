# SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
#
# SPDX-License-Identifier: MIT

REUSE := $(shell command -v reuse 2> /dev/null)
GOLANGCI := $(shell command -v golangci-lint 2> /dev/null)

all: test testcover lint reuse build

test:
	go test -race -coverprofile=coverage.out ./...

testcover: test
	go tool cover -func=coverage.out

lint:
ifndef GOLANGCI
	$(error "golangci is required to run lint checks on this projects")
endif
	golangci-lint run

reuse:
ifdef REUSE
	reuse lint
endif

build:
	go build -o bin/ ./cmd/git-mirror-me

.SILENT: clean

clean:
	rm -rf bin *.out
