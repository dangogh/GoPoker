.PHONY: all test cover clean

test:
	go test -coverprofile=coverage.out ./...


cover: test
	go tool cover -html=coverage.out


clean:
	rm -rf coverage.out
