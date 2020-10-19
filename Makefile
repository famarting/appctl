
appctl: go.mod $(shell find -type f -name '*.go')
	go build -o appctl

install: appctl
	cp appctl $(HOME)/bin

.PHONY: install