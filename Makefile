# The binary to build (just the basename).
BIN := puffer

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)

all: build

.PHONY: build

build: $(BIN)
	go build -o $(BIN) puffer.go

run: build
	sudo ./puffer watch

version:
	@echo $(VERSION)

clean:
	rm $(BIN)





