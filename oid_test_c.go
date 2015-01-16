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
	void *vp_str_to_oid,
	void *vp_release_buffer,
	void *vp_display_status)
{
	OM_uint32 maj, min, s, message_context;
	gss_buffer_desc buffer;
	gss_buffer_t b = &buffer;
	gss_OID_desc oid;
	gss_OID newoid;
	char context[512];

	OM_uint32 (*fp_oid_to_str)(OM_uint32 *, gss_OID, gss_buffer_t) = vp_oid_to_str;
	OM_uint32 (*fp_str_to_oid)(OM_uint32 *, gss_buffer_t, gss_OID*) = vp_str_to_oid;
	OM_uint32 (*fp_release_buffer)(OM_uint32 *, gss_buffer_t) = vp_release_buffer;
	OM_uint32 (*fp_display_status)(OM_uint32 *, OM_uint32, int, const gss_OID, OM_uint32 *, gss_buffer_t) = vp_display_status;

	// On the Mac (and other BSD systems?) oid_to_str produces output "1 2 840 113554 1 2 1 1"
	// without the surrounding {}, and excluding the trailing 0. On Linux it seems to produce
	// "{ 1 2 840 113554 1 2 1 1 }" including the \0.
	// str_to_oid always expects the {} form.
	static struct {
		char *to_str;
		int to_str_len;
		char *from_str;
		gss_OID_desc oid;
	} *t, tests[] = {
		// GSS_C_NT_USER_NAME
		{
#if defined(__APPLE__) && defined(__MACH__)
			"1 2 840 113554 1 2 1 1",
			22,
			"NOT APPLICABLE",
#else
			"{ 1 2 840 113554 1 2 1 1 }",
			27,
			"{ 1 2 840 113554 1 2 1 1 }",
#endif
			{ 10, "\x2A\x86\x48\x86\xF7\x12\x01\x02\x01\x01" }
		}
	};

	t = tests;
	int ntests = sizeof(tests)/sizeof(*tests);
	int i;
	for (i=0; i < ntests; i++, t++) {
		snprintf(context, sizeof(context), "%s: testing '%s':", __func__, t->from_str);

		maj = fp_oid_to_str(&min, &t->oid, b);
		if ((int) maj != GSS_S_COMPLETE || min != 0) {
			fprintf(stderr, "%s gss_oid_to_str: got statuses %x,%x, expected 0,0\n", context,
				(int) maj, (int) min);
			goto ERROR_OUT;
		}
		if (b->length != t->to_str_len) {
			fprintf(stderr, "%s gss_oid_to_str: got b->length:%d ('%.*s'), expected %d ('%s')\n", context,
				(int) b->length, (int) b->length, (char *) b->value,
				t->to_str_len, t->to_str);
			return 0;
		}
		if (memcmp(b->value, t->to_str, b->length) != 0) {
			fprintf(stderr, "%s gss_oid_to_str: got b->value:'%.*s', expected %s\n", context,
				(int) b->length, (char *) b->value, t->to_str);
			return 0;
		}
		fp_release_buffer(&min, b);

#if !defined(__APPLE__) || !defined(__MACH__)
		b->value = t->from_str;
		b->length = strlen(t->from_str)+1;
		maj = fp_str_to_oid(&min, b, &newoid);
		if (maj != GSS_S_COMPLETE || min != 0) {
			fprintf(stderr, "%s gss_str_to_oid: got statuses %x,%x, expected 0,0\n", context,
				maj, min);
			goto ERROR_OUT;
		}
		if (newoid->length != t->oid.length) {
			fprintf(stderr, "%s gss_str_to_oid: got newoid->length:%d, expected %d\n", context,
				(int) newoid->length, (int) t->oid.length);
			return 0;
		}
		if (memcmp(newoid->elements, t->oid.elements, newoid->length) != 0) {
			fprintf(stderr, "%s gss_str_to_oid: got newoid->elements:'%.*s', expected %.*s\n", context,
				(int) newoid->length, (char *) newoid->elements,
				(int) t->oid.length, (char *) t->oid.elements);
			return 0;
		}
#endif

	}

	return 1;

ERROR_OUT:
	message_context=0;
	s=maj;
	do {
		maj = fp_display_status(&min, s, GSS_C_GSS_CODE, GSS_C_NO_OID, &message_context, b);
		fprintf(stderr, "%.*s\n", (int) b->length, (char *) b->value);
		fp_release_buffer(&min, b);
	} while (message_context != 0);

	return 0;
}

*/
import "C"

func cOIDTest(l *Lib) bool {
	result := C.oid_to_str_test(
		l.Fp_gss_oid_to_str,
		l.Fp_gss_str_to_oid,
		l.Fp_gss_release_buffer,
		l.Fp_gss_display_status,
	)
	return result != 0
}
