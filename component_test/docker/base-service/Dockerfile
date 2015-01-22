FROM ubuntu:14.04
RUN apt-get -y update
RUN apt-get -y install libsasl2-modules-gssapi-mit libgssapi-krb5-2
ADD krb5.keytab /opt/go-gssapi-test-service/krb5.keytab
ADD krb5.conf.template /opt/go-gssapi-test-service/krb5.conf.template
ADD entrypoint.sh /opt/go-gssapi-test-service/entrypoint.sh
EXPOSE 80
ENTRYPOINT /opt/go-gssapi-test-service/entrypoint.sh
