// Copyright 2013-2015 Apcera Inc. All rights reserved.

package test

import (
	"fmt"
	"net/http"
	"os"

	"github.com/apcera/gssapi"
)

func Service(c *Context) error {
	if c.ServiceName == "" {
		return fmt.Errorf("Must provide a non-empty value for --service-name")
	}

	c.Debug(fmt.Sprintf("Starting service %q", c.ServiceName))

	nameBuf, err := c.MakeBufferString(c.ServiceName)
	if err != nil {
		return err
	}
	defer nameBuf.Release()

	name, err := nameBuf.Name(c.GSS_KRB5_NT_PRINCIPAL_NAME)
	if err != nil {
		return err
	}
	defer name.Release()

	cred, actualMechs, _, err := c.AcquireCred(name,
		gssapi.GSS_C_INDEFINITE, c.GSS_C_NO_OID_SET, gssapi.GSS_C_ACCEPT)
	actualMechs.Release()
	if err != nil {
		return err
	}
	defer cred.Release()

	keytab := os.Getenv("KRB5_KTNAME")
	if keytab == "" {
		keytab = "default /etc/krb5.keytab"
	}
	c.Debug(fmt.Sprintf("Acquired credentials using %v", keytab))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pass, err := filter(c, cred, w, r)
		if err != nil {
			//TODO: differentiate invalid tokens here and return a 403
			c.Err(fmt.Sprintf("ACCESS %d %q %q %q",
				http.StatusInternalServerError,
				r.Method,
				r.URL.String(), err))
			finalize(c, w, http.StatusInternalServerError, nil)
			return
		}
		if !pass {
			c.Warn(fmt.Sprintf(`ACCESS %d %q %q "no input token provided"`,
				http.StatusUnauthorized,
				r.Method,
				r.URL.String()))
			finalize(c, w, http.StatusUnauthorized, nil)
			return
		}
		w.Write([]byte("Hello!"))
		c.Info(fmt.Sprintf("ACCESS %d %q %q", http.StatusOK, r.Method, r.URL.String()))
	})

	err = http.ListenAndServe(c.ServiceAddress, nil)
	if err != nil {
		return err
	}

	return nil
}

func filter(c *Context,
	cred *gssapi.CredId, w http.ResponseWriter, r *http.Request) (
	pass bool, err error) {

	negotiate, inputToken := c.CheckSPNEGONegotiate(r.Header, "Authorization")

	// returning a 401 with a challenge, but no token will make the client
	// initiate security context and re-submit with a non-empty Authorization
	if !negotiate || inputToken.IsEmpty() {
		return false, nil
	}

	ctx, srcName, _, outputToken, _, _, delegatedCredHandle, err :=
		c.AcceptSecContext(c.GSS_C_NO_CONTEXT,
			cred, inputToken, c.GSS_C_NO_CHANNEL_BINDINGS)

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
