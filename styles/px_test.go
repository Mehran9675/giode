package styles

import (
	"testing"

	"gioui.org/unit"

	"github.com/mehran9675/giode/styles/properties"
)

func TestPx_ScalesLengthsKeepsSentinels(t *testing.T) {
	st := Styles{
		Width: -1, Height: 10, Gap: 4,
		Padding:     properties.UniformInset(8),
		BorderRadii: properties.TopCorners(6),
		BoxShadow:   properties.BoxShadow{Blur: 4, Color: "black"},
		FontSize:    14,
	}.Px(unit.Metric{PxPerDp: 2})
	if st.Width != -1 {
		t.Errorf("Width -1 (expand) changed to %d", st.Width)
	}
	if st.Height != 20 || st.Gap != 8 || st.Padding.Left != 16 || st.BorderRadii.TopLeft != 12 || st.BoxShadow.Blur != 8 {
		t.Errorf("lengths not doubled: %+v", st)
	}
	if st.FontSize != 14 {
		t.Errorf("FontSize (sp) scaled to %v", st.FontSize)
	}
}

func TestMerge_Variadic(t *testing.T) {
	got := Merge(Styles{Color: "red", Gap: 1}, Styles{Gap: 2}, Styles{Color: "transparent"})
	if got.Gap != 2 || got.Color != "transparent" {
		t.Errorf("Merge = %+v, want later overrides to win, transparent included", got)
	}
	if got := (Styles{Gap: 1}).With(Styles{Width: 5}); got.Gap != 1 || got.Width != 5 {
		t.Errorf("With = %+v", got)
	}
}
