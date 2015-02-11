// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

import (
	"bytes"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	l, err := testLoad()
	if err != nil {
		t.Error(err)
		return
	}
	defer l.Unload()

	b := l.NewBuffer(false)
	defer b.Release()

	if b == nil {
		t.Errorf("Got nil, expected non-nil")
		return
	}
	if b.Lib != l {
		t.Errorf("b.Lib didn't get set correctly, got %p, expected %p", b.Lib, l)
		return
	}
	if b.C_gss_buffer_t == nil {
		t.Errorf("Got nil buffer, expected non-nil")
		return
	}
	if b.String() != "" {
		t.Errorf(`String(): got %q, expected ""`, b.String())
		return
	}
}

// Also tests MakeBufferBytes, implicitly
func TestMakeBufferString(t *testing.T) {
	l, err := testLoad()
	if err != nil {
		t.Error(err)
		return
	}
	defer l.Unload()

	test := "testing"
	b := l.MakeBufferString(test)
	defer b.Release()

	if b == nil {
		t.Errorf("Got nil, expected non-nil")
		return
	}
	if b.Lib != l {
		t.Errorf("b.Lib didn't get set correctly, got %p, expected %p", b.Lib, l)
		return
	}
	if b.String() != test {
		t.Errorf("Got %q, expected %q", b.String(), test)
		return
	} else if !bytes.Equal(b.Bytes(), []byte(test)) {
		t.Fatalf("Got '%v'; expected '%v'", b.Bytes(), []byte(test))
	}
}
