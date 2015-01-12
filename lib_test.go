// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

import (
	"testing"
)

func TestLoadLib(t *testing.T) {
	l, err := LoadLib()

	if err != nil {
		t.Fatal(err)
	}

	if l.Fp_gss_export_name == nil {
		t.Error("Fp_gss_export_name did not get initialized")
		return
	}

	// TODO: maybe use reflect to enumerate all Fp's

	defer l.Unload()
}
