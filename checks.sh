#!/usr/bin/env bash

# TODO: Choose a real task runner

set -ex

# Lint
golint -set_exit_status ./...

# Vet
go vet ./...

# Unit Tests
go test -v -cover ./reach/...

# Build
go build -a -v -tags netgo -o "/dev/null"

set +ex

printf "\n\nALL CHECKS PASS!\n"
