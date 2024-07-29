FRONT_END_PORT=8081
AUTH_BINARY=authApp
MAIL_BINARY=mailApp
BROKER_BINARY=brokerApp
LOGGER_BINARY=loggerApp
LISTENER_BINARY=listenerApp
FRONT_END_BINARY=frontEndApp

## up_build: stops docker-compose(if running), builds all projects and starts docker compose
up_build: down build_frontend build_broker build_auth build_logger build_mail build_listener
	@echo "Building (when required) and starting docker images..."
	docker-compose up -d --build
	@echo "Docker images built and started!"

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## down: stops docker compose
down:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Docker Stopped!"

## build_frontend: builds the front-end binary as a linux executable
build_frontend:
	@echo "Building front-end binary..."
	cd ./front-end && env GOOS=linux CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Front-end binary generated!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Broker binary generated!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ./auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Auth binary generated!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Logger binary generated!"

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ./mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Mail binary generated!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ./listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Listener binary generated!"
