FRONT_END_BINARY=frontApp
FRONT_END_PORT=80

## build_front: builds the front-end binary
build_front:
	@echo "Building front-end binary..."
	cd ./front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Front-end build completed!"

## start: starts the front-end server
start_front: build_front
	@echo "Starting front-end server"	
	cd ./front-end && ./${FRONT_END_BINARY}

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
