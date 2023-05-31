# Building Tradepipe

This document will help guide you through understanding the build process.

## Requirements

- Docker

  **Note**: You will need to check if Docker CLI plugin buildx is properly installed (`docker-buildx` file should be present in `~/.docker/cli-plugins`). You can install buildx according to the [instructions](https://github.com/docker/buildx/blob/master/README.md#installing).

## Overview

It is possible to build tradepipe using local golang, but we have a build process that runs in a Docker container. This simplifies initial set up.


## Scripts

The following scripts are found in [build/](.) directory.
