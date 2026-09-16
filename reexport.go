package giode

import (
	"github.com/mehran9675/giode/components/accordion"
	"github.com/mehran9675/giode/components/appshell"
	"github.com/mehran9675/giode/components/autocomplete"
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
	"github.com/mehran9675/giode/components/pagination"
	"github.com/mehran9675/giode/components/progress"
	"github.com/mehran9675/giode/components/scroll"
	"github.com/mehran9675/giode/components/slider"
	"github.com/mehran9675/giode/components/spinner"
	"github.com/mehran9675/giode/components/tabs"
	"github.com/mehran9675/giode/components/title"
	"github.com/mehran9675/giode/components/tooltip"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/fonts"
	"github.com/mehran9675/giode/router"
	"github.com/mehran9675/giode/state"
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

type Mode = properties.Mode
type Palette = properties.Palette

type AppShellConfig = appshell.Config
type AppShellSlot = appshell.Slot

type ButtonVariant = button.Variant

// State is a minimal global-state primitive: an observable value
// readable/writable from anywhere, including a goroutine outside the
// render loop (a network response, a timer), not just from within a
// Layout call. Create it with NewState.
type State[T any] = state.State[T]

// NewState returns a State holding initial.
func NewState[T any](initial T) *State[T] {
	return state.New(initial)
}

const (
	VariantContained = button.VariantContained
	VariantOutlined  = button.VariantOutlined
	VariantGhost     = button.VariantGhost
)

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

	Table       = elements.Table
	TableHeader = elements.TableHeader
	TableRow    = elements.TableRow
	TableCell   = elements.TableCell
)

// --- Components -----------------------------------------------------------

var (
	Button = button.New
	Tabs   = tabs.New
	Scroll = scroll.New
	Drawer = drawer.New
	Dialog = dialog.New
	Menu   = menu.New

	Progress   = progress.Progress
	Spinner    = spinner.Spinner
	Image      = imageutil.New
	Tooltip    = tooltip.New
	Accordion  = accordion.New
	Pagination = pagination.New

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
	// AutoComplete creates a text field with a filtered suggestion
	// list.
	AutoComplete func(placeholder string, options []string, st ...Styles) *autocomplete.AutoComplete
}

// Input groups the form-control constructors that collect user
// input: Input.Text, Input.Checkbox, Input.Slider, Input.Select and
// Input.AutoComplete.
var Input = InputAPI{
	Text:         input.New,
	Checkbox:     checkbox.New,
	Slider:       slider.New,
	Select:       dropdown.New,
	AutoComplete: autocomplete.New,
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

	// LoadIconSVG parses a whole .svg file's bytes (viewBox and one or
	// more <path> elements), unlike IconSVG/IconSVGCustom, which take
	// raw path data. Call .Icon(st...) on the result to render it.
	LoadIconSVG = icon.LoadSVG
	// MustLoadIconSVG is LoadIconSVG, panicking on error - convenient
	// for bytes known at compile time, e.g. via go:embed.
	MustLoadIconSVG = icon.MustLoadSVG
)

type IconDocument = icon.Document

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

// --- Theming ---------------------------------------------------------------

const (
	ModeLight = properties.ModeLight
	ModeDark  = properties.ModeDark
)

var (
	// SetMode switches the active color mode, changing what every
	// "theme.*" color token (e.g. Styles{Background: "theme.surface"})
	// resolves to on the next frame.
	SetMode = properties.SetMode
	// CurrentMode returns the active color mode (ModeLight by
	// default).
	CurrentMode = properties.CurrentMode
	// SetPalette replaces the palette used for mode m, changing what
	// every "theme.*" color token resolves to while that mode is
	// active. Start from DefaultLightPalette/DefaultDarkPalette to
	// override just a few fields.
	SetPalette = properties.SetPalette
)

// DefaultLightPalette and DefaultDarkPalette are the built-in
// palettes for ModeLight/ModeDark, in effect until SetPalette
// overrides one of them.
var (
	DefaultLightPalette = properties.DefaultLightPalette
	DefaultDarkPalette  = properties.DefaultDarkPalette
)

// --- AppShell ---------------------------------------------------------------

var (
	// AppShell returns a full-window layout: Header/Navbar/Aside/Footer
	// (each optional, via AppShellConfig) are fixed to their respective
	// edges, and Main fills the remaining space, offset by whichever of
	// them are present.
	AppShell = appshell.New
	// Header, Navbar, Aside and Footer build an AppShellConfig field -
	// AppShellConfig{Header: giode.Header(60, ...), ...}.
	Header = appshell.Header
	Navbar = appshell.Navbar
	Aside  = appshell.Aside
	Footer = appshell.Footer
	// Main wraps an AppShell's main content area.
	Main = appshell.Main
	// Section groups content inside a Navbar/Aside as a scrollable
	// area.
	Section = appshell.Section
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
