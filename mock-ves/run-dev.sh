#!/bin/bash

export APP_NAME=mock-ves
export DOCKER_PORT=30417
export APP_PORT=30417

docker build . -t $APP_NAME
docker run -p $DOCKER_PORT:$APP_PORT --name $APP_NAME $APP_NAME
