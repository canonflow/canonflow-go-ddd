.PHONY: new-domain

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