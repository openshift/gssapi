// Copyright 2013-2015 Apcera Inc. All rights reserved.

// +build darwin linux

package gssapi

//#cgo LDFLAGS: -L/opt/local/lib -Wl,-search_paths_first -lgssapi_krb5 -lkrb5 -lk5crypto -lcom_err
//import "C"
// Uncomment the above line to directly link against the libraries, to avoid
// the dlopen layer.

/*
#cgo linux LDFLAGS: -ldl

#include <gssapi/gssapi.h>
#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

// matches log, not fmt
type Printer interface {
	Print(a ...interface{})
}

// ftable fields will be initialized to the corresponding function pointers from the
// GSSAPI library. They must be of form Fp_function_name (Capital 'F' so that
// we can use reflect.
type ftable struct {
	// buffer.go
	Fp_gss_release_buffer unsafe.Pointer
	Fp_gss_import_name    unsafe.Pointer
	Fp_gss_str_to_oid     unsafe.Pointer

	// context.go
	Fp_gss_init_sec_context      unsafe.Pointer
	Fp_gss_accept_sec_context    unsafe.Pointer
	Fp_gss_delete_sec_context    unsafe.Pointer
	Fp_gss_process_context_token unsafe.Pointer
	Fp_gss_context_time          unsafe.Pointer
	Fp_gss_inquire_context       unsafe.Pointer
	Fp_gss_wrap_size_limit       unsafe.Pointer
	Fp_gss_export_sec_context    unsafe.Pointer
	Fp_gss_import_sec_context    unsafe.Pointer

	// credential.go
	Fp_gss_acquire_cred unsafe.Pointer
	Fp_gss_release_cred unsafe.Pointer

	// name.go
	Fp_gss_canonicalize_name      unsafe.Pointer
	Fp_gss_compare_name           unsafe.Pointer
	Fp_gss_display_name           unsafe.Pointer
	Fp_gss_duplicate_name         unsafe.Pointer
	Fp_gss_export_name            unsafe.Pointer
	Fp_gss_inquire_mechs_for_name unsafe.Pointer
	Fp_gss_inquire_names_for_mech unsafe.Pointer
	Fp_gss_release_name           unsafe.Pointer

	// oid.go
	// Fp_gss_oid_equal  unsafe.Pointer
	Fp_gss_oid_to_str unsafe.Pointer

	// oid_set.go
	Fp_gss_create_empty_oid_set unsafe.Pointer
	Fp_gss_add_oid_set_member   unsafe.Pointer
	Fp_gss_release_oid_set      unsafe.Pointer
	Fp_gss_test_oid_set_member  unsafe.Pointer

	// status.go
	Fp_gss_display_status unsafe.Pointer

	// krb5_keytab.go -- where does this come from?
	// Fp_gsskrb5_register_acceptor_identity unsafe.Pointer
}

// Lib encapsulates both the GSSAPI and the library dlopen()'d for it.
type Lib struct {
	Printer

	handle unsafe.Pointer

	ftable
}

const (
	ENVVAR_GSSAPILIB          = "APC_GSSAPI_LIB"
	GSSAPILIB_DEFAULT_MIT     = "libgssapi_krb5"
	GSSAPILIB_DEFAULT_HEIMDAL = "libgssapi"

	fpPrefix = "Fp_"
)

func LibPath(path string, useHeimdal bool, useMIT bool) string {
	switch {
	case path != "" && !useMIT && !useHeimdal:
		return path

	case path == "" && useMIT && !useHeimdal:
		return appendExt(GSSAPILIB_DEFAULT_MIT)

	case path == "" && !useMIT && useHeimdal:
		return appendExt(GSSAPILIB_DEFAULT_HEIMDAL)

	case path == "" && runtime.GOOS == "freebsd":
		return appendExt(GSSAPILIB_DEFAULT_HEIMDAL)

	case path == "" && !useMIT && !useHeimdal:
		if envLib := os.Getenv(ENVVAR_GSSAPILIB); envLib != "" {
			return envLib
		}
		return appendExt(GSSAPILIB_DEFAULT_MIT)
	}
	return ""
}

func LoadDefaultLib() (*Lib, error) {
	lib, err := LoadLib(LibPath("", false, false))
	if err != nil {
		return nil, err
	}

	lib.Printer = log.New(os.Stderr, "gssapi", log.LstdFlags)

	return lib, nil
}

func LoadLib(path string) (*Lib, error) {
	// We get the error in a separate call, so we need to lock OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	lib_cs := C.CString(path)
	defer C.free(unsafe.Pointer(lib_cs))
	// we don't use RTLD_FIRST, it might be the case that the GSSAPI lib
	// delegates symbols to other libs it links against (eg, Kerberos)
	dlhandle := C.dlopen(lib_cs, C.RTLD_NOW|C.RTLD_LOCAL)
	if dlhandle == nil {
		return nil, fmt.Errorf("%s", C.GoString(C.dlerror()))
	}

	lib := &Lib{
		handle: dlhandle,
	}

	err := lib.populateFunctions()
	if err != nil {
		lib.Unload()
		return nil, err
	}

	return lib, nil
}

func (lib *Lib) Unload() error {
	if lib == nil || lib.handle == nil {
		return nil
	}

	//TODO: is runtime.LockOSThread needed here?
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	i := C.dlclose(lib.handle)
	if i == -1 {
		return fmt.Errorf("%s", C.GoString(C.dlerror()))
	}

	lib.handle = nil
	return nil
}

func appendExt(path string) string {
	ext := ".so"
	if runtime.GOOS == "darwin" {
		ext = ".dylib"
	}
	if !strings.HasSuffix(path, ext) {
		path += ext
	}
	return path
}

func (lib *Lib) populateFunctions() error {
	libT := reflect.TypeOf(lib.ftable)
	functionsV := reflect.ValueOf(lib).Elem().FieldByName("ftable")

	n := libT.NumField()
	for i := 0; i < n; i++ {
		// Get the field name, and make sure it's an Fp_.
		f := libT.FieldByIndex([]int{i})

		if !strings.HasPrefix(f.Name, fpPrefix) {
			return fmt.Errorf(
				"Unexpected: field %q does not start with %q",
				f.Name, fpPrefix)
		}

		// Resolve the symbol.
		cfname := C.CString(f.Name[len(fpPrefix):])
		v := C.dlsym(lib.handle, cfname)
		C.free(unsafe.Pointer(cfname))
		if v == nil {
			return fmt.Errorf("%s", C.GoString(C.dlerror()))
		}

		// Save the value into the struct
		functionsV.FieldByIndex([]int{i}).SetPointer(v)
	}

	return nil
}

func (lib *Lib) Print(a ...interface{}) {
	if lib == nil || lib.Printer == nil {
		return
	}
	lib.Printer.Print(a...)
}
