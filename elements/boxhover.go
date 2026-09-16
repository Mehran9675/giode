package elements

import (
	"gioui.org/gesture"
	"gioui.org/layout"

	"github.com/mehran9675/giode/styles"
)

// OnHover sets styles applied while the pointer is over the element,
// merged on top of its base styles for both layout and painting:
//
//	card := giode.Box(base, children...).OnHover(giode.Styles{
//		Background: "#1e40af",
//		Cursor:     giode.CursorPointer,
//	})
//
// Hover state lives in the element, so create it once and reuse it
// every frame; rebuilding it in the view function resets the hover on
// each frame.
func (b *BoxElement) OnHover(st styles.Styles) *BoxElement {
	b.hover = new(gesture.Hover)
	b.hoverSt = st
	return b
}

// Hovered reports whether the pointer was over the element on the last
// frame.
func (b *BoxElement) Hovered() bool {
	return b.hovered
}

// resolveStyles processes hover events and returns the styles to lay
// out and paint with: the base styles, with the hover styles merged
// on top while hovered.
func (b *BoxElement) resolveStyles(gtx layout.Context) styles.Styles {
	if b.hover == nil {
		return b.st
	}
	b.hovered = b.hover.Update(gtx.Source)
	if !b.hovered {
		return b.st
	}
	return styles.Merge(b.st, b.hoverSt)
}
