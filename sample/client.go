package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/apcera/gssapi"
)

// Client attempts to connect to the configured server on the security context
// provided.
func Client(c *Context) error {
	if c.Lib == nil {
		lib, err := gssapi.LoadLib()
		if err != nil {
			return err
		}
		c.Lib = lib
	}

	t, err := c.Lib.NewSPNEGOTransport(c.ServiceName)
	if err != nil {
		return err
	}

	client := http.Client{
		Transport: t,
	}

	u := c.ServiceAddress + "/" + c.RequestPath
	fmt.Printf("CLIENT WANTS: GET %s\n\n", u)

	resp, err := client.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	fmt.Println("\n<- CLIENT RECEIVED:")
	fmt.Println(string(out))
	return nil
}
