// Copyright 2013 Apcera Inc. All rights reserved.

// +build darwin linux

package gssapi

//#cgo LDFLAGS: -L/opt/local/lib -Wl,-search_paths_first -lgssapi_krb5 -lkrb5 -lk5crypto -lcom_err
//import "C"
// Uncomment the above line to directly link against the libraries, to avoid
// the dlopen layer.

/*
#include <gssapi/gssapi.h>
*/
import "C"

import (
	"sync"
	"unsafe"
)

// Encapsulates both the GSSAPI and the library dlopen()'d for it.
type GssapiLib struct {
	handle      unsafe.Pointer
	populate    sync.Once
	populateErr error

	// fp_ == function pointer, resolved at load time
	fp_gss_release_buffer unsafe.Pointer
	fp_gss_display_status unsafe.Pointer
	// name.go
	fp_gss_canonicalize_name      unsafe.Pointer
	fp_gss_compare_name           unsafe.Pointer
	fp_gss_display_name           unsafe.Pointer
	fp_gss_duplicate_name         unsafe.Pointer
	fp_gss_export_name            unsafe.Pointer
	fp_gss_import_name            unsafe.Pointer
	fp_gss_inquire_mechs_for_name unsafe.Pointer
	fp_gss_inquire_names_for_mech unsafe.Pointer
	fp_gss_release_name           unsafe.Pointer
}

type GssBuffer struct {
	lib    *GssapiLib
	buffer C.gss_buffer_t
}

func (buf GssBuffer) Bytes() []byte {
	return C.GoBytes(buf.buffer.value, C.int(buf.buffer.length))
}

func (buf GssBuffer) Release() Status {
	var min C.OM_uint32
	maj := buf.lib.gss_release_buffer(&min, buf.buffer)
	return NewStatus(maj, min)
}

func (buf GssBuffer) String() string {
	return C.GoStringN((*C.char)(buf.buffer.value), C.int(buf.buffer.length))
}
