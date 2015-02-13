// Copyright 2013-2015 Apcera Inc. All rights reserved.

// +build clienttest

package test

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"testing"
)

func TestLoadAndRequest(t *testing.T) {
	ch := make(chan error, 1)
	oneRequest(t, ch, c.DebugLog)
	err := <-ch
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkLoadAndRequest(b *testing.B) {
	ch := make(chan error, b.N)

	for i := 0; i < b.N; i++ {
		go oneRequest(b, ch, false)
	}

	for i := 0; i < b.N; i++ {
		err := <-ch
		if err != nil {
			b.Fatal(err)
		}
	}
}

func oneRequest(tb testing.TB, ch chan error, verbose bool) {
	var err error
	defer func() {
		ch <- err
	}()

	lib, err := loadlib(verbose, "go-gssapi-test-client")
	if err != nil {
		tb.Fatal(err)
	}
	defer lib.Unload()

	transport, err := lib.NewSPNEGOTransport(c.ServiceName)
	if err != nil {
		return
	}

	client := http.Client{
		Transport: transport,
	}

	u := c.ServiceAddress + "/" // c.RequestPath
	if !strings.HasPrefix(u, "http://") {
		u = "http://" + u
	}
	lib.Debug("CLIENT WANTS: GET ", u)

	resp, err := client.Get(u)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return
	}
	lib.Debug("<- CLIENT RECEIVED:\n", string(out), "\n")

	if resp.StatusCode != http.StatusOK || !strings.Contains(string(out), "Hello!") {
		err = fmt.Errorf("Test failed: unexpected response: code:%v, body:\n%s", resp.StatusCode, string(out))
	}
}
