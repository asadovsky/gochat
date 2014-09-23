export SHELL := /bin/bash -euo pipefail
export GOPATH := $(shell pwd)

src/code.google.com/p/go.net:
	go get -d code.google.com/p/go.net/websocket

serve: src/code.google.com/p/go.net
	cd src/main && go run main.go

clean:
	rm -rf bin pkg

lint:
	tools/lint.sh

.PHONY: serve clean lint
