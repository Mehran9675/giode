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

var (
	shaper     *text.Shaper
	collection []font.FontFace
)

// Register replaces the font collection used to shape text.
func Register(faces []font.FontFace) {
	collection = faces
	shaper = text.NewShaper(text.WithCollection(faces))
}

// Collection returns the registered font collection, in the order the
// shaper loaded it.
func Collection() []font.FontFace {
	return collection
}

// UseGofont embeds the Go fonts for portable text rendering.
func UseGofont() {
	Register(gofont.Collection())
}

// Shaper returns the text shaper, or nil if no fonts were registered.
func Shaper() *text.Shaper {
	return shaper
}
