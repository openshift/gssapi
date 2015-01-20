sed -e "s/KDC_ADDRESS/$KDC_PORT_88_TCP_ADDR:$KDC_PORT_88_TCP_PORT/g" \
	/opt/go-gssapi-test-service/krb5.conf.template \
	> /opt/go-gssapi-test-service/krb5.conf

exec /opt/go-gssapi-test-service/service \
	-mode=service \
	-service-name=service.user/kdc.example.com \
	-service-address=:80 \
	-gssapi-path=/usr/lib/x86_64-linux-gnu/libgssapi_krb5.so.2 \
	-krb5-config=/opt/go-gssapi-test-service/krb5.conf \
	-krb5-ktname=/opt/go-gssapi-test-service/krb5.keytab
