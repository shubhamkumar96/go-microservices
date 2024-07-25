FRONT_END_BINARY=frontApp
FRONT_END_PORT=80
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp

## up_all: starts all the services (front-end & docker-services)
up_all: up start_front

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose(if running), builds all projects and starts docker compose
up_build: build_broker build_auth
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up -d --build
	@echo "Docker images built and started!"

## down: stops docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Docker Stopped!"

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

## build_front: builds the front-end binary
build_front:
	@echo "Building front-end binary..."
	cd ./front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Front-end build completed!"

## start: starts the front-end server
start_front: build_front stop_port
	@echo "Starting front-end server"	
	cd ./front-end && ./${FRONT_END_BINARY} &

## stop: stops the front-end server
stop_front:
	@echo "Stopping front-end server..."
	$(MAKE) stop_port PORT=$(FRONT_END_PORT)
	@echo "Stopped front-end server!"

## TO STOP THE PORT BEING USED

# Default values for variables
PORT ?= 80

# Command to find the PID of the process using the port
FIND_PID := lsof -t -i :$(PORT)

# Command to kill the process using the port
KILL_PID := $(FIND_PID) | xargs kill -9

# Stop the process using the port, whose value will be passed while calling 'stop_port'.
stop_port:
	@echo "Checking for processes using port $(PORT)..."
	-@if [ "$$($(FIND_PID))" ]; then \
		echo "Stopping process using port $(PORT)"; \
		$(KILL_PID); \
	else \
		echo "No process using port $(PORT)"; \
	fi

# Declaring .PHONY targets to force execution
.PHONY: stop_port
