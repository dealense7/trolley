# Makefile for Goose migrations using config/config.yml

# Use yq to extract DB config from YAML
DB_USER := $(shell yq e '.db.user' config/config.yaml)
DB_PASS := $(shell yq e '.db.password' config/config.yaml)
DB_HOST := $(shell yq e '.db.host' config/config.yaml)
DB_PORT := $(shell yq e '.db.port' config/config.yaml)
DB_NAME := $(shell yq e '.db.name' config/config.yaml)

# Construct MySQL DSN
DB_URL := $(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

MIGRATIONS_DIR := ./migrations

.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  make migrate              # Run all pending migrations"
	@echo "  make rollback             # Rollback last migration"
	@echo "  make status               # Show current migration status"
	@echo "  make create_migration NAME=<name>   # Create new empty migration file"

# Run all pending migrations
.PHONY: migrate
migrate:
	echo "$(DB_URL)"
	goose -dir $(MIGRATIONS_DIR) mysql "$(DB_URL)" up

# Run all pending migrations
.PHONY: migrate_fresh
migrate_fresh:
	@echo "Rolling back all migrations..."
	goose -dir $(MIGRATIONS_DIR) mysql "$(DB_URL)" down-to 0
	@echo "Running all migrations..."
	goose -dir $(MIGRATIONS_DIR) mysql "$(DB_URL)" up

# Rollback last migration
.PHONY: rollback
rollback:
	goose -dir $(MIGRATIONS_DIR) mysql "$(DB_URL)" down

# Show current migration status
.PHONY: status
status:
	goose -dir $(MIGRATIONS_DIR) mysql "$(DB_URL)" status

# Create a new migration
# Usage: make create_migration NAME=create_users_table
.PHONY: create_migration
create_migration:
ifndef NAME
	$(error NAME is not set. Usage: make create NAME=create_table)
endif
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql
	@echo "Migration file created in $(MIGRATIONS_DIR)"