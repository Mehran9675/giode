// Package checkbox provides a stateful checkbox with an optional
// label.
package checkbox

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Checkbox is a stateful checkbox with an optional label. Create it
// once with New and lay it out every frame.
type Checkbox struct {
	label    string
	bool     widget.Bool
	onChange func(checked bool)
	st       styles.Styles
	iconName string
	iconSt   styles.Styles
}

// New returns a Checkbox. label is optional: pass "" for a bare
// checkbox with no label. The styles argument is optional: Color is
// the label color, Background is the checked fill, BorderColor is
// the unchecked border, BorderWidth and BorderRadius the border, and
// Width/Height size the box (default 18x18).
func New(label string, st ...styles.Styles) *Checkbox {
	c := &Checkbox{label: label, iconName: "check"}
	if len(st) > 0 {
		c.st = st[0]
	}
	return c
}

// Icon replaces the material icon shown when checked (default
// "check"). The styles argument is optional and overrides the icon's
// Color (default white) and size (defaults to the box size).
func (c *Checkbox) Icon(name string, st ...styles.Styles) *Checkbox {
	c.iconName = name
	if len(st) > 0 {
		c.iconSt = st[0]
	}
	return c
}

// Checked reports the current state.
func (c *Checkbox) Checked() bool {
	return c.bool.Value
}

// SetChecked sets the state without invoking OnChange.
func (c *Checkbox) SetChecked(v bool) {
	c.bool.Value = v
}

// OnChange sets the handler invoked when the state changes. It
// returns c for chaining.
func (c *Checkbox) OnChange(fn func(checked bool)) *Checkbox {
	c.onChange = fn
	return c
}

// Layout lays out and updates the checkbox.
func (c *Checkbox) Layout(gtx layout.Context) layout.Dimensions {
	if c.bool.Update(gtx) && c.onChange != nil {
		c.onChange(c.bool.Value)
		gtx.Execute(op.InvalidateCmd{})
	}
	st := styles.Merge(styles.Styles{Color: defaultText}, c.st)
	return c.bool.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if c.bool.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}
		box := c.box()
		if c.label == "" {
			return box.Layout(gtx)
		}
		return elements.Row(styles.Styles{Align: properties.AlignCenter, Gap: 8},
			box,
			elements.Text(c.label, styles.Styles{Color: st.Color, FontSize: st.FontSize}),
		).Layout(gtx)
	})
}
