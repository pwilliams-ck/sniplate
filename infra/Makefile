SNIPLATE_BINARY=sniplateApp

# Default values for the environment variables
DEFAULT_PORT=4200
DEFAULT_ENV=development
DEFAULT_TLS=false
DEFAULT_LOG=false
DEFAULT_DSN=postgres://postgres:password@postgres:5432/sniplate?sslmode=disable

## start: Stops docker-compose (if running), builds all projects and starts docker compose with flags.
start: build_sniplate
	@echo "Stopping docker images (if running...)"
	@PORT=$${PORT:-${DEFAULT_PORT}} \
	 ENV=$${ENV:-${DEFAULT_ENV}} \
	 TLS=$${TLS:-${DEFAULT_TLS}} \
	 LOG=$${LOG:-${DEFAULT_LOG}} \
	 DSN=$${DSN:-${DEFAULT_DSN}} \
	 docker-compose down
	@echo "Building (when required) and starting docker images..."
	@PORT=$${PORT:-${DEFAULT_PORT}} \
	 ENV=$${ENV:-${DEFAULT_ENV}} \
	 TLS=$${TLS:-${DEFAULT_TLS}} \
	 LOG=$${LOG:-${DEFAULT_LOG}} \
	 DSN=$${DSN:-${DEFAULT_DSN}} \
	 docker-compose up --build -d sniplate postgres
	@echo "Docker images built and started!"

## stop: Stop docker-compose.
stop:
	@echo "Stopping docker compose..."
	@PORT=$${PORT:-${DEFAULT_PORT}} \
	 ENV=$${ENV:-${DEFAULT_ENV}} \
	 TLS=$${TLS:-${DEFAULT_TLS}} \
	 LOG=$${LOG:-${DEFAULT_LOG}} \
	 DSN=$${DSN:-${DEFAULT_DSN}} \
	 docker-compose down
	@echo "Done!"

## migrate_up: Runs database migrations.
migrate_up:
	@echo "Running database migrations..."
	@PORT=$${PORT:-${DEFAULT_PORT}} \
	 ENV=$${ENV:-${DEFAULT_ENV}} \
	 TLS=$${TLS:-${DEFAULT_TLS}} \
	 LOG=$${LOG:-${DEFAULT_LOG}} \
	 DSN=$${DSN:-${DEFAULT_DSN}} \
	 docker-compose up -d migrate
	@echo "Database migrations completed!"

## migrate_up: Runs database migrations.
migrate_stop:
	@echo "Running database migrations..."
	@PORT=$${PORT:-${DEFAULT_PORT}} \
	 ENV=$${ENV:-${DEFAULT_ENV}} \
	 TLS=$${TLS:-${DEFAULT_TLS}} \
	 LOG=$${LOG:-${DEFAULT_LOG}} \
	 DSN=$${DSN:-${DEFAULT_DSN}} \
	 docker-compose down migrate
	@echo "Database migrations completed!"

## build_sniplate: Builds the sniplate binary as a linux executable.
build_sniplate:
	@echo "Building sniplate binary..."
	cd ../../sniplate && env GOOS=linux GOARCH=amd64 go build -o ${SNIPLATE_BINARY} ./cmd/api
	@echo "Done!"

