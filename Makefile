-include .env
.SILENT:

DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

tidy:
	@echo "Go Mod Tidy"
	@go mod tidy
	@go mod vendor

run:
	@go run cmd/main.go

# database migrations
# 
# create migration name=users
create-migration:
	@migrate create -ext sql -dir migrations -seq $(name)
#
# up all migrations
migrateup:
	@migrate -path migrations -database "$(DB_URL)" -verbose up
#
# up migration last one
migrateup1:
	@migrate -path migrations -database "$(DB_URL)" -verbose up 1
#
# down migrations all
migratedown:
	@migrate -path migrations -database "$(DB_URL)" -verbose down
#
# down the migration last
migratedown1:
	@migrate -path migrations -database "$(DB_URL)" -verbose down 1
