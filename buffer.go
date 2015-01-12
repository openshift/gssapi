// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

/*
#include <gssapi.h>
#include <string.h>

OM_uint32
wrap_any_gss_one_buffer(void *fp,
	OM_uint32 *minor_status,
	gss_buffer_t buf)
{
	return ((OM_uint32(*)(
		OM_uint32*,
		gss_buffer_t))fp) (minor_status, buf);
}

OM_uint32
wrap_gss_import_name(void *fp,
	OM_uint32 *minor_status,
	const gss_buffer_t input_name_buffer,
	const gss_OID input_name_type,
	gss_name_t *output_name)
{
	return ((OM_uint32(*)(
		OM_uint32 *,
		const gss_buffer_t,
		const gss_OID,
		gss_name_t *)) fp) (
			minor_status,
			input_name_buffer,
			input_name_type,
			output_name);
}

OM_uint32
wrap_gss_str_to_oid(void *fp,
	OM_uint32 *minor_status,
	gss_buffer_t oid_str,
	gss_OID *oid)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		gss_buffer_t,
		gss_OID*)) fp)(
			minor_status, oid_str, oid);
}

int
wrap_gss_buffer_equal(
	gss_buffer_t b1,
	gss_buffer_t b2)
{
	return
		b1 != NULL &&
		b2 != NULL &&
		b1->length == b2->length &&
		(memcmp(b1->value,b2->value,b1->length) == 0);
}

int
wrap_gss_buffer_empty(
	gss_buffer_t b)
{
	return b == NULL || b->length == 0;
}

*/
import "C"

import (
	"unsafe"
)

// NewBuffer returns an uninitialized (empty) Buffer
func (lib *Lib) NewBuffer(releasable bool) *Buffer {
	return &Buffer{
		Lib:            lib,
		C_gss_buffer_t: C.gss_buffer_t(&C.gss_buffer_desc{}),
		releasable:     releasable,
	}
}

func (lib *Lib) GSS_C_NO_BUFFER() *Buffer {
	return &Buffer{
		Lib: lib,
		// C_gss_buffer_t: C.GSS_C_NO_BUFFER,
		releasable: true,
	}
}

// MakeBufferBytes makes a Buffer encapsulating a byte slice
func (lib *Lib) MakeBufferBytes(content []byte) *Buffer {
	return &Buffer{
		Lib: lib,
		C_gss_buffer_t: C.gss_buffer_t(
			&C.gss_buffer_desc{
				length: C.size_t(len(content)),
				value:  (unsafe.Pointer)(&content[0]),
			}),
	}
}

// MakeBufferBytes makes a Buffer encapsulating the contents of a string
func (lib *Lib) MakeBufferString(content string) *Buffer {
	return lib.MakeBufferBytes([]byte(content))
}

// Release safely frees the contents of a Buffer. C.gss_buffer_t (and thus
// our Buffer) can come from Go or from the GSSAPI library; Those coming
// from GSSAPI must have been wrapped by us, so all API wrappers must set the
// releasable flag.
func (b *Buffer) Release() error {
	if b == nil || !b.releasable || b.C_gss_buffer_t == nil {
		return nil
	}

	var min C.OM_uint32
	maj := C.wrap_any_gss_one_buffer(b.Fp_gss_release_buffer, &min, b.C_gss_buffer_t)
	err := b.MakeError(maj, min).GoError()
	if err != nil {
		return err
	}

	b.releasable = false
	return nil
}

// Bytes returns teh contents of a Buffer as a byte slice
func (b Buffer) Bytes() []byte {
	return C.GoBytes(b.C_gss_buffer_t.value, C.int(b.C_gss_buffer_t.length))
}

// String returns the contents of a Buffer as a string
func (b Buffer) String() string {
	return C.GoStringN((*C.char)(b.C_gss_buffer_t.value), C.int(b.C_gss_buffer_t.length))
}

// Name converts a Buffer representing a name into a Name (internal
// opaque representation) using the specified nametype
func (b Buffer) Name(nametype *OID) (*Name, error) {
	var min C.OM_uint32
	var result C.gss_name_t

	maj := C.wrap_gss_import_name(b.Fp_gss_import_name, &min,
		b.C_gss_buffer_t, nametype.C_gss_OID, &result)
	err := b.MakeError(maj, min).GoError()
	if err != nil {
		return nil, err
	}

	return &Name{
		Lib:          b.Lib,
		C_gss_name_t: result,
	}, nil
}

// OID converts the buffer to an OID. Note that the OID is allocated by the C
// code, The only way to release the memory is to add it to an OIDSet, and then
// Release the entire OIDSet
func (b *Buffer) OID() (oid *OID, err error) {
	var min C.OM_uint32

	oid = &OID{
		Lib: b.Lib,
	}

	maj := C.wrap_gss_str_to_oid(b.Fp_gss_str_to_oid,
		&min, b.C_gss_buffer_t, &oid.C_gss_OID)

	err = b.MakeError(maj, min).GoError()
	if err != nil {
		return nil, err
	}

	return oid, nil
}

func (b *Buffer) Equal(other *Buffer) bool {
	isEqual := C.wrap_gss_buffer_equal(b.C_gss_buffer_t, other.C_gss_buffer_t)
	return isEqual != 0
}

func (b *Buffer) IsEmpty() bool {
	if b == nil {
		return true
	}

	isEmpty := C.wrap_gss_buffer_empty(b.C_gss_buffer_t)
	return isEmpty != 0
}
