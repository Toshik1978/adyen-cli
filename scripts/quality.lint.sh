#!/usr/bin/env bash

GOLANGCI_LINT_VERSION="1.56.1"

if ! [ -x "$(command -v golangci-lint)" ]; then
  curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v$GOLANGCI_LINT_VERSION
fi
golangci-lint run ./...
