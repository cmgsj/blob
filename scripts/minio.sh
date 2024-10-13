#!/bin/bash

blob server start \
	--minio-address='localhost:9000' \
	--minio-access-key='root' \
	--minio-secret-key='password' \
	--minio-secure='false' \
	--minio-bucket='cmgsj' \
	--minio-object-prefix=''
