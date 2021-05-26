SHELL := /bin/bash

.PHONY: test test-watch coverage-examine upgrade

test:
	go test ./... -cover -coverprofile=coverage.coverprofile

test-watch:
	ginkgo watch -r -trace

coverage-examine: test
	go tool cover -html=coverage.coverprofile

upgrade:
	go-mod-upgrade
