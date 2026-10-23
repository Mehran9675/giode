package input

import (
	"image"
	"testing"

	"gioui.org/f32"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/styles"
)

type harness struct {
	r   input.Router
	in  *Input
	max image.Point
}

func (h *harness) frame() layout.Dimensions {
	var ops op.Ops
	gtx := layout.Context{Ops: &ops, Source: h.r.Source(), Constraints: layout.Constraints{Max: h.max}}
	dims := h.in.Layout(gtx)
	h.r.Frame(&ops)
	return dims
}

func (h *harness) tap(x, y float32, src pointer.Source) {
	pos := f32.Pt(x, y)
	h.r.Queue(
		pointer.Event{Kind: pointer.Press, Source: src, Position: pos, Buttons: pointer.ButtonPrimary},
		pointer.Event{Kind: pointer.Release, Source: src, Position: pos},
	)
	h.frame()
	h.frame()
}

func TestInput_TapAnywhereFocuses(t *testing.T) {
	for _, tc := range []struct {
		name string
		x, y float32
	}{
		// Inside the padding, left of the text.
		{"padding", 4, 20},
		// Far right of an empty field, past the old 80dp text area.
		{"blank", 380, 20},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := &harness{in: New("Name", styles.Styles{Height: 40}), max: image.Pt(400, 400)}
			h.frame()
			h.tap(tc.x, tc.y, pointer.Touch)
			if !h.in.Focused() {
				t.Errorf("tap at (%v, %v) didn't focus the field", tc.x, tc.y)
			}
		})
	}
}

func TestInput_EditorFillsWidth(t *testing.T) {
	in := New("Name")
	in.SetText("ab")
	var ops op.Ops
	gtx := layout.Context{Ops: &ops, Constraints: layout.Constraints{Max: image.Pt(300, 400)}}
	if got := (&editorEl{input: in}).Layout(gtx).Size.X; got != 300 {
		t.Errorf("editor width = %d, want 300 (the whole field)", got)
	}

	// Without a width to fill (a horizontal scroll), it hugs the text.
	gtx.Constraints.Max.X = 1e6
	if got := (&editorEl{input: in}).Layout(gtx).Size.X; got >= 300 {
		t.Errorf("unbounded editor width = %d, want the text's width", got)
	}
}

func TestInput_SecondaryClickDoesNotFocus(t *testing.T) {
	h := &harness{in: New("Name", styles.Styles{Height: 40}), max: image.Pt(400, 400)}
	h.frame()
	pos := f32.Pt(4, 20)
	h.r.Queue(
		pointer.Event{Kind: pointer.Press, Source: pointer.Mouse, Position: pos, Buttons: pointer.ButtonSecondary},
		pointer.Event{Kind: pointer.Release, Source: pointer.Mouse, Position: pos},
	)
	h.frame()
	h.frame()
	if h.in.Focused() {
		t.Error("secondary click in the padding focused the field")
	}
}

func TestInput_KeyboardHint(t *testing.T) {
	in := New("URL").Keyboard(KeyboardURL)
	if in.editor.InputHint != key.HintURL {
		t.Errorf("hint = %v, want URL", in.editor.InputHint)
	}
	in.Password(true)
	if in.editor.InputHint != key.HintPassword {
		t.Errorf("password hint = %v, want Password", in.editor.InputHint)
	}
	in.Password(false)
	if in.editor.InputHint != key.HintURL {
		t.Errorf("hint after Password(false) = %v, want URL", in.editor.InputHint)
	}
}
