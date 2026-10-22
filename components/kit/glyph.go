package kit

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// Glyph is one of the few small indicator shapes components draw
// themselves (chevrons, a check mark) so they don't depend on any
// icon set. Apps bring their own icons through the icon package.
type Glyph int

const (
	GlyphChevronDown Glyph = iota
	GlyphChevronUp
	GlyphChevronLeft
	GlyphChevronRight
	GlyphCheck
	GlyphClose
	GlyphMinus
	GlyphPlus
	// GlyphMaximize is a window's maximize button: a square.
	GlyphMaximize
	// GlyphRestore is a maximized window's restore button: two
	// overlapping squares.
	GlyphRestore
	// GlyphCopy is a copy-to-clipboard indicator: two stacked pages.
	GlyphCopy
)

// glyphLines holds each glyph as polylines on a 24x24 grid.
var glyphLines = map[Glyph][][]f32.Point{
	GlyphChevronDown:  {{{X: 6, Y: 9}, {X: 12, Y: 15}, {X: 18, Y: 9}}},
	GlyphChevronUp:    {{{X: 6, Y: 15}, {X: 12, Y: 9}, {X: 18, Y: 15}}},
	GlyphChevronLeft:  {{{X: 15, Y: 6}, {X: 9, Y: 12}, {X: 15, Y: 18}}},
	GlyphChevronRight: {{{X: 9, Y: 6}, {X: 15, Y: 12}, {X: 9, Y: 18}}},
	GlyphCheck:        {{{X: 5, Y: 12.5}, {X: 10, Y: 17.5}, {X: 19, Y: 7}}},
	GlyphClose:        {{{X: 6, Y: 6}, {X: 18, Y: 18}}, {{X: 18, Y: 6}, {X: 6, Y: 18}}},
	GlyphMinus:        {{{X: 5, Y: 12}, {X: 19, Y: 12}}},
	GlyphPlus:         {{{X: 5, Y: 12}, {X: 19, Y: 12}}, {{X: 12, Y: 5}, {X: 12, Y: 19}}},
	GlyphMaximize:     {{{X: 5, Y: 5}, {X: 19, Y: 5}, {X: 19, Y: 19}, {X: 5, Y: 19}, {X: 5, Y: 5}}},
	GlyphCopy: {
		{{X: 9, Y: 9}, {X: 19, Y: 9}, {X: 19, Y: 20}, {X: 9, Y: 20}, {X: 9, Y: 9}},
		{{X: 5, Y: 15}, {X: 5, Y: 4}, {X: 15, Y: 4}, {X: 15, Y: 6}},
	},
	GlyphRestore: {
		{{X: 4, Y: 8}, {X: 16, Y: 8}, {X: 16, Y: 20}, {X: 4, Y: 20}, {X: 4, Y: 8}},
		{{X: 8, Y: 8}, {X: 8, Y: 4}, {X: 20, Y: 4}, {X: 20, Y: 16}, {X: 16, Y: 16}},
	},
}

// DrawGlyph paints g in col at size dp, with round caps and joins, and
// returns its dimensions.
func DrawGlyph(gtx layout.Context, g Glyph, size unit.Dp, col color.NRGBA) layout.Dimensions {
	px := gtx.Dp(size)
	sz := gtx.Constraints.Constrain(image.Pt(px, px))
	if sz.X <= 0 || sz.Y <= 0 {
		return layout.Dimensions{Size: sz}
	}
	scale := float32(min(sz.X, sz.Y)) / 24
	defer op.Affine(f32.AffineId().Scale(f32.Point{}, f32.Pt(scale, scale))).Push(gtx.Ops).Pop()
	const width = 2.2
	for _, line := range glyphLines[g] {
		var p clip.Path
		p.Begin(gtx.Ops)
		p.MoveTo(line[0])
		for _, pt := range line[1:] {
			p.LineTo(pt)
		}
		st := clip.Stroke{Path: p.End(), Width: width}.Op().Push(gtx.Ops)
		paint.ColorOp{Color: col}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		st.Pop()
		for _, pt := range line {
			var d clip.Path
			d.Begin(gtx.Ops)
			d.MoveTo(f32.Pt(pt.X+width/2, pt.Y))
			d.ArcTo(pt, pt, 2*math.Pi)
			d.Close()
			ds := clip.Outline{Path: d.End()}.Op().Push(gtx.Ops)
			paint.ColorOp{Color: col}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			ds.Pop()
		}
	}
	return layout.Dimensions{Size: sz}
}
