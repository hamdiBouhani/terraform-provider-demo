# terraform-provider-demo Makefile
BINARY_NAME := terraform-provider-demo
MODULE      := github.com/hamdiBouhani/terraform-provider-demo
GO          ?= go
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the provider binary
	$(GO) build -o $(BINARY_NAME) .

.PHONY: install
install: build ## Build and print the path to use in .terraformrc
	@echo ""
	@echo "Provider built at: $(CURDIR)/$(BINARY_NAME)"
	@echo "Add this to your .terraformrc:"
	@echo ""
	@echo "  provider_installation {"
	@echo "    dev_overrides {"
	@echo "      \"hamdiBouhani/demo\" = \"$(CURDIR)\""
	@echo "    }"
	@echo "    direct {}"
	@echo "  }"

.PHONY: fmt
fmt: ## Format Go code
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: test
test: ## Run tests
	$(GO) test ./... -v

.PHONY: tidy
tidy: ## Tidy go.mod / go.sum
	$(GO) mod tidy

.PHONY: check
check: fmt vet build ## Run fmt, vet, and build

.PHONY: clean
clean: ## Remove build artifacts
	rm -f $(BINARY_NAME)
	rm -rf examples/.terraform examples/.terraform.lock.hcl examples/terraform.tfstate*

.PHONY: plan
plan: build ## Build provider and run terraform plan in examples/
	cd examples && TF_CLI_CONFIG_FILE=../.terraformrc terraform init -upgrade && \
		TF_CLI_CONFIG_FILE=../.terraformrc terraform plan

.PHONY: apply
apply: build ## Build provider and run terraform apply in examples/
	cd examples && TF_CLI_CONFIG_FILE=../.terraformrc terraform apply -auto-approve

.PHONY: destroy
destroy: ## Destroy the example resources
	cd examples && TF_CLI_CONFIG_FILE=../.terraformrc terraform destroy -auto-approve