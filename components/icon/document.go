package icon

import (
	"encoding/xml"
	"fmt"
	"image"
	"strconv"
	"strings"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Document is a parsed SVG document: its viewBox and the fill path
// data from every <path> element found (including ones nested in
// <g> groups). Use LoadSVG/MustLoadSVG to parse one and Icon to
// render it - unlike SVG/SVGCustom, which take raw path data, LoadSVG
// takes a whole .svg file's bytes.
//
// Only fill geometry is supported: every path in the document paints
// with the same styles.Color, like SVG/SVGCustom. Stroke-only icon
// sets (fill="none" with a visible stroke) aren't rendered, since
// that needs outline-to-fill conversion this parser doesn't do; a
// per-path fill color from the source file also isn't preserved, to
// keep icons themeable the same way as every other icon in the
// package (via styles.Color).
type Document struct {
	viewX, viewY float32
	viewW, viewH float32
	paths        []string
}

// LoadSVG parses a full SVG document: the viewBox (or width/height,
// when viewBox is absent) and one fill path per <path d="..."> found,
// recursing into <g> groups. Other elements (rect, circle, text, ...)
// are ignored. It returns an error for malformed XML, or a document
// with neither a viewBox nor width and height attributes.
func LoadSVG(data []byte) (*Document, error) {
	var root xmlSVG
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("giode: parse svg: %w", err)
	}

	d := &Document{}
	if vb, ok := parseViewBox(root.ViewBox); ok {
		d.viewX, d.viewY, d.viewW, d.viewH = vb[0], vb[1], vb[2], vb[3]
	} else if w, h, ok := parseSize(root.Width, root.Height); ok {
		d.viewW, d.viewH = w, h
	} else {
		return nil, fmt.Errorf("giode: parse svg: no viewBox or width/height")
	}
	d.paths = collectPaths(root.Paths, root.Groups)
	return d, nil
}

// MustLoadSVG is LoadSVG, panicking on error - convenient for an
// icon's bytes known at compile time, e.g. via go:embed.
func MustLoadSVG(data []byte) *Document {
	d, err := LoadSVG(data)
	if err != nil {
		panic(err)
	}
	return d
}

// Icon renders the document as an icon element, like SVG/SVGCustom:
// size comes from styles.Width/Height (default 24px) and color from
// styles.Color (default black).
func (d *Document) Icon(st ...styles.Styles) elements.Element {
	return &svgDocEl{doc: d, st: first(st)}
}

type xmlSVG struct {
	ViewBox string     `xml:"viewBox,attr"`
	Width   string     `xml:"width,attr"`
	Height  string     `xml:"height,attr"`
	Paths   []xmlPath  `xml:"path"`
	Groups  []xmlGroup `xml:"g"`
}

type xmlGroup struct {
	Paths  []xmlPath  `xml:"path"`
	Groups []xmlGroup `xml:"g"`
}

type xmlPath struct {
	D string `xml:"d,attr"`
}

// collectPaths gathers every path's d attribute, recursing into
// groups, skipping empty ones.
func collectPaths(paths []xmlPath, groups []xmlGroup) []string {
	var out []string
	for _, p := range paths {
		if p.D != "" {
			out = append(out, p.D)
		}
	}
	for _, g := range groups {
		out = append(out, collectPaths(g.Paths, g.Groups)...)
	}
	return out
}

// parseViewBox parses "minX minY width height", numbers separated by
// whitespace and/or commas.
func parseViewBox(s string) (v [4]float32, ok bool) {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == ',' || r == '\t' || r == '\n'
	})
	if len(fields) != 4 {
		return v, false
	}
	for i, f := range fields {
		n, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return v, false
		}
		v[i] = float32(n)
	}
	return v, true
}

// parseSize parses SVG width/height attributes, tolerating a trailing
// unit like "24px".
func parseSize(w, h string) (width, height float32, ok bool) {
	pw, okW := parseLength(w)
	ph, okH := parseLength(h)
	return pw, ph, okW && okH
}

func parseLength(s string) (float32, bool) {
	s = strings.TrimSpace(s)
	for _, suffix := range []string{"px", "pt"} {
		s = strings.TrimSuffix(s, suffix)
	}
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, false
	}
	return float32(n), true
}

// svgDocEl renders every path of a parsed Document, scaled from its
// viewBox, filled with a single color - the multi-path equivalent of
// svgEl.
type svgDocEl struct {
	doc *Document
	st  styles.Styles
}

func (s *svgDocEl) Layout(gtx layout.Context) layout.Dimensions {
	sizeX := int(s.st.Width)
	sizeY := int(s.st.Height)
	if sizeX <= 0 {
		sizeX = 24
	}
	if sizeY <= 0 {
		sizeY = 24
	}
	size := gtx.Constraints.Constrain(image.Pt(sizeX, sizeY))
	if size.X <= 0 || size.Y <= 0 {
		return layout.Dimensions{}
	}

	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	if s.doc.viewW > 0 && s.doc.viewH > 0 {
		scale := f32.AffineId().
			Offset(f32.Pt(-s.doc.viewX, -s.doc.viewY)).
			Scale(f32.Point{}, f32.Pt(float32(size.X)/s.doc.viewW, float32(size.Y)/s.doc.viewH))
		defer op.Affine(scale).Push(gtx.Ops).Pop()
	}

	var p clip.Path
	p.Begin(gtx.Ops)
	for _, d := range s.doc.paths {
		parsePath(&p, d)
	}
	paint.FillShape(gtx.Ops, properties.ResolveColor(s.st.Color), clip.Outline{Path: p.End()}.Op())
	return layout.Dimensions{Size: size}
}
