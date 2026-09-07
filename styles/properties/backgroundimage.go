package properties

import (
	"image"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// BackgroundImage is an image painted behind an element's content
// (CSS background-image), scaled according to Fit. A nil image is not
// painted.
type BackgroundImage = image.Image

// BackgroundFit controls how BackgroundImage is scaled.
type BackgroundFit = Fit

// PaintBackgroundImage paints img into rect, scaled according to fit.
func PaintBackgroundImage(ops *op.Ops, rect image.Rectangle, img image.Image, fit Fit) {
	if img == nil || rect.Dx() <= 0 || rect.Dy() <= 0 {
		return
	}
	isz := img.Bounds().Size()
	if isz.X <= 0 || isz.Y <= 0 {
		return
	}
	var dst image.Rectangle
	switch fit.Mode() {
	case FitNone:
		dst = image.Rectangle{Min: rect.Min, Max: rect.Min.Add(isz)}
	case FitStretch:
		dst = rect
	case FitCover:
		scale := max(float32(rect.Dx())/float32(isz.X), float32(rect.Dy())/float32(isz.Y))
		w := int(float32(isz.X) * scale)
		h := int(float32(isz.Y) * scale)
		dst = centered(rect, w, h)
	default: // contain
		scale := min(float32(rect.Dx())/float32(isz.X), float32(rect.Dy())/float32(isz.Y))
		w := int(float32(isz.X) * scale)
		h := int(float32(isz.Y) * scale)
		dst = centered(rect, w, h)
	}
	if dst.Dx() <= 0 || dst.Dy() <= 0 {
		return
	}
	defer clip.Rect(rect).Push(ops).Pop()
	sx := float32(dst.Dx()) / float32(isz.X)
	sy := float32(dst.Dy()) / float32(isz.Y)
	trans := f32.AffineId().Scale(f32.Point{}, f32.Pt(sx, sy)).Offset(f32.Pt(float32(dst.Min.X), float32(dst.Min.Y)))
	defer op.Affine(trans).Push(ops).Pop()
	paint.NewImageOp(img).Add(ops)
	paint.PaintOp{}.Add(ops)
}

func centered(rect image.Rectangle, w, h int) image.Rectangle {
	x := rect.Min.X + (rect.Dx()-w)/2
	y := rect.Min.Y + (rect.Dy()-h)/2
	return image.Rect(x, y, x+w, y+h)
}
