==========================
GSSAPI Go Language Roadmap
==========================

Rather than reimplement mechanisms for Go, this wraps libraries provided for
the C language bindings, so we conform to values specified in [RFC2744][].

 * [RFC2744][]: Generic Security Service API Version 2 : C-bindings
 * [RFC5587][]: Extended Generic Security Service Mechanism Inquiry APIs

Methods with names starting `gss_` are internal to this module and not
exported; they often map directly to the C bindings, with much finangling
going on.

The exported methods are based on sub-types and are shortened.
For instance, the `(*GssapiLib).gss_release_buffer()` method is internal,
whereas the `(GssBuffer).Release()` method is exposed, acting upon a
`GssBuffer`, which is the Go-exported type which ties a `C.gss_buffer_t` into
the housekeeping information needed.


Code relating to dynamic library loading is mostly in [dynload.go](dynload.go).

Our own top-level types and generic handling is in [gssapi.go](gssapi.go).

Code relating to errors, statuses and constants relating to those is in [status.go](status.go).

Most constants of the form `GSS_C_*` are in [consts.go](consts.go).

Code relating to credential handling, per [RFC2744][] table 2-1, is in [credential.go](credential.go).

Code relating to contexts, per [RFC2744][] table 2-2, is in [context.go](context.go).

Code relating to per-message handling, per [RFC2744][] table 2-3, is in [message.go](message.go).

Code relating to name manipulation, per [RFC2744][] table 2-4, is in [name.go](name.go).

Handling of routines classified as miscellaneous, per [RFC2744][] table 2-5, is in [misc.go](misc.go).

[RFC2744]: http://www.ietf.org/rfc/rfc2744.txt
[RFC5587]: http://www.ietf.org/rfc/rfc5587.txt
