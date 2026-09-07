package spinner

import (
	"image"
	"math"
	"time"

	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/styles"
)

// spinnerEl renders the rotating arc.
type spinnerEl struct {
	st styles.Styles
}

// Layout draws one frame of the animation and requests the next.
func (s *spinnerEl) Layout(gtx layout.Context) layout.Dimensions {
	sizeX := int(s.st.Width)
	if sizeX <= 0 {
		sizeX = 32
	}
	sizeY := int(s.st.Height)
	if sizeY <= 0 {
		sizeY = sizeX
	}
	size := gtx.Constraints.Constrain(image.Pt(sizeX, sizeY))
	if size.X <= 0 || size.Y <= 0 {
		return layout.Dimensions{}
	}
	radius := float32(size.X) / 2
	if float32(size.Y)/2 < radius {
		radius = float32(size.Y) / 2
	}
	c := s.st.Color
	if c.A == 0 {
		c = defaultColor
	}

	defer op.Offset(image.Pt(size.X/2, size.Y/2)).Push(gtx.Ops).Pop()

	dt := float32((time.Duration(gtx.Now.UnixNano()) % time.Second).Seconds())
	startAngle := dt * math.Pi * 2
	endAngle := startAngle + math.Pi*1.5

	defer clipArc(gtx.Ops, startAngle, endAngle, radius).Push(gtx.Ops).Pop()
	paintFill(gtx.Ops, c)

	// Keep animating while visible.
	gtx.Execute(op.InvalidateCmd{})
	return layout.Dimensions{Size: size}
}
