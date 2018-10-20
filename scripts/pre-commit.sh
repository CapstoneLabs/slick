#!/bin/sh

GOFMT_FILES=$(env GO111MODULE=on gofmt -l .)
if [ -n "${GOFMT_FILES}" ]; then
  printf >&2 'gofmt failed for the following files:\n%s\n\nplease run "gofmt -w ." on your changes before committing.\n' "${GOFMT_FILES}"
  exit 1
fi

GOVET_ERRORS=$(env GO111MODULE=on go tool vet *.go 2>&1)
if [ -n "${GOVET_ERRORS}" ]; then
  printf >&2 'go vet failed for the following reasons:\n%s\n\nplease run "go tool vet *.go" on your changes before committing.\n' "${GOVET_ERRORS}"
  exit 1
fi

if [ -z "${NOTEST}" ]; then
  env GO111MODULE=on go test -v -short ./...
fi