package slider

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// paintTrack draws the background track.
func paintTrack(gtx layout.Context, size image.Point, track color.NRGBA) {
	c := resolve(track, defaultTrack)
	trackRect := image.Rect(0, (height-trackH)/2, size.X, (height-trackH)/2+trackH)
	defer clip.UniformRRect(trackRect, trackH/2).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, c)
}

// paintFill draws the filled portion up to v.
func paintFill(gtx layout.Context, size image.Point, v float32, fill color.NRGBA) {
	c := resolve(fill, defaultFill)
	w := int(v * float32(size.X))
	if w <= trackH/2 {
		return
	}
	fillRect := image.Rect(0, (height-trackH)/2, w, (height-trackH)/2+trackH)
	defer clip.UniformRRect(fillRect, trackH/2).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, c)
}

// paintThumb draws the draggable thumb.
func paintThumb(gtx layout.Context, size image.Point, v float32, thumb color.NRGBA) {
	c := resolve(thumb, defaultThumb)
	cx := padding + int(v*float32(size.X-2*padding))
	r := image.Rect(cx-thumbD/2, (height-thumbD)/2, cx+thumbD/2, (height-thumbD)/2+thumbD)
	defer clip.Ellipse(r).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, c)
}
