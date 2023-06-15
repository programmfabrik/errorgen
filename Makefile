all: build

build:
	go build

clean:
	rm -f resources.go example/example.go errorgen

test: build
	 cd example && go generate && go test

.PHONY: test all build clean