#!/bin/bash

APP_NAME=focus
podman run -d \
--name ${APP_NAME:?}-db \
-e POSTGRES_DB=${APP_NAME:?} \
-e POSTGRES_USER=${APP_NAME:?} \
-e POSTGRES_PASSWORD=password \
-v ${APP_NAME:?}-data:/var/lib/postgresql \
-p 5432:5432 \
docker.io/library/postgres:latest