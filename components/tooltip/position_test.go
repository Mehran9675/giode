package tooltip

import (
	"image"
	"testing"

	"github.com/mehran9675/giode/components/scroll"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/internal/uitest"
	"github.com/mehran9675/giode/styles"
)

// The trigger sits in a Row inside a Stack inside a scrolled Scroll:
// three kinds of containers that place children after laying them out.
func TestTooltipShowsBelowTriggerInNestedContainers(t *testing.T) {
	trigger := elements.Box(styles.Styles{Width: 100, Height: 30})
	tip := New(trigger, "A helpful hint")
	sc := scroll.New(func() elements.Element {
		return elements.Stack(styles.Styles{Gap: 20},
			elements.Box(styles.Styles{Height: 300}),
			elements.Row(styles.Styles{Gap: 16, Wrap: "nowrap"},
				elements.Box(styles.Styles{Width: 120, Height: 30}),
				tip,
			),
			elements.Box(styles.Styles{Height: 900}),
		)
	})
	h := uitest.New(func() elements.Element {
		return elements.Box(styles.Styles{Width: -1, Height: -1}, sc)
	})
	h.Frame()
	sc.ScrollTo(200)
	h.Frames(2)

	// The trigger is at x 136 (120 + 16) and y 120 (300 + 20 - 200).
	h.Move(image.Pt(150, 130))
	hint, ok := h.Find(elements.SemanticText, "A helpful hint")
	if !ok {
		t.Fatal("tooltip didn't show")
	}
	if hint.Bounds.Min.X < 136 || hint.Bounds.Min.X > 160 || hint.Bounds.Min.Y < 150 || hint.Bounds.Min.Y > 180 {
		t.Errorf("tooltip text at %v, want just below the trigger at (136, 120)-(236, 150)", hint.Bounds)
	}
}
