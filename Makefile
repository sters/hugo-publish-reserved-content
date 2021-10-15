
export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

TOOLS=$(shell cat tools/tools.go | egrep '^\s_ '  | awk '{ print $$2 }')

.PHONY: bootstrap-tools
bootstrap-tools:
	@echo "Installing: " $(TOOLS)
	@go install $(TOOLS)

.PHONY: run
run:
	go run cmd/command/main.go

.PHONY: lint
lint:
	$(GOBIN)/golangci-lint run -v ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: cover
cover:
	go test -v -cover -race ./...

.PHONY: tidy
tidy:
	go mod tidy
