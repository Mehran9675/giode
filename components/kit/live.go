package kit

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op"
)

// Change tracking.
//
// Giode repaints only what changed: an elements.Memo caches the ops of
// its subtree and replays them while nothing it depends on changed.
// A widget inside a memo that changes on its own - an animation, a
// focused text field, a widget that just handled an event - must tell
// the enclosing memos so they lay the subtree out again on the next
// frame instead of replaying a stale picture. Every built-in component
// does this; a custom component does it by calling RequestFrame (for
// animations) or MarkLive (for anything else that changed).
//
// Layout runs on a single goroutine, so the scope stack below needs no
// locking.

var liveStack []*bool

// PushLiveScope opens a change-tracking scope: until the returned pop
// function runs, MarkLive sets *flag. Scopes nest; MarkLive marks
// every open scope, since a change inside a nested memo is also a
// change inside each memo enclosing it. Elements that cache output
// (elements.Memo) use it; ordinary components don't need to.
func PushLiveScope(flag *bool) (pop func()) {
	liveStack = append(liveStack, flag)
	n := len(liveStack)
	return func() {
		liveStack = liveStack[:n-1]
	}
}

// MarkLive reports that the widget being laid out changed, or will
// change without further input, so cached output enclosing it must not
// be reused on the next frame.
func MarkLive() {
	for _, f := range liveStack {
		*f = true
	}
}

// RequestFrame asks for another frame as soon as possible and marks
// the enclosing caches live. Call it from an animating widget, or after
// a handler changed state the current frame already read.
func RequestFrame(gtx layout.Context) {
	gtx.Execute(op.InvalidateCmd{})
	MarkLive()
}

// RequestFrameAt asks for a frame at t and marks the enclosing caches
// live. Use it for timed changes (a tooltip delay, a caret blink)
// instead of redrawing every frame until then.
func RequestFrameAt(gtx layout.Context, t time.Time) {
	gtx.Execute(op.InvalidateCmd{At: t})
	MarkLive()
}
