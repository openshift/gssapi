#!/bin/bash -eu

export PATH=$PATH:$GOROOT/bin

cat /tmp/krb5.conf.template \
        | sed -e "s/KDC_ADDRESS/$KDC_PORT_88_TCP_ADDR:$KDC_PORT_88_TCP_PORT/g" \
        | sed -e "s/DOMAIN_NAME/${DOMAIN_NAME}/g" \
        | sed -e "s/REALM_NAME/${REALM_NAME}/g" \
	> /opt/go-gssapi-test-service/krb5.conf

(cd /opt/go-gssapi-test-service && go build github.com/apcera/gssapi/component_test/service)

exec /opt/go-gssapi-test-service/service \
	-service-name=${SERVICE_NAME} \
	-service-address=:80 \
	-gssapi-path=/usr/lib/x86_64-linux-gnu/libgssapi_krb5.so.2 \
	-krb5-config=/opt/go-gssapi-test-service/krb5.conf \
	-krb5-ktname=/opt/go-gssapi-test-service/krb5.keytab
