package properties

import (
	"image"
	"image/color"
	"os"
	"runtime"
	"testing"

	"gioui.org/f32"
	"gioui.org/gpu/headless"
	"gioui.org/op"
	"gioui.org/op/paint"
)

// render paints with fn into a w×h image, or skips without a GPU.
func render(t *testing.T, w, h int, fn func(ops *op.Ops)) *image.RGBA {
	t.Helper()
	if runtime.GOOS == "linux" && os.Getenv("EGL_PLATFORM") == "" {
		// Mesa renders without a display this way.
		t.Setenv("EGL_PLATFORM", "surfaceless")
	}
	win, err := headless.NewWindow(w, h)
	if err != nil {
		t.Skipf("no offscreen renderer: %v", err)
	}
	defer win.Release()
	var ops op.Ops
	paint.Fill(&ops, color.NRGBA{255, 255, 255, 255})
	fn(&ops)
	if err := win.Frame(&ops); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if err := win.Screenshot(img); err != nil {
		t.Fatal(err)
	}
	return img
}

// maxDiff returns the largest channel difference between a and b.
func maxDiff(a, b *image.RGBA) (int, image.Point) {
	worst, at := 0, image.Point{}
	for i := range a.Pix {
		d := int(a.Pix[i]) - int(b.Pix[i])
		if d < 0 {
			d = -d
		}
		if d > worst {
			worst = d
			px := i / 4
			at = image.Pt(px%a.Rect.Dx(), px/a.Rect.Dx())
		}
	}
	return worst, at
}

func fillCol(ops *op.Ops, col color.NRGBA) func() {
	return func() {
		paint.ColorOp{Color: col}.Add(ops)
		paint.PaintOp{}.Add(ops)
	}
}

func TestFillRRect_MatchesPath(t *testing.T) {
	col := color.NRGBA{R: 30, G: 90, B: 200, A: 140} // translucent shows overlaps
	r := image.Rect(10, 12, 190, 107)
	for _, c := range []Corners{UniformCorners(16), {TopLeft: 30, TopRight: 4, BottomRight: 0, BottomLeft: 47}, UniformCorners(200), {}} {
		// Under a scale the pieces show seams, which is what exact is for.
		for _, scale := range []float32{1} {
			draw := func(exact bool) func(*op.Ops) {
				return func(ops *op.Ops) {
					if scale != 1 {
						defer op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Push(ops).Pop()
					}
					FillRRect(ops, r, c, exact, fillCol(ops, col))
				}
			}
			want := render(t, 300, 200, draw(true))
			got := render(t, 300, 200, draw(false))
			if d, at := maxDiff(want, got); d > 2 {
				t.Errorf("corners %+v, scale %v: pieces differ from the path by %d at %v", c, scale, d, at)
			}
		}
	}
}

func TestStrokeRRect_MatchesPath(t *testing.T) {
	col := color.NRGBA{R: 200, G: 40, B: 40, A: 160}
	r := image.Rect(10, 12, 190, 107)
	for _, w := range []int{1, 3, 12, 60} {
		for _, c := range []Corners{UniformCorners(16), {TopLeft: 30, TopRight: 2, BottomLeft: 47}, {}} {
			want := render(t, 200, 120, func(ops *op.Ops) { StrokeRRect(ops, r, c, w, true, fillCol(ops, col)) })
			got := render(t, 200, 120, func(ops *op.Ops) { StrokeRRect(ops, r, c, w, false, fillCol(ops, col)) })
			if d, at := maxDiff(want, got); d > 2 {
				t.Errorf("width %d, corners %+v: pieces differ from the path by %d at %v", w, c, d, at)
			}
		}
	}
}

func TestRoundedCornersClear(t *testing.T) {
	r := image.Rect(0, 0, 100, 100)
	c := UniformCorners(20)
	if !RoundedCornersClear(r, image.Rect(8, 8, 92, 92), c) {
		t.Error("content inset 8 of a 20 radius is inside the arcs")
	}
	if RoundedCornersClear(r, image.Rect(2, 2, 98, 98), c) {
		t.Error("content inset 2 reaches the cut-away corners")
	}
	if !RoundedCornersClear(r, image.Rect(0, 25, 100, 75), c) {
		t.Error("content beside the corner squares is clear")
	}
	if !RoundedCornersClear(r, r, Corners{}) {
		t.Error("no radius: always clear")
	}
}

// BenchmarkRoundedCards renders a screen of rounded, bordered cards,
// as one path each (exact) and as pieces.
func BenchmarkRoundedCards(b *testing.B) {
	if runtime.GOOS == "linux" && os.Getenv("EGL_PLATFORM") == "" {
		b.Setenv("EGL_PLATFORM", "surfaceless")
	}
	for _, exact := range []bool{true, false} {
		name := "pieces"
		if exact {
			name = "path"
		}
		b.Run(name, func(b *testing.B) {
			win, err := headless.NewWindow(1080, 2400)
			if err != nil {
				b.Skipf("no offscreen renderer: %v", err)
			}
			defer win.Release()
			var ops op.Ops
			// Reading a pixel back waits for the GPU to finish.
			px := image.NewRGBA(image.Rect(0, 0, 1080, 2400))
			for i := 0; i < b.N; i++ {
				ops.Reset()
				for y := 0; y < 2400; y += 120 {
					r := image.Rect(24, y+8, 1056, y+112)
					c := UniformCorners(24)
					FillRRect(&ops, r, c, exact, fillCol(&ops, color.NRGBA{240, 240, 245, 255}))
					StrokeRRect(&ops, r, c, 3, exact, fillCol(&ops, color.NRGBA{200, 200, 210, 255}))
				}
				if err := win.Frame(&ops); err != nil {
					b.Fatal(err)
				}
				if err := win.Screenshot(px); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func drawImage(ops *op.Ops, img image.Image, at image.Point) {
	defer op.Offset(at).Push(ops).Pop()
	paint.NewImageOp(img).Add(ops)
	paint.PaintOp{}.Add(ops)
}

// benchScreen renders a phone screen of cards drawn by draw.
func benchScreen(b *testing.B, draw func(ops *op.Ops, r image.Rectangle)) {
	if runtime.GOOS == "linux" && os.Getenv("EGL_PLATFORM") == "" {
		b.Setenv("EGL_PLATFORM", "surfaceless")
	}
	win, err := headless.NewWindow(1080, 2400)
	if err != nil {
		b.Skipf("no offscreen renderer: %v", err)
	}
	defer win.Release()
	px := image.NewRGBA(image.Rect(0, 0, 1080, 2400))
	var ops op.Ops
	for i := 0; i < b.N; i++ {
		ops.Reset()
		for y := 0; y < 2400; y += 240 {
			draw(&ops, image.Rect(60, y+40, 1020, y+200))
		}
		win.Frame(&ops)
		win.Screenshot(px)
	}
}
