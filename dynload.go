// Copyright 2013 Apcera Inc. All rights reserved.

// +build darwin linux

package gssapi

// This file provides for library dynamic loading support

/*
#cgo linux LDFLAGS: -ldl

#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"

import (
	"os"
	"runtime"
	"strings"
	"unsafe"
)

const (
	ENVVAR_GSSAPILIB          = "APC_GSSAPI_LIB"
	GSSAPILIB_DEFAULT_MIT     = "libgssapi_krb5"
	GSSAPILIB_DEFAULT_HEIMDAL = "libgssapi"
)

type DlError string

func (e DlError) Error() string { return string(e) }

// Beware that if you call this, you should lock the OS thread before
// calling the function which errored
func GetDlError() error {
	return DlError(C.GoString(C.dlerror()))
}

func (gss *GssapiLib) DlSym(symbol string) (unsafe.Pointer, error) {
	cSym := C.CString(symbol)
	defer func() { C.free(unsafe.Pointer(cSym)) }()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	v := C.dlsym(gss.handle, cSym)
	if v == nil {
		return nil, GetDlError()
	}
	return v, nil
}

func LoadGssapiLibrary() (*GssapiLib, error) {
	libname := GSSAPILIB_DEFAULT_MIT
	libExt := ".so"
	switch runtime.GOOS {
	case "darwin":
		libExt = ".dylib"
	case "freebsd":
		libname = GSSAPILIB_DEFAULT_HEIMDAL
	}
	if envLib := os.Getenv(ENVVAR_GSSAPILIB); envLib != "" {
		libname = envLib
	}

	if !strings.HasSuffix(libname, libExt) {
		libname += libExt
	}

	// We get the error in a separate call, so we need to lock OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	lib_cs := C.CString(libname)
	defer C.free(unsafe.Pointer(lib_cs))
	// we don't use RTLD_FIRST, it might be the case that the GSSAPI lib
	// delegates symbols to other libs it links against (eg, Kerberos)
	dlhandle := C.dlopen(lib_cs, C.RTLD_NOW|C.RTLD_LOCAL)
	if dlhandle == nil {
		return nil, GetDlError()
	}

	lib := &GssapiLib{handle: dlhandle}
	err := lib.Populate()
	if err != nil {
		lib.LibraryClose()
		return nil, err
	}
	return lib, nil
}

func (gss *GssapiLib) LibraryClose() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	i := C.dlclose(gss.handle)
	if i == -1 {
		return GetDlError()
	}
	gss.handle = nil
	return nil
}
