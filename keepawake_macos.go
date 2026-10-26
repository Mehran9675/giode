//go:build darwin && !ios

package giode

/*
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#include <stdlib.h>
#include <IOKit/pwr_mgt/IOPMLib.h>

static IOPMAssertionID giodeAssert(int display, const char *reason) {
	CFStringRef r = CFStringCreateWithCString(NULL, reason, kCFStringEncodingUTF8);
	IOPMAssertionID id = 0;
	CFStringRef type = display ? kIOPMAssertionTypeNoDisplaySleep : kIOPMAssertionTypePreventUserIdleSystemSleep;
	if (IOPMAssertionCreateWithName(type, kIOPMAssertionLevelOn, r, &id) != kIOReturnSuccess) {
		id = 0;
	}
	CFRelease(r);
	return id;
}
*/
import "C"

import (
	"context"
	"errors"
	"sync"
	"unsafe"
)

func (a *App) haptic(HapticKind) {}

var screenAssertion struct {
	sync.Mutex
	id C.IOPMAssertionID
}

func assert(display bool, reason string) (C.IOPMAssertionID, error) {
	cs := C.CString(reason)
	defer C.free(unsafe.Pointer(cs))
	d := C.int(0)
	if display {
		d = 1
	}
	id := C.giodeAssert(d, cs)
	if id == 0 {
		return 0, errors.New("giode: IOPMAssertionCreateWithName failed")
	}
	return id, nil
}

func (a *App) keepScreenOn(ctx context.Context, on bool) error {
	screenAssertion.Lock()
	defer screenAssertion.Unlock()
	switch {
	case on && screenAssertion.id == 0:
		id, err := assert(true, "Keeping the screen on")
		if err != nil {
			return err
		}
		screenAssertion.id = id
	case !on && screenAssertion.id != 0:
		C.IOPMAssertionRelease(screenAssertion.id)
		screenAssertion.id = 0
	}
	return nil
}

func (a *App) keepAlive(reason string) (func(), error) {
	id, err := assert(false, reason)
	if err != nil {
		return func() {}, err
	}
	var once sync.Once
	return func() { once.Do(func() { C.IOPMAssertionRelease(id) }) }, nil
}
