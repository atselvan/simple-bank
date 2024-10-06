.PHONY: dep test app-up app-down migrateup migratedown sqlc

# Load the environment variables from the config.env file
include config.env
export $(shell sed 's/=.*//' config.env)

dep:
	@echo "Install/Update dependencies"
	go get -u -t
	go mod tidy

test:
	@echo "Running tests"
	go test ./...

# Bring up all elements of the application
app-up:
	@echo "Starting the application"
	@docker-compose up -d

# Bring down all elements of the application
app-down:
	@echo "Stopping the application"
	@docker-compose down

# Initialize database migration scripts
migrate-init:
	@echo "Running migration init"
	@migrate create -ext sql -dir db/migration -seq init_schema

# Migrate the database up
migrate-up:
	@echo "Running migration up"
	@migrate -path db/migration -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

# Migrate the database down
migrate-down:
	@echo "Running migration down"
	@migrate -path db/migration -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

sqlc:
	@echo "Running sqlc"
	@sqlc generate
