// Copyright 2013-2015 Apcera Inc. All rights reserved.

//+build servicetest

package test

// test the credentials APIs with a keytab, configured against a real KDC

import (
	"testing"

	"github.com/apcera/gssapi"
)

// Assumes that the keytab is
func TestCredential(t *testing.T) {
	if !c.RunAsService {
		t.Skip()
	}

	if c.ServiceName == "" {
		t.Fatal("Need a --service-name")
	}

	nameBuf, err := c.MakeBufferString(c.ServiceName)
	if err != nil {
		t.Fatal(err)
	}
	defer nameBuf.Release()

	name, err := nameBuf.Name(c.GSS_KRB5_NT_PRINCIPAL_NAME())
	if err != nil {
		t.Fatal(err)
	}
	defer name.Release()

	cred, actualMechs, _, err := c.AcquireCred(name,
		gssapi.GSS_C_INDEFINITE, c.GSS_C_NO_OID_SET(), gssapi.GSS_C_ACCEPT)
	defer cred.Release()
	defer actualMechs.Release()
	if err != nil {
		t.Fatal(err)
	}
}
