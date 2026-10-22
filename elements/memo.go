package elements

import (
	"reflect"
	"sync/atomic"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/state"
	"github.com/mehran9675/giode/styles/properties"
)

// MemoCache caches the painted output of a subtree across frames, so
// the subtree is only rebuilt and laid out again when something it
// depends on changed. Everything else replays the ops recorded last
// time, which costs next to nothing regardless of how large the
// subtree is.
//
// A cached subtree is rebuilt when any of these changed since it was
// recorded:
//
//   - the deps passed to Layout (compared with ==, or
//     reflect.DeepEqual for values that aren't comparable);
//   - any State the subtree read while being built (tracked
//     automatically);
//   - the layout constraints or the display density;
//   - the theme (SetMode/SetPalette);
//   - a widget inside asked for it: an animation, a focused field, a
//     widget that handled an event (see kit.MarkLive);
//   - the pointer is over the subtree, so hover, press and scroll are
//     always handled live.
//
// Values read from plain variables (not States) must be passed as
// deps, or the cache won't notice them change.
//
// Create a MemoCache once and lay it out every frame, or use Memo for
// a cache keyed by an identifier inside the view function. A cache
// laid out twice in the same frame lays the second use out uncached.
type MemoCache struct {
	ops   op.Ops
	call  op.CallOp
	dims  layout.Dimensions
	valid bool

	constraints layout.Constraints
	metric      unit.Metric
	theme       uint64
	epoch       uint64
	deps        []any
	reads       *state.Tracker

	// live is set (through kit.MarkLive) by widgets that changed
	// while the subtree was laid out.
	live bool
	// inside tracks whether the pointer is over the subtree.
	inside bool
	// portal is set when the subtree registered portals, which can't
	// be replayed from cache.
	portal bool
	// nodes are the semantic nodes recorded with the subtree, relative
	// to its origin.
	nodes []SemanticNode

	frame uint64

	// Stats, for tests and diagnostics.
	hits, misses uint64
}

// NewMemoCache returns an empty cache.
func NewMemoCache() *MemoCache {
	return new(MemoCache)
}

// Invalidate drops the cached output, forcing a rebuild on the next
// layout.
func (c *MemoCache) Invalidate() {
	c.valid = false
}

// Stats returns how many layouts were served from cache (hits) and how
// many rebuilt the subtree (misses).
func (c *MemoCache) Stats() (hits, misses uint64) {
	return c.hits, c.misses
}

// Layout lays out build's element, or replays the cached result when
// nothing it depends on changed.
func (c *MemoCache) Layout(gtx layout.Context, build func() Element, deps ...any) layout.Dimensions {
	if c.frame == frameNumber && frameNumber != 0 {
		// Already used this frame: its ops are referenced by the
		// earlier use and must not be re-recorded.
		return build().Layout(gtx)
	}
	c.frame = frameNumber

	// Pointer activity over the subtree means its widgets must see the
	// events, so it's laid out live.
	active := false
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: c,
			Kinds:  pointer.Enter | pointer.Leave | pointer.Move | pointer.Press | pointer.Release | pointer.Drag | pointer.Cancel,
		})
		if !ok {
			break
		}
		pe, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		active = true
		switch pe.Kind {
		case pointer.Leave, pointer.Cancel:
			c.inside = false
		default:
			c.inside = true
		}
	}

	reuse := c.valid && !active && !c.inside && !c.live && !c.portal &&
		c.constraints == gtx.Constraints &&
		c.metric == gtx.Metric &&
		c.theme == properties.ThemeVersion() &&
		c.epoch == cacheEpoch.Load() &&
		depsEqual(c.deps, deps) &&
		!c.reads.Changed()

	if reuse {
		c.hits++
		origin := AbsolutePosition()
		if semantics != nil {
			for _, n := range c.nodes {
				n.Bounds = n.Bounds.Add(origin)
				*semantics = append(*semantics, n)
			}
		}
		c.call.Add(gtx.Ops)
	} else {
		c.misses++
		c.record(gtx, build, deps)
	}

	// Observe the pointer over the subtree without blocking anything
	// beneath it.
	ps := pointer.PassOp{}.Push(gtx.Ops)
	cl := clip.Rect{Max: c.dims.Size}.Push(gtx.Ops)
	event.Op(gtx.Ops, c)
	cl.Pop()
	ps.Pop()
	return c.dims
}

// record rebuilds and lays out the subtree into the cache's own ops.
func (c *MemoCache) record(gtx layout.Context, build func() Element, deps []any) {
	c.ops.Reset()
	cgtx := gtx
	cgtx.Ops = &c.ops

	c.live = false
	popLive := kit.PushLiveScope(&c.live)
	reads := state.Track()
	portals := len(portalQueue)
	semStart := 0
	if semantics != nil {
		semStart = len(*semantics)
	}

	macro := op.Record(&c.ops)
	c.dims = build().Layout(cgtx)
	c.call = macro.Stop()

	reads.Stop()
	popLive()

	c.reads = reads
	c.portal = len(portalQueue) > portals
	c.nodes = c.nodes[:0]
	if semantics != nil {
		origin := AbsolutePosition()
		for _, n := range (*semantics)[semStart:] {
			n.Bounds = n.Bounds.Sub(origin)
			c.nodes = append(c.nodes, n)
		}
	}
	c.constraints = gtx.Constraints
	c.metric = gtx.Metric
	c.theme = properties.ThemeVersion()
	c.epoch = cacheEpoch.Load()
	c.deps = append(c.deps[:0], deps...)
	c.valid = true

	c.call.Add(gtx.Ops)
}

func depsEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x == nil || y == nil {
			if x != y {
				return false
			}
			continue
		}
		tx, ty := reflect.TypeOf(x), reflect.TypeOf(y)
		if tx != ty {
			return false
		}
		if tx.Comparable() {
			if x != y {
				return false
			}
			continue
		}
		if !reflect.DeepEqual(x, y) {
			return false
		}
	}
	return true
}

// cacheEpoch invalidates every cache at once when advanced.
var cacheEpoch atomic.Uint64

// InvalidateCaches forces every MemoCache and Memo to rebuild on the
// next frame. Call it after changing something caches can't observe -
// a plain variable not passed as a dep, a component driven
// programmatically (Button.Click, Input.SetText) - when the change
// must show up inside a cached subtree. It is safe to call from any
// goroutine; it doesn't request a frame by itself.
func InvalidateCaches() {
	cacheEpoch.Add(1)
}

// --- Keyed memos -------------------------------------------------------------

// frameNumber counts frames; BeginFrame advances it. Keyed memos not
// used during a frame are dropped at the next EndFrame.
var (
	frameNumber uint64
	memos       = map[any]*MemoCache{}
)

// BeginFrame starts a new frame. The app calls it before laying out
// the view; component authors don't need to.
func BeginFrame() {
	frameNumber++
}

// EndFrame finishes a frame, releasing keyed memos (see Memo) that
// weren't laid out during it. The app calls it after every frame.
func EndFrame() {
	for k, c := range memos {
		if c.frame != frameNumber {
			delete(memos, k)
		}
	}
}

// Memo returns an element that caches build's output under key and
// rebuilds it only when deps, the States it reads, its constraints or
// the theme change, or when a widget inside asks for it - see
// MemoCache for the full rules. Use it inside the view function to
// skip rebuilding parts of the UI that didn't change:
//
//	giode.Memo("sidebar", func() giode.Element {
//		return sidebar(user.Name)
//	}, user.Name)
//
// key identifies the cache across frames and must be unique among the
// memos laid out in a frame (any comparable value: a string, an ID, a
// pointer). A key not laid out for a frame is released.
func Memo(key any, build func() Element, deps ...any) Element {
	return &memoEl{key: key, build: build, deps: deps}
}

type memoEl struct {
	key   any
	build func() Element
	deps  []any
}

func (m *memoEl) Layout(gtx layout.Context) layout.Dimensions {
	c, ok := memos[m.key]
	if !ok {
		c = new(MemoCache)
		memos[m.key] = c
	}
	return c.Layout(gtx, m.build, m.deps...)
}
