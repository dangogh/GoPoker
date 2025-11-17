.PHONY: test cover clean

test:
	go test ./... -coverprofile=coverage.out

cover: test
	go tool cover -html=coverage.out

clean:
	rm -f coverage.out
