help: ## Show this help.
	echo "help"

lint: ## init
	golangci-lint run ./...

test: ## tests
	go test ./...
