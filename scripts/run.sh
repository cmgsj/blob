#!/bin/bash

set -euo pipefail

export MINIO_ROOT_USER='root'
export MINIO_ROOT_PASSWORD='password'

go run ./cmd/blob server \
    --driver='minio=http://localhost:9000/test'
