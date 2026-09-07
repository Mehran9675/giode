// Package fonts holds the text shaper shared by all Giode elements.
//
// Register a font collection once before running the app, or call
// UseGofont to embed the Go fonts. When no fonts are registered the
// shaper still falls back to system fonts where the platform supports
// them.
package fonts

import (
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/text"
)

var shaper *text.Shaper

// Register replaces the font collection used to shape text.
func Register(faces []font.FontFace) {
	shaper = text.NewShaper(text.WithCollection(faces))
}

// UseGofont embeds the Go fonts for portable text rendering.
func UseGofont() {
	Register(gofont.Collection())
}

// Shaper returns the text shaper, or nil if no fonts were registered.
func Shaper() *text.Shaper {
	return shaper
}
