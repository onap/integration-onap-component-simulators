#!/bin/bash

export APP_NAME=mock-ves

docker rm -f ${APP_NAME}
docker rmi ${APP_NAME}

