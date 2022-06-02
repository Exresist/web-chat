#!/usr/bin/env bash

set -euo pipefail

export IMAGE_TAG="${CI_COMMIT_TAG:-$CI_COMMIT_SHORT_SHA}"

export IMAGE="${CI_REGISTRY_IMAGE}"
