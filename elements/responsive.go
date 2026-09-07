package elements

import "gioui.org/layout"

// Breakpoint is a simplified media query: it matches when the
// available width is at most Width. A Width of zero always matches and
// can serve as the fallback when listed last.
type Breakpoint struct {
	Width int
	View  func() Element
}

type responsiveEl struct {
	bps []Breakpoint
}

// Responsive renders different views depending on the available
// width. Breakpoints are evaluated in order and the first match wins,
// so list the narrowest first and a zero-width fallback last:
//
//	giode.Responsive(
//		giode.Breakpoint{Width: 600, View: phoneView},
//		giode.Breakpoint{Width: 1024, View: tabletView},
//		giode.Breakpoint{View: desktopView},
//	)
//
// Only the matching view is built each frame. Width is measured in
// the layout pixel space of the parent container.
func Responsive(bps ...Breakpoint) Element {
	return &responsiveEl{bps: bps}
}

func (r *responsiveEl) Layout(gtx layout.Context) layout.Dimensions {
	w := gtx.Constraints.Max.X
	for _, bp := range r.bps {
		if bp.Width == 0 || w <= bp.Width {
			return bp.View().Layout(gtx)
		}
	}
	return layout.Dimensions{}
}
