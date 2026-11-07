SHELL := bash

DIST_DIR := build/bin
DIST_PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64
VERSION ?= $(shell git describe --tags --dirty 2>/dev/null | sed 's/^v//')
TARGET_GOOS ?= $(shell go env GOOS)
GOEXE := $(shell go env GOEXE)
GOLANGCI_LINT := build/tools/golangci-lint$(GOEXE)
GOVULNCHECK := build/tools/govulncheck$(GOEXE)
GO_TOOL := go tool -modfile=tools/go.mod
TEST_PACKAGES := ./cmd/... ./internal/... ./tests/unit ./tests/contract ./tests/golden
TEST_REPORTS := --junitfile build/reports/tests.xml --jsonfile build/reports/tests.json

.PHONY: build dist tools test race coverage integration lint fmt fmt-check modules workflows vuln check golden tidy

build:
	go build -trimpath -o bin/herdr-scope ./cmd/herdr-scope

dist:
	rm -rf $(DIST_DIR)
	mkdir -p $(DIST_DIR)
	for p in $(DIST_PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; \
		out=$(DIST_DIR)/herdr-scope-$$os-$$arch; \
		echo "$$out"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath \
			-ldflags "-s -w -X main.version=$(VERSION)" -o $$out ./cmd/herdr-scope || exit 1; \
	done

tools: $(GOLANGCI_LINT) $(GOVULNCHECK)
	$(GO_TOOL) gotestsum --version
	$(GO_TOOL) actionlint -version

$(GOLANGCI_LINT): .golangci-lint-version
	mkdir -p build/tools
	curl -fsSL --retry 3 https://raw.githubusercontent.com/golangci/golangci-lint/25c5369a93aacb8cccea0d46e1569ca4a194564d/install.sh -o build/tools/install-golangci-lint.sh
	sh build/tools/install-golangci-lint.sh -b build/tools "$$(tr -d '\r\n' < .golangci-lint-version)"

$(GOVULNCHECK): tools/go.mod tools/go.sum
	mkdir -p build/tools
	go -C tools build -o ../$(GOVULNCHECK) golang.org/x/vuln/cmd/govulncheck

test:
	mkdir -p build/reports
	$(GO_TOOL) gotestsum $(TEST_REPORTS) -- -count=1 -timeout=5m $(TEST_PACKAGES)

race:
	mkdir -p build/reports
	$(GO_TOOL) gotestsum $(TEST_REPORTS) -- -race -count=1 -timeout=5m $(TEST_PACKAGES)

coverage:
	mkdir -p build/reports
	$(GO_TOOL) gotestsum $(TEST_REPORTS) -- -race -count=1 -timeout=5m \
		-covermode=atomic -coverpkg=./... -coverprofile=build/reports/coverage.out $(TEST_PACKAGES)
	go tool cover -func=build/reports/coverage.out > build/reports/coverage.txt
	tail -n 1 build/reports/coverage.txt

integration:
	mkdir -p build/reports
	HERDR_SCOPE_INTEGRATION=1 $(GO_TOOL) gotestsum \
		--junitfile build/reports/integration.xml --jsonfile build/reports/integration.json \
		-- -race -count=1 -timeout=5m ./tests/integration

lint: $(GOLANGCI_LINT)
	GOOS=$(TARGET_GOOS) $(GOLANGCI_LINT) run --config .golangci.yaml ./...

fmt: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) fmt --config .golangci.yaml

fmt-check: $(GOLANGCI_LINT)
	@diff=$$($(GOLANGCI_LINT) fmt --config .golangci.yaml --diff) && test -z "$$diff" || { printf '%s\n' "$$diff"; exit 1; }

modules:
	go mod tidy -diff
	go mod verify
	go -C tools mod tidy -diff
	go -C tools mod verify

workflows:
	$(GO_TOOL) actionlint

vuln: $(GOVULNCHECK)
	GOOS=$(TARGET_GOOS) $(GOVULNCHECK) ./...

check: fmt-check modules workflows lint vuln coverage

golden:
	go test ./tests/golden -update -count=1

tidy:
	go mod tidy
	go -C tools mod tidy
