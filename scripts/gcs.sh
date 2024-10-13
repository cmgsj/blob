#!/bin/bash

export STORAGE_EMULATOR_HOST="http://localhost:4443"

blob server start --gcs-uri='http://localhost:4443/storage/v1/cmgsj'

# List buckets
# curl -sSL "${STORAGE_EMULATOR_HOST}/storage/v1/b" | jq

# Get bucket
# curl -sSL "${STORAGE_EMULATOR_HOST}/storage/v1/b/cmgsj" | jq

# List bucket objects
# curl -sSL "${STORAGE_EMULATOR_HOST}/storage/v1/b/cmgsj/o" | jq -r '.items.[].name'

# Get bucket objects
# curl -sSL "${STORAGE_EMULATOR_HOST}/storage/v1/b/cmgsj/o/blobs/foo" | jq
