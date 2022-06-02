#!/bin/bash
set -euo pipefail errexit

export HELM_STARTER="helm-scaffolds/general-backend"

# Generating a helm chart from custom (our) starter
helm create ${CI_PROJECT_NAME} --starter $HELM_STARTER

# Linting helm chart alongside with repo-based values
mv ${CI_PROJECT_NAME}/values.yaml values.yaml
helm lint ${CI_PROJECT_NAME} --values ${CI_ENVIRONMENT_NAME}.values.yaml
mv values.yaml ${CI_PROJECT_NAME}/values.yaml

# If all is ok with checkings, it's time do make a deploy now
helm upgrade --install --atomic ${CI_PROJECT_NAME} \
  --namespace=${KUBE_NAMESPACE} ${CI_PROJECT_NAME} \
  --set image.tag=${IMAGE_TAG} \
  --values ${CI_ENVIRONMENT_NAME}.values.yaml
