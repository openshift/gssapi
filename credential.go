// Copyright 2013 Apcera Inc. All rights reserved.

package gssapi

// This file provides GSSCredential methods

/*
#include <gssapi/gssapi.h>

OM_uint32
wrap_gss_acquire_cred(void *fp,
	OM_uint32 * minor_status,
	const gss_name_t desired_name,
	OM_uint32 time_req,
	const gss_OID_set desired_mechs,
	gss_cred_usage_t cred_usage,
	gss_cred_id_t * output_cred_handle,
	gss_OID_set * actual_mechs,
	OM_uint32 * time_rec)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		const gss_name_t,
		OM_uint32,
		const gss_OID_set,
		gss_cred_usage_t,
		gss_cred_id_t *,
		gss_OID_set *,
		OM_uint32 *)
	) fp)(
		minor_status,
		desired_name,
		time_req,
		desired_mechs,
		cred_usage,
		output_cred_handle,
		actual_mechs,
		time_rec);
}

OM_uint32
wrap_gss_release_cred(void *fp,
	OM_uint32 * minor_status,
	gss_cred_id_t * cred_handle)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		gss_cred_id_t *)
	) fp)(
		minor_status,
		cred_handle);
}

*/
import "C"

import (
	// "unsafe"
	"time"
)

func (lib *Lib) NewCredId() *CredId {
	return &CredId{
		Lib: lib,
	}
}

func (lib *Lib) GSS_C_NO_CREDENTIAL() *CredId {
	return lib.NewCredId()
}

// Note that the actualMechs MUST be released by the caller
func (lib *Lib) AcquireCred(desiredName *Name, timeReq time.Duration,
	desiredMechs *OIDSet, credUsage CredUsage) (outputCredHandle *CredId,
	actualMechs *OIDSet, timeRec time.Duration, err error) {

	min := C.OM_uint32(0)
	actualMechs = lib.NewOIDSet()
	outputCredHandle = lib.NewCredId()
	timerec := C.OM_uint32(0)

	maj := C.wrap_gss_acquire_cred(lib.Fp_gss_acquire_cred,
		&min,
		desiredName.C_gss_name_t,
		C.OM_uint32(timeReq.Seconds()),
		desiredMechs.C_gss_OID_set,
		C.gss_cred_usage_t(credUsage),
		&outputCredHandle.C_gss_cred_id_t,
		&actualMechs.C_gss_OID_set,
		&timerec)

	err = lib.MakeError(maj, min).GoError()
	if err != nil {
		return nil, nil, 0, err
	}

	return outputCredHandle, actualMechs, time.Duration(timerec) * time.Second, nil
}

func (c *CredId) Release() error {
	if c == nil || c.C_gss_cred_id_t == nil {
		return nil
	}

	min := C.OM_uint32(0)
	maj := C.wrap_gss_release_cred(c.Fp_gss_release_cred,
		&min,
		&c.C_gss_cred_id_t)

	return c.MakeError(maj, min).GoError()
}

// TODO: gss_acquire_cred
// TODO: gss_add_cred
// TODO: gss_inquire_cred
// TODO: gss_inquire_cred_by_mech
