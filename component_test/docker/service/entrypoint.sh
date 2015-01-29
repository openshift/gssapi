#!/bin/sh -eu

export PATH=$PATH:$GOROOT/bin

sed -e "s/KDC_ADDRESS/$KDC_PORT_88_TCP_ADDR:$KDC_PORT_88_TCP_PORT/g" \
	/opt/go-gssapi-test-service/krb5.conf.template \
	> /opt/go-gssapi-test-service/krb5.conf

(cd /opt/go-gssapi-test-service && go build github.com/apcera/gssapi/component_test/service)

exec /opt/go-gssapi-test-service/service \
	-service-name=HTTP/service.s.gogssapi.com \
	-service-address=:80 \
	-gssapi-path=/usr/lib/x86_64-linux-gnu/libgssapi_krb5.so.2 \
	-krb5-config=/opt/go-gssapi-test-service/krb5.conf \
	-krb5-ktname=/opt/go-gssapi-test-service/krb5.keytab
