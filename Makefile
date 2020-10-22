all: build

build:
	go generate && go build

clean:
	rm -f resources.go

test: build
	 cd example && go generate && go test

.PHONY: test all build clean