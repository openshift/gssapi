// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

// Side-note: gss_const_name_t is defined in RFC5587 as a bug-fix over RFC2744,
// since "const gss_name_t foo" says that the foo pointer is const, not the item
// pointed to is const.  Ideally, we'd be able to detect that, or have a macro
// which indicates availability of the 5587 extensions.  Instead, we're stuck with
// the ancient system GSSAPI headers on MacOS not supporting this.
//
// Choosing between "correctness" on the target platform and losing that for others,
// I've chosen to pull in /opt/local/include for MacPorts on MacOS; that should get
// us a functioning type; it's a pointer, at the ABI level the typing doesn't matter,
// so once we compile we're good.  If modern (correct) headers are available in other
// locations, just add them to the search path for the relevant OS below.
//
// Using "MacPorts" on MacOS gives us: -I/opt/local/include
// Using "brew" on MacOS gives us: -I/usr/local/opt/heimdal/include

/*
#cgo darwin CFLAGS: -I/opt/local/include -I/usr/local/opt/heimdal/include
#include <gssapi.h>
#include <stdio.h>

OM_uint32
wrap_gss_display_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t input_name,
	gss_buffer_t output_name_buffer,
	gss_OID *output_name_type)
{
	return ((OM_uint32(*)(
		OM_uint32 *, const gss_name_t, gss_buffer_t, gss_OID *)
	)fp)(
		minor_status, input_name, output_name_buffer, output_name_type);
}

OM_uint32
wrap_gss_compare_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t name1,
	const gss_name_t name2,
	int * name_equal)
{
	return ((OM_uint32(*)(
		OM_uint32 *, const gss_name_t, const gss_name_t, int *)
	)fp)(
		minor_status, name1, name2, name_equal);
}

OM_uint32
wrap_gss_release_name(void *fp,
	OM_uint32 *minor_status,
	gss_name_t *input_name)
{
	return ((OM_uint32(*)(
		OM_uint32 *, gss_name_t *)
	)fp)(
		minor_status, input_name);
}

OM_uint32
wrap_gss_inquire_mechs_for_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t input_name,
	gss_OID_set *mech_types)
{
	return ((OM_uint32(*)(
		OM_uint32 *, const gss_name_t, gss_OID_set *)
	)fp)(
		minor_status, input_name, mech_types);
}

OM_uint32
wrap_gss_inquire_names_for_mech(void *fp,
	OM_uint32 *minor_status,
	const gss_OID mechanism,
	gss_OID_set * name_types)
{
	return ((OM_uint32(*)(
		OM_uint32 *, const gss_OID, gss_OID_set *)
	)fp)(
		minor_status, mechanism, name_types);
}

OM_uint32
wrap_gss_canonicalize_name(void *fp,
	OM_uint32 *minor_status,
	gss_const_name_t input_name,
	const gss_OID mech_type,
	gss_name_t *output_name)
{
	return ((OM_uint32(*)(
		OM_uint32 *, gss_const_name_t, const gss_OID, gss_name_t *)
	)fp)(
		minor_status, input_name, mech_type, output_name);
}

OM_uint32
wrap_gss_export_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t input_name,
	gss_buffer_t exported_name)
{
	OM_uint32 maj;

	maj = ((OM_uint32(*)(
		OM_uint32 *, const gss_name_t, gss_buffer_t)
	)fp)(
		minor_status, input_name, exported_name);

	return maj;
}

OM_uint32
wrap_gss_duplicate_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t src_name,
	gss_name_t *dest_name)
{
	return ((OM_uint32(*)(
		OM_uint32 *, const gss_name_t, gss_name_t *)
	)fp)(
		minor_status, src_name, dest_name);
}

*/
import "C"

func (lib *Lib) NewName() *Name {
	return &Name{
		Lib: lib,
	}
}

func (lib *Lib) GSS_C_NO_NAME() *Name {
	return lib.NewName()
}

// Name-Types.  These are standardized in the RFCs.  The library requires that
// a given name be usable for resolution, but it's typically a macro, there's
// no guarantee about the name exported from the library.  But since they're
// static, and well-defined, we can just define them ourselves.

// RFC2744-mandated values, mapping from as-near-as-possible to cut&paste
func (lib *Lib) GSS_C_NT_USER_NAME() *OID {
	return lib.MakeOIDString("\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x01")
}
func (lib *Lib) GSS_C_NT_MACHINE_UID_NAME() *OID {
	return lib.MakeOIDString("\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x02")
}
func (lib *Lib) GSS_C_NT_STRING_UID_NAME() *OID {
	return lib.MakeOIDString("\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x03")
}
func (lib *Lib) GSS_C_NT_HOSTBASED_SERVICE_X() *OID {
	return lib.MakeOIDString("\x2b\x06\x01\x05\x06\x02")
}
func (lib *Lib) GSS_C_NT_HOSTBASED_SERVICE() *OID {
	return lib.MakeOIDString("\x2a\x86\x48\x86\xf7\x12\x01\x02\x01\x04")
}
func (lib *Lib) GSS_C_NT_ANONYMOUS() *OID {
	return lib.MakeOIDString("\x2b\x06\x01\x05\x06\x03") // original had \01
}
func (lib *Lib) GSS_C_NT_EXPORT_NAME() *OID {
	return lib.MakeOIDString("\x2b\x06\x01\x05\x06\x04")
}

// from gssapi_krb5.h: This name form shall be represented by the Object
// Identifier {iso(1) member-body(2) United States(840) mit(113554) infosys(1)
// gssapi(2) krb5(2) krb5_name(1)}.  The recommended symbolic name for this
// type is "GSS_KRB5_NT_PRINCIPAL_NAME".
func (lib *Lib) GSS_KRB5_NT_PRINCIPAL_NAME() *OID {
	return lib.MakeOIDString("\x2a\x86\x48\x86\xf7\x12\x01\x02\x02\x01")
}

// Release frees the memory associated with an internal representation of the
// name.
func (n *Name) Release() error {
	if n == nil || n.C_gss_name_t == nil {
		return nil
	}
	var min C.OM_uint32
	maj := C.wrap_gss_release_name(n.Fp_gss_release_name, &min, &n.C_gss_name_t)
	err := n.MakeError(maj, min).GoError()
	if err == nil {
		n.C_gss_name_t = nil
	}
	return err
}

// Equal tests 2 names for semantic equality (refer to the same entity)
func (n Name) Equal(other Name) (equal bool, err error) {
	var min C.OM_uint32
	var isEqual C.int

	maj := C.wrap_gss_compare_name(n.Fp_gss_compare_name, &min,
		n.C_gss_name_t, other.C_gss_name_t, &isEqual)
	err = n.MakeError(maj, min).GoError()
	if err != nil {
		return false, err
	}

	return isEqual != 0, nil
}

// Display "allows an application to obtain a textual representation of an
// opaque internal-form name for display purposes"
func (n Name) Display() (name string, oid C.gss_OID, err error) {
	var min C.OM_uint32
	b := n.NewBuffer(true)

	maj := C.wrap_gss_display_name(n.Fp_gss_display_name, &min,
		n.C_gss_name_t, b.C_gss_buffer_t, &oid)

	err = n.MakeError(maj, min).GoError()
	if b.C_gss_buffer_t == nil {
		return "", nil, err
	}
	defer b.Release()

	return b.String(), oid, err
}

// Go-friendly version of Display ("" on error)
func (n Name) String() string {
	s, _, _ := n.Display()
	return s
}

// Canonicalize returns a copy of this name, canonicalized for the specified
// mechanism
func (n Name) Canonicalize(mech_type *OID) (canonical *Name, err error) {
	canonical = &Name{
		Lib: n.Lib,
	}

	var min C.OM_uint32
	maj := C.wrap_gss_canonicalize_name(n.Fp_gss_canonicalize_name, &min,
		n.C_gss_name_t, mech_type.C_gss_OID, &canonical.C_gss_name_t)
	err = n.MakeError(maj, min).GoError()
	if err != nil {
		return nil, err
	}

	return canonical, nil
}

// Duplicate creates a new independent imported name; after this, both the original and
// the duplicate will need to be .Released().
func (n *Name) Duplicate() (duplicate *Name, err error) {
	duplicate = &Name{
		Lib: n.Lib,
	}

	var min C.OM_uint32
	maj := C.wrap_gss_duplicate_name(n.Fp_gss_duplicate_name, &min,
		n.C_gss_name_t, &duplicate.C_gss_name_t)
	err = n.MakeError(maj, min).GoError()
	if err != nil {
		return nil, err
	}

	return duplicate, nil
}

// Export makes a text (Buffer) version from an internal representation
func (n *Name) Export() (b *Buffer, err error) {
	b = n.NewBuffer(true)

	var min C.OM_uint32
	maj := C.wrap_gss_export_name(n.Fp_gss_export_name, &min,
		n.C_gss_name_t, b.C_gss_buffer_t)
	err = n.MakeError(maj, min).GoError()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// InquireMechs returns the set of mechanisms supported by the GSS-API
// implementation that may be able to process the specified name
func (n *Name) InquireMechs() (oids *OIDSet, err error) {
	oidset := n.NewOIDSet()
	if err != nil {
		return nil, err
	}

	var min C.OM_uint32
	maj := C.wrap_gss_inquire_mechs_for_name(n.Fp_gss_inquire_mechs_for_name, &min,
		n.C_gss_name_t, &oidset.C_gss_OID_set)
	err = n.MakeError(maj, min).GoError()
	if err != nil {
		return nil, err
	}

	return oidset, nil
}

/*
TODO: provide API for gss_inquire_mechs_for_name and gss_inquire_names_for_mech

OM_uint32
wrap_gss_inquire_mechs_for_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t input_name,
	gss_OID_set *mech_types

OM_uint32
wrap_gss_inquire_names_for_mech(void *fp,
	OM_uint32 *minor_status,
	const gss_OID mechanism,
	gss_OID_set * name_types

*/
