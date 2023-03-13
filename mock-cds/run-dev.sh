#!/bin/bash

export APP_NAME=mock-cds-app
export DOCKER_PORT=8080
export APP_PORT=8080

docker build -t $APP_NAME .
docker run -d -p $DOCKER_PORT:$APP_PORT $APP_NAME