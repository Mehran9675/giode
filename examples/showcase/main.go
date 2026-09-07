// Command showcase demonstrates the Giode component suite: routing,
// a global right-click context menu, a drawer, media-query-like
// breakpoints, and the input, checkbox, slider, tabs, dropdown,
// progress, spinner, dialog, scroll and icon components.
package main

import (
	"fmt"
	"image"
	"log"

	"gioui.org/layout"

	"github.com/mehran9675/giode"
)

const (
	bg     = "#0f172a"
	card   = "#1e293b"
	border = "#334155"
	text   = "#f8fafc"
	muted  = "#94a3b8"
	accent = "#3b82f6"
)

func main() {
	giode.UseGofont()

	// --- Context menu: opens with a right-click anywhere in the
	// window and stays within the window bounds.
	menu := giode.NewMenu().
		Item("New", func() { log.Println("menu: new") }).
		IconItem("edit", "Rename", func() { log.Println("menu: rename") }).
		IconItem("copy", "Duplicate", func() { log.Println("menu: duplicate") }).
		Separator().
		IconItem("trash", "Delete", func() { log.Println("menu: delete") }).
		DisabledItem("Locked").
		Separator().
		IconItem("logout", "Quit", func() { log.Println("menu: quit") })

	// --- Stateful handles, created once.
	drawer := giode.NewDrawer(giode.SideLeft, 260, giode.Styles{Background: giode.HexColor(card)})
	dialog := giode.NewDialog(giode.Styles{})
	search := giode.NewInput("Search...")
	notify := giode.NewCheckbox("Enable notifications")
	notify.SetChecked(true)
	volume := giode.NewSlider()
	sizeSel := giode.NewSelect("Small", "Medium", "Large")
	scroller := giode.NewScroll()

	// --- Routing.
	settingsName := giode.NewInput("Your name")
	router := giode.NewRouter(
		giode.Route{Path: "/", View: func() giode.Element { return homePage(search, notify, volume, sizeSel, scroller) }},
		giode.Route{Path: "/settings", View: func() giode.Element { return settingsPage(settingsName) }},
	)

	openDrawer := giode.NewButton("Drawer").OnClick(func() { drawer.Toggle() })
	openDialog := giode.NewButton("Dialog").OnClick(func() { dialog.Open() })

	// --- Window action handles.
	centerBtn := giode.NewButton("Center")
	minimizeBtn := giode.NewButton("Minimize")
	maximizeBtn := giode.NewButton("Maximize")
	fullscreenBtn := giode.NewButton("Fullscreen")

	app := giode.New(giode.Config{
		Title:      "Giode Showcase",
		Size:       image.Pt(1000, 700),
		Resizable:  true,
		MinSize:    image.Pt(480, 360),
		Frameless:  true,
		Draggable:  true,
		Background: giode.HexColor(bg),
	})
	app.SetContextMenu(menu)

	centerBtn.OnClick(func() { app.Center() })
	minimizeBtn.OnClick(func() { app.Minimize() })
	maximizeBtn.OnClick(func() { app.Maximize() })
	fullscreenBtn.OnClick(func() { app.SetFullscreen(true) })

	if err := app.Run(func() giode.Element {
		return giode.Box(
			giode.Styles{Width: -1, Height: -1},
			giode.Box(
				giode.Styles{
					Padding:   giode.UniformInset(24),
					Color:     giode.HexColor(text),
					Display:   giode.DisplayFlex,
					Direction: giode.DirColumn,
					Gap:       16,
				},
				giode.H1("Giode Showcase", giode.Styles{Color: giode.HexColor(text)}),
				giode.Text("Right-click anywhere for the context menu. Drag empty space to move the window.", giode.Styles{Color: giode.HexColor(muted)}),
				giode.Flex(
					giode.Styles{Gap: 12, Align: giode.AlignCenter},
					giode.Box(giode.Styles{Padding: giode.SymmetricInset(16, 8), Background: giode.HexColor(accent), BorderRadius: 8}, router.Link("/", "Home")),
					giode.Box(giode.Styles{Padding: giode.SymmetricInset(16, 8), Background: giode.HexColor(card), BorderRadius: 8}, router.Link("/settings", "Settings")),
					giode.Box(giode.Styles{Width: 16}),
					openDrawer,
					openDialog,
				),
				giode.Divider(),
				router,
				giode.Text("Resize the window: the box below switches layout at 600px.", giode.Styles{Color: giode.HexColor(muted)}),
				giode.Responsive(
					giode.Breakpoint{Width: 600, View: func() giode.Element {
						return giode.Box(giode.Styles{Padding: giode.UniformInset(16), Background: giode.HexColor(card), BorderRadius: 8},
							giode.Text("Narrow layout (<= 600px)", giode.Styles{Color: giode.HexColor(accent)}))
					}},
					giode.Breakpoint{View: func() giode.Element {
						return giode.Box(giode.Styles{Padding: giode.UniformInset(16), Background: giode.HexColor(card), BorderRadius: 8},
							giode.Text("Wide layout (> 600px)", giode.Styles{Color: giode.HexColor(accent)}))
					}},
				),
				giode.Text("Window actions (frameless window):", giode.Styles{Color: giode.HexColor(muted)}),
				giode.Flex(giode.Styles{Gap: 8},
					centerBtn,
					minimizeBtn,
					maximizeBtn,
					fullscreenBtn,
				),
				giode.Box(giode.Styles{Height: 16}),
			),
			giode.Raw(func(gtx layout.Context) layout.Dimensions {
				return drawer.Layout(gtx, giode.Box(
					giode.Styles{Padding: giode.UniformInset(16), Gap: 12, Display: giode.DisplayFlex, Direction: giode.DirColumn, Color: giode.HexColor(text)},
					giode.H4("Navigation", giode.Styles{Color: giode.HexColor(text)}),
					giode.Text("Drawer item", giode.Styles{Color: giode.HexColor(muted)}),
					giode.Text("Another item", giode.Styles{Color: giode.HexColor(muted)}),
				))
			}),
			giode.Raw(func(gtx layout.Context) layout.Dimensions {
				return dialog.Layout(gtx, giode.Box(
					giode.Styles{Gap: 12, Display: giode.DisplayFlex, Direction: giode.DirColumn, Color: giode.HexColor(text)},
					giode.H4("Dialog", giode.Styles{Color: giode.HexColor(text)}),
					giode.Text("Click outside to dismiss.", giode.Styles{Color: giode.HexColor(muted)}),
				))
			}),
		)
	}); err != nil {
		log.Fatal(err)
	}
}

// homePage builds the home view: form controls and a scrolled list.
func homePage(search *giode.Input, notify *giode.Checkbox, volume *giode.Slider, sizeSel *giode.Select, scroller *giode.Scroll) giode.Element {
	// The dropdown's own type; re-exported constructors keep the API
	// single-import.
	sel := sizeSel

	return giode.Flex(
		giode.Styles{Gap: 16, Align: giode.AlignStart},
		giode.Box(
			giode.Styles{
				Padding:      giode.UniformInset(20),
				Background:   giode.HexColor(card),
				BorderRadius: 12,
				BorderWidth:  1,
				BorderColor:  giode.HexColor(border),
				BoxShadow:    giode.BoxShadow{Y: 4, Blur: 12, Color: giode.HexColor("#00000060")},
				Width:        360,
				Display:      giode.DisplayFlex,
				Direction:    giode.DirColumn,
				Gap:          16,
				Color:        giode.HexColor(text),
			},
			giode.Box(
				giode.Styles{Position: giode.PositionRelative},
				giode.H3("Controls", giode.Styles{Color: giode.HexColor(text)}),
				giode.Box(
					giode.Styles{
						Position:     giode.PositionAbsolute,
						Top:          -8,
						Right:        -60,
						ZIndex:       10,
						Padding:      giode.SymmetricInset(8, 2),
						Background:   giode.HexColor("#ef4444"),
						BorderRadius: 99,
						Color:        giode.HexColor(text),
					},
					giode.Text("3", giode.Styles{Color: giode.HexColor(text), FontSize: 11}),
				),
			),
			giode.Text("Search", giode.Styles{Color: giode.HexColor(muted)}),
			search,
			notify,
			giode.Text(fmt.Sprintf("Size: %s", sel.Value()), giode.Styles{Color: giode.HexColor(muted)}),
			sel,
			giode.Text(fmt.Sprintf("Volume: %.0f%%", volume.Value()*100), giode.Styles{Color: giode.HexColor(muted)}),
			volume,
			giode.Progress(volume.Value(), giode.Styles{Color: giode.HexColor(accent), Background: giode.HexColor(border)}),
			giode.Box(giode.Styles{Gap: 8, Display: giode.DisplayFlex, Align: giode.AlignCenter},
				giode.Spinner(giode.Styles{Width: 20, Height: 20}),
				giode.Text("Loading...", giode.Styles{Color: giode.HexColor(muted)}),
			),
		),
		giode.Box(giode.Styles{FlexGrow: 1, Height: 360},
			giode.Raw(func(gtx layout.Context) layout.Dimensions {
				var items []giode.Element
				for i := 0; i < 40; i++ {
					items = append(items, giode.Box(giode.Styles{
						Padding:      giode.UniformInset(12),
						Background:   giode.HexColor(card),
						BorderRadius: 8,
						Margin:       giode.Inset{Bottom: 8},
					}, giode.Text(fmt.Sprintf("Item %d", i+1), giode.Styles{Color: giode.HexColor(text)})))
				}
				return scroller.Layout(gtx, giode.Box(giode.Styles{}, items...))
			}),
		),
	)
}

func settingsPage(name *giode.Input) giode.Element {
	save := giode.NewButton("Save").OnClick(func() {
		log.Printf("saved: %q", name.Text())
	})
	return giode.Box(
		giode.Styles{
			Padding:      giode.UniformInset(20),
			Background:   giode.HexColor(card),
			BorderRadius: 12,
			Width:        360,
			Display:      giode.DisplayFlex,
			Direction:    giode.DirColumn,
			Gap:          12,
			Color:        giode.HexColor(text),
		},
		giode.H3("Settings", giode.Styles{Color: giode.HexColor(text)}),
		name,
		giode.Box(giode.Styles{Align: giode.AlignEnd, Display: giode.DisplayFlex, Direction: giode.DirColumn},
			save,
		),
	)
}
