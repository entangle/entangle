PACKAGES := \
	entangle/token \
	entangle/source \
	entangle/parser \
	entangle/lexer \
	entangle/declarations \
	entangle/utils \
	entangle/term
SOURCE := $(wildcard $(addsuffix /*.go, $(addprefix src/, $(PACKAGES))))

export GOPATH=$(shell pwd)

all: entangle

entangle: $(SOURCE)
	@go build -v -o bin/entangle cmds/entangle

test: all
	@go test -v $(PACKAGES)

format:
	@gofmt -l -w $(SOURCE)

clean:
	@rm -rf bin pkg

.PHONY: test clean
