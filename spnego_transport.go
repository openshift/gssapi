// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

import (
	"encoding/base64"
	"fmt"
	"net/http"
	_ "net/http/httputil"
	"strings"
)

type SPNEGOTransport struct {
	http.RoundTripper
	*Lib

	serviceName   *Name
	ctx           *CtxId
	authorization *Buffer
}

func (lib *Lib) NewSPNEGOTransport(serviceName string) (*SPNEGOTransport, error) {
	t := &SPNEGOTransport{
		Lib:          lib,
		RoundTripper: &http.Transport{},
	}

	namebuf, err := lib.MakeBufferString(serviceName)
	if err != nil {
		return nil, err
	}
	defer namebuf.Release()

	name, err := namebuf.Name(lib.GSS_KRB5_NT_PRINCIPAL_NAME)
	if err != nil {
		return nil, err
	}
	t.serviceName = name

	lib.Debug("New SPNEGOTransport for ", t.serviceName)
	return t, nil
}

func (t *SPNEGOTransport) Release() error {
	if t == nil {
		return nil
	}

	if err := t.serviceName.Release(); err != nil {
		return err
	}
	if err := t.ctx.Release(); err != nil {
		return err
	}
	if err := t.authorization.Release(); err != nil {
		return err
	}

	return nil
}

func (t *SPNEGOTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// verify that req does not contain an Authorization header of its own
	if req.Header.Get("Authorization") != "" {
		return nil, fmt.Errorf("Authorization header must not be set when using SPNEGOTransport")
	}

	// try HEAD request until get a clean 200, a cacheable authorization, or a failure

	negotiate := false
	sendToken := t.GSS_C_NO_BUFFER
	receiveToken := t.GSS_C_NO_BUFFER

	for {
		if !t.authorization.IsEmpty() {
			break
		}

		rcopy := *req
		rcopy.Method = "HEAD"
		resp, negotiate, receiveToken, err = t.doRoundTrip(&rcopy, sendToken)
		if err != nil {
			return nil, err
		}
		// out, _ := httputil.DumpResponse(resp, true)
		// t.Debug("<- SPNEGO RECEIVED:\n", string(out), "\n")
		resp.Body.Close()

		if !negotiate {
			break
		}

		// other outputs can be safely ignored, no need to release
		t.ctx, _, sendToken, _, _, err = t.InitSecContext(
			t.GSS_C_NO_CREDENTIAL,
			t.ctx, t.serviceName, t.GSS_C_NO_OID, 0, 0,
			t.GSS_C_NO_CHANNEL_BINDINGS, receiveToken)
		receiveToken.Release()
		if err != nil {
			e, ok := err.(*Error)
			if ok {
				ok = e.Major.ContinueNeeded()
			}
			if !ok {
				t.ctx.DeleteSecContext()
				sendToken.Release()
				return nil, err
			}
		} else {
			t.authorization = sendToken
		}
	}

	resp, _, _, err = t.doRoundTrip(req, t.authorization)
	return resp, err

	//FIXME: process the response and invalidate the token if needed, also update gssapi context if needed
}

func (t *SPNEGOTransport) doRoundTrip(req *http.Request, inputToken *Buffer) (
	resp *http.Response, negotiate bool, outputToken *Buffer, err error) {

	t.AddSPNEGONegotiate(req.Header, "Authorization", inputToken)

	// out, _ := httputil.DumpRequest(req, true)
	// t.Debug("-> SPNEGO SEND:\n", string(out), "\n")

	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, false, nil, err
	}
	negotiate, outputToken = t.CheckSPNEGONegotiate(resp.Header, "WWW-Authenticate")

	return resp, negotiate, outputToken, nil
}

func (lib *Lib) AddSPNEGONegotiate(h Header, name string, token *Buffer) {
	if name == "" {
		return
	}

	v := "Negotiate"
	if !token.IsEmpty() {
		data := token.Bytes()
		v = v + " " + base64.StdEncoding.EncodeToString(data)
	}
	h.Set(name, v)
}

func (lib *Lib) CheckSPNEGONegotiate(h Header, name string) (present bool, token *Buffer) {
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
