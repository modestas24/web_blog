include .env

MIGRATE_PATH=./cmd/migrate/migrations
DB_URL=$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)
ARGUMENTS=$(filter-out $@,$(MAKECMDGOALS))


.PHONY: migrate-status
migrate-status:
	@goose $(DB_DRIVER) $(DB_URL) -dir $(MIGRATE_PATH) status

.PHONY: migrate-create
migrate-create:
	@goose $(DB_DRIVER) $(DB_URL) -dir $(MIGRATE_PATH) create $(ARGUMENTS) sql -no-versioning

.PHONY: migrate-up
migrate-up:
	@goose $(DB_DRIVER) $(DB_URL) -dir $(MIGRATE_PATH) up -no-versioning

.PHONY: migrate-down
migrate-down:
	@goose $(DB_DRIVER) $(DB_URL) -dir $(MIGRATE_PATH) down -no-versioning

.PHONY: migrate-reset
migrate-reset:
	@goose $(DB_DRIVER) $(DB_URL) -dir $(MIGRATE_PATH) reset -no-versioning

.PHONY: migrate-validate
migrate-validate:
	@goose $(DB_DRIVER) $(DB_URL) -dir $(MIGRATE_PATH) validate

.PHONY: migrate-fix
migrate-fix:
	@goose $(DB_DRIVER) $(DB_URL) -dir $(MIGRATE_PATH) fix

.PHONY: seed
seed:
	@go run ./cmd/migrate/seed/main.go

.PHONY: swag-init
swag-init:
	@swag init -g ./main/main.go -d cmd,internal && swag fmt