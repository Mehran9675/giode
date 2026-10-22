package icon

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// SVG renders SVG path data (commands M, L, H, V, C, S, Q, T, A, Z in
// absolute or relative form) filled with styles.Color, scaled from a
// 24x24 viewBox. The path is parsed once and cached, so passing the
// same string every frame is cheap.
func SVG(pathData string, st ...styles.Styles) elements.Element {
	return &svgEl{path: compileCached(pathData), viewW: 24, viewH: 24, st: first(st)}
}

// SVGCustom renders SVG path data scaled from a custom viewBox.
func SVGCustom(pathData string, viewBox image.Point, st ...styles.Styles) elements.Element {
	return &svgEl{path: compileCached(pathData), viewW: float32(viewBox.X), viewH: float32(viewBox.Y), st: first(st)}
}

// SVGStroke renders SVG path data as a stroke of the given width (in
// viewBox units) instead of a fill, scaled from a 24x24 viewBox - the
// format of outline icon sets like Lucide or Feather.
func SVGStroke(pathData string, width float32, st ...styles.Styles) elements.Element {
	return &svgEl{path: compileCached(pathData), viewW: 24, viewH: 24, st: first(st), stroke: width}
}

// svgEl renders a compiled path scaled from its viewBox.
type svgEl struct {
	path         compiled
	viewW, viewH float32
	st           styles.Styles
	stroke       float32
}

// Layout fills (or strokes) the path into the element's size.
func (s *svgEl) Layout(gtx layout.Context) layout.Dimensions {
	size, ok := iconSize(gtx, s.st)
	if !ok {
		return layout.Dimensions{}
	}
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	defer op.Affine(viewBoxTransform(0, 0, s.viewW, s.viewH, size)).Push(gtx.Ops).Pop()

	col := paint.ColorOp{Color: properties.ResolveColor(s.st.Color)}
	if s.stroke > 0 {
		strokePath(gtx.Ops, s.path, s.stroke, true, col)
		return layout.Dimensions{Size: size}
	}
	var p clip.Path
	p.Begin(gtx.Ops)
	s.path.add(&p)
	defer clip.Outline{Path: p.End()}.Op().Push(gtx.Ops).Pop()
	col.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}
