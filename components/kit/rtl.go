package kit

import (
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

// RTL reports whether gtx lays out right to left: an Arabic, Persian,
// Hebrew, ... locale.
func RTL(gtx layout.Context) bool {
	return gtx.Locale.Direction == system.RTL
}

// MirrorRTL flips what's drawn next, and the pointer positions it
// receives, horizontally across an area width pixels wide when gtx
// reads right to left, until the returned function runs. Drawing
// without text - a slider's track, a progress bar, a switch's knob -
// then runs from the right with no other change. Text drawn under it
// would read backwards.
func MirrorRTL(gtx layout.Context, width int) (pop func()) {
	if !RTL(gtx) {
		return func() {}
	}
	t := f32.AffineId().Scale(f32.Point{}, f32.Pt(-1, 1)).Offset(f32.Pt(float32(width), 0))
	return op.Affine(t).Push(gtx.Ops).Pop
}

// GlyphForward resolves GlyphChevronForward for gtx: the chevron
// pointing right, or left right to left.
func GlyphForward(gtx layout.Context) Glyph {
	if RTL(gtx) {
		return GlyphChevronLeft
	}
	return GlyphChevronRight
}

// GlyphBack resolves GlyphChevronBack for gtx.
func GlyphBack(gtx layout.Context) Glyph {
	if RTL(gtx) {
		return GlyphChevronRight
	}
	return GlyphChevronLeft
}
