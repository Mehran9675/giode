package properties

import (
	"strconv"
	"strings"
	"time"
)

// Transition animates the change between an element's base and hover
// styles (CSS transition). Colors (Background, Color, BorderColor) and
// Opacity interpolate; other properties switch immediately. The zero
// value disables it.
type Transition struct {
	Duration time.Duration
	Easing   Easing
}

// Easing maps linear progress in 0-1 to eased progress, like a CSS
// timing function. Besides the constants, any CSS cubic-bezier() works:
//
//	properties.Easing("cubic-bezier(0.34, 1.56, 0.64, 1)")
type Easing string

const (
	// EaseInOut is the default easing: slow start and end.
	EaseInOut Easing = "ease-in-out"
	// Ease is CSS's default timing function: a quick start and gentle
	// end.
	Ease       Easing = "ease"
	EaseLinear Easing = "linear"
	// EaseIn starts slowly and accelerates.
	EaseIn Easing = "ease-in"
	// EaseOut starts fast and decelerates.
	EaseOut Easing = "ease-out"
	// EaseOutBack overshoots the end slightly, then settles.
	EaseOutBack Easing = "ease-out-back"
	// EaseInBack pulls back slightly before moving.
	EaseInBack Easing = "ease-in-back"
	// EaseOutBounce bounces at the end, like a dropped ball.
	EaseOutBounce Easing = "ease-out-bounce"
)

// Apply eases t, clamped to 0-1. Back easings may return values
// slightly outside 0-1 by design.
func (e Easing) Apply(t float32) float32 {
	t = max(0, min(1, t))
	switch e {
	case EaseLinear:
		return t
	case Ease:
		return cubicBezier(0.25, 0.1, 0.25, 1, t)
	case EaseIn:
		return cubicBezier(0.42, 0, 1, 1, t)
	case EaseOut:
		return cubicBezier(0, 0, 0.58, 1, t)
	case EaseOutBack:
		const c1, c3 = 1.70158, 2.70158
		u := t - 1
		return 1 + c3*u*u*u + c1*u*u
	case EaseInBack:
		const c1, c3 = 1.70158, 2.70158
		return c3*t*t*t - c1*t*t
	case EaseOutBounce:
		const n, d = 7.5625, 2.75
		switch {
		case t < 1/d:
			return n * t * t
		case t < 2/d:
			t -= 1.5 / d
			return n*t*t + 0.75
		case t < 2.5/d:
			t -= 2.25 / d
			return n*t*t + 0.9375
		default:
			t -= 2.625 / d
			return n*t*t + 0.984375
		}
	case EaseInOut, "":
		return cubicBezier(0.42, 0, 0.58, 1, t)
	}
	if x1, y1, x2, y2, ok := parseCubicBezier(string(e)); ok {
		return cubicBezier(x1, y1, x2, y2, t)
	}
	return cubicBezier(0.42, 0, 0.58, 1, t)
}

// parseCubicBezier parses "cubic-bezier(x1, y1, x2, y2)".
func parseCubicBezier(s string) (x1, y1, x2, y2 float32, ok bool) {
	if !strings.HasPrefix(s, "cubic-bezier(") {
		return
	}
	args, ok := funcArgs(s)
	if !ok || len(args) != 4 {
		return 0, 0, 0, 0, false
	}
	var v [4]float32
	for i, a := range args {
		f, err := strconv.ParseFloat(a, 32)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		v[i] = float32(f)
	}
	// x coordinates must stay in 0-1 for the curve to be a function.
	if v[0] < 0 || v[0] > 1 || v[2] < 0 || v[2] > 1 {
		return 0, 0, 0, 0, false
	}
	return v[0], v[1], v[2], v[3], true
}

// cubicBezier evaluates a CSS cubic-bezier timing function at x = t:
// it solves the curve's x(s) = t for s with Newton's method (falling
// back to bisection) and returns y(s).
func cubicBezier(x1, y1, x2, y2, t float32) float32 {
	bez := func(a, b, s float32) float32 {
		// Bernstein form with P0 = 0 and P3 = 1.
		u := 1 - s
		return 3*u*u*s*a + 3*u*s*s*b + s*s*s
	}
	dbez := func(a, b, s float32) float32 {
		u := 1 - s
		return 3*u*u*a + 6*u*s*(b-a) + 3*s*s*(1-b)
	}
	s := t
	for i := 0; i < 8; i++ {
		x := bez(x1, x2, s) - t
		if x > -1e-5 && x < 1e-5 {
			return bez(y1, y2, s)
		}
		d := dbez(x1, x2, s)
		if d > -1e-6 && d < 1e-6 {
			break
		}
		s -= x / d
	}
	lo, hi := float32(0), float32(1)
	s = t
	for i := 0; i < 32; i++ {
		x := bez(x1, x2, s)
		if x > -1e-5+t && x < 1e-5+t {
			break
		}
		if x < t {
			lo = s
		} else {
			hi = s
		}
		s = (lo + hi) / 2
	}
	return bez(y1, y2, s)
}

// LerpColor blends from a to b by t in 0-1. Unset colors blend as
// transparent versions of the other color, so fading a background in
// from nothing doesn't pass through black.
func LerpColor(a, b Color, t float32) Color {
	if t <= 0 {
		return a
	}
	if t >= 1 || a == b {
		return b
	}
	ca, cb := CalcColor(a), CalcColor(b)
	if a == "" {
		ca = cb
		ca.A = 0
	}
	if b == "" {
		cb = ca
		cb.A = 0
	}
	l := func(x, y uint8) uint8 { return uint8(float32(x) + (float32(y)-float32(x))*t + 0.5) }
	ca.R, ca.G, ca.B, ca.A = l(ca.R, cb.R), l(ca.G, cb.G), l(ca.B, cb.B), l(ca.A, cb.A)
	return CalcColorReverse(ca)
}
