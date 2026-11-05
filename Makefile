-include .env

DATABASE_URL?=SA
MIGRATIONS_DIR=migration

.PHONY: migrate-down migrate-up migrate-new migrate-force migrate-fix migrate-status

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-new:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)

migrate-fix:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $$(migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version 2>&1 | grep -oE '^[0-9]+' || echo 0)

migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version