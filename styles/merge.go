package styles

import "github.com/mehran9675/giode/styles/properties"

// Merge overlays the non-zero values of over on top of base and
// returns the result. It is used by components to apply defaults
// before user styles.
func Merge(base, over Styles) Styles {
	if over.Display != "" {
		base.Display = over.Display
	}
	if over.Direction != "" {
		base.Direction = over.Direction
	}
	if over.Align != "" {
		base.Align = over.Align
	}
	if over.Justify != "" {
		base.Justify = over.Justify
	}
	if over.Gap != 0 {
		base.Gap = over.Gap
	}
	if over.FlexGrow != 0 {
		base.FlexGrow = over.FlexGrow
	}
	if over.Width != 0 {
		base.Width = over.Width
	}
	if over.Height != 0 {
		base.Height = over.Height
	}
	if over.WidthPct != 0 {
		base.WidthPct = over.WidthPct
	}
	if over.HeightPct != 0 {
		base.HeightPct = over.HeightPct
	}
	if over.Padding != (properties.Inset{}) {
		base.Padding = over.Padding
	}
	if over.Margin != (properties.Inset{}) {
		base.Margin = over.Margin
	}
	if over.Background.A != 0 {
		base.Background = over.Background
	}
	if over.BackgroundImage != nil {
		base.BackgroundImage = over.BackgroundImage
	}
	if over.BackgroundFit != "" {
		base.BackgroundFit = over.BackgroundFit
	}
	if over.BorderWidth != 0 {
		base.BorderWidth = over.BorderWidth
	}
	if over.BorderColor.A != 0 {
		base.BorderColor = over.BorderColor
	}
	if over.BorderRadius != 0 {
		base.BorderRadius = over.BorderRadius
	}
	if over.Opacity != 0 {
		base.Opacity = over.Opacity
	}
	if over.Color.A != 0 {
		base.Color = over.Color
	}
	if over.FontSize != 0 {
		base.FontSize = over.FontSize
	}
	if over.FontWeight != "" {
		base.FontWeight = over.FontWeight
	}
	if over.TextAlign != "" {
		base.TextAlign = over.TextAlign
	}
	if over.LineHeight != 0 {
		base.LineHeight = over.LineHeight
	}
	if over.MaxLines != 0 {
		base.MaxLines = over.MaxLines
	}
	if over.Cursor != "" {
		base.Cursor = over.Cursor
	}
	if over.Fit != "" {
		base.Fit = over.Fit
	}
	if over.Position != "" {
		base.Position = over.Position
	}
	if over.Top != 0 {
		base.Top = over.Top
	}
	if over.Right != 0 {
		base.Right = over.Right
	}
	if over.Bottom != 0 {
		base.Bottom = over.Bottom
	}
	if over.Left != 0 {
		base.Left = over.Left
	}
	if over.ZIndex != 0 {
		base.ZIndex = over.ZIndex
	}
	if over.MinWidth != 0 {
		base.MinWidth = over.MinWidth
	}
	if over.MaxWidth != 0 {
		base.MaxWidth = over.MaxWidth
	}
	if over.MinHeight != 0 {
		base.MinHeight = over.MinHeight
	}
	if over.MaxHeight != 0 {
		base.MaxHeight = over.MaxHeight
	}
	if over.BoxShadow.Color.A != 0 {
		base.BoxShadow = over.BoxShadow
	}
	if over.AspectRatio != 0 {
		base.AspectRatio = over.AspectRatio
	}
	if over.Overflow != "" {
		base.Overflow = over.Overflow
	}
	if over.TextDecoration != "" {
		base.TextDecoration = over.TextDecoration
	}
	if over.FontStyle != "" {
		base.FontStyle = over.FontStyle
	}
	if over.FontFamily != "" {
		base.FontFamily = over.FontFamily
	}
	return base
}
