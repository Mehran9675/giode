package elements

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Layout resolves the box size and paints the background stack
// (shadow, fill, gradient, image, border) with the content replayed on
// top.
func (b *BoxElement) Layout(gtx layout.Context) layout.Dimensions {
	if b.press != nil {
		// OnClick: the Pressable lays the box out over its hit area.
		b.press.cursor = b.st.Cursor
		b.press.content = func(PressState) Element { return boxBody{b} }
		return b.press.Layout(gtx)
	}
	return b.layoutBody(gtx)
}

// boxBody lays out a clickable box's own content inside its Pressable.
type boxBody struct{ b *BoxElement }

func (bb boxBody) Layout(gtx layout.Context) layout.Dimensions { return bb.b.layoutBody(gtx) }

func (b *BoxElement) layoutBody(gtx layout.Context) layout.Dimensions {
	st := b.resolveStyles(gtx)
	if st.Display == properties.DisplayNone {
		return layout.Dimensions{}
	}
	st = st.Px(gtx.Metric)

	// Margin carves space out of the constraints; the painted content
	// is then offset inside it.
	margin := st.Margin
	gtx.Constraints = insetConstraints(gtx.Constraints, margin)

	// Fixed and expanding sizes.
	if st.Width != 0 {
		gtx.Constraints.Min.X, gtx.Constraints.Max.X = st.Width.Constrain(gtx.Constraints.Min.X, gtx.Constraints.Max.X)
	}
	if st.Height != 0 {
		gtx.Constraints.Min.Y, gtx.Constraints.Max.Y = st.Height.Constrain(gtx.Constraints.Min.Y, gtx.Constraints.Max.Y)
	}

	// Clamps, percentage sizes (which take precedence) and aspect
	// ratio.
	applyMinMax(&gtx.Constraints, st)
	applyPct(&gtx.Constraints, layout.Horizontal, int(st.WidthPct))
	applyPct(&gtx.Constraints, layout.Vertical, int(st.HeightPct))
	applyAspect(&gtx.Constraints, st)

	flow, absolute := b.split()

	// Lay the content out into a macro so the background can be
	// painted underneath it once its size is known. Children's own
	// Layout calls happen inside this block, so the absolute-position
	// tracker (see absolute.go) must be pushed around it, not around
	// the later replay below - by the time that runs, every child has
	// already finished laying itself out.
	macro := op.Record(gtx.Ops)
	padding := st.Padding
	cgtx := gtx
	cgtx.Constraints = insetConstraints(gtx.Constraints, padding)
	contentOff := image.Pt(margin.Left+padding.Left, margin.Top+padding.Top)
	popAbs := PushOffset(contentOff)
	popColor := PushColor(st.Color)
	dims := b.layoutFlex(cgtx, st, flow)
	if len(absolute) > 0 {
		b.layoutAbsolute(cgtx, absolute, dims.Size)
	}
	popColor()
	popAbs()
	content := macro.Stop()

	size := dims.Size
	size.X += padding.Left + padding.Right
	size.Y += padding.Top + padding.Bottom
	outer := image.Pt(size.X+margin.Left+margin.Right, size.Y+margin.Top+margin.Bottom)

	// A hidden box keeps its space but paints nothing and takes no
	// input: the recorded content is simply dropped.
	if st.Visibility == properties.VisibilityHidden {
		return layout.Dimensions{Size: outer}
	}

	// Opacity 0 means unset (fully opaque); it applies to the
	// background and the content.
	opacity := float32(st.Opacity)
	if opacity == 0 {
		opacity = 1
	}
	if opacity < 1 {
		defer paint.PushOpacity(gtx.Ops, opacity).Pop()
	}

	// Background stack, offset into the margin area. The shadow
	// paints first so it stays behind the box.
	defer op.Offset(image.Pt(margin.Left, margin.Top)).Push(gtx.Ops).Pop()
	if !st.Transform.IsIdentity() {
		defer op.Affine(st.Transform.Affine(size)).Push(gtx.Ops).Pop()
	}
	corners := st.Corners()
	box := image.Rectangle{Max: size}
	properties.PaintBoxShadowCorners(gtx.Ops, box, corners, st.BoxShadow)
	properties.PaintBackgroundCorners(gtx.Ops, size, st.Background, corners)
	properties.PaintGradientCorners(gtx.Ops, size, st.BackgroundGradient, corners)
	if st.BackgroundImage != nil {
		cl := properties.RRect(box, corners).Push(gtx.Ops)
		properties.PaintBackgroundImage(gtx.Ops, box, st.BackgroundImage, st.BackgroundFit)
		cl.Pop()
	}

	properties.PaintBorderCorners(gtx.Ops, size, st.BorderWidth, st.BorderColor, corners)

	// Replay the content at its padded position. When overflow is
	// hidden it's clipped to the padding box - inside the border, with
	// the corner radii shrunk to match, like CSS - so children can't
	// cover the border or spill past rounded corners.
	var clipStack clip.Stack
	hidden := st.Overflow.Mode() == properties.OverflowHidden
	if hidden {
		bw := int(st.BorderWidth)
		shrink := func(v int) int { return max(0, v-bw) }
		inner := properties.Corners{TopLeft: shrink(corners.TopLeft), TopRight: shrink(corners.TopRight), BottomRight: shrink(corners.BottomRight), BottomLeft: shrink(corners.BottomLeft)}
		clipStack = properties.RRect(box.Inset(bw), inner).Push(gtx.Ops)
	}
	off := op.Offset(image.Pt(padding.Left, padding.Top)).Push(gtx.Ops)
	content.Add(gtx.Ops)
	off.Pop()
	if hidden {
		clipStack.Pop()
	}

	// Register the hover area over the painted box. A pass-through
	// region sees pointer events without blocking children above it.
	if b.hover != nil {
		ps := pointer.PassOp{}.Push(gtx.Ops)
		cl := clip.Rect(box).Push(gtx.Ops)
		b.hover.Add(gtx.Ops)
		if b.hovered && st.Cursor != "" {
			st.Cursor.CursorFor().Add(gtx.Ops)
		}
		cl.Pop()
		ps.Pop()
	}

	return layout.Dimensions{Size: outer}
}

// applyMinMax clamps the resolved constraints with the min/max size
// properties.
func applyMinMax(c *layout.Constraints, st styles.Styles) {
	if st.MinWidth > 0 && c.Min.X < int(st.MinWidth) {
		c.Min.X = int(st.MinWidth)
	}
	if st.MaxWidth > 0 && c.Max.X > int(st.MaxWidth) {
		c.Max.X = int(st.MaxWidth)
	}
	if st.MinHeight > 0 && c.Min.Y < int(st.MinHeight) {
		c.Min.Y = int(st.MinHeight)
	}
	if st.MaxHeight > 0 && c.Max.Y > int(st.MaxHeight) {
		c.Max.Y = int(st.MaxHeight)
	}
	if c.Max.X < c.Min.X {
		c.Max.X = c.Min.X
	}
	if c.Max.Y < c.Min.Y {
		c.Max.Y = c.Min.Y
	}
}

// applyPct sizes the axis to pct percent of the available space.
func applyPct(c *layout.Constraints, axis layout.Axis, pct int) {
	if pct <= 0 {
		return
	}
	if pct > 100 {
		pct = 100
	}
	if axis == layout.Horizontal {
		if c.Max.X >= unbounded {
			return
		}
		v := c.Max.X * pct / 100
		c.Min.X, c.Max.X = v, v
	} else {
		if c.Max.Y >= unbounded {
			return
		}
		v := c.Max.Y * pct / 100
		c.Min.Y, c.Max.Y = v, v
	}
}

// applyAspect derives the hugging axis from the resized axis when an
// aspect ratio is set.
func applyAspect(c *layout.Constraints, st styles.Styles) {
	if st.AspectRatio <= 0 {
		return
	}
	ratio := float32(st.AspectRatio)
	if st.Width != 0 && st.Height == 0 && c.Max.X > 0 && c.Max.X < unbounded {
		h := int(float32(c.Max.X) / ratio)
		c.Min.Y, c.Max.Y = h, h
	} else if st.Height != 0 && st.Width == 0 && c.Max.Y > 0 && c.Max.Y < unbounded {
		w := int(float32(c.Max.Y) * ratio)
		c.Min.X, c.Max.X = w, w
	}
}
