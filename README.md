# gssapi

The gssapi package is a Golang wrapper around [RFC 2743](https://www.ietf.org/rfc/rfc2743.txt),
the Generic Security Services Application Programming Interface. (GSSAPI)

## Uses

We use it to authenticate clients with our authentication server. Clients talk
to a Kerberos or Active Directory Domain Controller to retrieve a Kerberos
service ticket, which we verify with a keytab on our authentication server.

When a user logs into Kerberos using `kinit`, they get a Kerberos TGT. During
Kerberos authentication, that TGT is used to retrieve a Service Ticket from the
Domain Controller. GSSAPI lets us authenticate without having to know where or
in what form the TGT is stored.

What do you use it for? Let us know!

## TODO

See our [TODO doc](TODO.md) on stuff you can do to help. We welcome
contributions!

## Verified platforms

We've tested that we can authenticate against:

- Heimdal Kerberos
- Active Directory

We suspect we can authenticate against:

- MIT Kerberos

We definitely cannot authenticate with:

- Windows clients (because Windows uses SSPI instead of GSSAPI as the library
  interface)
