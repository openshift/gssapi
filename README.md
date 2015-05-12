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

## Testing

Tests in the main `gssapi` repository can be run using the built-in `go test`.

To run an integrated test against a live Heimdal Kerberos Domain Controller,
`cd test` and bring up [Docker](https://www.docker.com/), (or
[boot2docker](http://boot2docker.io/)). Then, run `./run-heimdal.sh`. This will
run some go tests using three Docker images: a client, a service, and a domain
controller. The service will receive a generated keytab file, and the client
will point to the domain controller for authentication.

**NOTE:** to run Docker tests, your `GOROOT` environment variable MUST be set.

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
