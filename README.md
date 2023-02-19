# TRADEPIPE

This is a microservice for the private API of the Trade Republic online brokerage. I am not affiliated with Trade Republic Bank GmbH.

Inspired by https://github.com/marzzzello/pytr

## Features

### Overview

- [GRPC server](#grpc-server) - [protobuffer](./api/proto/tradepipe.proto)
- [HTTP server](#http-server) - [openapi](./api/openapi/openapi.yaml)
- [Single command](#single-command)

### GRPC server
#### Usage
```
# Build the binary 
$ make build

# Run the command 
$ ./build/bin/tradepipe-<GOOS>-<GOHOSTARCH>-<VERSION>  --grpc 
```

Use the generate client from https://github.com/Sannrox/tradepipe/grpc

For example take a look at this [fakeclient](./helper/testhelpers/fakegrpcclient/fake_client.go)
### HTTP server 
#### Usage 

```
# Build the binary 
$ make build

# Run the command 
$ ./build/bin/tradepipe-<GOOS>-<GOHOSTARCH>-<VERSION>  --http
```
Use the [openapi-spec](./api/openapi/openapi.yaml) to build/generate a client


### Single command

#### Usage

```
# Build the binary 
$ make build

# Run the command 
$ ./build/bin/tradepipe-<GOOS>-<GOHOSTARCH>-<VERSION>  <TR-NUMBER> <TR-PIN>

# Need to verify with 2FA
Enter Token: <token>
```