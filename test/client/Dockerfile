FROM ubuntu:14.04

RUN apt-get -y update
RUN apt-get -y install libsasl2-modules-gssapi-heimdal libgssapi-krb5-2 krb5-user
ADD client /opt/go-gssapi-test-client/client
ADD krb5.conf.template /opt/go-gssapi-test-client/krb5.conf.template
ADD entrypoint.sh /opt/go-gssapi-test-client/entrypoint.sh

ENTRYPOINT /opt/go-gssapi-test-client/entrypoint.sh
