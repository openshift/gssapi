package main

import (
	"fmt"
	"github.com/apcera/gssapi"
	"net/http"
	"net/http/httptest"
)

func Server(c *Context) (string, error) {
	if c.ServiceName == "" {
		return "", fmt.Errorf("Must provide a non-empty value for --service-name")
	}

	nameBuf := c.MakeBufferString(c.ServiceName)
	defer nameBuf.Release()
	name, err := nameBuf.Name(c.GSS_KRB5_NT_PRINCIPAL_NAME())
	defer name.Release()
	if err != nil {
		return "", err
	}

	cred, actualMechs, _, err := c.AcquireCred(name,
		gssapi.GSS_C_INDEFINITE, c.GSS_C_NO_OID_SET(), gssapi.GSS_C_ACCEPT)
	actualMechs.Release()
	if err != nil {
		return "", err
	}
	defer cred.Release()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass, err := filter(c, cred, w, r)
		if err != nil {
			finalize(c, w, http.StatusInternalServerError, nil)
			return
		}
		if !pass {
			return
		}
		w.Write([]byte("Hello!"))
	}))
	return ts.URL, nil
}

func filter(c *Context,
	cred *gssapi.CredId, w http.ResponseWriter, r *http.Request) (
	pass bool, err error) {

	negotiate, inputToken := c.CheckSPNEGONegotiate(r.Header, "Authorization")

	// returning a 401 with a challenge, but no token will make the client
	// initiate security context and re-submit with a non-empty Authorization
	if !negotiate || inputToken.IsEmpty() {
		finalize(c, w, http.StatusUnauthorized, nil)
		return false, nil
	}

	ctx, srcName, _, outputToken, _, _, delegatedCredHandle, err :=
		c.AcceptSecContext(c.GSS_C_NO_CONTEXT(),
			cred, inputToken, gssapi.GSS_C_NO_CHANNEL_BINDINGS)

	//TODO: special case handling of GSS_S_CONTINUE_NEEDED
	// but it doesn't change the logic, still fail
	if err != nil {
		return false, err
	}
	srcName.Release()
	delegatedCredHandle.Release()
	ctx.DeleteSecContext()

	if !outputToken.IsEmpty() {
		c.AddSPNEGONegotiate(w.Header(), "WWW-Authenticate", outputToken)
	}

	return true, nil
}

func finalize(c *Context, w http.ResponseWriter, code int, token *gssapi.Buffer) {
	c.AddSPNEGONegotiate(w.Header(), "WWW-Authenticate", token)
	w.WriteHeader(code)
	return
}
