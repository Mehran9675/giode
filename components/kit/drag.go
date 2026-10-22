package kit

import (
	"image"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

// blockDragAction is a system.Action value that isn't any of Gio's
// named window actions (gioui.org/io/system's Minimize..Move
// constants). Gio resolves overlapping system.ActionInputOp
// registrations to whichever is registered most recently (see
// gioui.org/io/input's TestPointerSystemAction/"uses topmost action
// op"), and no platform backend special-cases an unrecognized action
// value - so registering this one over a widget's own area is a safe,
// portable "not draggable" marker, distinct from claiming any real OS
// action.
const blockDragAction = system.Action(1 << 20)

// BlockDragHere marks the current clip area as not draggable, so a
// press inside it is delivered as a normal click instead of being
// hijacked into a native window-move by an ancestor's whole-window
// system.ActionInputOp(system.ActionMove) registration (see the root
// package's frameless-drag support). Call it from inside a clip scope
// that's already open for the widget's own pointer/gesture
// registration (Slider, Scroll's viewport and thumb, ...).
func BlockDragHere(gtx layout.Context) {
	system.ActionInputOp(blockDragAction).Add(gtx.Ops)
}

// DragArea marks a size-d area at the current paint position as a
// window drag handle: pressing it and moving drags a frameless window,
// like an OS title bar. On Windows the area also behaves as the
// system caption, so double-clicking it maximizes and dragging it to a
// screen edge snaps the window. Register it before laying out the
// handle's contents: widgets inside carve themselves out with
// BlockDrag/BlockDragHere, which only wins when registered later.
func DragArea(gtx layout.Context, size image.Point) {
	cl := clip.Rect{Max: size}.Push(gtx.Ops)
	system.ActionInputOp(system.ActionMove).Add(gtx.Ops)
	cl.Pop()
}

// BlockDrag is BlockDragHere for a widget that doesn't already have a
// clip scope open at the point its own area needs to be marked (most
// widget.Clickable-based components - Button, Checkbox, Tabs, Menu
// rows, Dropdown's field and options, ...): it pushes its own
// clip.Rect matching size, at the current paint position, marks it,
// then pops. Call right after laying the widget out, passing its own
// measured size.
func BlockDrag(gtx layout.Context, size image.Point) {
	cl := clip.Rect{Max: size}.Push(gtx.Ops)
	BlockDragHere(gtx)
	cl.Pop()
}
