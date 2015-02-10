// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

/*
#include <gssapi/gssapi.h>

void helper_gss_OID_desc_set_bytes(gss_OID_desc *oid , int l, void *p) {
	oid->length = l;
	oid->elements = p;
}

void helper_gss_OID_desc_get_bytes(gss_OID_desc *oid , int *l, char **p) {
	*l = oid->length;
	*p = oid->elements;
}

int
wrap_gss_oid_equal(void *fp, gss_OID oid1, gss_OID oid2)
{
	return ((int(*) (gss_OID, gss_OID)) fp)(oid1, oid2);
}

*/
import "C"

import (
	"fmt"
	"unsafe"
)

func (lib *Lib) NewOID() *OID {
	return &OID{
		Lib: lib,
	}
}

func (lib *Lib) GSS_C_NO_OID() *OID {
	return lib.NewOID()
}

// MakeOIDBytes makes an OID encapsulating a byte slice. Note that it does not duplicate
// the data, points to it directly
func (lib *Lib) MakeOIDBytes(data []byte) *OID {
	oid := &OID{
		Lib:       lib,
		C_gss_OID: &C.gss_OID_desc{},
		data:      data,
	}

	// because of the alignment issues I can't access o.oid's fields from go,
	// so invoking a C function to do the same as:
	// o.oid.length = C.OM_uint32(len(oid))
	// o.oid.elements = unsafe.Pointer(&oid[0])
	C.helper_gss_OID_desc_set_bytes(oid.C_gss_OID, C.int(len(oid.data)), unsafe.Pointer(&oid.data[0]))

	return oid
}

func (lib *Lib) MakeOIDString(data string) *OID {
	return lib.MakeOIDBytes([]byte(data))
}

func (oid OID) Bytes() []byte {
	var l C.int
	var p *C.char

	C.helper_gss_OID_desc_get_bytes(oid.C_gss_OID, &l, &p)

	return C.GoBytes(unsafe.Pointer(p), l)
}

// Used in testing to check that we get back something reasonable form the C
// converted form
func DebugStringCGssOIDDesc(oid C.gss_OID_desc) string {
	var l C.int
	var p *C.char

	C.helper_gss_OID_desc_get_bytes(&oid, &l, &p)

	return fmt.Sprintf(`{%d %p}:"%x"`, l, p, C.GoStringN(p, l))
}
