LOCAL_BIN:=$(CURDIR)/bin

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

start: ### Start service with integration tests
	chmod +x start.sh
	./start.sh
.PHONY: start

compose-up: ### Start main service
	docker-compose up --build -d banknote_exchange && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

install-all: ### Install all tools
	make install-linter
	make install-swag
.PHONY: install-all

install-linter: ### Install golangci-lint
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2
.PHONY: install-linter

lint: ### Check by golangci linter
	$(LOCAL_BIN)/golangci-lint run
.PHONY: lint

test: ### Run unit tests
	go test -v ./internal/...
.PHONY: test

install-swag: ### Install swag
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@v1.16.3

swag: ### Generate swag docs
	$(LOCAL_BIN)/swag init -g cmd/app/main.go
.PHONY: swag