// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

type SPNEGOTransport struct {
	serviceName *Name

	http.RoundTripper
	*Lib

	// Authorization value to use in the requests. Note that if requests are not challenged, that information is not cached
	ctx           *CtxId
	authorization *Buffer
}

func (lib *Lib) NewSPNEGOTransportRoundTripper(rt http.RoundTripper, serviceName string) (http.RoundTripper, error) {
	namebuf := lib.MakeBufferString(serviceName)
	defer namebuf.Release()

	name, err := namebuf.Name(lib.GSS_KRB5_NT_PRINCIPAL_NAME())
	if err != nil {
		return nil, err
	}

	t := &SPNEGOTransport{
		Lib:          lib,
		serviceName:  name,
		RoundTripper: rt,
	}

	return t, nil
}

func (lib *Lib) NewSPNEGOTransport(serviceName string) (http.RoundTripper, error) {
	return lib.NewSPNEGOTransportRoundTripper(&http.Transport{}, serviceName)
}

func (t *SPNEGOTransport) Release() {
	if t == nil {
		return
	}

	t.serviceName.Release()
}

func (t *SPNEGOTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// verify that req does not contain an Authorization header of its own
	if req.Header.Get("Authorization") != "" {
		return nil, fmt.Errorf("Authorization header must not be set when using SPNEGOTransport")
	}

	// try HEAD request until get a clean 200, a cacheable authorization, or a failure

	negotiate := false
	sendToken := t.GSS_C_NO_BUFFER()
	receiveToken := t.GSS_C_NO_BUFFER()

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
		out, _ := httputil.DumpResponse(resp, true)
		fmt.Println("<- SPNEGO RECEIVED:")
		fmt.Printf("%s\n\n", string(out))
		resp.Body.Close()

		if !negotiate {
			break
		}

		// other outputs can be safely ignored, no need to release
		// fmt.Println("DEBUG: preparing to initiate security context")
		t.ctx, _, sendToken, _, _, err = t.InitSecContext(
			t.GSS_C_NO_CREDENTIAL(),
			t.ctx, t.serviceName, t.GSS_C_NO_OID(), 0, 0,
			GSS_C_NO_CHANNEL_BINDINGS, receiveToken)
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

	out, _ := httputil.DumpRequest(req, true)
	fmt.Println("-> SPNEGO SEND:")
	fmt.Printf("%s\n\n", string(out))

	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, false, nil, err
	}
	negotiate, outputToken = t.CheckSPNEGONegotiate(resp.Header, "WWW-Authenticate")

	return resp, negotiate, outputToken, nil
}

func (lib *Lib) AddSPNEGONegotiate(h http.Header, name string, token *Buffer) {
	if name == "" {
		return
	}

	v := "Negotiate"
	if !token.IsEmpty() {
		v = v + " " + base64.StdEncoding.EncodeToString(token.Bytes())
	}
	h.Set(name, v)
}

func (lib *Lib) CheckSPNEGONegotiate(h http.Header, name string) (present bool, token *Buffer) {
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
		token = lib.MakeBufferBytes(tbytes)
	}

	return present, token
}
