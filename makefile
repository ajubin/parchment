
.PHONY: help run deploy build

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	
run: ## Run the project locally
	@echo "Running Parchment locally..."
	@go run main.go
 
deploy: ## Deploy the project to rasbperry pi
	@./deploy.sh

build: ## Build the package as an executable
	@go build -o parchment main.go
