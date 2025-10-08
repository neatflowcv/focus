#!/bin/bash
set -euo pipefail

APP_NAME=focus
podman stop ${APP_NAME:?}-db
podman rm ${APP_NAME:?}-db
podman volume rm ${APP_NAME:?}-data