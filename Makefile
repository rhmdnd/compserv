OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

MIGRATE_VERSION := v4.15.2
GOLANGCI_LINT_VERSION := v1.46.2

BUILDS_DIR := builds
TOOLS_DIR := tools
MIGRATE?=

.PHONY: $(BUILDS_DIR)
$(BUILDS_DIR):
	mkdir -p $(BUILDS_DIR)

.PHONY: $(TOOLS_DIR)
$(TOOLS_DIR):
	mkdir -p $(TOOLS_DIR)

.PHONY: build
build: $(BUILDS_DIR)
	go build -o $(BUILDS_DIR) cmd/compserv-server.go

.PHONY: test
test:
	go test -v ./...

.PHONY: test-migrate
test-migrate: $(TOOLS_DIR)/migrate
	MIGRATE=$(MIGRATE) migrations/test.sh

.PHONY: verify
verify: verify-go-lint

.PHONY: verify-go-lint
verify-go-lint: $(TOOLS_DIR)/golangci-lint ## Verify the golang code by linting
	# we use go 1.17 because golangci-lint still has issues with 1.18
	GL_DEBUG=gocritic $(TOOLS_DIR)/golangci-lint --go 1.17 run

MIGRATE = ./$(TOOLS_DIR)/migrate
$(TOOLS_DIR)/migrate: $(TOOLS_DIR) ## Download migrate locally if necessary.
ifeq (,$(wildcard $(MIGRATE)))
ifeq (,$(shell which migrate 2>/dev/null))
	@{ \
	set -e ;\
	curl -sSL https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.$(OS)-$(ARCH).tar.gz | tar xz migrate -O > $(MIGRATE) ;\
	chmod u+x $(MIGRATE) ;\
	}
else
MIGRATE = $(shell which migrate)
endif
endif

$(TOOLS_DIR)/golangci-lint:
	export \
		VERSION=$(GOLANGCI_LINT_VERSION) \
		URL=https://raw.githubusercontent.com/golangci/golangci-lint \
		BINDIR=$(TOOLS_DIR) && \
	curl -sfL $$URL/$$VERSION/install.sh | sh -s $$VERSION
	$(TOOLS_DIR)/golangci-lint version
	$(TOOLS_DIR)/golangci-lint linters

.PHONY: deploy
deploy:
	kubectl apply -k kustomize
