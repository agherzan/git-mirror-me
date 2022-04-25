# SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
#
# SPDX-License-Identifier: MIT

INEFFASSIGN := $(shell command -v ineffassign 2> /dev/null)
GOCYCLO := $(shell command -v gocyclo 2> /dev/null)
MISSPELL := $(shell command -v misspell 2> /dev/null)
REUSE := $(shell command -v reuse 2> /dev/null)

all: testcover checks build

test:
	go test -race -coverprofile=coverage.out ./...

testcover: test
	go tool cover -func=coverage.out

checks: vet fmt reuse
ifdef INEFFASSIGN
	ineffassign ./...
endif
ifdef GOCYCLO
	gocyclo -ignore "_test" -over 15 .
endif
ifdef MISSPELL
	misspell .
endif

vet:
	go vet ./...

fmt:
	@echo "gofmt checks..."; fails="$$(go list -f '{{.Dir}}' ./... | \
		xargs -L1 gofmt -l | grep -v /vendor/)"; \
		if [ -n "$$fails" ]; then \
		for fail in $$fails; do \
			gofmt -s -d $$fail; \
		done; \
		false; \
	fi;
	@echo "gofmt checks done."

reuse:
ifdef REUSE
	reuse lint
endif

build:
	go build -o bin/ ./cmd/git-mirror-me

.SILENT: clean

clean:
	rm -rf bin *.out
