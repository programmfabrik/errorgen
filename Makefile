all: build

build:
	go generate && go build

test: build
	 cd example && go generate && go test

.PHONY: test all build