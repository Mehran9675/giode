package giode

import (
	"github.com/mehran9675/giode/components/button"
	"github.com/mehran9675/giode/components/checkbox"
	"github.com/mehran9675/giode/components/dialog"
	"github.com/mehran9675/giode/components/drawer"
	"github.com/mehran9675/giode/components/dropdown"
	"github.com/mehran9675/giode/components/icon"
	"github.com/mehran9675/giode/components/imageutil"
	"github.com/mehran9675/giode/components/input"
	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/components/menu"
	"github.com/mehran9675/giode/components/progress"
	"github.com/mehran9675/giode/components/scroll"
	"github.com/mehran9675/giode/components/slider"
	"github.com/mehran9675/giode/components/spinner"
	"github.com/mehran9675/giode/components/tabs"
	"github.com/mehran9675/giode/components/title"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/fonts"
	"github.com/mehran9675/giode/router"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// --- Types ---------------------------------------------------------------

type Element = elements.Element
type Styles = styles.Styles
type Inset = properties.Inset
type BoxShadow = properties.BoxShadow
type ScrollBar = properties.ScrollBar

type Breakpoint = elements.Breakpoint
type Tab = tabs.Tab
type Side = drawer.Side

type Router = router.Router
type Route = router.Route

// --- Elements ------------------------------------------------------------

var (
	Box   = elements.Box
	Row   = elements.Row
	Stack = elements.Stack
	Text  = elements.Text
	Raw   = elements.Raw

	Divider    = elements.Divider
	Responsive = elements.Responsive
)

// --- Components -----------------------------------------------------------

var (
	Button = button.New
	Tabs   = tabs.New
	Scroll = scroll.New
	Drawer = drawer.New
	Dialog = dialog.New
	Menu   = menu.New

	Progress = progress.Progress
	Spinner  = spinner.Spinner
	Image    = imageutil.New

	H1 = title.H1
	H2 = title.H2
	H3 = title.H3
	H4 = title.H4
	H5 = title.H5
	H6 = title.H6
)

// InputAPI groups the constructors for components whose purpose is
// collecting a value from the user.
type InputAPI struct {
	// Text creates a single-line text field.
	Text func(placeholder string, st ...Styles) *input.Input
	// Checkbox creates a checkbox. label is optional: pass "" for a
	// bare checkbox with no label.
	Checkbox func(label string, st ...Styles) *checkbox.Checkbox
	// Slider creates a horizontal slider in the range [0, 1].
	Slider func(st ...Styles) *slider.Slider
	// Select creates a dropdown choosing among options.
	Select func(st Styles, options ...string) *dropdown.Select
}

// Input groups the form-control constructors that collect user
// input: Input.Text, Input.Checkbox, Input.Slider and Input.Select.
var Input = InputAPI{
	Text:     input.New,
	Checkbox: checkbox.New,
	Slider:   slider.New,
	Select:   dropdown.New,
}

const (
	SideLeft  = drawer.Left
	SideRight = drawer.Right
)

// --- Icons ----------------------------------------------------------------

var (
	IconMaterial  = icon.Material
	IconSVG       = icon.SVG
	IconSVGCustom = icon.SVGCustom
	IconImage     = icon.Image
	IconIconVG    = icon.IconVG
)

// --- Fonts ----------------------------------------------------------------

var (
	UseGofont     = fonts.UseGofont
	RegisterFonts = fonts.Register
)

// --- Colors and insets -----------------------------------------------------

var (
	HexColor = properties.HexColor

	UniformInset   = properties.UniformInset
	SymmetricInset = properties.SymmetricInset
)

// --- Routing ---------------------------------------------------------------

var NewRouter = router.New

// --- Component kit -----------------------------------------------------------

var (
	Scale    = kit.Scale
	Hovered  = kit.Hovered
	Pressed  = kit.Pressed
	Alpha    = kit.Alpha
	ClampPos = kit.ClampPos
)

// OverlayMargin is the default distance kept between an overlay and
// the window edges.
const OverlayMargin = kit.Margin

// --- Style constants ---------------------------------------------------------

const (
	FlexDirectionRow    = properties.FlexDirectionRow
	FlexDirectionColumn = properties.FlexDirectionColumn

	WrapNoWrap = properties.WrapNoWrap
	WrapWrap   = properties.WrapWrap

	AlignStart  = properties.AlignStart
	AlignCenter = properties.AlignCenter
	AlignEnd    = properties.AlignEnd

	JustifyStart        = properties.JustifyStart
	JustifyCenter       = properties.JustifyCenter
	JustifyEnd          = properties.JustifyEnd
	JustifySpaceBetween = properties.JustifySpaceBetween
	JustifySpaceAround  = properties.JustifySpaceAround
	JustifySpaceEvenly  = properties.JustifySpaceEvenly

	TextAlignStart  = properties.TextAlignStart
	TextAlignCenter = properties.TextAlignCenter
	TextAlignEnd    = properties.TextAlignEnd

	FontWeightThin      = properties.FontWeightThin
	FontWeightLight     = properties.FontWeightLight
	FontWeightNormal    = properties.FontWeightNormal
	FontWeightMedium    = properties.FontWeightMedium
	FontWeightBold      = properties.FontWeightBold
	FontWeightExtraBold = properties.FontWeightExtraBold
	FontWeightBlack     = properties.FontWeightBlack

	FontStyleNormal = properties.FontStyleNormal
	FontStyleItalic = properties.FontStyleItalic

	CursorDefault    = properties.CursorDefault
	CursorNone       = properties.CursorNone
	CursorText       = properties.CursorText
	CursorPointer    = properties.CursorPointer
	CursorCrosshair  = properties.CursorCrosshair
	CursorGrab       = properties.CursorGrab
	CursorGrabbing   = properties.CursorGrabbing
	CursorNotAllowed = properties.CursorNotAllowed
	CursorWait       = properties.CursorWait

	FitNone    = properties.FitNone
	FitContain = properties.FitContain
	FitCover   = properties.FitCover
	FitStretch = properties.FitStretch

	PositionStatic   = properties.PositionStatic
	PositionRelative = properties.PositionRelative
	PositionAbsolute = properties.PositionAbsolute

	TextDecorationNone        = properties.TextDecorationNone
	TextDecorationUnderline   = properties.TextDecorationUnderline
	TextDecorationLineThrough = properties.TextDecorationLineThrough

	OverflowVisible = properties.OverflowVisible
	OverflowHidden  = properties.OverflowHidden
)
