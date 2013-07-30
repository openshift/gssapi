// Copyright 2013 Apcera Inc. All rights reserved.

// +build ignore

package main

import (
	"fmt"
	"os"

	"gssapi"
)

func main() {
	handle, err := gssapi.LoadGssapiLibrary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load GSSAPI driver: %s\n", err)
		os.Exit(1)
	}
	v, err := handle.DlSym("gss_acquire_cred")
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"symbol resolution (of '%s') failed: %s\n",
			"gss_acquire_cred", err)
		handle.LibraryClose()
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "gss_acquire_cred resolved: %#+v\n", v)
	handle.LibraryClose()
	fmt.Println("Done.")
}
