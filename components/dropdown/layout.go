package dropdown

import (
	"image"

	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Layout lays out the field and, when open, the list just below it,
// drawn above everything else (see elements.RegisterPopup). Clicking
// anywhere outside the open list closes it.
//
// Keyboard: with the field focused (Tab), Enter, Space or the arrow
// keys open the list; in the open list Up and Down move the
// highlight, Enter or Space picks it and Escape closes the list.
func (s *Select) Layout(gtx layout.Context) layout.Dimensions {
	st := styles.Merge(styles.Styles{
		Background:   properties.CalcColorReverse(defaultBackground),
		Color:        properties.CalcColorReverse(defaultText),
		BorderWidth:  1,
		BorderColor:  properties.CalcColorReverse(defaultBorder),
		BorderRadius: 6,
		Padding:      properties.SymmetricInset(12, 8),
	}, s.st)

	if s.field.Clicked(gtx) {
		s.setOpen(!s.open)
		kit.RequestFrame(gtx)
	}
	s.handleKeys(gtx)
	switch s.focus {
	case focusList:
		gtx.Execute(key.FocusCmd{Tag: s})
	case focusField:
		gtx.Execute(key.FocusCmd{Tag: &s.field})
	}
	s.focus = focusNone
	if s.open {
		// The open panel is a portal; keep enclosing caches live so
		// it's registered every frame.
		kit.MarkLive()
	}

	// Field.
	fieldDims := s.field.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if s.st.Width > 0 {
			w := gtx.Dp(unit.Dp(s.st.Width))
			gtx.Constraints.Min.X, gtx.Constraints.Max.X = w, w
		}
		fieldSt := st
		if s.field.Hovered() {
			fieldSt.BorderColor = kit.Hovered(st.BorderColor)
			pointer.CursorPointer.Add(gtx.Ops)
		}
		// With a fixed width the label fills the space left of the
		// chevron, pushing it to the right edge. The row must not wrap:
		// FlexGrow is ignored on wrapped lines.
		rowSt := styles.Styles{Align: properties.AlignCenter, Wrap: properties.WrapNoWrap}
		labelSt := styles.Styles{FlexGrow: 1}
		if s.st.Width > 0 {
			rowSt.Width = -1
			labelSt.Width = -1
		}
		row := elements.Row(rowSt,
			elements.Box(labelSt, elements.Text(s.Value(), styles.Styles{Color: st.Color, FontSize: st.FontSize})),
			elements.Box(styles.Styles{Width: 4}),
			elements.Glyph(kit.GlyphChevronDown, styles.Styles{Color: st.Color, Width: 16}),
		)
		return elements.Box(fieldSt, row).Layout(gtx)
	})
	// Keep the field clickable even inside a frameless window's
	// whole-window drag area - see button.Layout's identical comment.
	kit.BlockDrag(gtx, fieldDims.Size)
	// The open list takes keyboard focus under this tag. It passes
	// pointer events through so the field below stays clickable.
	ps := pointer.PassOp{}.Push(gtx.Ops)
	cl := clip.Rect{Max: fieldDims.Size}.Push(gtx.Ops)
	event.Op(gtx.Ops, s)
	cl.Pop()
	ps.Pop()
	elements.RecordSemantic(gtx, elements.SemanticNode{Kind: elements.SemanticSelect, Label: s.Value(), Value: s.Value(), Target: s}, fieldDims.Size)

	if s.open && len(s.options) > 0 {
		fieldWidth := fieldDims.Size.X
		anchor := elements.AbsolutePosition().Add(image.Pt(0, fieldDims.Size.Y+4))
		elements.RegisterPopup(gtx, anchor, s, func() elements.Element {
			return elements.Raw(func(gtx layout.Context) layout.Dimensions {
				return s.dropdown(gtx, st, fieldWidth)
			})
		}, elements.PopupOptions{OnDismiss: s.Close})
	}
	return fieldDims
}

// handleKeys processes keyboard input for the field and the open list.
func (s *Select) handleKeys(gtx layout.Context) {
	names := []key.Name{key.NameUpArrow, key.NameDownArrow, key.NameReturn, key.NameEnter, key.NameSpace, key.NameEscape}
	filters := []event.Filter{key.FocusFilter{Target: s}}
	for _, n := range names {
		filters = append(filters, key.Filter{Focus: s, Name: n})
	}
	// On the focused field, the arrow keys open the list too (Enter and
	// Space already click it).
	filters = append(filters,
		key.Filter{Focus: &s.field, Name: key.NameUpArrow},
		key.Filter{Focus: &s.field, Name: key.NameDownArrow})
	for {
		ev, ok := gtx.Event(filters...)
		if !ok {
			return
		}
		e, ok := ev.(key.Event)
		if !ok || e.State != key.Press || len(s.options) == 0 {
			continue
		}
		kit.RequestFrame(gtx)
		if !s.open {
			if e.Name != key.NameEscape {
				s.setOpen(true)
			}
			continue
		}
		switch e.Name {
		case key.NameUpArrow:
			s.active = max(0, s.active-1)
		case key.NameDownArrow:
			s.active = min(len(s.options)-1, s.active+1)
		case key.NameReturn, key.NameEnter, key.NameSpace:
			s.Choose(s.active)
		case key.NameEscape:
			s.Close()
		}
	}
}
