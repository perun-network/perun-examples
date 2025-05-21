#!/bin/bash

set -e

DETACHED=0

case "$1" in
standalone)
    echo "Using standalone network"
    ARGS="--local"
    ;;
standalone-detached)
    echo "Using standalone network (detached mode)"
    ARGS="--local"
    DETACHED=1
    ;;
futurenet)
    echo "Using Futurenet network"
    ARGS="--futurenet"
    ;;
*)
    echo "Usage: $0 standalone|standalone-detached|futurenet"
    exit 1
    ;;
esac

# This is set to the quickstart `soroban-dev` image annointed as the release 
# for a given Soroban Release. See: https://soroban.stellar.org/docs/reference/releases 
# QUICKSTART_SOROBAN_DOCKER_SHA=stellar/quickstart:soroban-dev
QUICKSTART_SOROBAN_DOCKER_SHA=stellar/quickstart:testing@sha256:274395daab6fa8033b9213f152d56699358917fb01d7c7e95392a37fc00c9d01

shift

# Create Docker network if it doesn't exist
echo "Creating docker soroban network"
(docker network inspect soroban-network -f '{{.Id}}' 2>/dev/null) \
  || docker network create soroban-network

# Clean up existing soroban-preview container if present
echo "Searching for a previous soroban-preview docker container"
containerID=$(docker ps --filter="name=soroban-preview" --all --quiet)
if [[ ${containerID} ]]; then
    echo "Start removing soroban-preview container."
    docker rm --force soroban-preview
    echo "Finished removing soroban-preview container."
else
    echo "No previous soroban-preview container was found"
fi

# Start the soroban-preview container in detached mode
currentDir=$(pwd)
docker run -d \
  --volume "${currentDir}:/workspace" \
  --memory="2g" \
  --cpus="2.0" \
  --name soroban-preview \
  -p 8001:8000 \
  --ipc=host \
  --network soroban-network \
  soroban-preview:10

# Start the stellar quickstart container
if [[ "$DETACHED" == "1" ]]; then
  echo "Running stellar quickstart container in detached mode"
  docker run -d \
    --name stellar \
    --pull always \
    --network soroban-network \
    -p 8000:8000 \
    "$QUICKSTART_SOROBAN_DOCKER_SHA" \
    $ARGS \
    --enable-soroban-rpc \
    "$@"
else
  echo "Running stellar quickstart container in foreground"
  docker run --rm \
    --name stellar \
    --pull always \
    --network soroban-network \
    -p 8000:8000 \
    "$QUICKSTART_SOROBAN_DOCKER_SHA" \
    $ARGS \
    --enable-soroban-rpc \
    "$@"
fi
