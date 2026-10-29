//go:build windows && !race

// Not under the race detector: fyne.io/systray's quit writes a flag its
// tray thread reads without synchronization, which -race reports
// whenever a tray is stopped.

package giode

import (
	"testing"
	"time"
)

// TestTray_Windows starts the tray and waits for it to be ready.
func TestTray_Windows(t *testing.T) {
	a := New(Config{})
	tr := a.Tray().SetIcon(testPNG(t)).SetTooltip("tip").SetMenu(TrayAction("Open", nil))
	for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
		tr.mu.Lock()
		ready := tr.ready
		tr.mu.Unlock()
		if ready {
			tr.Remove()
			return
		}
	}
	t.Fatal("the tray never became ready")
}
