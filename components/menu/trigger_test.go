package menu

import (
	"image"
	"testing"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/internal/uitest"
	"github.com/mehran9675/giode/styles"
)

func TestTriggerMenuOpensBelowTriggerInsideRowAndStack(t *testing.T) {
	m := New().Item("New", func() {}).Item("Open", func() {})
	btn := elements.NewPressable(func(elements.PressState) elements.Element {
		return elements.Box(styles.Styles{Width: 80, Height: 30})
	}).Label("File").OnClick(func(elements.ClickEvent) {})
	trigger := m.Trigger(btn)
	h := uitest.New(func() elements.Element {
		return elements.Stack(styles.Styles{Gap: 20},
			elements.Box(styles.Styles{Height: 100}),
			elements.Row(styles.Styles{Gap: 16, Wrap: "nowrap"},
				elements.Box(styles.Styles{Width: 200, Height: 30}),
				trigger,
			),
			m,
		)
	})
	h.Frames(2)
	file, ok := h.Find(elements.SemanticButton, "File")
	if !ok {
		t.Fatal("trigger not found")
	}
	h.Click(file.Bounds.Min.Add(image.Pt(5, 5)))
	item, ok := h.Find(elements.SemanticMenuItem, "Open")
	if !ok {
		t.Fatal("menu didn't open")
	}
	if item.Bounds.Min.X < file.Bounds.Min.X || item.Bounds.Min.X > file.Bounds.Min.X+20 ||
		item.Bounds.Min.Y < file.Bounds.Max.Y || item.Bounds.Min.Y > file.Bounds.Max.Y+80 {
		t.Errorf("menu item at %v, want just below the trigger at %v", item.Bounds, file.Bounds)
	}
}
