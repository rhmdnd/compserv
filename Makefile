OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

MIGRATE_VERSION := v4.15.2
GOLANGCI_LINT_VERSION := v1.46.2
KUBECTL_VERSION := v1.24.1

BUILDS_DIR := builds
TOOLS_DIR := tools
MIGRATE?=
KUBECTL = ./$(TOOLS_DIR)/kubectl

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
	go test -v ./pkg/...
	go test -v ./cmd/...

.PHONY: test-migrate
test-migrate: $(TOOLS_DIR)/migrate
	MIGRATE=$(MIGRATE) migrations/test.sh

.PHONY: test-database-integration
test-database-integration:
	./utils/run_integration_tests.sh

.PHONY: verify
verify: verify-go-lint

# Find all bash scripts by relying on the file extension and pass them to the
# linter, but ignore anything in vendor/.
.PHONY: bash-lint
bash-lint:
	shellcheck -x $(shell find . -type f -name '*.sh' -not -path "./vendor/*")

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
deploy: $(TOOLS_DIR)/kubectl
	$(KUBECTL) apply -k kustomize

$(TOOLS_DIR)/kubectl: $(TOOLS_DIR)
# Check if tools/kubectl exists - if it does then the default value provided
# above will work.
ifeq (,$(wildcard $(KUBECTL)))
# If tools/kubectl doesn't exist, check if the binary exists somewhere else in
# the path and use that. Otherwise, if we get back an empty string here we need
# to download a copy of kubectl and put it in the tools/ directory.
ifeq (,$(shell which kubectl 2>/dev/null))
	@{ \
	set -e ;\
	curl -L --output $(KUBECTL) "https://dl.k8s.io/release/$(KUBECTL_VERSION)/bin/linux/amd64/kubectl" ;\
	chmod u+x $(KUBECTL) ;\
	}
else
KUBECTL = $(shell which kubectl)
endif
endif
