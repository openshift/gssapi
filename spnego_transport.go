// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (lib *Lib) AddSPNEGONegotiate(h http.Header, name string, token *Buffer) {
	if name == "" {
		return
	}

	v := "Negotiate"
	if token.Length() != 0 {
		data := token.Bytes()
		v = v + " " + base64.StdEncoding.EncodeToString(data)
	}
	h.Set(name, v)
}

func (lib *Lib) CheckSPNEGONegotiate(h http.Header, name string) (present bool, token *Buffer) {
	var err error
	defer func() {
		if err != nil {
			lib.Debug(fmt.Sprintf("CheckSPNEGONegotiate: %v", err))
		}
	}()

	v := h.Get(name)
	if len(v) == 0 || !strings.HasPrefix(v, "Negotiate") {
		return false, nil
	}

	present = true
	tbytes, err := base64.StdEncoding.DecodeString(
		strings.TrimSpace(v[len("Negotiate"):]))
	if err != nil {
		return false, nil
	}
	if len(tbytes) > 0 {
		token, err = lib.MakeBufferBytes(tbytes)
		if err != nil {
			return false, nil
		}
	}

	return present, token
}
