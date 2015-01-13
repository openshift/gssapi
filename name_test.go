// Copyright 2014 Apcera Inc. All rights reserved.

package gssapi

import (
	"testing"
)

// Tests importing exporting names
func TestNameImportExport(t *testing.T) {
	l, err := LoadLib()
	if err != nil {
		t.Error(err)
		return
	}
	defer l.Unload()

	names := []string{
		`test@corp.apcera.com`,
		`test@corp.ApCera.com`,
		`test@CORP.APCERA.COM`,
	}

	makeName := func(n string) (name *Name) {
		b := l.MakeBufferString(n)
		if b == nil {
			t.Errorf("%q: Got nil, expected non-nil", n)
			return nil
		}
		defer b.Release()

		//name, err := b.Name(l.GSS_KRB5_NT_PRINCIPAL_NAME())
		name, err := b.Name(l.GSS_C_NT_HOSTBASED_SERVICE())
		if err != nil {
			t.Errorf("%q: Got error %q, expected nil", n, err.Error())
			return nil
		}
		if name == nil {
			t.Errorf("%q: Got nil, expected non-nil", n)
			return nil
		}
		return name
	}

	// Make the reference name
	n0 := makeName(names[0])
	defer n0.Release()

	// Make sure we can have the krb mechanism, and normalize the reference
	// name using it
	mechs, err := n0.InquireMechs()
	if err != nil {
		//TODO: need a better test for OS X since this InquireMechs doesn't
		// seem to work
		t.Skip(`Couldn't get mechs for`, names[0], `, error:`, err.Error())
		return
	}

	// This OID seems to be an avalable merch on linux
	kerbOID, err := l.MakeBufferString("{ 1 2 840 113554 1 2 2 }\x00").OID()
	if err != nil {
		t.Errorf("Got error %q, expected nil", err.Error())
		return
	}

	contains, err := mechs.Contains(kerbOID)
	if err != nil {
		t.Errorf("Got error %q, expected nil", err.Error())
		return
	}
	if !contains {
		t.Errorf("Expected true")
		return
	}

	makeNames := func(n string) (
		name *Name, canonical *Name, display string, exported *Buffer) {

		name = makeName(n)
		if name == nil {
			return nil, nil, "", nil
		}

		origDisplay, _, err := name.Display()
		if err != nil {
			t.Errorf("Got error %q, expected nil", err.Error())
			return nil, nil, "", nil
		}
		if origDisplay != n {
			t.Errorf("Got %q, expected %q", origDisplay, n)
			return nil, nil, "", nil
		}

		canonical, err = name.Canonicalize(kerbOID)
		if err != nil {
			t.Errorf("Got error %q, expected nil", err.Error())
			return nil, nil, "", nil
		}
		if canonical == nil {
			t.Errorf("Got nil, expected non-nil")
			return nil, nil, "", nil
		}

		display, _, err = canonical.Display()
		if err != nil {
			t.Errorf("Got error %q, expected nil", err.Error())
			return nil, nil, "", nil
		}

		exported, err = canonical.Export()
		if err != nil {
			t.Errorf("Got error %q, expected nil", err.Error())
			return nil, nil, "", nil
		}
		if exported == nil {
			t.Errorf("Got nil, expected non-nil")
			return nil, nil, "", nil
		}

		return name, canonical, display, exported
	}

	n0, _, d0, e0 := makeNames(names[0])
	if n0 == nil {
		t.Errorf("Got nil, expected non-nil")
		return
	}

	for _, n := range names {
		n, _, d, e := makeNames(n)
		if n == nil {
			t.Errorf("%s: Got nil, expected non-nil", n)
			return
		}
		if d != d0 {
			t.Errorf("%s: Got %q, expected %q", n, d, d0)
			return
		}
		if !e.Equal(e0) {
			t.Errorf("%s: Got %q, expected %q", n, e.String(), e0.String())
			return
		}
	}
}
