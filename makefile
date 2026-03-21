.PHONY: new-domain migrate-create migrate-up migrate-down migrate-version migrate-force migrate-drop

include .env
export

new-domain:
ifndef DOMAIN
	$(error DOMAIN is not set. Usage: make new-domain DOMAIN=<your-domain> [MESSAGING=false])
endif

	@echo "Creating new domain: $(DOMAIN)"
	@mkdir -p internal/domain/$(DOMAIN)/repository
	@mkdir -p internal/domain/$(DOMAIN)/dto
	@mkdir -p internal/domain/$(DOMAIN)/model
	@mkdir -p internal/domain/$(DOMAIN)/delivery
	@mkdir -p internal/domain/$(DOMAIN)/delivery/http
	@mkdir -p internal/domain/$(DOMAIN)/usecase

ifeq ($(MESSAGING),true)
	@mkdir -p internal/domain/$(DOMAIN)/gateway/messaging
	@mkdir -p internal/domain/$(DOMAIN)/delivery/messaging
	@echo "package messaging" > internal/domain/$(DOMAIN)/gateway/messaging/$(DOMAIN)_producer.go
	@echo "package messaging" > internal/domain/$(DOMAIN)/delivery/messaging/$(DOMAIN)_consumer.go
	@echo "package model" > internal/domain/$(DOMAIN)/model/$(DOMAIN)_event.go
# 	@mkdir -p internal/contract
# 	@mkdir -p pkg/injector
# 	@if [ ! -f internal/contract/event.go ]; then \
# 		echo "package contract\n\ntype Event interface {\n\tGetId() string\n}" > internal/contract/event.go; \
# 		echo "package contract\n\ntype ProducerContract interface {\n\tGetTopic() *string\n\tSend(event Event) error\n}" > internal/contract/producer.go; \
# 	fi
endif

	@echo "package repository" > internal/domain/$(DOMAIN)/repository/$(DOMAIN)_repository.go
	@echo "package repository" > internal/domain/$(DOMAIN)/repository/$(DOMAIN)_repository_impl.go

	@echo "package dto" > internal/domain/$(DOMAIN)/dto/$(DOMAIN).go

	@echo "package model" > internal/domain/$(DOMAIN)/model/$(DOMAIN).go

	@echo "package http" > internal/domain/$(DOMAIN)/delivery/http/route.go
	@echo "package http" > internal/domain/$(DOMAIN)/delivery/http/handler.go

	@echo "package usecase" > internal/domain/$(DOMAIN)/usecase/$(DOMAIN)_usecase.go
	@echo "package usecase" > internal/domain/$(DOMAIN)/usecase/$(DOMAIN)_usecase_impl.go

ifeq ($(MESSAGING),true)
	@echo "Domain $(DOMAIN) with messaging created successfully."
else
	@echo "Domain $(DOMAIN) created successfully."
endif

migrate-create:
ifndef NAME
	$(error NAME is not set. Usage: make migrate-create NAME=<migration_name>)
endif

ifndef DB_DRIVER
	$(error DB_DRIVER is not set. Make sure DB_DRIVER is defined in your .env file)
endif
	@echo "Creating migration: $(NAME)"
	@mkdir -p migrations/$(DB_DRIVER)
	migrate create -ext sql -dir ./migrations/$(DB_DRIVER) -seq $(NAME)
	@echo "Migration created in ./migrations/$(DB_DRIVER)/"

migrate-up:
	go run cmd/migration/main.go up

migrate-down:
	go run cmd/migration/main.go down

migrate-version:
	go run cmd/migration/main.go version

migrate-force:
ifndef VERSION
	$(error VERSION is not set. Usage: make migrate-force VERSION=<version>)
endif
	@echo "WARNING: You are about to force migration to version $(VERSION)."
	@read -p "Are you sure? [y/N]: " confirm; \
	if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	@read -p "This may cause data inconsistency. Type 'force' to confirm: " final; \
	if [ "$$final" != "force" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	go run cmd/migration/main.go force $(VERSION)

migrate-drop:
	@echo "WARNING: You are about to DROP ALL migrations. This action is IRREVERSIBLE."
	@read -p "Are you sure? [y/N]: " confirm; \
	if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	@read -p "Type 'drop' to confirm: " final; \
	if [ "$$final" != "drop" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	go run cmd/migration/main.go drop