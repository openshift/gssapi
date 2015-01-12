// Copyright 2014-2015 Apcera Inc. All rights reserved.

package gssapi

// This file provides GSSKRB5 methods

/*
#include <gssapi/gssapi.h>
#include <stdlib.h>

OM_uint32
wrap_gsskrb5_register_acceptor_identity(void *fp,
	const char *identity)
{
	return ((OM_uint32(*) (
		const char *identity)
	) fp)(
		identity);
}

*/
import "C"

import (
// "unsafe"
)

func (lib *Lib) KRB5RegisterAcceptorIdentity(filename string) (err error) {
	/*
		fn := C.CString(filename)
		defer C.free(unsafe.Pointer(fn))
		maj := C.wrap_gsskrb5_register_acceptor_identity(lib.Fp_gsskrb5_register_acceptor_identity,
			fn)

		err = lib.MakeError(maj, 0).GoError()
		if err != nil {
			return err
		}
	*/
	panic("Not implemented")
	return nil
}
