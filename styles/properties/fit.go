package properties

// Fit describes how an image is scaled into its box.
type Fit string

const (
	// FitNone draws the image at its natural size.
	FitNone Fit = "none"
	// FitContain scales the image as large as possible without
	// cropping.
	FitContain Fit = "contain"
	// FitCover scales the image to cover the whole box, cropping the
	// overflow.
	FitCover Fit = "cover"
	// FitStretch scales the image to fill the box exactly.
	FitStretch Fit = "stretch"
)

// Mode resolves the value, defaulting to contain.
func (f Fit) Mode() Fit {
	if f == "" {
		return FitContain
	}
	return f
}
