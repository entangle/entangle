PACKAGES := \
	token \
	source \
	parser \
	lexer \
	declarations \
	utils \
	term
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
