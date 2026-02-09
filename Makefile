include .env
export

DB_URL=postgres://$(SHARED_DB_USER):$(SHARED_DB_PASSWORD)@localhost:$(SHARED_DB_PORT)/$(SHARED_DB)?sslmode=disable

migrate-up:
	migrate -path ./schema -database $(DB_URL) up

migrate-down:
	migrate -path ./schema -database $(DB_URL) down 1