package elements

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Layout resolves the box size and paints the background stack
// (shadow, fill, image, border) with the content replayed on top.
func (b *boxEl) Layout(gtx layout.Context) layout.Dimensions {
	st := b.st

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
	// painted underneath it once its size is known.
	macro := op.Record(gtx.Ops)
	padding := st.Padding
	cgtx := gtx
	cgtx.Constraints = insetConstraints(gtx.Constraints, padding)
	dims := b.layoutFlex(cgtx, flow)
	if len(absolute) > 0 {
		b.layoutAbsolute(cgtx, absolute, dims.Size)
	}
	content := macro.Stop()

	size := dims.Size
	size.X += padding.Left + padding.Right
	size.Y += padding.Top + padding.Bottom

	// Opacity and Visibility both mean unset at 0 (fully visible); it
	// applies to the background and the content.
	opacity := float32(st.Opacity)
	if opacity == 0 {
		opacity = 1
	}
	visibility := float32(st.Visibility)
	if visibility == 0 {
		visibility = 1
	}
	opacity *= visibility
	if opacity < 1 {
		opStack := paint.PushOpacity(gtx.Ops, opacity)
		defer opStack.Pop()
	}

	// Background stack, offset into the margin area. The shadow
	// paints first so it stays behind the box.
	bgOff := op.Offset(image.Pt(margin.Left, margin.Top)).Push(gtx.Ops)
	properties.PaintBoxShadow(gtx.Ops, image.Rectangle{Max: size}, st.BorderRadius, st.BoxShadow)
	var clipStack clip.Stack
	if st.Overflow.Mode() == properties.OverflowHidden {
		clipStack = clip.Rect{Max: size}.Push(gtx.Ops)
	}
	properties.PaintBackground(gtx.Ops, size, st.Background, st.BorderRadius)
	properties.PaintBackgroundImage(gtx.Ops, image.Rectangle{Max: size}, st.BackgroundImage, st.BackgroundFit)
	properties.PaintBorder(gtx.Ops, size, st.BorderWidth, st.BorderColor, st.BorderRadius)
	if st.Overflow.Mode() == properties.OverflowHidden {
		clipStack.Pop()
	}
	bgOff.Pop()

	// Replay the content at its padded position.
	off := op.Offset(image.Pt(margin.Left+padding.Left, margin.Top+padding.Top)).Push(gtx.Ops)
	content.Add(gtx.Ops)
	off.Pop()

	return layout.Dimensions{
		Size: image.Pt(size.X+margin.Left+margin.Right, size.Y+margin.Top+margin.Bottom),
	}
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
