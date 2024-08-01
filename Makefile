FRONT_END_PORT=8081
AUTH_BINARY=authApp
MAIL_BINARY=mailApp
BROKER_BINARY=brokerApp
LOGGER_BINARY=loggerApp
LISTENER_BINARY=listenerApp
FRONT_END_BINARY=frontEndApp


## For Prod - full step for service build generation, docker-image creation, and pushing docker image for Prod EC2 Instance(linux/amd64)
## And finally running the build_all to generate the builds for local machine, that we have overwritten in previous steps.
full_prod_linux_amd64: build_prod_linux_amd64 build_docker_images_prod_linux_amd64 push_docker_images_prod_linux_amd64	build_all

## For Prod - generating build for linux_amd64 machine, to generate linux/amd64 docker image to run in EC2 Instance(linux/amd64)
build_prod_linux_amd64:
	@echo "Building (linux/amd64) compatible build..."
	cd ./front-end && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${FRONT_END_BINARY} ./cmd/web
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${BROKER_BINARY} ./cmd/api
	cd ./auth-service && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${AUTH_BINARY} ./cmd/api
	cd ./logger-service && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${LOGGER_BINARY} ./cmd/api
	cd ./mail-service && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${MAIL_BINARY} ./cmd/api
	cd ./listener-service && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${LISTENER_BINARY} .
	@echo "Build for (linux/amd64) Completed."

## For Prod - build docker images for EC2 Instance(linux/amd64)
build_docker_images_prod_linux_amd64:
	@echo "Creating docker-images..."
	cd ./ && docker buildx build --platform linux/amd64 -f caddy.production.dockerfile -t shub96/micro-caddy-production-linux-amd64 .
	cd ./auth-service && docker buildx build --platform linux/amd64 -f auth-service.dockerfile -t shub96/auth-service-production-linux-amd64 .
	cd ./broker-service && docker buildx build --platform linux/amd64 -f broker-service.dockerfile -t shub96/broker-service-production-linux-amd64 .
	cd ./front-end && docker buildx build --platform linux/amd64 -f front-end.dockerfile -t shub96/front-end-production-linux-amd64 .
	cd ./listener-service && docker buildx build --platform linux/amd64 -f listener-service.dockerfile -t shub96/listener-service-production-linux-amd64 .
	cd ./logger-service && docker buildx build --platform linux/amd64 -f logger-service.dockerfile -t shub96/logger-service-production-linux-amd64 .
	cd ./mail-service && docker buildx build --platform linux/amd64 -f mail-service.dockerfile -t shub96/mail-service-production-linux-amd64 .
	@echo "Created docker-images."

## For Prod - push docker images to docker-hub for EC2 Instance(linux/amd64)
push_docker_images_prod_linux_amd64:
	@echo "Pushing docker-images to docker-hub..."
	docker login
	docker push shub96/micro-caddy-production-linux-amd64
	docker push shub96/auth-service-production-linux-amd64
	docker push shub96/broker-service-production-linux-amd64
	docker push shub96/front-end-production-linux-amd64
	docker push shub96/listener-service-production-linux-amd64
	docker push shub96/logger-service-production-linux-amd64
	docker push shub96/mail-service-production-linux-amd64
	@echo "Pushed docker-images to docker-hub."


## up_build: For local - stops docker-compose(if running), builds all projects and starts docker compose
up_build: down build_all
	@echo "Building (when required) and starting docker images..."
	docker-compose up -d --build
	@echo "Docker images built and started!"

## build_all: generate build for all services
build_all: build_frontend build_broker build_auth build_logger build_mail build_listener	

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
