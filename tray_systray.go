//go:build windows || (darwin && !ios) || (linux && !android)

package giode

import (
	"runtime"

	"fyne.io/systray"
)

const traySupported = true

// trayStart starts the platform tray; t.onReady runs once it exists.
func trayStart(t *Tray) {
	startSystray(t.onReady)
}

// trayApply shows s in the tray.
func trayApply(t *Tray, s trayState) {
	if len(s.icon) > 0 {
		setTrayIcon(s.icon, s.template)
	}
	systray.SetTitle(s.title)
	systray.SetTooltip(s.tooltip)
	if s.customClick || runtime.GOOS != "darwin" {
		systray.SetOnTapped(t.clicked)
	} else {
		// macOS menu bar items open their menu on a click.
		systray.SetOnTapped(nil)
	}
	systray.ResetMenu()
	addTrayItems(t.app, nil, s.items)
}

// addTrayItems adds items to the menu, or to parent's submenu.
func addTrayItems(a *App, parent *systray.MenuItem, items []TrayItem) {
	for _, it := range items {
		if it.Separator {
			if parent == nil {
				systray.AddSeparator()
			} else {
				parent.AddSeparator()
			}
			continue
		}
		var mi *systray.MenuItem
		switch {
		case parent == nil && it.Checkable:
			mi = systray.AddMenuItemCheckbox(it.Label, it.Tooltip, it.Checked)
		case parent == nil:
			mi = systray.AddMenuItem(it.Label, it.Tooltip)
		case it.Checkable:
			mi = parent.AddSubMenuItemCheckbox(it.Label, it.Tooltip, it.Checked)
		default:
			mi = parent.AddSubMenuItem(it.Label, it.Tooltip)
		}
		if it.Disabled {
			mi.Disable()
		}
		if len(it.Items) > 0 {
			addTrayItems(a, mi, it.Items)
		}
		if fn := it.OnClick; fn != nil {
			// The channel closes when the menu is rebuilt.
			go func() {
				for range mi.ClickedCh {
					a.Post(fn)
				}
			}()
		}
	}
}
