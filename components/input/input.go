package input

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Input is a stateful text field, single-line by default. Create it
// once with New and lay it out every frame.
type Input struct {
	editor       widget.Editor
	placeholder  string
	label        string
	st           styles.Styles
	onSubmit     func(text string)
	onChange     func(text string)
	requestFocus bool
	focused      bool
}

// New returns an Input with the given placeholder. The styles
// argument is optional.
func New(placeholder string, st ...styles.Styles) *Input {
	i := &Input{placeholder: placeholder}
	if len(st) > 0 {
		i.st = st[0]
	}
	i.editor.SingleLine = true
	i.editor.Submit = true
	return i
}

// Text returns the current content.
func (i *Input) Text() string {
	return i.editor.Text()
}

// SetText replaces the content without invoking OnChange.
func (i *Input) SetText(s string) {
	i.editor.SetText(s)
}

// Submit sets the handler invoked when Enter is pressed in a
// single-line field. It returns i for chaining.
func (i *Input) Submit(fn func(text string)) *Input {
	i.onSubmit = fn
	return i
}

// OnChange sets the handler invoked after every edit. It returns i for
// chaining.
func (i *Input) OnChange(fn func(text string)) *Input {
	i.onChange = fn
	return i
}

// Password masks the content with bullets. It returns i for chaining.
func (i *Input) Password(on bool) *Input {
	if on {
		i.editor.Mask = '•'
	} else {
		i.editor.Mask = 0
	}
	return i
}

// Multiline lets the field hold several lines; Enter inserts a newline
// instead of submitting. It returns i for chaining.
func (i *Input) Multiline(on bool) *Input {
	i.editor.SingleLine = !on
	i.editor.Submit = !on
	return i
}

// ReadOnly makes the content selectable but not editable. It returns i
// for chaining.
func (i *Input) ReadOnly(on bool) *Input {
	i.editor.ReadOnly = on
	return i
}

// MaxLength limits the content to n runes; 0 removes the limit. It
// returns i for chaining.
func (i *Input) MaxLength(n int) *Input {
	i.editor.MaxLen = n
	return i
}

// Label sets the accessible name reported to automation (the semantic
// tree, the MCP server); the placeholder is used when unset. It
// returns i for chaining.
func (i *Input) Label(label string) *Input {
	i.label = label
	return i
}

// ChangeText replaces the content and invokes OnChange, as typing
// would.
func (i *Input) ChangeText(s string) {
	i.editor.SetText(s)
	if i.onChange != nil {
		i.onChange(s)
	}
}

// SubmitText invokes the Submit handler with the current content, as
// pressing Enter would.
func (i *Input) SubmitText() {
	if i.onSubmit != nil {
		i.onSubmit(i.editor.Text())
	}
}

// Focus requests keyboard focus.
func (i *Input) Focus() {
	i.requestFocus = true
}

// Focused reports whether the field had keyboard focus in the last
// frame.
func (i *Input) Focused() bool {
	return i.focused
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
		switch e := ev.(type) {
		case widget.SubmitEvent:
			if i.onSubmit != nil {
				i.onSubmit(e.Text)
				kit.RequestFrame(gtx)
			}
		case widget.ChangeEvent:
			if i.onChange != nil {
				i.onChange(i.editor.Text())
				kit.RequestFrame(gtx)
			}
		}
	}
	i.focused = gtx.Source.Focused(&i.editor)
	if i.focused {
		// A focused field receives key events and blinks its caret, so
		// cached ancestors must keep laying it out.
		kit.MarkLive()
	}

	st := styles.Merge(styles.Styles{
		Padding:      properties.SymmetricInset(12, 8),
		BorderWidth:  1,
		BorderRadius: 6,
		BorderColor:  "theme.border",
		Color:        "theme.text",
	}, i.st)
	if st.Width == 0 {
		st.Width = -1
	}

	dims := elements.Box(st, &editorEl{input: i, st: st}).Layout(gtx)
	label := i.label
	if label == "" {
		label = i.placeholder
	}
	value := i.editor.Text()
	if i.editor.Mask != 0 {
		value = ""
	}
	elements.RecordSemantic(gtx, elements.SemanticNode{Kind: elements.SemanticInput, Label: label, Value: value, Target: i}, dims.Size)
	return dims
}
