# Style properties

Every element takes a `giode.Styles` value. Fields left at their zero value are unset, and the
element uses its default.

```go
card := giode.Styles{
	FlexDirection: giode.FlexDirectionColumn,
	Gap:           12,
	Padding:       giode.UniformInset(16),
	Background:    "theme.surface",
	BorderRadius:  12,
	BoxShadow:     giode.BoxShadow{Y: 4, Blur: 12, Color: "#00000060"},
}
```

Combine styles with `giode.Merge(base, overrides...)` or `base.With(override)`: later values
replace earlier ones where they're set.

## Units

Lengths (sizes, padding, margins, gaps, radii, border widths, offsets, shadows) are in
**dp**, density-independent pixels. On a standard screen 1dp is one pixel; on a 2× HiDPI
screen it's two, so the UI keeps its physical size. Font sizes are in sp, which scale the
same way. In your own elements, `st.Px(gtx.Metric)` converts a `Styles` value to pixels.

## Colors

A color is a string in any CSS syntax:

- hex: `"#fff"`, `"#1e293b"`, `"#00000060"` (with alpha)
- functions: `"rgb(59 130 246)"`, `"rgba(59, 130, 246, 0.5)"`, `"hsl(217 91% 60%)"`
- names: any of the 148 CSS colors, such as `"tomato"`, plus `"transparent"`
- theme tokens: `"theme.primary"`, `"theme.surface"`, ... (see below)

An empty string means unset. A string that isn't a valid color paints white, so the
mistake is visible; check user-provided colors with `giode.ParseColor`.

**Inheritance.** A box's `Color` carries down to the text, icons and glyphs inside it that
don't set their own, like CSS. `"currentColor"` refers to it explicitly.

### Theme

Theme tokens resolve against the active palette, so switching modes restyles everything
that uses them:

| Token | Use |
| --- | --- |
| `theme.background` | Window background |
| `theme.surface` | Cards, panels, inputs |
| `theme.text`, `theme.muted` | Text, secondary text |
| `theme.border` | Borders and dividers |
| `theme.primary`, `theme.onprimary` | Accent color, and text drawn on it |
| `theme.secondary`, `theme.accent` | Other accents |
| `theme.success`, `theme.warning`, `theme.danger`, `theme.info` | Status colors |

Add an alpha with a slash: `"theme.primary/0.5"` or `"theme.primary/50%"`.

```go
giode.SetMode(giode.ModeDark) // or giode.ToggleMode()
giode.SetPalette(giode.ModeDark, giode.Palette{Primary: "#22c55e"}) // unset fields keep the defaults
```

Both are safe to call from any goroutine and redraw the window.

## Layout

| Property | Values | Notes |
| --- | --- | --- |
| `FlexDirection` | `FlexDirectionRow` (default), `FlexDirectionColumn` | Direction children are laid out in. |
| `Align` | `AlignStart` (default), `AlignCenter`, `AlignEnd` | Placement across that direction. With the default, children stretch to fill a fixed-size container, like CSS. |
| `Justify` | `JustifyStart` (default), `JustifyCenter`, `JustifyEnd`, `JustifySpaceBetween`, `JustifySpaceAround`, `JustifySpaceEvenly` | Placement along that direction. |
| `Wrap` | `WrapWrap`, `WrapNoWrap` | Whether children flow onto new lines. Rows wrap by default, stacks don't. |
| `Gap` | dp | Space between children, and between wrapped lines. |
| `FlexGrow` | ≥ 0 | Share of the parent's spare space this element takes. Ignored in wrapping rows. |
| `Width`, `Height` | dp; `0` fits the content, `-1` fills | |
| `WidthPct`, `HeightPct` | 1–100 | Percentage of the available space. Wins over `Width`/`Height`. |
| `MinWidth`, `MaxWidth`, `MinHeight`, `MaxHeight` | dp | Limits on the final size. |
| `AspectRatio` | width ÷ height | Derives the free side from the sized one. |
| `Padding`, `Margin` | `UniformInset(v)`, `SymmetricInset(h, v)`, `Inset{Top, Right, Bottom, Left}` | dp. |
| `Position` | `PositionStatic` (default), `PositionRelative`, `PositionAbsolute` | See below. |
| `Top`, `Right`, `Bottom`, `Left` | dp | Offsets of absolutely positioned children. |
| `ZIndex` | integer | Higher paints on top and gets input first. Unlike CSS, it also reorders flex children. |
| `Overflow` | `OverflowVisible` (default), `OverflowHidden` | Hidden clips content inside the border, following rounded corners. |
| `Display` | `DisplayNone` | Removes the element: no space, nothing painted. |
| `Visibility` | `VisibilityHidden` | Keeps the space but paints nothing and takes no input. |

## Background, border and effects

| Property | Values | Notes |
| --- | --- | --- |
| `Background` | color | |
| `BackgroundGradient` | `giode.Linear(angle, colors...)` | Linear gradient with any number of colors. CSS angles: 90 is left to right, 180 (the default) top to bottom. Use `LinearGradient{Angle, Stops}` for uneven stops. |
| `BackgroundImage`, `BackgroundFit` | `image.Image`; `FitContain` (default), `FitCover`, `FitStretch`, `FitNone` | |
| `BorderWidth`, `BorderColor` | dp, color | |
| `BorderRadius` | dp | Same radius on every corner. |
| `BorderRadii` | `UniformCorners(r)`, `TopCorners(r)`, `BottomCorners(r)`, `Corners{TopLeft, TopRight, BottomRight, BottomLeft}` | Per-corner radii; wins over `BorderRadius`. |
| `BoxShadow` | `BoxShadow{X, Y, Blur, Spread, Color}` | dp. |
| `Opacity` | 0–1 | Applies to the element and its children. `0` means unset, so fully opaque; hide with `Visibility` instead. |
| `Transform` | `Transform{Rotate, Scale, ScaleX, ScaleY, TranslateX, TranslateY, OriginX, OriginY}` | Rotation in degrees, translation in dp, origin as a fraction of the size (center by default). Affects painting and input, not layout. |
| `Transition` | `StyleTransition{Duration, Easing}` | Animates `Background`, `Color`, `BorderColor`, shadow color and `Opacity` between a box's normal and `OnHover` styles. |

Easings: `EaseInOut` (default), `Ease`, `EaseIn`, `EaseOut`, `EaseLinear`, `EaseOutBack`,
`EaseInBack`, `EaseOutBounce`, or any CSS curve such as
`giode.Easing("cubic-bezier(0.4, 0, 0.2, 1)")`.

## Text

| Property | Values | Notes |
| --- | --- | --- |
| `Color` | color | Inherited from the enclosing box; black if none sets it. |
| `FontSize` | sp | 14 by default. |
| `FontWeight` | `FontWeightThin` … `FontWeightBlack` | |
| `FontStyle` | `FontStyleNormal`, `FontStyleItalic` | |
| `FontFamily` | typeface name | Must be a registered font. |
| `TextAlign` | `TextAlignStart` (default), `TextAlignCenter`, `TextAlignEnd` | |
| `LineHeight` | sp | Distance between lines. |
| `MaxLines` | integer | Limits the number of lines. |
| `TextOverflow` | `TextOverflowClip` (default), `TextOverflowEllipsis` | Ellipsis implies a single line unless `MaxLines` is set. |
| `TextTransform` | `TextTransformUppercase`, `TextTransformLowercase`, `TextTransformCapitalize` | |
| `TextDecoration` | `TextDecorationUnderline`, `TextDecorationLineThrough` | |

## Other

| Property | Values | Notes |
| --- | --- | --- |
| `Cursor` | `CursorPointer`, `CursorText`, `CursorGrab`, `CursorNotAllowed`, ... | Cursor while the pointer is over the element (with `OnHover`) or a component. |
| `Fit` | `FitContain`, `FitCover`, `FitStretch`, `FitNone` | How `Image` scales. |
| `ScrollBar` | `ScrollBar{Width, TrackColor, ThumbColor, Radius, MinThumbLength}` | Scrollbar of `Scroll`, `VirtualList` and `TreeView`. |

## Absolute positioning

A box child with `Position: PositionAbsolute` leaves the normal flow and is placed by `Top`
and `Left` (or `Right` and `Bottom`) inside its parent's padding. It paints above the other
children and doesn't affect the parent's size.

```go
giode.Box(giode.Styles{}, avatar,
	giode.Box(giode.Styles{
		Position: giode.PositionAbsolute, Top: 2, Right: 2, // 0 would mean "unset"
		Background: "theme.danger", BorderRadius: 99, Padding: giode.SymmetricInset(6, 2),
	}, giode.Text("3", giode.Styles{Color: "white", FontSize: 11})),
)
```

Only boxes can be positioned this way; wrap other elements in a `Box`.

## Not supported

CSS grid, floats, `position: sticky`, `white-space: nowrap`, keyframe animations and
pseudo-elements have no equivalent. Use [Transition](components/transition.md) for
show/hide animations.
