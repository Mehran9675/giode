//go:build darwin

package i18n

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>

static void giodePreferredLanguage(char *buf, int size) {
	buf[0] = 0;
	CFArrayRef langs = CFLocaleCopyPreferredLanguages();
	if (langs == NULL) return;
	if (CFArrayGetCount(langs) > 0) {
		CFStringRef l = (CFStringRef)CFArrayGetValueAtIndex(langs, 0);
		CFStringGetCString(l, buf, size, kCFStringEncodingUTF8);
	}
	CFRelease(langs);
}
*/
import "C"

import "unsafe"

func systemLocale() string {
	var buf [64]C.char
	C.giodePreferredLanguage(&buf[0], C.int(len(buf)))
	return C.GoString((*C.char)(unsafe.Pointer(&buf[0])))
}
