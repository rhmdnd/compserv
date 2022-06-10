OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

MIGRATE_VERSION := v4.15.2

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

.PHONY: $(TOOLS_DIR)/migrate
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
