BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
IDENTITY_BINARY=identityApp
TRANSACTION_BINARY=transactionApp
LISTENER_BINARY=listenerApp
## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_identity build_listener build_transaction
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_identity: builds the identity binary as a linux executable
build_identity:
	@echo "Building logger binary..."
	cd ./identity-service && env GOOS=linux CGO_ENABLED=0 go build -o ${IDENTITY_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ./auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building auth binary..."
	cd ./listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} ./
	@echo "Done!"
## build_listener: builds the listener binary as a linux executable
build_transaction:
	@echo "Building auth binary..."
	cd ./transaction-service && env GOOS=linux CGO_ENABLED=0 go build -o ${TRANSACTION_BINARY} ./cmd/api
	@echo "Done!"

