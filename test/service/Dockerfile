FROM ubuntu:14.04

RUN apt-get -y update
RUN apt-get -y install libsasl2-modules-gssapi-heimdal libgssapi-krb5-2
ADD service /opt/go-gssapi-test-service/service
ADD krb5.keytab /etc/krb5.keytab
ADD krb5.conf.template /opt/go-gssapi-test-service/krb5.conf.template
ADD entrypoint.sh /opt/go-gssapi-test-service/entrypoint.sh

EXPOSE 80

# service startup
ENTRYPOINT /opt/go-gssapi-test-service/entrypoint.sh
