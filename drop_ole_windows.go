//go:build windows

package giode

import (
	"sync"
	"syscall"
	"unsafe"

	"gioui.org/f32"
	"golang.org/x/sys/windows"
)

// Drag enter, over and leave on Windows: an OLE IDropTarget, a COM
// object built here by hand - a vtable of Go callbacks - registered on
// the window's thread, where OLE must be initialized.

var (
	ole32                = windows.NewLazySystemDLL("ole32.dll")
	procOleInitialize    = ole32.NewProc("OleInitialize")
	procRegisterDragDrop = ole32.NewProc("RegisterDragDrop")
	procRevokeDragDrop   = ole32.NewProc("RevokeDragDrop")
	procReleaseStgMedium = ole32.NewProc("ReleaseStgMedium")
	procScreenToClient   = user32.NewProc("ScreenToClient")
)

const (
	cfHDROP          = 15
	dvaspectContent  = 1
	tymedHGlobal     = 1
	dropEffectNone   = 0
	dropEffectCopy   = 1
	sOK              = 0
	eNoInterface     = 0x80004002
	dataObjGetData   = 3 // IDataObject vtable slots
	dataObjQueryData = 5
)

var (
	iidIUnknown    = windows.GUID{Data1: 0x00000000, Data4: [8]byte{0xC0, 0, 0, 0, 0, 0, 0, 0x46}}
	iidIDropTarget = windows.GUID{Data1: 0x00000122, Data4: [8]byte{0xC0, 0, 0, 0, 0, 0, 0, 0x46}}
)

type formatEtc struct {
	format uint16
	ptd    uintptr
	aspect uint32
	index  int32
	tymed  uint32
}

type stgMedium struct {
	tymed   uint32
	handle  uintptr
	release uintptr
}

// oleDropTarget is the COM object: its first word points to the
// vtable.
type oleDropTarget struct {
	vtbl  *[7]uintptr
	hwnd  uintptr
	app   *App
	files bool // the drag carries files
}

var (
	oleVtbl    *[7]uintptr
	oleVtblMu  sync.Mutex
	oleTargets = map[uintptr]*oleDropTarget{} // keeps them alive
)

// registerOLEDrop runs on the window's thread.
func registerOLEDrop(a *App, hwnd uintptr) {
	// POINTL arrives by value: one register on 64-bit systems, which
	// the callbacks rely on.
	if unsafe.Sizeof(uintptr(0)) != 8 {
		return
	}
	procOleInitialize.Call(0)
	oleVtblMu.Lock()
	if oleVtbl == nil {
		oleVtbl = &[7]uintptr{
			windows.NewCallback(oleQueryInterface),
			windows.NewCallback(oleAddRef),
			windows.NewCallback(oleRelease),
			windows.NewCallback(oleDragEnter),
			windows.NewCallback(oleDragOver),
			windows.NewCallback(oleDragLeave),
			windows.NewCallback(oleDrop),
		}
	}
	t := &oleDropTarget{vtbl: oleVtbl, hwnd: hwnd, app: a}
	oleTargets[hwnd] = t
	oleVtblMu.Unlock()
	if hr, _, _ := procRegisterDragDrop.Call(hwnd, uintptr(unsafe.Pointer(t))); hr != sOK {
		a.logf(LogLevelWarning, "giode: drag events unavailable: RegisterDragDrop failed (0x%x)", hr)
		oleVtblMu.Lock()
		delete(oleTargets, hwnd)
		oleVtblMu.Unlock()
	}
}

func revokeOLEDrop(hwnd uintptr) {
	oleVtblMu.Lock()
	_, ok := oleTargets[hwnd]
	delete(oleTargets, hwnd)
	oleVtblMu.Unlock()
	if ok {
		procRevokeDragDrop.Call(hwnd)
	}
}

func oleQueryInterface(this *oleDropTarget, riid *windows.GUID, out **oleDropTarget) uintptr {
	if *riid == iidIUnknown || *riid == iidIDropTarget {
		*out = this
		return sOK
	}
	*out = nil
	return eNoInterface
}

// The targets live as long as their windows (oleTargets), so the
// reference count isn't kept.
func oleAddRef(this *oleDropTarget) uintptr  { return 1 }
func oleRelease(this *oleDropTarget) uintptr { return 1 }

// clientPos converts a POINTL in screen pixels to window pixels.
func (t *oleDropTarget) clientPos(pt uintptr) f32.Point {
	p := [2]int32{int32(uint32(pt)), int32(uint32(pt >> 32))}
	procScreenToClient.Call(t.hwnd, uintptr(unsafe.Pointer(&p)))
	return f32.Pt(float32(p[0]), float32(p[1]))
}

func hdropFormat() formatEtc {
	return formatEtc{format: cfHDROP, aspect: dvaspectContent, index: -1, tymed: tymedHGlobal}
}

// comCall calls method slot of a COM object.
func comCall(obj unsafe.Pointer, slot int, args ...uintptr) uintptr {
	vtbl := *(*unsafe.Pointer)(obj)
	fn := *(*uintptr)(unsafe.Add(vtbl, uintptr(slot)*unsafe.Sizeof(uintptr(0))))
	r, _, _ := syscall.SyscallN(fn, append([]uintptr{uintptr(obj)}, args...)...)
	return r
}

func (t *oleDropTarget) effect(out *uint32) uintptr {
	*out = dropEffectNone
	if t.files {
		*out = dropEffectCopy
	}
	return sOK
}

func oleDragEnter(t *oleDropTarget, data unsafe.Pointer, keys, pt uintptr, effect *uint32) uintptr {
	f := hdropFormat()
	t.files = data != nil && comCall(data, dataObjQueryData, uintptr(unsafe.Pointer(&f))) == sOK
	if t.files {
		t.app.deliverFileDrag(FileDragEnter, t.clientPos(pt))
	}
	return t.effect(effect)
}

func oleDragOver(t *oleDropTarget, keys, pt uintptr, effect *uint32) uintptr {
	if t.files {
		t.app.deliverFileDrag(FileDragOver, t.clientPos(pt))
	}
	return t.effect(effect)
}

func oleDragLeave(t *oleDropTarget) uintptr {
	if t.files {
		t.files = false
		t.app.deliverFileDrag(FileDragLeave, f32.Point{})
	}
	return sOK
}

func oleDrop(t *oleDropTarget, data unsafe.Pointer, keys, pt uintptr, effect *uint32) uintptr {
	if !t.files || data == nil {
		return t.effect(effect)
	}
	t.files = false
	pos := t.clientPos(pt)
	f := hdropFormat()
	var m stgMedium
	var paths []string
	if comCall(data, dataObjGetData, uintptr(unsafe.Pointer(&f)), uintptr(unsafe.Pointer(&m))) == sOK {
		if m.tymed == tymedHGlobal && m.handle != 0 {
			paths = readDrop(m.handle)
		}
		procReleaseStgMedium.Call(uintptr(unsafe.Pointer(&m)))
	}
	t.app.deliverFileDrop(paths, pos)
	*effect = dropEffectCopy
	return sOK
}
