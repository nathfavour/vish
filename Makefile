.PHONY: build clean

build:
	mkdir -p bin
	go build -ldflags="-s -w" -o bin/vish ./cmd/vish

clean:
	rm -rf bin
