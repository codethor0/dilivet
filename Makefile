MODULE := github.com/codethor0/ml-dsa-debug-whitepaper
PKG_ALL := ./...
PKG_CLEAN := ./code/clean

BIN := mldsa
RELEASE_DIR := dist
PLATFORMS := darwin/arm64 darwin/amd64 linux/arm64 linux/amd64

.PHONY: all test bench run lint fmt vet tidy cover ci-local release clean

all: test

test:
	go test $(PKG_ALL) -v

bench:
	go test $(PKG_CLEAN) -bench Benchmark -benchmem -count=1

run:
	go build -v ./cmd/mldsa
	./$(BIN) -mode $(MODE) -msg "$(MSG)"

lint:
	@command -v golangci-lint >/dev/null 2>&1 || $(MAKE) install-lint
	golangci-lint run

install-lint:
	@echo "Installing golangci-lint..."
	@if command -v brew >/dev/null 2>&1; then \
		brew list golangci-lint >/dev/null 2>&1 || brew install golangci-lint; \
	else \
		GO111MODULE=on GOBIN="$(HOME)/.local/bin" go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.2; \
	fi

fmt:
	go fmt $(PKG_ALL)

vet:
	go vet $(PKG_ALL)

tidy:
	go mod tidy

cover:
	go test $(PKG_ALL) -coverprofile=coverage.out
	@echo "Coverage summary:" && go tool cover -func=coverage.out | tail -n 1

ci-local: fmt vet test

release:
	rm -rf $(RELEASE_DIR)
	mkdir -p $(RELEASE_DIR)
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; \
		echo "Building for $$os/$$arch..."; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(RELEASE_DIR)/$(BIN)-$$os-$$arch ./cmd/mldsa; \
	done
	@ls -lh $(RELEASE_DIR)

clean:
	rm -rf $(RELEASE_DIR) coverage.out
