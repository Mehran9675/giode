package main

import (
	"fmt"

	"github.com/mehran9675/giode"
)

// buttonPage holds the button and its counter, created once so the
// click handler and count persist across frames.
type buttonPage struct {
	count int
	btn   giode.Element
}

func NewButtonPage() *buttonPage {
	p := &buttonPage{}
	p.btn = giode.Button("Button", giode.Styles{
		Width: 50, Background: "gray", BorderRadius: 10,
		BoxShadow: giode.BoxShadow{X: 1, Y: 1, Blur: 4, Color: "black"},
	}).OnClick(func() { p.count++ })
	return p
}

func (p *buttonPage) View() giode.Element {
	return giode.Stack(
		boxStyle,
		giode.Text(fmt.Sprintf("Count: %d", p.count), textStyle),
		p.btn,
	)
}
