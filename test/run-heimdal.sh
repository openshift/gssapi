#!/bin/bash -eu

REUSE_DOCKER_IMAGES="1" \
SERVICE_LOG_FILTER="true" \
EXT_KDC_HOST="" \
EXT_KDC_PORT="" \
KEYTAB_FILE="" \
SERVICE_NAME="HTTP/auth.www.xample.test" \
REALM_NAME="XAMPLE.TEST" \
DOMAIN_NAME="xample.test" \
USER_NAME="testuser" \
USER_PASSWORD="P@ssword!" \
CLIENT_IN_CONTAINER="" \
        ./run.sh


