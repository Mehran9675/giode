package giode

import (
	"context"
	"runtime"
	"sync"

	"golang.org/x/sys/windows"
)

var procSetThreadExecutionState = windows.NewLazySystemDLL("kernel32.dll").NewProc("SetThreadExecutionState")

const (
	esContinuous      = 0x80000000
	esSystemRequired  = 0x00000001
	esDisplayRequired = 0x00000002
)

// awake holds the execution state flags requested by KeepScreenOn and
// KeepAlive. SetThreadExecutionState applies to the calling thread, so
// one locked thread owns it.
var awake struct {
	sync.Mutex
	screen  bool
	holds   int
	updates chan uint32
}

func setExecutionState() {
	flags := uint32(esContinuous)
	if awake.screen {
		flags |= esDisplayRequired | esSystemRequired
	}
	if awake.holds > 0 {
		flags |= esSystemRequired
	}
	if awake.updates == nil {
		awake.updates = make(chan uint32, 1)
		go func() {
			runtime.LockOSThread()
			for f := range awake.updates {
				procSetThreadExecutionState.Call(uintptr(f))
			}
		}()
	}
	select {
	case <-awake.updates:
	default:
	}
	awake.updates <- flags
}

func (a *App) haptic(HapticKind) {}

func (a *App) keepScreenOn(ctx context.Context, on bool) error {
	awake.Lock()
	defer awake.Unlock()
	awake.screen = on
	setExecutionState()
	return nil
}

func (a *App) keepAlive(reason string) (func(), error) {
	awake.Lock()
	defer awake.Unlock()
	awake.holds++
	setExecutionState()
	var once sync.Once
	return func() {
		once.Do(func() {
			awake.Lock()
			defer awake.Unlock()
			awake.holds--
			setExecutionState()
		})
	}, nil
}
