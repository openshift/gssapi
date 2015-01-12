// Copyright 2013-2015 Apcera Inc. All rights reserved.

// Wrappers for the main gssapi types, all in one file for consistency.

package gssapi

/*
#include <gssapi/gssapi.h>
*/
import "C"

// Struct types. The structs themselves are allocated in Go and are therefore
// GCed, the contents may come from Go or from C/gssapi calls, and therefore
// must be explicitly released.  Calling the Release method once per instance
// is safe.

type Buffer struct {
	*Lib
	C_gss_buffer_t C.gss_buffer_t

	// indicates if the contents of the buffer must be released with
	// gss_release_buffer
	releasable bool
}

type Name struct {
	*Lib
	C_gss_name_t C.gss_name_t
}

// OID is the wrapper for gss_OID_desc type. OIDs are not released explicitly,
// only as part of OIDSet
type OID struct {
	*Lib
	C_gss_OID C.gss_OID
	data      []byte
}

type OIDSet struct {
	*Lib
	C_gss_OID_set C.gss_OID_set
}

type CredId struct {
	*Lib
	C_gss_cred_id_t C.gss_cred_id_t
}

type CtxId struct {
	*Lib
	C_gss_ctx_id_t C.gss_ctx_id_t
}

// Aliases for the simple types
type CredUsage C.gss_cred_usage_t // C.int
type ChannelBindingAddressFamily uint32

// A struct pointer technically, but not really used yet, and it's a static,
// non-releaseable struct so this may suffice
type ChannelBindings C.gss_channel_bindings_t
