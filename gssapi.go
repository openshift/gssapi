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
	handle unsafe.Pointer
	populate sync.Once
	populateErr error

	// fp_ == function pointer, resolved at load time
	fp_gss_release_buffer unsafe.Pointer
}

type GssBuffer struct {
	lib *GssapiLib
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
