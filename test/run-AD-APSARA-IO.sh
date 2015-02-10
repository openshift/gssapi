#!/bin/bash -eu

#TODO need a pubic user for the client

tmp_keytab_file="/tmp/keytab$(date +%s)"
echo "BQIAAABHAAIACUFQU0FSQS5JTwAESFRUUAAFdGFzdHkAAAABAAAAAAMAEgAg1wBAGJBWc122iNwyNJOtbWq8OIhoS2NzCv9PKfLkFnQ=" | base64 --decode >$tmp_keytab_file
REUSE_DOCKER_IMAGES="1" \
SERVICE_LOG_FILTER="true" \
EXT_KDC_HOST="ad1.apsara.io" \
EXT_KDC_PORT="88" \
KEYTAB_FILE="$tmp_keytab_file" \
SERVICE_NAME="HTTP/tasty" \
REALM_NAME="APSARA.IO" \
DOMAIN_NAME="apsara.io" \
USER_NAME="lev" \
USER_PASSWORD="P@ssword!" \
CLIENT_IN_CONTAINER="" \
        ./run.sh
rm $tmp_keytab_file
