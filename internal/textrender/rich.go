package textrender

// The rich text layout follows gioui.org/x/styledtext (Unlicense OR
// MIT), collecting the shaped segments instead of painting them, so a
// layout can be cached and drawn every frame without shaping.

import (
	"hash/maphash"
	"image"
	"image/color"
	"math"
	"sync"
	"unicode/utf8"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"golang.org/x/image/math/fixed"
)

// Span is a run of text with its own style.
type Span struct {
	Font    font.Font
	Size    unit.Sp
	Content string
}

// RichParams are the paragraph-wide settings of rich text.
type RichParams struct {
	Alignment       text.Alignment
	LineHeight      unit.Sp
	LineHeightScale float32
}

// Segment is a piece of one span on one line.
type Segment struct {
	// Span indexes the spans the text was laid out from.
	Span int
	// Offset is the segment's top left corner in the text; Size its
	// size; Ascent its baseline from its top.
	Offset image.Point
	Size   image.Point
	Ascent int
	// Glyphs are the segment's glyphs, relative to Offset.
	Glyphs []text.Glyph
	// Runes are the span's runes the segment covers, from its start.
	RuneStart, RuneEnd int
}

// RichLaid is laid out rich text.
type RichLaid struct {
	Segments []Segment
	Dims     layout.Dimensions
	keys     []uint64
	used     uint64
}

var richLayouts struct {
	sync.Mutex
	m map[uint64]*richEntry
}

type richEntry struct {
	spans []Span
	*RichLaid
}

// LayoutRich lays out spans as styledtext.TextStyle would, from cache
// when the spans, settings and constraints are the same as before.
func LayoutRich(sh *text.Shaper, gtx layout.Context, p RichParams, spans []Span) *RichLaid {
	key := richKey(gtx, p, spans)
	now := frame.Load()
	richLayouts.Lock()
	if e, ok := richLayouts.m[key]; ok && sameSpans(e.spans, spans) {
		e.used = now
		richLayouts.Unlock()
		return e.RichLaid
	}
	richLayouts.Unlock()
	l := layoutRich(sh, gtx, p, spans)
	l.used = now
	richLayouts.Lock()
	if richLayouts.m == nil {
		richLayouts.m = make(map[uint64]*richEntry)
	}
	richLayouts.m[key] = &richEntry{spans: append([]Span(nil), spans...), RichLaid: l}
	richLayouts.Unlock()
	return l
}

func richKey(gtx layout.Context, p RichParams, spans []Span) uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	var b [8]byte
	w := func(v uint64) { putU64(b[:], v); h.Write(b[:]) }
	w(uint64(gtx.Constraints.Max.X))
	w(uint64(gtx.Constraints.Max.Y))
	w(uint64(math.Float32bits(gtx.Metric.PxPerSp)))
	w(uint64(p.Alignment))
	w(uint64(math.Float32bits(float32(p.LineHeight))))
	w(uint64(math.Float32bits(p.LineHeightScale)))
	h.WriteString(gtx.Locale.Language)
	for _, s := range spans {
		h.WriteString(string(s.Font.Typeface))
		w(uint64(s.Font.Weight)<<16 | uint64(s.Font.Style))
		w(uint64(math.Float32bits(float32(s.Size))))
		w(uint64(len(s.Content)))
		h.WriteString(s.Content)
	}
	return h.Sum64()
}

func sameSpans(a, b []Span) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Draw paints the segments, each in the color colorOf returns for its
// span.
func (l *RichLaid) Draw(ops *op.Ops, sh *text.Shaper, colorOf func(span int) color.NRGBA) {
	if len(l.keys) != len(l.Segments) {
		l.keys = make([]uint64, len(l.Segments))
		for i, s := range l.Segments {
			l.keys[i] = runHash(s.Glyphs)
		}
	}
	for i, s := range l.Segments {
		st := op.Offset(s.Offset).Push(ops)
		drawRunKeyed(ops, sh, s.Glyphs, colorOf(s.Span), l.keys[i])
		st.Pop()
	}
}

type spanState struct {
	Span
	idx       int
	runeStart int
}

type spanResult struct {
	glyphs           []text.Glyph
	width, height    int
	ascent           int
	runes            int
	multiLine        bool
	endedWithNewline bool
}

func shapeSpan(sh *text.Shaper, gtx layout.Context, p RichParams, maxWidth int, span spanState, truncate bool) (spanResult, richIterator) {
	maxLines := 0
	if truncate {
		maxLines = 1
	}
	sh.LayoutString(text.Parameters{
		Font:            span.Font,
		PxPerEm:         fixed.I(gtx.Sp(span.Size)),
		MaxLines:        maxLines,
		MaxWidth:        maxWidth,
		Truncator:       "​",
		Locale:          gtx.Locale,
		WrapPolicy:      text.WrapWords,
		LineHeight:      fixed.I(gtx.Sp(p.LineHeight)),
		LineHeightScale: p.LineHeightScale,
	}, span.Content)
	ti := richIterator{viewport: image.Rectangle{Max: gtx.Constraints.Max}, maxLines: 1}
	var res spanResult
	for g, ok := sh.NextGlyph(); ok; g, ok = sh.NextGlyph() {
		visibleOrBefore := ti.processGlyph(g)
		if ti.visible {
			if !ti.init {
				ti.firstX = g.X
				ti.init = true
			}
			g.X -= ti.firstX
			res.glyphs = append(res.glyphs, g)
		}
		if !visibleOrBefore {
			break
		}
	}
	return res, ti
}

func layoutSpan(sh *text.Shaper, gtx layout.Context, p RichParams, maxWidth int, span spanState) spanResult {
	res, ti := shapeSpan(sh, gtx, p, maxWidth, span, true)
	runes := ti.runes
	multiLine := runes < utf8.RuneCountInString(span.Content)
	endedWithNewline := ti.hasNewline
	if multiLine {
		i := 0
		for n := 0; n < runes; n++ {
			_, sz := utf8.DecodeRuneInString(span.Content[i:])
			i += sz
		}
		first, _ := utf8.DecodeRuneInString(span.Content[i:])
		if first == '\n' {
			endedWithNewline = true
			runes++
		} else if runes == 0 {
			// Nothing fit on the line and it wasn't a hard newline: wrap
			// without truncation to find the word that doesn't fit.
			res, ti = shapeSpan(sh, gtx, p, maxWidth, span, false)
			runes = ti.runes
			multiLine = runes < utf8.RuneCountInString(span.Content)
			endedWithNewline = ti.hasNewline
		}
	}
	res.width, res.height = ti.bounds.Dx(), ti.bounds.Dy()
	res.ascent = ti.baseline
	res.runes = runes
	res.multiLine = multiLine
	res.endedWithNewline = endedWithNewline
	return res
}

func layoutRich(sh *text.Shaper, gtx layout.Context, p RichParams, in []Span) *RichLaid {
	spans := make([]spanState, len(in))
	for i, s := range in {
		spans[i] = spanState{Span: s, idx: i}
	}
	scale := p.LineHeightScale
	lineHeightPx := gtx.Sp(p.LineHeight)
	if scale == 0 {
		scale = 1.2
	}
	var (
		lineDims image.Point
		overall  image.Point
		line     []Segment
		out      = &RichLaid{}
	)
	for i := 0; i < len(spans); i++ {
		span := spans[i]
		maxWidth := gtx.Constraints.Max.X - lineDims.X
		res := layoutSpan(sh, gtx, p, maxWidth, span)
		forceToNextLine := lineDims.X > 0 && res.width > maxWidth
		if !forceToNextLine {
			line = append(line, Segment{
				Span:      span.idx,
				Offset:    image.Point{X: lineDims.X},
				Size:      image.Point{X: res.width, Y: res.height},
				Ascent:    res.ascent,
				Glyphs:    res.glyphs,
				RuneStart: span.runeStart,
				RuneEnd:   span.runeStart + res.runes,
			})
			lineDims.X += res.width
			lineDims.Y = max(lineDims.Y, res.height)
			overall.X = max(overall.X, lineDims.X)
		}
		if res.multiLine || res.endedWithNewline || i == len(spans)-1 || forceToNextLine {
			if len(line) > 0 {
				last := line[len(line)-1]
				lineWidth := last.Offset.X + last.Size.X
				pad := 0
				if lineWidth < gtx.Constraints.Max.X {
					switch p.Alignment {
					case text.Middle:
						pad = (gtx.Constraints.Max.X - lineWidth) / 2
					case text.End:
						pad = gtx.Constraints.Max.X - lineWidth
					}
				}
				for _, s := range line {
					s.Offset.X += pad
					s.Offset.Y = overall.Y
					out.Segments = append(out.Segments, s)
				}
			}
			line = line[:0]
			h := lineDims.Y
			if p.LineHeight != 0 {
				h = lineHeightPx
			}
			overall.Y += int(float32(h) * scale)
			lineDims = image.Point{}
		}
		if res.multiLine && !forceToNextLine {
			spans = append(spans, spanState{})
			copy(spans[i+2:], spans[i+1:])
			byteLen := 0
			for n := 0; n < res.runes; n++ {
				_, sz := utf8.DecodeRuneInString(span.Content[byteLen:])
				byteLen += sz
			}
			span.Content = span.Content[byteLen:]
			span.runeStart += res.runes
			spans[i+1] = span
		} else if forceToNextLine {
			i--
		}
	}
	out.Dims = layout.Dimensions{Size: gtx.Constraints.Constrain(overall)}
	return out
}

// richIterator follows styledtext's textIterator.
type richIterator struct {
	viewport   image.Rectangle
	maxLines   int
	linesSeen  int
	init       bool
	firstX     fixed.Int26_6
	hasNewline bool
	padding    image.Rectangle
	bounds     image.Rectangle
	runes      int
	visible    bool
	first      bool
	baseline   int
}

func (it *richIterator) processGlyph(g text.Glyph) bool {
	logical := image.Rectangle{
		Min: image.Pt(g.X.Floor(), int(g.Y)-g.Ascent.Ceil()),
		Max: image.Pt((g.X + g.Advance).Ceil(), int(g.Y)+g.Descent.Ceil()),
	}
	if g.Flags&text.FlagTruncator != 0 {
		if it.runes == 0 {
			it.hasNewline = true
		}
		it.bounds.Min.Y = min(it.bounds.Min.Y, logical.Min.Y)
		it.bounds.Max.Y = max(it.bounds.Max.Y, logical.Max.Y)
		it.visible = false
		return false
	}
	it.runes += int(g.Runes)
	it.hasNewline = it.hasNewline || (g.Flags&text.FlagLineBreak > 0 && g.Flags&text.FlagParagraphBreak > 0)
	if it.maxLines > 0 {
		if g.Flags&text.FlagLineBreak != 0 {
			it.linesSeen++
		}
		if it.linesSeen == it.maxLines && g.Flags&text.FlagParagraphBreak != 0 {
			it.visible = false
			return false
		}
	}
	if d := g.Bounds.Min.X.Floor(); d < it.padding.Min.X {
		it.padding.Min.X = d
	}
	if d := (g.Bounds.Max.X - g.Advance).Ceil(); d > it.padding.Max.X {
		it.padding.Max.X = d
	}
	if !it.first {
		it.first = true
		it.baseline = int(g.Y)
		it.bounds = logical
	}
	above := logical.Max.Y < it.viewport.Min.Y
	below := logical.Min.Y > it.viewport.Max.Y
	left := logical.Max.X < it.viewport.Min.X
	right := logical.Min.X > it.viewport.Max.X
	it.visible = !above && !below && !left && !right
	if it.visible {
		it.bounds = grow(it.bounds, logical)
	}
	return !below
}
