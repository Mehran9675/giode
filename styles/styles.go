package styles

import "github.com/mehran9675/giode/styles/properties"

// Styles is the set of CSS-like properties applied to an element.
// Zero values mean "unset"; elements fall back to their defaults.
type Styles struct {
	// FlexDirection is the main axis of a flex container: "row" or
	// "column".
	FlexDirection properties.FlexDirection
	// Align positions children on the flex cross axis.
	Align properties.Align
	// Justify distributes children on the flex main axis.
	Justify properties.Justify
	// Wrap controls whether children continue onto a new line when
	// they overflow the main axis. Unset defaults to wrapping on Row
	// containers and not wrapping on Stack (column) containers.
	Wrap properties.Wrap
	// Gap is the spacing between flex children, and between wrapped
	// lines.
	Gap properties.Gap
	// FlexGrow is the share of spare space a flex child takes.
	FlexGrow properties.FlexGrow
	// Width and Height size the element: 0 hugs content, -1 expands,
	// >0 is fixed.
	Width  properties.Width
	Height properties.Height
	// WidthPct and HeightPct size the element as a percentage of the
	// available space and take precedence over Width and Height.
	WidthPct  properties.WidthPct
	HeightPct properties.HeightPct
	// Padding is the space between the edge and the content.
	Padding properties.Padding
	// Margin is the space outside the painted area.
	Margin properties.Margin
	// Background is the fill color behind the content.
	Background properties.Background
	// BackgroundImage is an image behind the content, scaled with
	// BackgroundFit.
	BackgroundImage properties.BackgroundImage
	BackgroundFit   properties.BackgroundFit
	// BorderWidth and BorderColor draw a border around the element.
	BorderWidth properties.BorderWidth
	BorderColor properties.BorderColor
	// BorderRadius rounds the background corners.
	BorderRadius properties.BorderRadius
	// Opacity is the element transparency in the range 0-1.
	Opacity properties.Opacity
	// Color is the text color.
	Color properties.Color
	// FontSize is the text size in sp.
	FontSize properties.FontSize
	// FontWeight is the text weight.
	FontWeight properties.FontWeight
	// TextAlign aligns text within its box.
	TextAlign properties.TextAlign
	// LineHeight is the distance between text baselines.
	LineHeight properties.LineHeight
	// MaxLines truncates text after the given number of lines.
	MaxLines properties.MaxLines
	// Cursor is the cursor shown while hovering.
	Cursor properties.Cursor
	// Fit controls how images scale into their box.
	Fit properties.Fit
	// Position selects static, relative or absolute placement.
	Position properties.Position
	// Top, Right, Bottom and Left offset absolutely positioned
	// elements from their parent box edges.
	Top    properties.Top
	Right  properties.Right
	Bottom properties.Bottom
	Left   properties.Left
	// ZIndex orders siblings: higher values paint above.
	ZIndex properties.ZIndex
	// MinWidth, MaxWidth, MinHeight and MaxHeight clamp the resolved
	// size.
	MinWidth  properties.MinWidth
	MaxWidth  properties.MaxWidth
	MinHeight properties.MinHeight
	MaxHeight properties.MaxHeight
	// BoxShadow paints a drop shadow behind the element.
	BoxShadow properties.BoxShadow
	// ScrollBar customizes the scrollbar a Scroll container draws.
	ScrollBar properties.ScrollBar
	// AspectRatio derives the unresized axis from the resized one.
	AspectRatio properties.AspectRatio
	// Overflow clips content to the box when set to hidden.
	Overflow properties.Overflow
	// TextDecoration decorates text.
	TextDecoration properties.TextDecoration
	// FontStyle selects normal or italic text.
	FontStyle properties.FontStyle
	// FontFamily selects the text typeface.
	FontFamily properties.FontFamily
}
