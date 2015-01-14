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
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"
)

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
	handle unsafe.Pointer

	ftable

	// to populate the functions
	populate    sync.Once
	populateErr error

	unload    sync.Once
	unloadErr error
}

// A GSSLibrarian is an interface that defines minimal functionality for SPNEGO
// and credential issuance using GSSAPI.
type GSSLibrarian interface {
	GSSCredentialer

	GSSSPNEGOtiator
}

// A GSSCredentialer allows for the acquisition of a credential.
type GSSCredentialer interface {
	// AcquireCred acquires a GSSAPI credential.
	AcquireCred(*Name, time.Duration, *OIDSet, CredUsage) (*CredId, *OIDSet, time.Duration, error)

	// MakeBufferString makes a pointer to a Buffer from a string.
	MakeBufferString(string) *Buffer

	// GSS_KRB5_NT_PRINCIPAL_NAME returns a name format.
	GSS_KRB5_NT_PRINCIPAL_NAME() string

	// GSS_C_NO_OID_SET returns an *OIDSet.
	GSS_C_NO_OID_SET() *OIDSet
}

// A GSSSPNEGOtiator handles SPNEGO communications.
type GSSSPNEGOtiator interface {
	// CheckSPNEGONegotiate handles negotiation based upon the presence of a
	// specified header.
	CheckSPNEGONegotiate(*http.Header, string) (bool, *Buffer)

	// AddSPNEGONegotiate adds a formatted token as a header to an HTTP writer.
	AddSPNEGONegotiate(*http.Header, string, *Buffer)

	// AcceptSecContext accepts a security context and attempts to validate
	// authorization.
	AcceptSecContext(*CtxId, *Buffer, ChannelBindings) (*CtxId, *Name, *OID,
		*Buffer, uint32, time.Duration, *CredId, error)

	// GSS_C_NO_CONTEXT returns a CtxId.
	GSS_C_NO_CONTEXT() *CtxId
}

const (
	ENVVAR_GSSAPILIB          = "APC_GSSAPI_LIB"
	GSSAPILIB_DEFAULT_MIT     = "libgssapi_krb5"
	GSSAPILIB_DEFAULT_HEIMDAL = "libgssapi"

	fpPrefix = "Fp_"
)

func LoadLib() (*Lib, error) {
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
		return nil, fmt.Errorf("%s", C.GoString(C.dlerror()))
	}

	lib := &Lib{
		handle: dlhandle,
	}

	lib.populate.Do(lib.populateFunctions)
	if lib.populateErr != nil {
		lib.Unload()
		return nil, lib.populateErr
	}

	return lib, nil
}

func (lib *Lib) Unload() error {
	if lib == nil || lib.handle == nil {
		return nil
	}

	lib.unload.Do(func() {
		// in case other threads do dl stuff...
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		i := C.dlclose(lib.handle)
		if i == -1 {
			lib.unloadErr = fmt.Errorf("%s", C.GoString(C.dlerror()))
			return
		}

		lib.handle = nil
		return
	})

	return lib.unloadErr
}

func (lib *Lib) populateFunctions() {

	libT := reflect.TypeOf(lib.ftable)
	functionsV := reflect.ValueOf(lib).Elem().FieldByName("ftable")

	// no one else to touch dl while we are working!
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	n := libT.NumField()
	for i := 0; i < n; i++ {
		// Get the field name, and make sure it's an Fp_.
		f := libT.FieldByIndex([]int{i})

		if !strings.HasPrefix(f.Name, fpPrefix) {
			lib.populateErr = fmt.Errorf(
				"Unexpected: field %q does not start with %q",
				f.Name, fpPrefix)
			return
		}

		// Resolve the symbol.
		cfname := C.CString(f.Name[len(fpPrefix):])
		v := C.dlsym(lib.handle, cfname)
		C.free(unsafe.Pointer(cfname))
		if v == nil {
			lib.populateErr = fmt.Errorf("%s", C.GoString(C.dlerror()))
			return
		}

		// Save the value into the struct
		functionsV.FieldByIndex([]int{i}).SetPointer(v)
	}
}
