GOLANGCI_LINT_VERSION := v2.4.0
GO_PATH := $(shell go env GOPATH)

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -race -count=1 ./...

.PHONY: cover
cover:
	go test -race -count=1 -coverprofile=coverage.out ./...

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: lint
lint: tools
	golangci-lint run ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: verify
verify:
	go mod verify

.PHONY: tools
tools:
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GO_PATH)/bin $(GOLANGCI_LINT_VERSION); \
	fi

.PHONY: clean
clean:
	rm -f coverage.out
