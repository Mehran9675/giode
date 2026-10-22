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
	"github.com/mehran9675/giode/components/iconbutton"
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
	"github.com/mehran9675/giode/components/titlebar"
	"github.com/mehran9675/giode/components/tooltip"
	"github.com/mehran9675/giode/components/transition"
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

type (
	Corners         = properties.Corners
	ColorStop       = properties.ColorStop
	LinearGradient  = properties.LinearGradient
	Transform       = properties.Transform
	StyleTransition = properties.Transition
	Easing          = properties.Easing
	Display         = properties.Display
	Visibility      = properties.Visibility
	TextTransform   = properties.TextTransform
	TextOverflow    = properties.TextOverflow

	// MemoCache caches a subtree's painted output across frames; see
	// Memo.
	PressState       = elements.PressState
	ClickEvent       = elements.ClickEvent
	ClickSource      = elements.ClickSource
	MouseButton      = elements.MouseButton
	PressableElement = elements.Pressable

	MemoCache    = elements.MemoCache
	SemanticNode = elements.SemanticNode
	SemanticKind = elements.SemanticKind
)
type Tab = tabs.Tab
type Side = drawer.Side

type Mode = properties.Mode
type Palette = properties.Palette

type AppShellConfig = appshell.Config
type AppShellSlot = appshell.Slot

type ButtonVariant = button.Variant
type IconButtonVariant = iconbutton.Variant

// State is a minimal global-state primitive: an observable value
// readable/writable from anywhere, including a goroutine outside the
// render loop (a network response, a timer), not just from within a
// Layout call. Create it with NewState.
type State[T any] = state.State[T]

// NewState returns a State holding initial.
func NewState[T any](initial T) *State[T] {
	return state.New(initial)
}

// Keyed returns the state stored under (scope, key), creating it with
// init the first time - persistent per-item state (a Pressable, a text
// field) for lists that change at runtime. Entries not asked for in a
// while are dropped, so removed items don't leak. Use stable keys
// (IDs), not indexes.
func Keyed[T any](scope, key any, init func() T) T {
	return elements.Keyed(scope, key, init)
}

// KeyedCache is a typed map of per-item state that drops entries not
// asked for in a while; hold one in a stateful component.
type KeyedCache[K comparable, V any] = elements.KeyedCache[K, V]

// NewKeyedCache returns an empty KeyedCache.
func NewKeyedCache[K comparable, V any]() *KeyedCache[K, V] {
	return elements.NewKeyedCache[K, V]()
}

var (
	// ForgetKeyed drops a Keyed entry right away.
	ForgetKeyed = elements.ForgetKeyed
	// SetKeyedRetention sets how many frames Keyed entries survive
	// without being asked for.
	SetKeyedRetention = elements.SetKeyedRetention
)

// NamedState returns a State registered under name, so tooling - the
// MCP server's state_list/state_set tools - can read and set it. T
// must be JSON-encodable for that.
func NamedState[T any](name string, initial T) *State[T] {
	return state.Named(name, initial)
}

const (
	VariantContained = button.VariantContained
	VariantOutlined  = button.VariantOutlined
	VariantGhost     = button.VariantGhost

	IconButtonContained = iconbutton.VariantContained
	IconButtonOutlined  = iconbutton.VariantOutlined
	IconButtonGhost     = iconbutton.VariantGhost
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

	// Pressable makes any content clickable; its content is rebuilt
	// with the current PressState so it can restyle itself.
	Pressable = elements.NewPressable

	// Memo caches the output of build under key, rebuilding only when
	// its deps, the States it reads, its constraints or the theme
	// change - see elements.MemoCache for the full rules.
	Memo = elements.Memo
	// NewMemoCache returns a cache to hold in a stateful component.
	NewMemoCache = elements.NewMemoCache
	// InvalidateCaches forces every Memo to rebuild on the next frame.
	InvalidateCaches = elements.InvalidateCaches

	// AfterLayout runs a function once the current frame's layout is
	// done, on the UI goroutine - for handlers that change the model
	// the frame is still reading.
	AfterLayout = elements.AfterLayout
	// DisplayScale returns the display's pixels per dp as of the last
	// frame; DpToPx and PxToDp convert with it. Safe from anywhere.
	DisplayScale = elements.DisplayScale
	DpToPx       = elements.DpToPx
	PxToDp       = elements.PxToDp

	Table       = elements.Table
	TableHeader = elements.TableHeader
	TableRow    = elements.TableRow
	TableCell   = elements.TableCell
)

// --- Components -----------------------------------------------------------

var (
	Button     = button.New
	IconButton = iconbutton.New
	Tabs       = tabs.New
	Scroll     = scroll.New
	Drawer     = drawer.New
	Dialog     = dialog.New
	Menu       = menu.New

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

// Giode bundles no icon set: load the icons your app uses from SVG
// files (or path data, IconVG, images) and, optionally, register them
// by name.
var (
	// IconSVG renders raw SVG path data on a 24x24 viewBox.
	IconSVG = icon.SVG
	// IconSVGCustom renders raw SVG path data on a custom viewBox.
	IconSVGCustom = icon.SVGCustom
	// IconSVGStroke strokes raw SVG path data (outline icon sets).
	IconSVGStroke = icon.SVGStroke
	IconImage     = icon.Image
	IconIconVG    = icon.IconVG
	LoadIconVG    = icon.LoadIconVG

	// LoadIconSVG parses a whole .svg file's bytes; call .Icon(st...)
	// on the result to render it.
	LoadIconSVG = icon.LoadSVG
	// MustLoadIconSVG is LoadIconSVG, panicking on error - convenient
	// for bytes known at compile time, e.g. via go:embed.
	MustLoadIconSVG = icon.MustLoadSVG
	// LoadIconSVGFile and LoadIconSVGFS read and parse an .svg file
	// from disk or from an fs.FS (such as an embed.FS).
	LoadIconSVGFile = icon.LoadSVGFile
	LoadIconSVGFS   = icon.LoadSVGFS

	// RegisterIcon, RegisterIconSVG, RegisterIconsFS and
	// RegisterIconsDir add icons to the default registry; Icon renders
	// one by name.
	RegisterIcon     = icon.Register
	RegisterIconSVG  = icon.RegisterSVG
	RegisterIconsFS  = icon.RegisterFS
	RegisterIconsDir = icon.RegisterDir
	LookupIcon       = icon.Lookup
	IconNames        = icon.Names
	// Icon renders a registered icon; an unknown name renders as empty
	// space of the icon's size.
	Icon = icon.Named
)

type (
	IconDocument = icon.Document
	IconSource   = icon.Source
	IconPath     = icon.Path
	IconRegistry = icon.Registry
)

// --- Fonts ----------------------------------------------------------------

var (
	UseGofont     = fonts.UseGofont
	RegisterFonts = fonts.Register
)

// --- Colors and insets -----------------------------------------------------

var (
	HexColor = properties.HexColor
	// ParseColor validates any CSS color (hex, rgb(), hsl(), named).
	ParseColor = properties.ParseColor
	// Linear builds a linear gradient for Styles.BackgroundGradient.
	Linear = properties.Linear

	UniformCorners = properties.UniformCorners
	TopCorners     = properties.TopCorners
	BottomCorners  = properties.BottomCorners

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
	// ToggleMode switches between light and dark mode.
	ToggleMode = properties.ToggleMode
	// CurrentPalette returns the palette of the active mode.
	CurrentPalette = properties.CurrentPalette
)

// DefaultLightPalette and DefaultDarkPalette are the built-in
// palettes for ModeLight/ModeDark, in effect until SetPalette
// overrides one of them.
var (
	DefaultLightPalette = properties.DefaultLightPalette
	DefaultDarkPalette  = properties.DefaultDarkPalette
)

// --- Title bar ----------------------------------------------------------------

// TitleBar returns a custom window title bar for a frameless window:
// a drag handle with the title, optional icon and menus, and
// minimize/maximize/close buttons. Pass the *App:
//
//	app := giode.New(giode.Config{Frameless: true, Resizable: true})
//	bar := giode.TitleBar(app, "My app")
var TitleBar = titlebar.New

type (
	TitleBarVariant = titlebar.Variant
	TitleBarButton  = titlebar.Button
	TitleBarWindow  = titlebar.Window
)

const (
	TitleBarWindows = titlebar.VariantWindows
	TitleBarMac     = titlebar.VariantMac

	TitleBarMinimize = titlebar.ButtonMinimize
	TitleBarMaximize = titlebar.ButtonMaximize
	TitleBarClose    = titlebar.ButtonClose
)

// App implements the window controls a TitleBar drives.
var _ titlebar.Window = (*App)(nil)

// --- Transitions ----------------------------------------------------------------

// Transition animates a component in and out. It takes the component,
// the effect (combine kinds with |), the duration and a CSS-style
// easing:
//
//	panel := giode.Transition(content, giode.FadeUp, giode.TransitionNormal, giode.EaseOut)
//	giode.Button("Toggle").OnClick(panel.Toggle)
var Transition = transition.New

type (
	TransitionKind    = transition.Kind
	TransitionElement = transition.Transition
)

const (
	Fade          = transition.Fade
	SlideUp       = transition.SlideUp
	SlideDown     = transition.SlideDown
	SlideLeft     = transition.SlideLeft
	SlideRight    = transition.SlideRight
	Enlarge       = transition.Enlarge
	Shrink        = transition.Shrink
	Rotate        = transition.Rotate
	Collapse      = transition.Collapse
	CollapseWidth = transition.CollapseWidth
	FadeUp        = transition.FadeUp
	Zoom          = transition.Zoom
	Pop           = transition.Pop

	TransitionFast   = transition.Fast
	TransitionNormal = transition.Normal
	TransitionSlow   = transition.Slow
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

	DisplayNone       = properties.DisplayNone
	VisibilityVisible = properties.VisibilityVisible
	VisibilityHidden  = properties.VisibilityHidden

	TextTransformNone       = properties.TextTransformNone
	TextTransformUppercase  = properties.TextTransformUppercase
	TextTransformLowercase  = properties.TextTransformLowercase
	TextTransformCapitalize = properties.TextTransformCapitalize

	TextOverflowClip     = properties.TextOverflowClip
	TextOverflowEllipsis = properties.TextOverflowEllipsis

	MousePrimary   = elements.MousePrimary
	MouseSecondary = elements.MouseSecondary
	MouseMiddle    = elements.MouseMiddle

	ClickPointer      = elements.ClickPointer
	ClickKeyboard     = elements.ClickKeyboard
	ClickProgrammatic = elements.ClickProgrammatic

	EaseInOut     = properties.EaseInOut
	Ease          = properties.Ease
	EaseLinear    = properties.EaseLinear
	EaseIn        = properties.EaseIn
	EaseOut       = properties.EaseOut
	EaseOutBack   = properties.EaseOutBack
	EaseInBack    = properties.EaseInBack
	EaseOutBounce = properties.EaseOutBounce
)
