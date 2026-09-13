# Style properties

Every element accepts a `giode.Styles{...}` struct. Zero values mean "unset": the element
falls back to its default. Each property is defined in its own file under
`styles/properties/`.

```go
card := giode.Styles{
	FlexDirection: giode.FlexDirectionColumn,
	Gap:           12,
	Padding:       giode.UniformInset(16),
	Background:    "#1e293b",
	BorderRadius:  12,
	BoxShadow:     giode.BoxShadow{Y: 4, Blur: 12, Color: "#00000060"},
}
```

## Colors

Colors are strings: CSS hex (`"#fff"`, `"#1e293b"`, `"#00000060"`) or one of the named
colors `white`, `black`, `gray`, `yellow`, `red`, `green`, `blue`, `magenta`, `cyan`,
`purple`. An empty string means unset (transparent); an unrecognized name renders white
so a typo is visible.

## Layout

| Property | Type | Values | Notes |
| --- | --- | --- | --- |
| `FlexDirection` | string | `FlexDirectionRow` (default), `FlexDirectionColumn` | Main axis (CSS flex-direction). |
| `Align` | string | `AlignStart`, `AlignCenter`, `AlignEnd` | Cross-axis alignment (align-items). |
| `Justify` | string | `JustifyStart` (default), `JustifyCenter`, `JustifyEnd`, `JustifySpaceBetween`, `JustifySpaceAround`, `JustifySpaceEvenly` | Main-axis distribution. |
| `Wrap` | string | `WrapWrap`, `WrapNoWrap` | Whether children continue onto new lines. Unset: Row wraps, Stack does not. |
| `Gap` | int | px | Space between children and between wrapped lines. |
| `FlexGrow` | int | ≥0 | Share of spare main-axis space. |
| `Width` / `Height` | int | `0` hug, `-1` expand, `>0` fixed px | |
| `WidthPct` / `HeightPct` | int | 1–100 | Percent of available space; takes precedence over `Width`/`Height`. |
| `MinWidth` / `MaxWidth` / `MinHeight` / `MaxHeight` | int | px | Clamp the resolved size. |
| `AspectRatio` | float | w/h | Derives the hugging axis from the resized axis. |
| `Padding` / `Margin` | Inset | px per side | `UniformInset(v)`, `SymmetricInset(h, v)` or `giode.Inset{Top, Right, Bottom, Left}`. |
| `Position` | string | `PositionStatic` (default), `PositionRelative`, `PositionAbsolute` | Absolute removes the element from the flow; see below. |
| `Top` / `Right` / `Bottom` / `Left` | int | px | Offsets of absolute children. |
| `ZIndex` | int | any | Stacking order of siblings. |
| `Overflow` | string | `OverflowVisible` (default), `OverflowHidden` | Hidden clips content to the box. |
| `ScrollBar` | struct | — | Customizes the scrollbar a `Scroll` draws: `Width`, `TrackColor`, `ThumbColor`, `Radius`, `MinThumbLength`. |

## Paint

| Property | Type | Values | Notes |
| --- | --- | --- | --- |
| `Background` | string | any color | Fill color. |
| `BackgroundImage` | `image.Image` | — | Image behind the content. |
| `BackgroundFit` | string | `FitContain` (default), `FitCover`, `FitStretch`, `FitNone` | Scaling of the background image. |
| `BorderWidth` / `BorderColor` | int / string | — | Border around the box. |
| `BorderRadius` | int | px | Rounds the box corners. |
| `BoxShadow` | struct | `{X, Y, Blur int; Color string}` | Drop shadow behind the box. |
| `Opacity` | float | 0–1 | Transparency of the element and its children; 0 means unset. |

## Text

| Property | Type | Values | Notes |
| --- | --- | --- | --- |
| `Color` | string | any color | Text color; empty resolves to black. |
| `FontSize` | sp | — | Default 14sp. |
| `FontWeight` | string | `FontWeightThin` … `FontWeightBlack` | Default normal. |
| `FontStyle` | string | `FontStyleNormal`, `FontStyleItalic` | |
| `FontFamily` | string | — | Typeface must exist in the registered fonts. |
| `TextAlign` | string | `TextAlignStart` (default), `TextAlignCenter`, `TextAlignEnd` | |
| `LineHeight` | sp | — | Distance between baselines. |
| `MaxLines` | int | — | Truncates with an ellipsis. |
| `TextDecoration` | string | `TextDecorationNone` (default), `TextDecorationUnderline`, `TextDecorationLineThrough` | |

## Interaction

| Property | Type | Values | Notes |
| --- | --- | --- | --- |
| `Cursor` | string | `CursorDefault`, `CursorPointer`, `CursorText`, `CursorGrab`, … | Cursor while hovering interactive components. |
| `Fit` | string | `FitContain`, `FitCover`, `FitStretch`, `FitNone` | Scaling of `Image` elements. |

## Absolute positioning

Children of a `Box` with `Position: PositionAbsolute` are removed from the flow and placed
at `Top`/`Left` (or anchored from `Right`/`Bottom`) of the parent's content box, painting
above the flow content. They do not affect the parent's size.

```go
giode.Box(giode.Styles{}, content,
	giode.Box(giode.Styles{
		Position: giode.PositionAbsolute,
		Top: 8, Right: 8, ZIndex: 10,
		Background: "#ef4444", BorderRadius: 99,
	}, giode.Text("3")),
)
```

Notes:

- `ZIndex` sorts siblings: higher values paint above and receive input first. Unlike CSS it
  also affects layout order in flex containers.
- Absolute positioning applies to `Box` children; wrap other elements in a `Box`.

## Not supported (by design)

Complex or platform-specific CSS: transitions, animations, transforms, box-sizing,
grid, float, `position: sticky`, `white-space` (Gio's label always wraps), pseudo
elements, media-query orientation. Raise an issue if one of these matters to you.
