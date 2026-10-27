package properties

import (
	"image"
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
)

// q places the control points of a cubic Bézier approximating a
// quarter circle, as a fraction of the radius from the arc's ends
// toward the corner.
const q = 4 * (math.Sqrt2 - 1) / 3

// FillRRect calls draw to paint r with rounded corners c. A rounded
// rectangle clip makes the GPU rasterize a path over the whole box, so
// instead the straight middle is clipped as plain rectangles (nearly
// free) and only the corners as small paths; a box without radii is a
// single rectangle. The pieces don't overlap, so translucent paint
// stays even. draw is called once per piece with its clip pushed and
// must paint in r's coordinates (a paint.PaintOp, say).
//
// Pieces meet on pixel boundaries, so they're seamless under offsets.
// Under a scale or rotation their antialiased edges could show
// hairlines; exact fills the box as one path for that case.
func FillRRect(ops *op.Ops, r image.Rectangle, c Corners, exact bool, draw func()) {
	if r.Empty() {
		return
	}
	c = clampCorners(c, r.Dx(), r.Dy())
	if c.IsZero() {
		st := clip.Rect(r).Push(ops)
		draw()
		st.Pop()
		return
	}
	if exact {
		st := RRect(r, c).Push(ops)
		draw()
		st.Pop()
		return
	}
	w, h := r.Dx(), r.Dy()
	// Horizontal bands between the corners' edges: within a band the
	// left and right ends don't change, so each band is one rectangle.
	// Radii are at most half the box, so a top and a bottom corner on
	// the same side never share a band.
	ys := [...]int{0, c.TopLeft, c.TopRight, h - c.BottomLeft, h - c.BottomRight, h}
	sortInts(ys[:])
	var prev image.Rectangle
	flush := func() {
		if !prev.Empty() {
			st := clip.Rect(prev.Add(r.Min)).Push(ops)
			draw()
			st.Pop()
		}
	}
	for i := 0; i+1 < len(ys); i++ {
		y0, y1 := ys[i], ys[i+1]
		if y0 == y1 {
			continue
		}
		left, right := 0, w
		switch {
		case y0 < c.TopLeft:
			left = c.TopLeft
		case y0 >= h-c.BottomLeft:
			left = c.BottomLeft
		}
		switch {
		case y0 < c.TopRight:
			right = w - c.TopRight
		case y0 >= h-c.BottomRight:
			right = w - c.BottomRight
		}
		band := image.Rect(left, y0, right, y1)
		if band.Empty() {
			continue
		}
		// Merge with the band above when the ends line up.
		if prev.Min.X == band.Min.X && prev.Max.X == band.Max.X && prev.Max.Y == band.Min.Y && !prev.Empty() {
			prev.Max.Y = band.Max.Y
			continue
		}
		flush()
		prev = band
	}
	flush()
	for _, k := range cornersOf(r, c) {
		if k.r == 0 {
			continue
		}
		rad := float32(k.r)
		var p clip.Path
		p.Begin(ops)
		p.MoveTo(k.pt(rad, rad))
		p.LineTo(k.pt(0, rad))
		p.CubeTo(k.pt(0, rad*(1-q)), k.pt(rad*(1-q), 0), k.pt(rad, 0))
		p.Close()
		st := clip.Outline{Path: p.End()}.Op().Push(ops)
		draw()
		st.Pop()
	}
}

// StrokeRRect calls draw to paint a border of width w along the inside
// of r with rounded corners c: the outer edge follows c and the inner
// edge is rounded by each radius less w. Like FillRRect, the straight
// sides are plain rectangles and only the corners are paths, so the
// cost doesn't grow with the box; exact strokes it as one path for
// boxes under a scale or rotation.
func StrokeRRect(ops *op.Ops, r image.Rectangle, c Corners, w int, exact bool, draw func()) {
	if w <= 0 || r.Empty() {
		return
	}
	c = clampCorners(c, r.Dx(), r.Dy())
	if 2*w >= min(r.Dx(), r.Dy()) {
		// All border.
		FillRRect(ops, r, c, exact, draw)
		return
	}
	if exact {
		// The outline less the inner rounded rectangle, as one path.
		inner := r.Inset(w)
		shrink := func(v int) int { return max(0, v-w) }
		ic := Corners{shrink(c.TopLeft), shrink(c.TopRight), shrink(c.BottomRight), shrink(c.BottomLeft)}
		var p clip.Path
		p.Begin(ops)
		rrectPath(&p, r, c, false)
		rrectPath(&p, inner, ic, true)
		st := clip.Outline{Path: p.End()}.Op().Push(ops)
		draw()
		st.Pop()
		return
	}
	ks := cornersOf(r, c)
	// s is the side of each corner's square piece.
	var s [4]int
	for i, k := range ks {
		s[i] = max(k.r, w)
	}
	tl, tr, br, bl := s[0], s[1], s[2], s[3]
	W, H := r.Dx(), r.Dy()
	for _, side := range [...]image.Rectangle{
		image.Rect(tl, 0, W-tr, w),   // top
		image.Rect(W-w, tr, W, H-br), // right
		image.Rect(bl, H-w, W-br, H), // bottom
		image.Rect(0, tl, w, H-bl),   // left
	} {
		if side.Empty() {
			continue
		}
		st := clip.Rect(side.Add(r.Min)).Push(ops)
		draw()
		st.Pop()
	}
	fw := float32(w)
	for i, k := range ks {
		if k.r == 0 {
			sq := image.Rectangle{Min: k.at, Max: k.at.Add(image.Pt(k.ux*w, k.uy*w))}.Canon()
			st := clip.Rect(sq).Push(ops)
			draw()
			st.Pop()
			continue
		}
		R, side := float32(k.r), float32(s[i])
		ri := max(0, R-fw)
		var p clip.Path
		p.Begin(ops)
		p.MoveTo(k.pt(0, side))
		p.LineTo(k.pt(0, R))
		p.CubeTo(k.pt(0, R*(1-q)), k.pt(R*(1-q), 0), k.pt(R, 0))
		p.LineTo(k.pt(side, 0))
		p.LineTo(k.pt(side, fw))
		if ri > 0 {
			// side == R here, so the inner arc runs from (R, w) to (w, R).
			p.CubeTo(k.pt(fw+ri*(1-q), fw), k.pt(fw, fw+ri*(1-q)), k.pt(fw, fw+ri))
		} else {
			p.LineTo(k.pt(fw, fw))
		}
		p.LineTo(k.pt(fw, side))
		p.Close()
		st := clip.Outline{Path: p.End()}.Op().Push(ops)
		draw()
		st.Pop()
	}
}

// corner is one corner of a box: its point, its radius, and the
// directions (±1) pointing into the box along x and y.
type corner struct {
	at     image.Point
	r      int
	ux, uy int
}

// pt returns the point a along x and b along y into the box from the
// corner.
func (k corner) pt(a, b float32) f32.Point {
	return f32.Pt(float32(k.at.X)+float32(k.ux)*a, float32(k.at.Y)+float32(k.uy)*b)
}

// cornersOf returns r's corners clockwise from the top left.
func cornersOf(r image.Rectangle, c Corners) [4]corner {
	return [4]corner{
		{at: r.Min, r: c.TopLeft, ux: 1, uy: 1},
		{at: image.Pt(r.Max.X, r.Min.Y), r: c.TopRight, ux: -1, uy: 1},
		{at: r.Max, r: c.BottomRight, ux: -1, uy: -1},
		{at: image.Pt(r.Min.X, r.Max.Y), r: c.BottomLeft, ux: 1, uy: -1},
	}
}

// rrectPath adds the outline of r with corners c to p as a closed
// subpath, clockwise, or counter-clockwise when reverse is set so it
// cuts a hole out of an enclosing clockwise one.
func rrectPath(p *clip.Path, r image.Rectangle, c Corners, reverse bool) {
	ks := cornersOf(r, c)
	// Each corner contributes the arc from its point on the previous
	// side to its point on the next side, going clockwise.
	type seg struct{ from, c1, c2, to f32.Point }
	var segs [4]seg
	for i, k := range ks {
		rad := float32(k.r)
		// Clockwise, the top left arc goes from its left side up to its
		// top side; the swap of axes alternates around the box.
		a, b := k.pt(0, rad), k.pt(rad, 0)
		c1, c2 := k.pt(0, rad*(1-q)), k.pt(rad*(1-q), 0)
		if i%2 == 1 {
			a, b = k.pt(rad, 0), k.pt(0, rad)
			c1, c2 = k.pt(rad*(1-q), 0), k.pt(0, rad*(1-q))
		}
		segs[i] = seg{a, c1, c2, b}
	}
	if !reverse {
		p.MoveTo(segs[0].from)
		for i, s := range segs {
			if i > 0 {
				p.LineTo(s.from)
			}
			p.CubeTo(s.c1, s.c2, s.to)
		}
		p.Close()
		return
	}
	p.MoveTo(segs[3].to)
	for i := 3; i >= 0; i-- {
		s := segs[i]
		if i < 3 {
			p.LineTo(s.to)
		}
		p.CubeTo(s.c2, s.c1, s.from)
	}
	p.Close()
}

// RoundedCornersClear reports whether content covering inner, inside a
// box r with corners c, stays clear of the area the rounding cuts
// away: then a plain rectangle clip of r clips it the same as the
// rounded one.
func RoundedCornersClear(r, inner image.Rectangle, c Corners) bool {
	c = clampCorners(c, r.Dx(), r.Dy())
	if c.IsZero() {
		return true
	}
	if inner.Empty() {
		return true
	}
	for _, k := range cornersOf(r, c) {
		if k.r == 0 {
			continue
		}
		// The content corner nearest this box corner.
		x, y := inner.Min.X, inner.Min.Y
		if k.ux < 0 {
			x = inner.Max.X
		}
		if k.uy < 0 {
			y = inner.Max.Y
		}
		dx := float64(k.ux) * float64(x-k.at.X)
		dy := float64(k.uy) * float64(y-k.at.Y)
		rad := float64(k.r)
		if dx >= rad || dy >= rad {
			continue // beside the corner's square
		}
		// Inside the square: clear only if within the arc's circle.
		if math.Hypot(rad-dx, rad-dy) > rad {
			return false
		}
	}
	return true
}

func sortInts(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j] < a[j-1]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}
