SHELL := /bin/bash
cwd := $(shell pwd)

test:
	go test ./... -cover -coverprofile=coverage.coverprofile

coverage-examine: test
	go tool cover -html=coverage.coverprofile