# TRADEPIPE

[![Test and validate](https://github.com/Sannrox/tradepipe/actions/workflows/test.yml/badge.svg)](https://github.com/Sannrox/tradepipe/actions/workflows/test.yml)

This is a microservice for the private API of the Trade Republic online brokerage. I am not affiliated with Trade Republic Bank GmbH.

Inspired by https://github.com/marzzzello/pytr


## Usage

If you only want to download your bank documents, it is enough to do the following to build the binary: `make tradepipe`

If you want to use the docker compose setup with the GRPC server and database you should follow these steps:

1. `./build/release-images.sh`
2. `./deployments/get-trade.sh`

**Note**: This `get-trade.sh` uses `docker system prune`

3. To start the compose file: `./deployments/trade-up.sh`
4. To stop the compose environment: `./deployments/trade-down.sh`
## Contributing

Please see [here](./CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

Feature requests are welcome.

## Roadmap

Depends on future feature requests



