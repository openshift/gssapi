#!/bin/bash -eu

REUSE_DOCKER_IMAGES="1" \
SERVICE_LOG_FILTER="true" \
EXT_KDC_IP="" \
EXT_KDC_PORT="" \
KEYTAB_FILE="" \
SERVICE_NAME="HTTP/auth.www.levtest.net" \
REALM_NAME="APSARA.IO" \
DOMAIN_NAME="www.levtest.net" \
USER_NAME="testuser" \
USER_PASSWORD="P@ssword!" \
CLIENT_IN_CONTAINER="" \
        ./run.sh


