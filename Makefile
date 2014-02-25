PACKAGES := \
	./token \
	./source \
	./parser \
	./lexer \
	./declarations \
	./utils
SOURCE := $(wildcard *.go $(addsuffix /*.go, $(PACKAGES)))

export GOPATH=$(shell pwd)

all: entangle

entangle: $(SOURCE)
	@go build -v -o entangle .

test: all
	@go test -v $(PACKAGES)

format:
	@gofmt -l -w $(SOURCE)

clean:
	@rm -rf entangle

.PHONY: test clean
