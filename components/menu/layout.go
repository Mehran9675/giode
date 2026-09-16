package menu

import (
	"image"

	"gioui.org/io/key"
	"gioui.org/layout"

	"github.com/mehran9675/giode/elements"
)

// Layout renders the open menu: a window-wide dismissal layer, plus
// the panel itself queued as a portal (see elements.RegisterPortal)
// so it paints above everything else regardless of where Layout is
// called from. A submenu (see SubMenu) skips the dismissal layer -
// its own parent row's hover state controls it instead.
func (m *Menu) Layout(gtx layout.Context) layout.Dimensions {
	if !m.open {
		return layout.Dimensions{}
	}

	if !m.isSubmenu {
		// Any primary click outside the panel closes it.
		if m.scrim.Clicked(gtx) {
			m.Close()
			return layout.Dimensions{}
		}
		m.scrim.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Max}
		})

		// Escape closes the menu. Focus is requested below, before the
		// events of this frame are delivered.
		for {
			ev, ok := gtx.Event(key.Filter{Focus: m, Name: key.NameEscape})
			if !ok {
				break
			}
			if e, ok := ev.(key.Event); ok && e.State == key.Press {
				m.Close()
				return layout.Dimensions{}
			}
		}
	}

	elements.RegisterPortal(gtx, image.Pt(int(m.pos.X), int(m.pos.Y)), func() elements.Element {
		return elements.Raw(m.panel)
	})

	if !m.isSubmenu {
		// Keep focus for the Escape handler while open.
		gtx.Execute(key.FocusCmd{Tag: m})
	}
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
