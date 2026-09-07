package icon

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// SVG renders an SVG path (commands M, L, H, V, C, S, Q, T, A, Z in
// absolute or relative form) scaled from a 24x24 viewBox.
func SVG(pathData string, st ...styles.Styles) elements.Element {
	return svg(pathData, 24, 24, first(st))
}

// SVGCustom renders an SVG path scaled from a custom viewBox.
func SVGCustom(pathData string, viewBox image.Point, st ...styles.Styles) elements.Element {
	return svg(pathData, viewBox.X, viewBox.Y, first(st))
}

// svg returns an element that fills an SVG path.
func svg(pathData string, viewW, viewH int, st styles.Styles) elements.Element {
	return &svgEl{
		path:  pathData,
		viewW: viewW,
		viewH: viewH,
		st:    st,
	}
}

// svgEl renders a filled SVG path scaled from its viewBox.
type svgEl struct {
	path  string
	viewW int
	viewH int
	st    styles.Styles
}

// Layout fills the path into the element's size.
func (s *svgEl) Layout(gtx layout.Context) layout.Dimensions {
	sizeX := int(s.st.Width)
	sizeY := int(s.st.Height)
	if sizeX <= 0 {
		sizeX = 24
	}
	if sizeY <= 0 {
		sizeY = 24
	}
	size := gtx.Constraints.Constrain(image.Pt(sizeX, sizeY))
	if size.X <= 0 || size.Y <= 0 {
		return layout.Dimensions{}
	}

	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	if s.viewW > 0 && s.viewH > 0 {
		scale := f32.AffineId().Scale(f32.Point{}, f32.Pt(float32(size.X)/float32(s.viewW), float32(size.Y)/float32(s.viewH)))
		defer op.Affine(scale).Push(gtx.Ops).Pop()
	}

	var p clip.Path
	p.Begin(gtx.Ops)
	parsePath(&p, s.path)
	paint.FillShape(gtx.Ops, properties.ResolveColor(s.st.Color), clip.Outline{Path: p.End()}.Op())
	return layout.Dimensions{Size: size}
}
