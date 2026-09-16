// Package slider provides a stateful horizontal slider.
package slider

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"github.com/mehran9675/giode/components/kit"

	"github.com/mehran9675/giode/styles"
)

// Slider is a stateful horizontal slider in the range [0, 1]. Create
// it once with New and lay it out every frame.
type Slider struct {
	value    float32
	drag     gesture.Drag
	hover    gesture.Hover
	onChange func(v float32)
	st       styles.Styles
}

// New returns a Slider. The styles argument is optional: Color is
// the fill, Background is the track, BorderColor is the thumb.
func New(st ...styles.Styles) *Slider {
	s := &Slider{value: 0}
	if len(st) > 0 {
		s.st = st[0]
	}
	return s
}

// Value returns the current value in the range [0, 1].
func (s *Slider) Value() float32 {
	return s.value
}

// SetValue sets the value without invoking OnChange.
func (s *Slider) SetValue(v float32) {
	s.value = clamp01(v)
}

// OnChange sets the handler invoked while dragging. It returns s for
// chaining.
func (s *Slider) OnChange(fn func(v float32)) *Slider {
	s.onChange = fn
	return s
}

// Layout lays out and updates the slider. It expands to the available
// width.
func (s *Slider) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	size := image.Pt(gtx.Constraints.Min.X, height)

	// The drag and hover regions must be clipped to the slider itself,
	// otherwise they would span whatever outer area it is laid out in
	// and steal events from siblings.
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	// Keep the slider draggable even inside a frameless window's
	// whole-window drag area (see the root package's layoutDrag):
	// without this, dragging the thumb would be hijacked into a
	// native window-move instead.
	kit.BlockDragHere(gtx)
	s.drag.Add(gtx.Ops)
	s.hover.Add(gtx.Ops)
	hovered := s.hover.Update(gtx.Source)
	changed := false
	for {
		ev, ok := s.drag.Update(gtx.Metric, gtx.Source, gesture.Horizontal)
		if !ok {
			break
		}
		switch ev.Kind {
		case pointer.Press, pointer.Drag:
			if ev.Priority == pointer.Grabbed {
				v := s.valueFromX(float32(ev.Position.X), float32(size.X))
				if v != s.value {
					s.value = v
					changed = true
				}
			}
		}
	}
	if changed && s.onChange != nil {
		s.onChange(s.value)
		gtx.Execute(op.InvalidateCmd{})
	}

	thumb := s.st.BorderColor
	switch {
	case s.drag.Dragging():
		thumb = kit.Pressed(resolve(thumb, defaultThumb))
		pointer.CursorGrabbing.Add(gtx.Ops)
	case hovered:
		thumb = kit.Hovered(resolve(thumb, defaultThumb))
		pointer.CursorGrab.Add(gtx.Ops)
	}

	paintTrack(gtx, size, s.st.Background)
	paintFill(gtx, size, s.value, s.st.Color)
	paintThumb(gtx, size, s.value, thumb)
	return layout.Dimensions{Size: size}
}

func (s *Slider) valueFromX(x, width float32) float32 {
	if width <= 2*padding {
		return 0
	}
	return clamp01((x - padding) / (width - 2*padding))
}

func clamp01(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
