#!/bin/bash

export APP_NAME=mock-cds-app

docker rm $(docker stop $(docker ps -a -q --filter ancestor="${APP_NAME}" --format="{{.ID}}"))
docker rmi ${APP_NAME}