.PHONY: test
test:
	go test -cover ./...
	go vet ./...

.PHONY: cover
cover:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

.DEFAULT_GOAL := test
