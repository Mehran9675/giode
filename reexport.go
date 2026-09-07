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

type Breakpoint = elements.Breakpoint
type Tab = tabs.Tab
type Side = drawer.Side

type Button = button.Button
type Input = input.Input
type Checkbox = checkbox.Checkbox
type Slider = slider.Slider
type Tabs = tabs.Tabs
type Select = dropdown.Select
type Scroll = scroll.Scroll
type Drawer = drawer.Drawer
type Dialog = dialog.Dialog
type Menu = menu.Menu

type Router = router.Router
type Route = router.Route

// --- Elements ------------------------------------------------------------

var (
	Box  = elements.Box
	Flex = elements.Flex
	Text = elements.Text
	Raw  = elements.Raw

	Divider    = elements.Divider
	Responsive = elements.Responsive
)

// --- Components -----------------------------------------------------------

var (
	NewButton   = button.New
	NewInput    = input.New
	NewCheckbox = checkbox.New
	NewSlider   = slider.New
	NewTabs     = tabs.New
	NewSelect   = dropdown.New
	NewScroll   = scroll.New
	NewDrawer   = drawer.New
	NewDialog   = dialog.New
	NewMenu     = menu.New

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
	HexColor = styles.Hex

	UniformInset   = properties.UniformInset
	SymmetricInset = properties.SymmetricInset
)

// --- Routing ---------------------------------------------------------------

var NewRouter = router.New

// --- Style constants ---------------------------------------------------------

const (
	DisplayBlock = properties.Block
	DisplayFlex  = properties.Flex
	DisplayNone  = properties.None

	DirRow    = properties.Row
	DirColumn = properties.Column

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
