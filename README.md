# go-microservices
This repo contains multiple micro-services (including 1 front-end service), Databases (MongoDB & Postgres), Messaging Queue (RabbitMQ), and Mail-Server (Mailpit) communicating with each other via REST, RPC, gRPC, & Messaging Queue.

Execute ```make up_build``` to run the 'docker-compose.yml', if you want to run the front-end and broker-service without using 'Caddy' in your local machine.

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
docker build -f front-end.dockerfile -t shub96/front-end:1.0.0 .
docker push shub96/front-end:1.0.0

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
2. Stop the existing docker services (if running) using ```make down```
3. Execute ```docker swarm init```
4. Execute ```docker stack deploy -c swarm.yml myapp```
5. Verify if all the services are up using ```docker service ls```
6. To scale the service execute ```docker service scale myapp_auth-service=3```
7. Stop entire docker-swarm using ```docker stack rm myapp```, removing manually from 'docker-desktop' will not work, as it will be getting created again and again if 'docker-swarm' is runnning.

## Update service in Docker-Swarm:
1. To deploy updated version of the docker service in Docker-Swarm, create updated build for that service using ```make build_auth```, than create new docker image for that service with new version tagged to it, and push it to docker-hub, similar to what we have done earlier.
2. You can also increase the current running instance of that service to at least 2, so that there is no downtime while updating the service in docker-swarm (in production), as docker-swarm updates one instance of the service at a time.
3. Use below command to update the service verison in the running docker-swarm:
``` docker service update --image shub96/auth-service:1.0.1 myapp_auth-service```, also don't forget to update the docker image version in your 'swarm.yml' file.

## Stopping Docker-Swarm:
- Stop a particular service in docker-swarm using ```docker service scale myapp_auth-service=0```, which scales down the instance of that service to '0', but does not removes it.
- Stop entire docker-swarm using ```docker stack rm myapp```, this will remove all the services in that docker-swarm. Removing manually from 'docker-desktop' will not work, as it will be getting created again and again if 'docker-swarm' is runnning.
- To entirely leave the docker-swarm execute ```docker swarm leave```, which will show you a warning message, so add the flag ```--force``` at the end of above command to force its close.


## Use of Caddy as reverse-proxy:
- Use the file 'caddy.dockerfile' & config 'Caddyfile', create a docker image, and push it to docker-hub, and post that use that image in your 'swarm.yml' file, to include it
  as part of your 'docker-swarm'.
  ```
    docker build -f caddy.dockerfile -t shub96/micro-caddy:1.0.0 .
    docker push shub96/micro-caddy:1.0.0
  ```

### Issue while deploying
- To deploy on your EC2 Instance(linux/amd64), use below command to generate the platform specific docker-image for all of your services that you have written, push the image to dockerhub and use this image in your 'swarm.production.yml' file, that is used for deploying to EC2 Instance. We are doing this because the previous images are built for the Mac Apple silicon chip, because of which, the images built using those builds does not works for running docker on EC2 Instance (which we have selected as a linux/amd64 OS). 

    ```make full_prod_linux_amd64```

This is where the mail Inbox can be accessed - [Inbox](http://node-1.s5m.in:8025/)

## How to Update the Changes on your EC2 Instance:
- Step-1: Do the required changes in your code.
- Step-2: Run ```make full_prod_linux_amd64``` in your local, which will create the updated docker-images and push to docker-hub.
- Step-3: In AWS, go to your EC2 instance, and connect to it using terminal.
- Step-4: Go to swarm folder using ```cd /swarm```
- Step-5: [We can skip this step] Stop the current running swarm using ```docker stack rm myapp```
- Step-6: Run the swarm again using ```docker stack deploy -c swarm.yml myapp```
By following above mentioned steps, your changes will be reflected on the deployed [front-end](https://swarm.s5m.in/)