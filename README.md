# go-microservices
This repo contains multiple micro-services (including 1 front-end service), Databases (MongoDB & Postgres), Messaging Queue (RabbitMQ), and Mail-Server (Mailpit) communicating with each other via REST, RPC, gRPC, & Messaging Queue.

## Command for generating grpc related code from 'logs.proto' file :
After getting inside the 'logger-service/logs' folder where the 'logs.proto' file is present,
exceute below command :
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto
```

## Command to build docker image and push to docker-hub:
After getting inside the folder where 'dockerfile' is present, run below command to build and tag docker image in local:
```
docker build -f logger-service.dockerfile -t dockerhub_username/logger-service:1.0.0 .
```

Login to your docker hub repo using below command, and entering your user-name & password
```
docker login
```
Now push the image to the DockerHub using below command:
```
docker push shub96/logger-service:1.0.0
```

Do similarly for all the services containing '.dockerfile' - 'auth-service', 'broker-service', 'listener-service', 'logger-service', 'mail-service'
```
docker build -f auth-service.dockerfile -t shub96/auth-service:1.0.0 .
docker push shub96/auth-service:1.0.0

docker build -f broker-service.dockerfile -t shub96/broker-service:1.0.0 .
docker push shub96/broker-service:1.0.0

docker build -f listener-service.dockerfile -t shub96/listener-service:1.0.0 .
docker push shub96/listener-service:1.0.0

docker build -f mail-service.dockerfile -t shub96/mail-service:1.0.0 .
docker push shub96/mail-service:1.0.0
```
## Deploy to Docker-Swarm:
Follow below steps to deploy to docker-swarm:
1. Create 'swarm.yml' file.
2. Stop the existing docker services (if running) using ```make down_all```
3. Execute ```docker swarm init```
4. Execute ```docker stack deploy -c swarm.yml myapp```
5. Verify if all the services are up using ```docker service ls```
6. To scale the service execute ```docker service scale myapp_auth-service=3```
7. Stop entire docker-swarm using ```docker stack rm myapp```, removing manually from 'docker-desktop' will not work, as it will be getting created again and again if 'docker-swarm' is runnning.

## Update service in Docker-Swarm:
1. To deploy updated version of the docker service in Docker-Swarm, just create new docker image for that service with new version tagged to it, and push it to docker-hub, similar to what we have done earlier.
2. You can also increase the current running instance of that service to at least 2, so that there is no downtime while updating the service in docker-swarm (in production), as docker-swarm updates one instance of the service at a time.
3. Use below command to update the service verison in the running docker-swarm:
``` docker service update --image shub96/auth-service:1.0.1 myapp_auth-service```

## Stopping Docker-Swarm:
- Stop a particular service in docker-swarm using ```docker service scale myapp_auth-service=0```, which scales down the instance of that service to '0', but does not removes it.
- Stop entire docker-swarm using ```docker stack rm myapp```, this will remove all the services in that docker-swarm. Removing manually from 'docker-desktop' will not work, as it will be getting created again and again if 'docker-swarm' is runnning.
- To entirely leave the docker-swarm execute ```docker swarm leave```, which will show you a warning message, so add the flag ```--force``` at the end of above command to force its close.


