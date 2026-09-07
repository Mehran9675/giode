package icon

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Image renders an image scaled according to styles.Fit.
func Image(img image.Image, st ...styles.Styles) elements.Element {
	return &imageEl{img: img, st: first(st)}
}

// imageEl adapts widget.Image to an element.
type imageEl struct {
	img image.Image
	st  styles.Styles
}

// Layout renders the image.
func (i *imageEl) Layout(gtx layout.Context) layout.Dimensions {
	if i.st.Width != 0 {
		gtx.Constraints.Min.X, gtx.Constraints.Max.X = int(i.st.Width), int(i.st.Width)
	}
	if i.st.Height != 0 {
		gtx.Constraints.Min.Y, gtx.Constraints.Max.Y = int(i.st.Height), int(i.st.Height)
	}
	var fit widget.Fit
	switch i.st.Fit.Mode() {
	case properties.FitNone:
		fit = widget.Unscaled
	case properties.FitCover:
		fit = widget.Cover
	case properties.FitStretch:
		fit = widget.Fill
	default:
		fit = widget.Contain
	}
	return widget.Image{
		Src: paint.NewImageOp(i.img),
		Fit: fit,
	}.Layout(gtx)
}
