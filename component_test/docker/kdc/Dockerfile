FROM ubuntu:14.04

RUN apt-get -y update
ADD krb5.conf /etc/krb5.conf
RUN apt-get -y install heimdal-kdc

EXPOSE	88

# Create keytab folder.
RUN mkdir /etc/docker-kdc

# Add kerberos principal/s.
RUN /bin/echo -e '\n\n\n\n\n\nP@ssword!\nP@ssword!\n' | kadmin -l add HTTP/service.s.gogssapi.com@TEST.GOGSSAPI.COM
RUN /bin/echo -e '\n\n\n\n\n\nP@ssword!\nP@ssword!\n' | kadmin -l add client.user@TEST.GOGSSAPI.COM

# Export keytab.
RUN kadmin -l ext_keytab -k /etc/docker-kdc/krb5.keytab HTTP/service.s.gogssapi.com@TEST.GOGSSAPI.COM

# KDC daemon startup.
#TODO -- what's in this config? Need to provide my own?
ENTRYPOINT ["/usr/lib/heimdal-servers/kdc", "--config-file=/etc/heimdal-kdc/kdc.conf", "-P 88"]

