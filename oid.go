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

OM_uint32
wrap_gss_oid_to_str(void *fp,
	OM_uint32 *minor_status,
	gss_OID oid,
	gss_buffer_t oid_str)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		gss_OID,
		gss_buffer_t)) fp)(
			minor_status, oid, oid_str);
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

// Equal compares this OID to another one, using the gss_oid_equal api
func (oid *OID) Equal(other *OID) bool {
	if oid == nil || other == nil {
		return false
	}

	//TODO: properly implement calling gss_oid_equal
	// isEqual := C.wrap_gss_oid_equal(oid.Fp_gss_oid_equal,
	//	oid.C_gss_OID, other.C_gss_OID)
	// return isEqual != 0

	return oid.String() == other.String()
}

// Buffer returns a string representation of the OID, as a gssapi.Buffer Unlike
// other wrapped gss types, OIDs to not contain a Lib backreference, the lib
// parameter provides that
func (oid OID) Buffer() (b *Buffer, err error) {
	b = oid.NewBuffer(true)

	var min C.OM_uint32
	maj := C.wrap_gss_oid_to_str(oid.Fp_gss_oid_to_str,
		&min, oid.C_gss_OID, b.C_gss_buffer_t)
	err = oid.MakeError(maj, min).GoError()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (oid OID) String() string {
	b, _ := oid.Buffer()
	defer b.Release()

	return b.String()
}

// Used in testing to check that we get back something reasonable form the C
// converted form
func DebugStringCGssOIDDesc(oid C.gss_OID_desc) string {
	var l C.int
	var p *C.char

	C.helper_gss_OID_desc_get_bytes(&oid, &l, &p)

	return fmt.Sprintf(`{%d %p}:"%x"`, l, p, C.GoStringN(p, l))
}
