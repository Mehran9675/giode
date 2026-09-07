// Package checkbox provides a stateful checkbox with a label.
package checkbox

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
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
}

// New returns a Checkbox with the given label.
func New(label string) *Checkbox {
	return &Checkbox{label: label}
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

// Styles replaces the checkbox styles: Color is the label color,
// Background is the checked box fill, BorderColor is the unchecked
// border. It returns c for chaining.
func (c *Checkbox) Styles(st styles.Styles) *Checkbox {
	c.st = st
	return c
}

// Layout lays out and updates the checkbox.
func (c *Checkbox) Layout(gtx layout.Context) layout.Dimensions {
	if c.bool.Update(gtx) && c.onChange != nil {
		c.onChange(c.bool.Value)
	}
	st := styles.Merge(styles.Styles{Color: defaultText}, c.st)
	return c.bool.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if c.bool.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}
		return elements.Flex(styles.Styles{Align: properties.AlignCenter, Gap: 8},
			c.box(),
			elements.Text(c.label, styles.Styles{Color: st.Color, FontSize: st.FontSize}),
		).Layout(gtx)
	})
}
