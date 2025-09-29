PKG := ./...
GO ?= go

.PHONY: all fmt vet lint test vuln tidy ci
all: fmt vet lint test vuln

fmt:
	$(GO) fmt $(PKG)

vet:
	$(GO) vet $(PKG)

lint:
	golangci-lint run

test:
	$(GO) test $(PKG) -race -coverprofile=coverage.out -covermode=atomic

vuln:
	govulncheck $(PKG)

tidy:
	$(GO) mod tidy

ci: tidy lint test vuln
