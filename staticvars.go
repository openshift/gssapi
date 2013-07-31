// Copyright 2013 Apcera Inc. All rights reserved.

package gssapi

/*
#include <gssapi/gssapi.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Name-Types
// These are static structs with variables exported by the library; they are
// standardized in the RFCs.  The library requires that a given name be usable
// for resolution, but it's typically a macro, there's no guarantee about the
// name exported from the library.  But since they're static, and well-defined,
// we can just define them ourselves.

type raw_gss_OID_desc []byte

// It's okay to reference the underlying storage, because these are package
// variables that won't disappear.
func (oid raw_gss_OID_desc) C_gss_OID_desc() C.gss_OID_desc {
	var t C.gss_OID_desc
	t.length = C.OM_uint32(len(oid))
	t.elements = unsafe.Pointer(&oid[0])
	return t
}

// Used in testing to check that we get back something reasonable form the C
// converted form
func DebugStringCGssOIDDesc(oid C.gss_OID_desc) string {
	return fmt.Sprintf("{%d %p}:%q",
		oid.length, oid.elements,
		C.GoStringN((*C.char)(oid.elements), C.int(oid.length)),
	)
}

var (
	// RFC2744-mandated values, mapping from as-near-as-possible to cut&paste
	GSS_C_NT_USER_NAME           = raw_gss_OID_desc("\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x01")
	GSS_C_NT_MACHINE_UID_NAME    = raw_gss_OID_desc("\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x02")
	GSS_C_NT_STRING_UID_NAME     = raw_gss_OID_desc("\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x03")
	GSS_C_NT_HOSTBASED_SERVICE_X = raw_gss_OID_desc("\x2b\x06\x01\x05\x06\x02")
	GSS_C_NT_HOSTBASED_SERVICE   = raw_gss_OID_desc("\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x04")
	GSS_C_NT_ANONYMOUS           = raw_gss_OID_desc("\x2b\x06\x01\x05\x06\x03") // original had \01
	GSS_C_NT_EXPORT_NAME         = raw_gss_OID_desc("\x2b\x06\x01\x05\x06\x04")
)
