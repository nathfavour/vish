.PHONY: build clean

build:
	mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/vish ./cmd/vish

clean:
	rm -rf bin
