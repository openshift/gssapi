// Copyright 2013-2015 Apcera Inc. All rights reserved.

// The OS X gsspi does not seem to support gss_str_to_oid, so don't run the tests there
//+build !darwin

package gssapi

import (
	"testing"
)

func TestOID_C(t *testing.T) {
	l, err := LoadLib()
	if err != nil {
		t.Error(err)
		return
	}
	defer l.Unload()

	if !cOIDTest(l) {
		t.Error("C test failed")
		return
	}
}

func TestOID(t *testing.T) {
	l, err := LoadLib()
	if err != nil {
		t.Error(err)
		return
	}
	defer l.Unload()

	data := []*OID{
		l.GSS_C_NT_USER_NAME(),
		l.GSS_C_NT_MACHINE_UID_NAME(),
		l.GSS_C_NT_STRING_UID_NAME(),
		l.GSS_C_NT_HOSTBASED_SERVICE_X(),
		l.GSS_C_NT_HOSTBASED_SERVICE(),
		l.GSS_C_NT_ANONYMOUS(),
		l.GSS_C_NT_EXPORT_NAME(),
		l.GSS_KRB5_NT_PRINCIPAL_NAME(),
	}

	testOne := func(oid *OID) {
		b, err := oid.Buffer()
		if err != nil {
			t.Error(err)
			return
		}
		defer b.Release()

		str := b.String()
		if str == "" {
			t.Errorf(`Got "" expected non-empty`)
			return
		}

		newoid, err := b.OID()
		if err != nil {
			t.Error(err)
			return
		}

		newstr := newoid.String()
		if str != newstr {
			t.Errorf(`Got %q expected %q`, newstr, str)
			return
		}

		equal := newoid.Equal(oid)
		if !equal {
			t.Errorf(`expected %q to be equal to %q using gss_oid_equal`, newstr, str)
			return
		}
	}

	for _, oid := range data {
		testOne(oid)
	}
}
