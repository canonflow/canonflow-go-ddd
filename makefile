.PHONY: new-domain

new-domain:
ifndef DOMAIN
	$(error DOMAIN is not set. Usage: make new-domain DOMAIN=<your-domain>)
endif

	@echo "Creating new domain: $(DOMAIN)"
	@mkdir -p internal/domain/$(DOMAIN)/repository
	@mkdir -p internal/domain/$(DOMAIN)/dto
	@mkdir -p internal/domain/$(DOMAIN)/model
	@mkdir -p internal/domain/$(DOMAIN)/delivery
	@mkdir -p internal/domain/$(DOMAIN)/delivery/http
	@mkdir -p internal/domain/$(DOMAIN)/usecase

	@echo "package repository" > internal/domain/$(DOMAIN)/repository/$(DOMAIN)_repository.go
	@echo "package repository" > internal/domain/$(DOMAIN)/repository/$(DOMAIN)_repository_impl.go

	@echo "package dto" > internal/domain/$(DOMAIN)/dto/$(DOMAIN).go

	@echo "package model" > internal/domain/$(DOMAIN)/model/$(DOMAIN).go

	@echo "package http" > internal/domain/$(DOMAIN)/delivery/http/route.go
	@echo "package http" > internal/domain/$(DOMAIN)/delivery/http/handler.go

	@echo "package usecase" > internal/domain/$(DOMAIN)/usecase/$(DOMAIN)_usecase.go
	@echo "package usecase" > internal/domain/$(DOMAIN)/usecase/$(DOMAIN)_usecase_impl.go

	@echo "Domain $(DOMAIN) created successfully."