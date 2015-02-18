// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

/*
#include <gssapi.h>
#include <stdlib.h>
#include <string.h>

const size_t gss_OID_size=sizeof(gss_OID_desc);

// Name-Types.  These are standardized in the RFCs.  The library requires that
// a given name be usable for resolution, but it's typically a macro, there's
// no guarantee about the name exported from the library.  But since they're
// static, and well-defined, we can just define them ourselves.

// RFC2744-mandated values, mapping from as-near-as-possible to cut&paste
const gss_OID_desc *_GSS_C_NT_USER_NAME           = & (gss_OID_desc) { 10, "\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x01" };
const gss_OID_desc *_GSS_C_NT_MACHINE_UID_NAME    = & (gss_OID_desc) { 10, "\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x02" };
const gss_OID_desc *_GSS_C_NT_STRING_UID_NAME     = & (gss_OID_desc) { 10, "\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x03" };
const gss_OID_desc *_GSS_C_NT_HOSTBASED_SERVICE_X = & (gss_OID_desc) {  6, "\x2b\x06\x01\x05\x06\x02" };
const gss_OID_desc *_GSS_C_NT_HOSTBASED_SERVICE   = & (gss_OID_desc) { 10, "\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x04" };
const gss_OID_desc *_GSS_C_NT_ANONYMOUS           = & (gss_OID_desc) {  6, "\x2b\x06\x01\x05\x06\x03" };  // original had \01
const gss_OID_desc *_GSS_C_NT_EXPORT_NAME         = & (gss_OID_desc) {  6, "\x2b\x06\x01\x05\x06\x04" };

// from gssapi_krb5.h: This name form shall be represented by the Object
// Identifier {iso(1) member-body(2) United States(840) mit(113554) infosys(1)
// gssapi(2) krb5(2) krb5_name(1)}.  The recommended symbolic name for this
// type is "GSS_KRB5_NT_PRINCIPAL_NAME".
const gss_OID_desc *_GSS_KRB5_NT_PRINCIPAL_NAME	 = & (gss_OID_desc) { 10, "\x2a\x86\x48\x86\xf7\x12\x01\x02\x02\x01" };

void helper_gss_OID_desc_free_elements(gss_OID oid) {
	free(oid->elements);
}

void helper_gss_OID_desc_set_elements(gss_OID oid, OM_uint32 l, void *p) {
	oid->length = l;
	oid->elements = p;
}

void helper_gss_OID_desc_get_elements(gss_OID oid, OM_uint32 *l, char **p) {
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

func (lib *Lib) GSS_C_NT_USER_NAME() *OID {
	return &OID{
		Lib:       lib,
		C_gss_OID: C._GSS_C_NT_USER_NAME,
	}
}

func (lib *Lib) GSS_C_NT_MACHINE_UID_NAME() *OID {
	return &OID{
		Lib:       lib,
		C_gss_OID: C._GSS_C_NT_MACHINE_UID_NAME,
	}
}

func (lib *Lib) GSS_C_NT_STRING_UID_NAME() *OID {
	return &OID{
		Lib:       lib,
		C_gss_OID: C._GSS_C_NT_MACHINE_UID_NAME,
	}
}

func (lib *Lib) GSS_C_NT_HOSTBASED_SERVICE_X() *OID {
	return &OID{
		Lib:       lib,
		C_gss_OID: C._GSS_C_NT_HOSTBASED_SERVICE_X,
	}
}

func (lib *Lib) GSS_C_NT_HOSTBASED_SERVICE() *OID {
	return &OID{
		Lib:       lib,
		C_gss_OID: C._GSS_C_NT_HOSTBASED_SERVICE,
	}
}

func (lib *Lib) GSS_C_NT_ANONYMOUS() *OID {
	return &OID{
		Lib:       lib,
		C_gss_OID: C._GSS_C_NT_ANONYMOUS,
	}
}

func (lib *Lib) GSS_C_NT_EXPORT_NAME() *OID {
	return &OID{
		Lib:       lib,
		C_gss_OID: C._GSS_C_NT_EXPORT_NAME,
	}
}

func (lib *Lib) GSS_KRB5_NT_PRINCIPAL_NAME() *OID {
	return &OID{
		Lib:       lib,
		C_gss_OID: C._GSS_KRB5_NT_PRINCIPAL_NAME,
	}
}

// MakeOIDBytes makes an OID encapsulating a byte slice. Note that it does not duplicate
// the data, points to it directly
func (lib *Lib) MakeOIDBytes(data []byte) (*OID, error) {
	oid := lib.NewOID()

	s := C.malloc(C.gss_OID_size) // s for struct
	if s == nil {
		return nil, ErrMallocFailed
	}
	C.memset(s, 0, C.gss_OID_size)

	l := C.size_t(len(data))
	e := C.malloc(l) // c for contents
	if e == nil {
		return nil, ErrMallocFailed
	}
	C.memcpy(e, (unsafe.Pointer)(&data[0]), l)

	oid.C_gss_OID = C.gss_OID(s)
	oid.alloc = allocMalloc

	// because of the alignment issues I can't access o.oid's fields from go,
	// so invoking a C function to do the same as:
	// oid.C_gss_OID.length = l
	// oid.C_gss_OID.elements = c
	C.helper_gss_OID_desc_set_elements(oid.C_gss_OID, C.OM_uint32(l), e)

	return oid, nil
}

func (lib *Lib) MakeOIDString(data string) (*OID, error) {
	return lib.MakeOIDBytes([]byte(data))
}

func (oid OID) Bytes() []byte {
	var l C.OM_uint32
	var p *C.char

	C.helper_gss_OID_desc_get_elements(oid.C_gss_OID, &l, &p)

	return C.GoBytes(unsafe.Pointer(p), C.int(l))
}

// Release safely frees the contents of an OID if it's allocated with malloc by
// MakeOIDBytes.
func (oid *OID) Release() error {
	if oid == nil || oid.C_gss_OID == nil {
		return nil
	}

	switch oid.alloc {
	case allocMalloc:
		// same as with get and set, use a C helper to free(oid.C_gss_OID.elements)
		C.helper_gss_OID_desc_free_elements(oid.C_gss_OID)
		C.free(unsafe.Pointer(oid.C_gss_OID))
		oid.C_gss_OID = nil
		oid.alloc = allocNone
	}

	return nil
}

// Used in testing to check that we get back something reasonable form the C
// converted form
func DebugStringCGssOIDDesc(oid C.gss_OID_desc) string {
	var l C.OM_uint32
	var p *C.char

	C.helper_gss_OID_desc_get_elements(&oid, &l, &p)

	return fmt.Sprintf(`{%d %p}:"%x"`, l, p, C.GoStringN(p, C.int(l)))
}
