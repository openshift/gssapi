// Copyright 2013-2015 Apcera Inc. All rights reserved.

package gssapi

/*
#include <gssapi.h>
#include <stdio.h>
#include <string.h>

// boolean 1 success, 0 failure
static int
oid_to_str_test(
	void *vp_oid_to_str,
	void *vp_str_to_oid)
{
	OM_uint32 maj, min;
	gss_buffer_desc buffer;
	gss_buffer_t b = &buffer;
	gss_OID_desc oid;
	gss_OID newoid;
	char context[512];

	OM_uint32 (*fp_oid_to_str)(OM_uint32 *, gss_OID, gss_buffer_t) = vp_oid_to_str;
	OM_uint32 (*fp_str_to_oid)(OM_uint32 *, gss_buffer_t, gss_OID*) = vp_str_to_oid;

	static struct {
		char *canonical;
		char *variant;
		gss_OID_desc oid;
	} *t, tests[] = {
		// GSS_C_NT_USER_NAME
		{
			"{ 1 2 840 113554 1 2 1 1 }",
			"1.2.840.113554.1.2.1.1",
			{ 10, "\x2A\x86\x48\x86\xF7\x12\x01\x02\x01\x01" }
		}
	};

	int ntests = sizeof(tests)/sizeof(*tests);
	int i;

	t = tests;
	for (i=0; i < ntests; i++, t++) {
		snprintf(context, sizeof(context), "%s: testing '%s':", __func__, t->canonical);

		maj = fp_oid_to_str(&min, &t->oid, b);
		if (maj != GSS_S_COMPLETE || min != 0) {
			fprintf(stderr, "%s gss_oid_to_str: got statuses %x,%x, expected 0,0", context,
				maj, min);
			return 0;
		}
		if (b->length != strlen(t->canonical)+1) {
			fprintf(stderr, "%s gss_oid_to_str: got b->length:%d, expected %d", context,
				(int) b->length, (int) strlen(t->canonical));
			return 0;
		}
		if (memcmp(b->value, t->canonical, b->length) != 0) {
			fprintf(stderr, "%s gss_oid_to_str: got b->value:'%.*s', expected %s", context,
				(int) b->length, (char *) b->value, t->canonical);
			return 0;
		}

		maj = fp_str_to_oid(&min, b, &newoid);
		if (maj != GSS_S_COMPLETE || min != 0) {
			fprintf(stderr, "%s gss_str_to_oid: got statuses %x,%x, expected 0,0", context,
				maj, min);
			return 0;
		}
		if (newoid->length != t->oid.length) {
			fprintf(stderr, "%s gss_str_to_oid: got newoid->length:%d, expected %d", context,
				(int) newoid->length, (int) t->oid.length);
			return 0;
		}
		if (memcmp(newoid->elements, t->oid.elements, newoid->length) != 0) {
			fprintf(stderr, "%s gss_str_to_oid: got newoid->elements:'%.*s', expected %.*s", context,
				(int) newoid->length, (char *) newoid->elements,
				(int) t->oid.length, (char *) t->oid.elements);
			return 0;
		}
	}

	return 1;
}

*/
import "C"

func cOIDTest(l *Lib) bool {
	result := C.oid_to_str_test(
		l.Fp_gss_oid_to_str,
		l.Fp_gss_str_to_oid)
	return result != 0
}
