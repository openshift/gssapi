// Copyright 2013 Apcera Inc. All rights reserved.

// GSS status and errors

package gssapi

/*
#include <gssapi/gssapi.h>
*/
import "C"

import (
	"fmt"
	"strings"
)

// Constant values are specified for C-language bindings in RFC 2744.
/*
"""
   These errors are encoded into the 32-bit GSS status code as follows:

      MSB                                                        LSB
      |------------------------------------------------------------|
      |  Calling Error | Routine Error  |    Supplementary Info    |
      |------------------------------------------------------------|
   Bit 31            24 23            16 15                       0
"""

Note that the first two fields hold integer consts, whereas Supplementary Info
is a bit-field.
*/

const (
	shiftCALLING = 24
	shiftROUTINE = 16
	maskCALLING  = 0xFF000000
	maskROUTINE  = 0x00FF0000
	maskSUPPINFO = 0x0000FFFF
)

// These are GSSAPI-defined:
type MajorStatus uint32

// These are mechanism-specific:
type MinorStatus uint32

// If Major is GSS_S_FAILURE then information will be in Minor
type Status struct {
	Major MajorStatus
	Minor MinorStatus
}

type StatusError struct {
	lib    *GssapiLib
	status Status
}

const (
	GSS_S_COMPLETE MajorStatus = 0

	GSS_S_CALL_INACCESSIBLE_READ  MajorStatus = 1 << shiftCALLING
	GSS_S_CALL_INACCESSIBLE_WRITE             = 2 << shiftCALLING
	GSS_S_CALL_BAD_STRUCTURE                  = 3 << shiftCALLING

	GSS_S_BAD_MECH             MajorStatus = 1 << shiftROUTINE
	GSS_S_BAD_NAME                         = 2 << shiftROUTINE
	GSS_S_BAD_NAMETYPE                     = 3 << shiftROUTINE
	GSS_S_BAD_BINDINGS                     = 4 << shiftROUTINE
	GSS_S_BAD_STATUS                       = 5 << shiftROUTINE
	GSS_S_BAD_MIC                          = 6 << shiftROUTINE
	GSS_S_BAD_SIG                          = 6 << shiftROUTINE // duplication deliberate
	GSS_S_NO_CRED                          = 7 << shiftROUTINE
	GSS_S_NO_CONTEXT                       = 8 << shiftROUTINE
	GSS_S_DEFECTIVE_TOKEN                  = 9 << shiftROUTINE
	GSS_S_DEFECTIVE_CREDENTIAL             = 10 << shiftROUTINE
	GSS_S_CREDENTIALS_EXPIRED              = 11 << shiftROUTINE
	GSS_S_CONTEXT_EXPIRED                  = 12 << shiftROUTINE
	GSS_S_FAILURE                          = 13 << shiftROUTINE
	GSS_S_BAD_QOP                          = 14 << shiftROUTINE
	GSS_S_UNAUTHORIZED                     = 15 << shiftROUTINE
	GSS_S_UNAVAILABLE                      = 16 << shiftROUTINE
	GSS_S_DUPLICATE_ELEMENT                = 17 << shiftROUTINE
	GSS_S_NAME_NOT_MN                      = 18 << shiftROUTINE

	field_GSS_S_CONTINUE_NEEDED = 1 << 0
	field_GSS_S_DUPLICATE_TOKEN = 1 << 1
	field_GSS_S_OLD_TOKEN       = 1 << 2
	field_GSS_S_UNSEQ_TOKEN     = 1 << 3
	field_GSS_S_GAP_TOKEN       = 1 << 4
)

// Equivalent to C GSS_CALLING_ERROR() macro
func (st MajorStatus) CallingError() MajorStatus {
	return st & maskCALLING
}

// Equivalent to C GSS_ROUTINE_ERROR() macro
func (st MajorStatus) RoutineError() MajorStatus {
	return st & maskROUTINE
}

// Equivalent to C GSS_SUPPLEMENTARY_INFO() macro
func (st MajorStatus) SupplementaryInfo() MajorStatus {
	return st & maskSUPPINFO
}

// Equivalent to C GSS_ERROR() macro
// Not 'Error' because that's special in Go conventions
func (st MajorStatus) IsError() bool {
	return st&(maskCALLING|maskROUTINE) != 0
}

func (st MajorStatus) ContinueNeeded() bool {
	return st&field_GSS_S_CONTINUE_NEEDED != 0
}

func (st MajorStatus) DuplicateToken() bool {
	return st&field_GSS_S_DUPLICATE_TOKEN != 0
}

func (st MajorStatus) OldToken() bool {
	return st&field_GSS_S_OLD_TOKEN != 0
}

func (st MajorStatus) UnseqToken() bool {
	return st&field_GSS_S_UNSEQ_TOKEN != 0
}

func (st MajorStatus) GapToken() bool {
	return st&field_GSS_S_GAP_TOKEN != 0
}

func NewStatus(major, minor C.OM_uint32) Status {
	return Status{
		Major: MajorStatus(major),
		Minor: MinorStatus(minor),
	}
}

func (lib *GssapiLib) CheckError(maybe Status) error {
	if !maybe.Major.IsError() {
		return nil
	}
	return &StatusError{lib: lib, status: maybe}
}

func (se *StatusError) Error() string {
	messages := make([]string, 0, 6)
	additional := make([]Status, 0, 2)
	buffer := GssBuffer{lib: se.lib}
	first := true
	context := C.OM_uint32(0)

	var inquiry C.OM_uint32
	var code_type int
	if se.status.Major.RoutineError() == GSS_S_FAILURE {
		inquiry = C.OM_uint32(se.status.Minor)
		code_type = GSS_C_MECH_CODE
	} else {
		inquiry = C.OM_uint32(se.status.Major)
		code_type = GSS_C_GSS_CODE
	}

	for first || context != C.OM_uint32(0) {
		first = false
		var (
			render_min C.OM_uint32
		)
		render_maj := se.lib.gss_display_status(&render_min,
			inquiry,
			code_type,
			GSS_C_NO_OID, // store a mech_type at the lib level?  Or context?
			&context,
			buffer.buffer,
		)
		resultStatus := NewStatus(render_maj, render_min)
		if resultStatus.Major.IsError() {
			additional = append(additional, resultStatus)
		}
		messages = append(messages, buffer.String())
		buffer.Release()
	}
	if len(additional) > 0 {
		messages = append(messages, fmt.Sprintf("additionally, %d conversions failed", len(additional)))
	}
	messages = append(messages, "")
	return strings.Join(messages, "\n")
}
