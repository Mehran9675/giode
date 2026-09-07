package drawer

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/elements"
)

const animDuration = 200 * time.Millisecond

// Layout renders the drawer: a fading scrim with the panel sliding
// from its side. content is rebuilt every frame while visible.
func (d *Drawer) Layout(gtx layout.Context, content elements.Element) layout.Dimensions {
	if !d.open && d.closeT.IsZero() {
		return layout.Dimensions{}
	}

	// Progress: 0 closed, 1 open.
	var progress float32
	if d.open {
		elapsed := gtx.Now.Sub(d.openT)
		if elapsed >= animDuration {
			progress = 1
		} else {
			progress = float32(elapsed) / float32(animDuration)
			gtx.Execute(op.InvalidateCmd{})
		}
	} else {
		elapsed := gtx.Now.Sub(d.closeT)
		if elapsed >= animDuration {
			d.closeT = time.Time{}
			return layout.Dimensions{}
		}
		progress = 1 - float32(elapsed)/float32(animDuration)
		gtx.Execute(op.InvalidateCmd{})
	}

	win := gtx.Constraints.Max
	width := d.width
	if width > win.X {
		width = win.X
	}

	// Scrim.
	if d.scrim.Clicked(gtx) {
		d.Close()
	}
	d.scrim.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		scrimColor := color.NRGBA{A: uint8(0.4 * 0xff * progress)}
		defer clip.Rect{Max: win}.Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, scrimColor)
		return layout.Dimensions{Size: win}
	})

	// Panel.
	var x int
	if d.side == Left {
		x = -width + int(float32(width)*progress)
	} else {
		x = win.X - int(float32(width)*progress)
	}
	off := op.Offset(image.Pt(x, 0)).Push(gtx.Ops)
	defer clip.Rect{Max: image.Pt(width, win.Y)}.Push(gtx.Ops).Pop()
	st := d.st
	if st.Background.A == 0 {
		st.Background = defaultBackground
	}
	panelGtx := gtx
	panelGtx.Constraints.Min.X = width
	panelGtx.Constraints.Max.X = width
	elements.Box(st, content).Layout(panelGtx)
	off.Pop()

	return layout.Dimensions{Size: win}
}
