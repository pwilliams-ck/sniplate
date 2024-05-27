SNIPLATE_BINARY=sniplateApp

# Default values for the environment variables
DEFAULT_PORT=4200
DEFAULT_ENV=development
DEFAULT_TLS=false
DEFAULT_LOG=false

## up: Starts all containers in the background without forcing build.
up:
	@echo "Starting Docker images..."
	@PORT=$${PORT:-${DEFAULT_PORT}} \
	 ENV=$${ENV:-${DEFAULT_ENV}} \
	 TLS=$${TLS:-${DEFAULT_TLS}} \
	 LOG=$${LOG:-${DEFAULT_LOG}} \
	 docker-compose up -d
	@echo "Docker images started!"

## up_build: Stops docker-compose (if running), builds all projects and starts docker compose with flags.
up_build: build_sniplate
	@echo "Stopping docker images (if running...)"
	@PORT=$${PORT:-${DEFAULT_PORT}} \
	 ENV=$${ENV:-${DEFAULT_ENV}} \
	 TLS=$${TLS:-${DEFAULT_TLS}} \
	 LOG=$${LOG:-${DEFAULT_LOG}} \
	 docker-compose down
	@echo "Building (when required) and starting docker images..."
	@PORT=$${PORT:-${DEFAULT_PORT}} \
	 ENV=$${ENV:-${DEFAULT_ENV}} \
	 TLS=$${TLS:-${DEFAULT_TLS}} \
	 LOG=$${LOG:-${DEFAULT_LOG}} \
	 docker-compose up --build -d
	@echo "Docker images built and started!"

## down: Stop docker-compose.
down:
	@echo "Stopping docker compose..."
	@PORT=$${PORT:-${DEFAULT_PORT}} \
	 ENV=$${ENV:-${DEFAULT_ENV}} \
	 TLS=$${TLS:-${DEFAULT_TLS}} \
	 LOG=$${LOG:-${DEFAULT_LOG}} \
	 docker-compose down
	@echo "Done!"

## build_sniplate: Builds the sniplate binary as a linux executable.
build_sniplate:
	@echo "Building sniplate binary..."
	cd ../../sniplate && env GOOS=linux GOARCH=amd64 go build -o ${SNIPLATE_BINARY} ./cmd/api
	@echo "Done!"

