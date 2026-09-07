// Package icon renders icons in the popular formats: SVG path data,
// the Material Design icon set, IconVG (Gio's native vector format)
// and plain images.
//
// All constructors are stateless and safe to call every frame. Size
// comes from styles.Width and styles.Height (default 24px) and color
// from styles.Color (default black).
package icon

import (
	"strconv"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

// Material returns a named icon from the built-in Material Design
// set. It panics if the name is unknown.
func Material(name string, st ...styles.Styles) elements.Element {
	pathData, ok := materialIcons[name]
	if !ok {
		panic("giode: unknown material icon " + strconv.Quote(name))
	}
	return svg(pathData, 24, 24, first(st))
}

// Names returns the names of all built-in material icons.
func Names() []string {
	names := make([]string, 0, len(materialIcons))
	for n := range materialIcons {
		names = append(names, n)
	}
	return names
}

func first(st []styles.Styles) styles.Styles {
	if len(st) > 0 {
		return st[0]
	}
	return styles.Styles{}
}
