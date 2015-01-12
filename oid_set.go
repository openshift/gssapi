// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

/*
#include <gssapi/gssapi.h>

OM_uint32
wrap_gss_create_empty_oid_set(void *fp,
	OM_uint32 *minor_status,
	gss_OID_set * set)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		gss_OID_set *)) fp)(
			minor_status,
			set);
}

OM_uint32
wrap_gss_release_oid_set(void *fp,
	OM_uint32 *minor_status,
	gss_OID_set * set)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		gss_OID_set *)) fp)(
			minor_status, set);
}

OM_uint32
wrap_gss_add_oid_set_member(void *fp,
	OM_uint32 *minor_status,
	const gss_OID member_oid,
	gss_OID_set * set)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		const gss_OID,
		gss_OID_set *)) fp)(
			minor_status, member_oid, set);
}

OM_uint32
wrap_gss_test_oid_set_member(void *fp,
	OM_uint32 *minor_status,
	const gss_OID member_oid,
	const gss_OID_set set,
	int * present)
{
	return ((OM_uint32(*) (
		OM_uint32 *,
		const gss_OID,
		const gss_OID_set,
		int *)) fp)(
			minor_status, member_oid, set, present);
}

gss_OID
oid_set_member(
	gss_OID_set set,
	int index)
{
	return set->elements + (index * sizeof(*set->elements));
}

*/
import "C"

import (
	"fmt"
)

func (lib *Lib) NewOIDSet() *OIDSet {
	return &OIDSet{
		Lib: lib,
		// C_gss_OID_set: (C.gss_OID_set)(unsafe.Pointer(nil)),
	}
}

func (lib Lib) GSS_C_NO_OID_SET() *OIDSet {
	return lib.NewOIDSet()
}

// CreateEmptyOIDSet makes an empty OIDSet
func (lib *Lib) CreateEmptyOIDSet() (s *OIDSet, err error) {
	return lib.CreateOIDSet()
}

// CreateOIDSet makes an OIDSet prepopulated with the given OIDs.
func (lib *Lib) CreateOIDSet(oids ...*OID) (s *OIDSet, err error) {
	s = &OIDSet{
		Lib: lib,
	}

	var min C.OM_uint32
	maj := C.wrap_gss_create_empty_oid_set(s.Fp_gss_create_empty_oid_set,
		&min, &s.C_gss_OID_set)
	err = s.MakeError(maj, min).GoError()
	if err != nil {
		return nil, err
	}

	err = s.Add(oids...)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Release frees all C memory associated with an OIDSet.
func (s *OIDSet) Release() (err error) {
	if s == nil {
		return nil
	}

	var min C.OM_uint32
	maj := C.wrap_gss_release_oid_set(s.Fp_gss_release_oid_set, &min, &s.C_gss_OID_set)
	return s.MakeError(maj, min).GoError()
}

// Add adds OIDs to an OIDSet.
func (s *OIDSet) Add(oids ...*OID) (err error) {
	var min C.OM_uint32
	for _, oid := range oids {
		maj := C.wrap_gss_add_oid_set_member(s.Fp_gss_add_oid_set_member,
			&min, oid.C_gss_OID, &s.C_gss_OID_set)
		err = s.MakeError(maj, min).GoError()
		if err != nil {
			return err
		}
	}

	return nil
}

// Contains (gss_test_oid_set_member) checks if an OID is present OIDSet.
func (s *OIDSet) Contains(oid *OID) (present bool, err error) {
	var min C.OM_uint32
	var isPresent C.int

	maj := C.wrap_gss_test_oid_set_member(s.Fp_gss_test_oid_set_member,
		&min, oid.C_gss_OID, s.C_gss_OID_set, &isPresent)
	err = s.MakeError(maj, min).GoError()
	if err != nil {
		return false, err
	}

	return isPresent != 0, nil
}

// Returns a specific OID from the set. The memory will be released when the
// set itself is released

func (s *OIDSet) Get(index int) (oid *OID) {

	if index < 0 || index >= int(s.C_gss_OID_set.count) {
		panic(fmt.Errorf("index %d out of bounds", index))
	}

	return &OID{
		C_gss_OID: C.oid_set_member(s.C_gss_OID_set, C.int(index)),
	}
}
