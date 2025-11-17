#!/bin/bash

set -euo pipefail

blob server \
	--minio-address='localhost:9000' \
	--minio-access-key='root' \
	--minio-secret-key='password' \
	--minio-bucket='test' \
	--minio-object-prefix='' \
	--minio-secure='false'
