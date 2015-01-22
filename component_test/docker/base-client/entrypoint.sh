#!/bin/sh -e

sed -e "s/KDC_ADDRESS/$KDC_PORT_88_TCP_ADDR:$KDC_PORT_88_TCP_PORT/g" \
	/opt/go-gssapi-test-client/krb5.conf.template \
	> /opt/go-gssapi-test-client/krb5.conf

export KRB5_CONFIG=/opt/go-gssapi-test-client/krb5.conf
echo P@ssword! | kinit client.user/kdc.example.com

/opt/go-gssapi-test-client/client.test \
	--test.bench=. \
	--test.v=false \
	--test.benchtime=5s \
	--service-name=service.user/kdc.example.com \
	--service-address=$SERVICE_PORT_80_TCP_ADDR:$SERVICE_PORT_80_TCP_PORT \
	--krb5-config=/opt/go-gssapi-test-client/krb5.conf \
	-gssapi-path=/usr/lib/x86_64-linux-gnu/libgssapi_krb5.so.2 2>&1

