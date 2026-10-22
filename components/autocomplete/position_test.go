package autocomplete

import (
	"image"
	"testing"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/internal/uitest"
	"github.com/mehran9675/giode/styles"
)

func TestSuggestionsOpenBelowFieldAndCloseOnOutsideClick(t *testing.T) {
	a := New("City", []string{"Paris", "Parma", "London"}, styles.Styles{Width: 200})
	h := uitest.New(func() elements.Element {
		return elements.Stack(styles.Styles{Gap: 20},
			elements.Box(styles.Styles{Height: 100}),
			elements.Row(styles.Styles{Gap: 16, Wrap: "nowrap"},
				elements.Box(styles.Styles{Width: 150, Height: 20}),
				a,
			),
		)
	})
	h.Frame()
	a.field.ChangeText("Par") // as if typed
	h.Frames(2)
	field, ok := h.Find(elements.SemanticInput, "City")
	if !ok {
		t.Fatal("field not found")
	}
	opt, ok := h.Find(elements.SemanticText, "Parma")
	if !ok {
		t.Fatal("suggestions didn't open")
	}
	if opt.Bounds.Min.X < field.Bounds.Min.X || opt.Bounds.Min.X > field.Bounds.Min.X+30 ||
		opt.Bounds.Min.Y < field.Bounds.Max.Y || opt.Bounds.Min.Y > field.Bounds.Max.Y+60 {
		t.Errorf("suggestion at %v, want just below the field at %v", opt.Bounds, field.Bounds)
	}

	h.Click(image.Pt(700, 500))
	if _, open := h.Find(elements.SemanticText, "Parma"); open {
		t.Error("click outside didn't hide the suggestions")
	}
}
