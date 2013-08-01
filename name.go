// Copyright 2013 Apcera Inc. All rights reserved.

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

/*
#cgo darwin CFLAGS: -I/opt/local/include
#include <gssapi/gssapi.h>

OM_uint32
wrap_gss_canonicalize_name(void *fp,
	OM_uint32 *minor_status,
	gss_const_name_t input_name,
	const gss_OID mech_type,
	gss_name_t *output_name
) {
	return (
		(
		OM_uint32(*)(OM_uint32 *, gss_const_name_t, const gss_OID, gss_name_t *)
		)fp
	)(minor_status, input_name, mech_type, output_name);
}

OM_uint32
wrap_gss_compare_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t name1,
	const gss_name_t name2,
	int * name_equal
) {
	return (
		(
		OM_uint32(*)(OM_uint32 *, const gss_name_t, const gss_name_t, int *)
		)fp
	)(minor_status, name1, name2, name_equal);
}

OM_uint32
wrap_gss_display_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t input_name,
	gss_buffer_t output_name_buffer,
	gss_OID *output_name_type
) {
	return (
		(
		OM_uint32(*)(OM_uint32 *, const gss_name_t, gss_buffer_t, gss_OID *)
		)fp
	)(minor_status, input_name, output_name_buffer, output_name_type);
}

OM_uint32
wrap_gss_duplicate_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t src_name,
	gss_name_t *dest_name
) {
	return (
		(
		OM_uint32(*)(OM_uint32 *, const gss_name_t, gss_name_t *)
		)fp
	)(minor_status, src_name, dest_name);
}

OM_uint32
wrap_gss_export_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t input_name,
	gss_buffer_t exported_name
) {
	return (
		(
		OM_uint32(*)(OM_uint32 *, const gss_name_t, gss_buffer_t)
		)fp
	)(minor_status, input_name, exported_name);
}

OM_uint32
wrap_gss_import_name(void *fp,
	OM_uint32 *minor_status,
	const gss_buffer_t input_name_buffer,
	const gss_OID input_name_type,
	gss_name_t *output_name
) {
	return (
		(
		OM_uint32(*)(OM_uint32 *, const gss_buffer_t, const gss_OID, gss_name_t *)
		)fp
	)(minor_status, input_name_buffer, input_name_type, output_name);
}

OM_uint32
wrap_gss_inquire_mechs_for_name(void *fp,
	OM_uint32 *minor_status,
	const gss_name_t input_name,
	gss_OID_set *mech_types
) {
	return (
		(
		OM_uint32(*)(OM_uint32 *, const gss_name_t, gss_OID_set *)
		)fp
	)(minor_status, input_name, mech_types);
}

OM_uint32
wrap_gss_inquire_names_for_mech(void *fp,
	OM_uint32 *minor_status,
	const gss_OID mechanism,
	gss_OID_set * name_types
) {
	return (
		(
		OM_uint32(*)(OM_uint32 *, const gss_OID, gss_OID_set *)
		)fp
	)(minor_status, mechanism, name_types);
}

OM_uint32
wrap_gss_release_name(void *fp,
	OM_uint32 *minor_status,
	gss_name_t *input_name
) {
	return (
		(
		OM_uint32(*)(OM_uint32 *, gss_name_t *)
		)fp
	)(minor_status, input_name);
}
*/
import "C"

// Called at library load time.
// Return false if could not resolve all symbols.
func (lib *GssapiLib) populateNameFunctions() bool {
	var ok bool
	if lib.fp_gss_canonicalize_name, ok = lib.symbolResolveOne("gss_canonicalize_name"); !ok {
		return false
	}
	if lib.fp_gss_compare_name, ok = lib.symbolResolveOne("gss_compare_name"); !ok {
		return false
	}
	if lib.fp_gss_display_name, ok = lib.symbolResolveOne("gss_display_name"); !ok {
		return false
	}
	if lib.fp_gss_duplicate_name, ok = lib.symbolResolveOne("gss_duplicate_name"); !ok {
		return false
	}
	if lib.fp_gss_export_name, ok = lib.symbolResolveOne("gss_export_name"); !ok {
		return false
	}
	if lib.fp_gss_import_name, ok = lib.symbolResolveOne("gss_import_name"); !ok {
		return false
	}
	if lib.fp_gss_inquire_mechs_for_name, ok = lib.symbolResolveOne("gss_inquire_mechs_for_name"); !ok {
		return false
	}
	if lib.fp_gss_inquire_names_for_mech, ok = lib.symbolResolveOne("gss_inquire_names_for_mech"); !ok {
		return false
	}
	if lib.fp_gss_release_name, ok = lib.symbolResolveOne("gss_release_name"); !ok {
		return false
	}

	return true
}

// Almost had Type suffix on the functions taking nametype as a parameter, with
// GSS_C_NO_OID wrappers that are slightly simpler, but removed that since the
// client of this library should know/understand name types and think through
// it carefully.

// Wrapper for the C gss_name_t, bundling in the library reference, with an
// exported type name.
type NameImported struct {
	lib *GssapiLib
	name C.gss_name_t
}

// Let GSSAPI free the storage associated with an internalized name.
func (name *NameImported) Release() error {
	var min C.OM_uint32
	maj := C.wrap_gss_release_name(name.lib.fp_gss_release_name, &min, &name.name)
	err := name.lib.MakeError(maj, min)
	if err == nil {
		name.name = GSS_C_NO_NAME
	}
	return err
}

// Closest to the C API: import a name in buffer format and a given type.
// This one exists as a lib method mostly for symmetry with NameImportString below.
// Normally, if you have a name as a GssBuffer, call Import() as a method on that.
func (lib *GssapiLib) NameImportBuffer(name GssBuffer, nametype C.gss_OID) (*NameImported, error) {
	var min C.OM_uint32
	var result C.gss_name_t
	maj := C.wrap_gss_import_name(lib.fp_gss_import_name, &min,
		name.buffer, nametype, &result)
	return &NameImported{lib: lib, name: result}, lib.MakeError(maj, min)
}

// Import a name from a string and a type.
func (lib *GssapiLib) NameImportString(name string, nametype C.gss_OID) (*NameImported, error) {
	bufName := lib.BufferString(name)
	return lib.NameImportBuffer(bufName, nametype)
}

// Given a GssBuffer representing a name, method to import it requiring only a
// type.
func (name GssBuffer) Import(nametype C.gss_OID) (*NameImported, error) {
	var min C.OM_uint32
	var result C.gss_name_t
	maj := C.wrap_gss_import_name(name.lib.fp_gss_import_name, &min,
		name.buffer, nametype, &result)
	return &NameImported{lib: name.lib, name: result}, name.lib.MakeError(maj, min)
}

// Compare two imported names, using Go-ish naming; "Compare" normally means
// more of a <,0,> comparator, where gss_compare_name only gives us an equality
// check.
func (name *NameImported) Equal(other *NameImported) bool {
	var min C.OM_uint32
	var isEqual C.int
	maj := C.wrap_gss_compare_name(name.lib.fp_gss_import_name, &min,
		name.name, other.name, &isEqual)
	if NewStatus(maj, min).IsError() {
		return false
	}
	// "non-zero - names refer to same entity"
	return isEqual != 0
}

// "Allows an application to obtain a textual representation of an opaque
// internal-form name for display purposes"
func (name *NameImported) Display() (string, error) {
	var min C.OM_uint32
	var cBuf C.gss_buffer_t
	var oid C.gss_OID
	maj := C.wrap_gss_display_name(name.lib.fp_gss_display_name, &min,
		name.name, cBuf, &oid)
	status := name.lib.MakeError(maj, min)
	if cBuf == GSS_C_NO_BUFFER {
		return "", status
	}
	buf := GssBuffer{lib: name.lib, buffer: cBuf}
	defer buf.Release()
	return buf.String(), status
}

// Assume that Display doesn't error?  Get to string-ify simply.
func (name *NameImported) String() string { s, _ := name.Display(); return s }

func (name *NameImported) Canonicalize(mech_type C.gss_OID) (*NameImported, error) {
	var min C.OM_uint32
	result := &NameImported{lib: name.lib}
	maj := C.wrap_gss_canonicalize_name(name.lib.fp_gss_canonicalize_name, &min,
		name.name, mech_type, &result.name)
	return result, name.lib.MakeError(maj, min)
}

// Create a new independent imported name; after this, both the original and
// the duplicate will need to be .Released().
func (name *NameImported) Duplicate() (*NameImported, error) {
	var min C.OM_uint32
	dup := &NameImported{lib: name.lib}
	maj := C.wrap_gss_duplicate_name(name.lib.fp_gss_duplicate_name, &min,
		name.name, &dup.name)
	return dup, name.lib.MakeError(maj, min)
}

func (name *NameImported) Export() (GssBuffer, error) {
	var min C.OM_uint32
	buf := GssBuffer{lib: name.lib}
	maj := C.wrap_gss_export_name(name.lib.fp_gss_export_name, &min,
		name.name, buf.buffer)
	return buf, name.lib.MakeError(maj, min)
}

/*
FIXME: provide API for:
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
