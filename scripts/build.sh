#!/bin/bash
docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY

docker build --progress auto \
    --build-arg CI_REGISTRY_USER=${CI_REGISTRY_USER} \
    --build-arg CI_REGISTRY_PASSWORD=${CI_REGISTRY_PASSWORD} \
    --build-arg CI_SERVER_HOST=${CI_SERVER_HOST} \
    --build-arg CI_SERVER_URL=${CI_SERVER_URL} \
    --tag ${IMAGE}:${IMAGE_TAG} .

docker push ${IMAGE}:${IMAGE_TAG}
