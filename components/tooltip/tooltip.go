// Package tooltip provides a stateful hover tooltip that wraps any
// element as its trigger.
package tooltip

import (
	"image"

	"gioui.org/layout"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

var (
	defaultBackground = "#0f172a"
	defaultText       = "#f8fafc"
)

// Tooltip is a stateful hover tooltip wrapping trigger: hovering
// trigger shows content in a small panel anchored just below it,
// rendered above everything else via elements.RegisterPortal. Create
// it once with New and lay it out every frame in trigger's place.
type Tooltip struct {
	trigger elements.Element
	content string
	st      styles.Styles
	hover   kit.HoverTracker
}

// New wraps trigger with a tooltip showing content on hover. The
// styles argument is optional: Background and Color style the panel,
// FontSize its text.
func New(trigger elements.Element, content string, st ...styles.Styles) *Tooltip {
	t := &Tooltip{trigger: trigger, content: content}
	if len(st) > 0 {
		t.st = st[0]
	}
	t.hover.Enable()
	return t
}

// SetContent replaces the tooltip text.
func (t *Tooltip) SetContent(content string) {
	t.content = content
}

// Hovered reports whether the trigger was hovered on the last frame.
func (t *Tooltip) Hovered() bool {
	return t.hover.Hovered()
}

// Layout lays out the trigger and, while hovered, queues the tooltip
// panel as a portal anchored just below it.
func (t *Tooltip) Layout(gtx layout.Context) layout.Dimensions {
	t.hover.Update(gtx)
	dims := t.trigger.Layout(gtx)
	t.hover.RegisterArea(gtx, image.Point{}, dims.Size, "")

	if t.hover.Hovered() {
		anchor := elements.AbsolutePosition().Add(image.Pt(0, dims.Size.Y+6))
		elements.RegisterPortal(gtx, anchor, func() elements.Element {
			return elements.Raw(t.panel)
		})
	}
	return dims
}

// panel renders the tooltip content, hugging it.
func (t *Tooltip) panel(gtx layout.Context) layout.Dimensions {
	st := styles.Merge(styles.Styles{
		Background:   defaultBackground,
		Color:        defaultText,
		Padding:      properties.SymmetricInset(10, 6),
		BorderRadius: 6,
		FontSize:     13,
	}, t.st)
	gtx.Constraints.Min = image.Point{}
	return elements.Box(st, elements.Text(t.content, styles.Styles{
		Color: st.Color, FontSize: st.FontSize,
	})).Layout(gtx)
}
