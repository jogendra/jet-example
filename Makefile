# Go build variables
BINARY_NAME=sfmc-content-fetcher
GO_CGO_ENABLED := 1
GO_BUILD_FLAGS=-ldflags "-s -w"

GO := go
GOIMPORTS := goimports

export CGO_ENABLED=$(GO_CGO_ENABLED)
export GOFLAGS=$(GO_LDFLAGS)

ifneq (,$(wildcard ./.env.example))
    include .env.example
    export
endif

.PHONY: deps
deps:
	@ $(GO) mod download

.PHONY: fmt
fmt:
	@ $(GO) fmt ./...
	@ $(GOIMPORTS) -w .

.PHONY: tidy
tidy:
	@ $(GO) mod tidy

.PHONY: build
build:
	@ $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/cli

.PHONY: run
run: build
	@./$(BINARY_NAME)

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: docker-build
docker-build:
	docker build -t sfmc-content-fetcher .

.PHONY: docker-run
docker-run: docker-build
	docker run -it --env-file .env.example sfmc-content-fetcher

# Help target
help:
	@echo "Targets:"
	@echo "  deps         Download dependencies"
	@echo "  fmt          Formatting"
	@echo "  tidy         Tidy"
	@echo "  deps         Download dependencies"
	@echo "  build        Builds the application"
	@echo "  run          Runs the application"
	@echo "  test         Runs the tests"
	@echo "  clean        Cleans the build binaries"
	@echo "  docker-build Builds the Docker image"
	@echo "  docker-run   Runs the application in a Docker container"
	@echo "  help         Displays this help message"