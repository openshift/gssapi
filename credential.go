// Copyright 2013 Apcera Inc. All rights reserved.

package gssapi

// This file provides GSSCredential methods

/*
#include <gssapi/gssapi.h>

void
wrap_gss_acquire_cred() {
}

*/
import "C"

func (gss *GssapiLib) GssAcquireCred() error {
	f, err := gss.DlSym("gss_acquire_cred")
	if err != nil {
		return err
	}
	_ = f
	C.wrap_gss_acquire_cred()
	return nil
}

// TODO: gss_acquire_cred
// TODO: gss_add_cred
// TODO: gss_inquire_cred
// TODO: gss_inquire_cred_by_mech
// TODO: gss_release_cred
