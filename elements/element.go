// Package elements provides the stateless building blocks of Giode:
// box, flex, text and raw. They are cheap to build and are meant to be
// created inside the view function on every frame. Stateful components
// (see components/) hold the persistent state interactive widgets
// need.
package elements

import "gioui.org/layout"

// Element is anything that can lay itself out in a frame.
type Element interface {
	Layout(gtx layout.Context) layout.Dimensions
}
