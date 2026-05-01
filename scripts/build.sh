#!/usr/bin/env bash

ENV=${1}
PROJECT_NAME=${2}
SERVICE_NAME=${3}
VERSION="${ENV}-$(date +%Y%m%d%H%M%S)"
PWD=$(cd "$(dirname "$0")" || exit; pwd)
export IMAGE_VERSION=${VERSION}
docker compose -p "${PROJECT_NAME}" -f "${PWD}"/../deploy/docker-compose.yml --compatibility up -d --build --force-recreate "${SERVICE_NAME}"
echo -e "\nservice: ${SERVICE_NAME}\nversion: ${VERSION}\n"
