package input

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Input is a stateful single-line text field. Create it once with New
// and lay it out every frame.
type Input struct {
	editor       widget.Editor
	placeholder  string
	st           styles.Styles
	onSubmit     func(text string)
	requestFocus bool
}

// New returns an Input with the given placeholder.
func New(placeholder string) *Input {
	i := &Input{placeholder: placeholder}
	i.editor.SingleLine = true
	i.editor.Submit = true
	return i
}

// Text returns the current content.
func (i *Input) Text() string {
	return i.editor.Text()
}

// SetText replaces the content.
func (i *Input) SetText(s string) {
	i.editor.SetText(s)
}

// Submit sets the handler invoked when Enter is pressed. It returns i
// for chaining.
func (i *Input) Submit(fn func(text string)) *Input {
	i.onSubmit = fn
	return i
}

// Styles replaces the field styles. It returns i for chaining.
func (i *Input) Styles(st styles.Styles) *Input {
	i.st = st
	return i
}

// Focus requests keyboard focus.
func (i *Input) Focus() {
	i.requestFocus = true
}

// Layout lays out and updates the field.
func (i *Input) Layout(gtx layout.Context) layout.Dimensions {
	if i.requestFocus {
		i.requestFocus = false
		gtx.Execute(key.FocusCmd{Tag: &i.editor})
	}
	for {
		ev, ok := i.editor.Update(gtx)
		if !ok {
			break
		}
		if e, ok := ev.(widget.SubmitEvent); ok && i.onSubmit != nil {
			i.onSubmit(e.Text)
		}
	}

	st := styles.Merge(styles.Styles{
		Padding:      properties.SymmetricInset(12, 8),
		BorderWidth:  1,
		BorderRadius: 6,
	}, i.st)
	if st.Width == 0 {
		st.Width = -1
	}

	return elements.Box(st, &editorEl{input: i, st: st}).Layout(gtx)
}
