package properties

import (
	"image"
	"image/color"
	"testing"

	"gioui.org/op"
)

// The nine slices draw the same shadow as blurring the whole box.
func TestBlurredShadow_SlicesMatchFullBlur(t *testing.T) {
	col := color.NRGBA{R: 10, G: 20, B: 40, A: 200}
	for _, c := range []Corners{UniformCorners(16), {TopLeft: 30, BottomRight: 6}, {}} {
		r := image.Rect(60, 50, 360, 250)
		sliced := render(t, 420, 300, func(ops *op.Ops) { paintBlurred(ops, r, c, 24, col) })
		full := render(t, 420, 300, func(ops *op.Ops) {
			m := 36
			img := blurredRRect(r.Dx(), r.Dy(), m, c, 24, col)
			drawImage(ops, img, r.Min.Sub(image.Pt(m, m)))
		})
		if d, at := maxDiff(full, sliced); d > 3 {
			t.Errorf("corners %+v: slices differ from the full blur by %d at %v", c, d, at)
		}
	}
}

func TestBlurredShadow_Cached(t *testing.T) {
	var ops op.Ops
	paintBlurred(&ops, image.Rect(0, 0, 400, 300), UniformCorners(12), 20, color.NRGBA{A: 100})
	shadowCache.Lock()
	n := len(shadowCache.m)
	shadowCache.Unlock()
	// Another size with the same radii and blur reuses the image.
	paintBlurred(&ops, image.Rect(0, 0, 900, 120), UniformCorners(12), 20, color.NRGBA{A: 100})
	shadowCache.Lock()
	defer shadowCache.Unlock()
	if len(shadowCache.m) != n {
		t.Errorf("cache grew from %d to %d for a different box size", n, len(shadowCache.m))
	}
}

func BenchmarkShadowedCards(b *testing.B) {
	benchScreen(b, func(ops *op.Ops, r image.Rectangle) {
		PaintBoxShadowCorners(ops, r, UniformCorners(24), BoxShadow{Y: 8, Blur: 24, Color: "#00000040"})
	})
}
