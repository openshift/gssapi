// Copyright 2013-2015 Apcera Inc. All rights reserved.

//+build servicetest

package test

// test the credentials APIs with a keytab, configured against a real KDC

import (
	"strings"
	"testing"

	"github.com/apcera/gssapi"
)

func prepareServiceName(t *testing.T) *gssapi.Name {
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

	name, err := nameBuf.Name(c.GSS_KRB5_NT_PRINCIPAL_NAME)
	if err != nil {
		t.Fatal(err)
	}
	if name.String() != c.ServiceName {
		t.Fatalf("name: got %q, expected %q", name.String(), c.ServiceName)
	}

	return name
}

func TestInquireMechsForName(t *testing.T) {
	name := prepareServiceName(t)
	defer name.Release()

	mechs, err := name.InquireMechs()
	if err != nil {
		t.Fatal(err)
	}
	defer mechs.Release()
	contains, _ := mechs.TestOIDSetMember(c.GSS_MECH_KRB5)
	if !contains {
		t.Fatalf("Expected mechs to contain %s, got %s",
			c.GSS_MECH_KRB5.DebugString(), mechs.DebugString())
	}
}

func TestCanonicalizeName(t *testing.T) {
	name := prepareServiceName(t)
	defer name.Release()

	name, err := name.Canonicalize(c.GSS_MECH_KRB5)
	if err != nil {
		t.Fatal(err)
	}
	defer name.Release()
	parts := strings.Split(name.String(), "@")
	if len(parts) != 2 || parts[0] != c.ServiceName {
		t.Fatalf("name: got %q, expected %q", name.String(), c.ServiceName+"@<domain>")
	}
}
