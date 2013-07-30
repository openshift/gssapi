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

// TODO:
//  * decide on the encapsulating data structures
//  * write the gss_* low-level functions
//  * write the exported functions
