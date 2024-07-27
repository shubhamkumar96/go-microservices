# go-microservices

## Command for generating grpc related code from 'logs.proto' file :
After getting inside the 'logger-service/logs' folder where the 'logs.proto' file is present,
exceute below command :

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto