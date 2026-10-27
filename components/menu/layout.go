package menu

import (
	"github.com/mehran9675/giode/components/kit"
	"image"

	"gioui.org/layout"

	"github.com/mehran9675/giode/elements"
)

// Layout renders the open menu: a window-wide dismissal layer, plus
// the panel itself queued as a portal (see elements.RegisterPortal)
// so it paints above everything else regardless of where Layout is
// called from. A submenu (see SubMenu) skips the dismissal layer -
// its own parent row's hover state controls it instead.
func (m *Menu) Layout(gtx layout.Context) layout.Dimensions {
	// Handles its input without an input pass; see kit.OpaqueInput.
	kit.OpaqueInput()
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

		// Back and Escape close the menu.
		elements.HandleDismiss(m.Close)
	}

	elements.RegisterPortal(gtx, image.Pt(int(m.pos.X), int(m.pos.Y)), func() elements.Element {
		return elements.Raw(m.panel)
	})
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
