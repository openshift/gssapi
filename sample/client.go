package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

// Client attempts to connect to the configured server on the security context
// provided.
func Client(c *Context) error {
	t, err := c.NewSPNEGOTransport(c.ServiceName)
	if err != nil {
		return err
	}

	client := http.Client{
		Transport: t,
	}

	u := c.ServiceAddress + c.RequestPath
	if !strings.HasPrefix(u, "http://") {
		u = "http://" + u
	}
	c.Print("CLIENT WANTS: GET ", u)

	resp, err := client.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	c.Print("<- CLIENT RECEIVED:\n", string(out), "\n")

	if resp.StatusCode != http.StatusOK || !strings.Contains(string(out), "Hello!") {
		return fmt.Errorf("Test failed")
	}

	return nil
}
