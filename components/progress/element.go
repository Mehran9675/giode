package progress

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/styles"
)

// progressEl renders the track and fill.
type progressEl struct {
	value float32
	st    styles.Styles
}

// Layout renders the bar.
func (p *progressEl) Layout(gtx layout.Context) layout.Dimensions {
	height := int(p.st.Height)
	if height <= 0 {
		height = 8
	}
	track := p.st.Background
	if track.A == 0 {
		track = defaultTrack
	}
	fill := p.st.Color
	if fill.A == 0 {
		fill = defaultFill
	}
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	size := image.Pt(gtx.Constraints.Min.X, height)

	defer clip.UniformRRect(image.Rectangle{Max: size}, height/2).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, track)

	v := p.value
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	fillW := int(v * float32(size.X))
	if fillW > height/2 {
		defer clip.UniformRRect(image.Rect(0, 0, fillW, height), height/2).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, fill)
	}
	return layout.Dimensions{Size: size}
}
