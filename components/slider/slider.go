// Package slider provides a stateful horizontal slider.
package slider

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/mehran9675/giode/styles"
)

// Slider is a stateful horizontal slider in the range [0, 1]. Create
// it once with New and lay it out every frame.
type Slider struct {
	value    float32
	drag     gesture.Drag
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

	// The drag region must be clipped to the slider itself, otherwise
	// it would span whatever outer area it is laid out in and steal
	// events from siblings.
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	s.drag.Add(gtx.Ops)
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

	paintTrack(gtx, size, properties.CalcColor(s.st.Background))
	paintFill(gtx, size, s.value, properties.CalcColor(s.st.Color))
	paintThumb(gtx, size, s.value, properties.CalcColor(s.st.BorderColor))
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
