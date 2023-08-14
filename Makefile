.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

run: ### Run
	go run .
.PHONY: run

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

test: ### run test
	go test -v -cover -race .
.PHONY: test