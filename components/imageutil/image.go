// Package imageutil renders plain images with fit modes.
package imageutil

import (
	"image"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

// New renders img as a stateless element. Width and Height fix the
// box size (the natural size is used when unset) and Fit controls the
// scaling: none, contain, cover or stretch.
func New(img image.Image, st ...styles.Styles) elements.Element {
	var s styles.Styles
	if len(st) > 0 {
		s = st[0]
	}
	return &imageEl{img: img, st: s}
}
