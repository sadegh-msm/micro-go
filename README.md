# micro-go
> microservice in golang 

## introduction
 Here is a mini project most for self-learning and working with `micro-services` and working with other technology's.

 this project has 5 different services that uses [gRPC](https://github.com/grpc/grpc-go) for more efficient and faster connections
 (first I use the http/rest and after it I wanted faster service, so I add gRPC so the http/rest handlers and endpoints are available to).
 I also used [rabbitMQ](https://www.rabbitmq.com) for message queue in case servers has more request that they can handle.
 
## up and running 
 to start working with these services in case you don't have the gRPC generated code you can run `protoc.sh` file in `project`
 directory, and you can pass the `log`, `shortner` or `auth` for creating the needed file for code to run.
 
 after generating the code go to `project` directory and use the make file to run any services you want (more data about make arguments in `Makefile` is included)
 for setting up all the services you can use `make up` command and wait for couple seconds to all configurations executes 
 and all services connect to databases and rabbitMQ start the work.
 
 after setting up the project you are good to go and send request to any services or send request to broker and specify the service you want to use
 and broker will handle the rest.
 
### database 
 in this project I used 3 different database for every service, [redis](https://github.com/redis/redis) for `url-shortner service`, 
 [mongodb](https://github.com/mongodb/mongo) for `logging service` and [postgres](https://www.postgresql.org) for `authentiction service`.
 
### models 
Project models are defined in `model` package. These models are used internally, but they can be used in response or request package.


