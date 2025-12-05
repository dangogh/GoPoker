.PHONY: all test cover clean build hands mcp-server

all: build

build: hands mcp-server

hands:
	go build -o bin/hands ./cmd/hands

mcp-server:
	go build -o bin/gopoker-mcp-server ./cmd/mcp-server

test:
	go test -coverprofile=coverage.out ./...


cover: test
	go tool cover -html=coverage.out


clean:
	rm -rf bin/ coverage.out
