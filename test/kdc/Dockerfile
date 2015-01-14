FROM ubuntu:14.04

RUN apt-get -y update

ADD krb5.conf /etc/krb5.conf

RUN apt-get -y install heimdal-kdc
RUN apt-get -y install libsasl2-modules-gssapi-heimdal

EXPOSE	88

# Create keytab folder.
RUN mkdir /etc/docker-kdc

# Add kerberos principal/s.
RUN /bin/echo -e '\n\n\n\n\n\nP@ssword!\nP@ssword!\n' | kadmin -l add service.user/kdc.example.com
RUN /bin/echo -e '\n\n\n\n\n\nP@ssword!\nP@ssword!\n' | kadmin -l add client.user/kdc.example.com

# Export keytab.
RUN kadmin -l ext_keytab -k /etc/docker-kdc/krb5.keytab service.user/kdc.example.com@EXAMPLE.COM

# KDC daemon startup.
ENTRYPOINT ["/usr/lib/heimdal-servers/kdc", "--config-file=/etc/heimdal-kdc/kdc.conf", "-P 88"]

