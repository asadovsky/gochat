export SHELL := /bin/bash -euo pipefail
export GOPATH := $(shell pwd)

src/code.google.com/p/go.net:
	go get code.google.com/p/go.net

src/code.google.com/p/go.text:
	go get code.google.com/p/go.text

serve: src/code.google.com/p/go.net src/code.google.com/p/go.text
	cd src/main && go run main.go

clean:
	rm -rf bin pkg

lint:
	tools/lint.sh

.PHONY: serve clean lint
