OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

MIGRATE_VERSION := v4.15.2

BUILDS_DIR := builds
TOOLS_DIR := tools

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

.PHONY: $(TOOLS_DIR)/migrate
$(TOOLS_DIR)/migrate: $(TOOLS_DIR)
	curl -sSL https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.$(OS)-$(ARCH).tar.gz | tar xz migrate -O > $(TOOLS_DIR)/migrate
	chmod u+x $@
