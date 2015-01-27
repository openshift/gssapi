#!/bin/sh -eu

# This script is used in the context of a docker VM when runnning the linux
# client test, and in the context of OS X when running on the Macintosh.  The
# following variables must be set (either via --link or explicitely)
#       KDC_PORT_88_TCP_ADDR
#       KDC_PORT_88_TCP_PORT
#               KDC address and port
#
#       SERVICE_PORT_80_TCP_ADDR
#       SERVICE_PORT_80_TCP_PORT
#               Http service address and port
#
#       KRB5_CONFIG_TEMPLATE
#       KRB5_CONFIG
#               The locations of the krb5.conf template, and where the
#               processed file must go
#
#       GSSAPI_PATH
#               The gssapi .so
#
#       TEST_DIR
#               The directory to build the client test app in

export PATH=$PATH:$GOROOT/bin

sed -e "s/KDC_ADDRESS/$KDC_PORT_88_TCP_ADDR:$KDC_PORT_88_TCP_PORT/g" \
	$KRB5_CONFIG_TEMPLATE > $KRB5_CONFIG

echo P@ssword! | kinit client.user@TEST.GOGSSAPI.COM >/dev/null

(cd $TEST_DIR && go test -c github.com/apcera/gssapi/component_test/client)

while ! echo exit | nc $SERVICE_PORT_80_TCP_ADDR $SERVICE_PORT_80_TCP_PORT >/dev/null; do
        echo "Waiting for service to start"
        sleep 1
done

$TEST_DIR/client.test \
	--test.bench=. \
	--test.v=false \
	--test.benchtime=5s \
	--service-name=HTTP/service.s.gogssapi.com@TEST.GOGSSAPI.COM \
	--service-address=$SERVICE_PORT_80_TCP_ADDR:$SERVICE_PORT_80_TCP_PORT \
	--krb5-config=$KRB5_CONFIG \
	--gssapi-path=$GSSAPI_PATH \
        2>&1
