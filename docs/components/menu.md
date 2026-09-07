# Context menu

Right-click menu that opens at the cursor and keeps itself within the window bounds,
flipping above the cursor when it would overflow the bottom edge.

## Global right-click

```go
menu := giode.NewMenu().
	Item("New", func() { ... }).
	IconItem("trash", "Delete", func() { ... }).
	Separator().
	DisabledItem("Locked")

app := giode.New(giode.Config{...})
app.SetContextMenu(menu) // right-click anywhere in the window opens it
```

The app registers the right-click area and renders the menu on top automatically.

## Manual control

```go
menu.OpenAt(f32.Pt(x, y)) // window pixel coordinates
menu.Toggle(f32.Pt(x, y))
menu.Close()
menu.Opened()
```

Without `SetContextMenu`, lay the menu out **last** in the view so it renders above the UI.

## API

| Member | Description |
| --- | --- |
| `NewMenu() *Menu` | Creates an empty menu. |
| `Item(label string, action func()) *Menu` | Clickable item. |
| `IconItem(iconName, label string, action func()) *Menu` | Item with a material icon. |
| `DisabledItem(label string) *Menu` | Grayed-out item. |
| `Separator() *Menu` | Horizontal line. |
| `Styles(st Styles) *Menu` | Panel styles: `Background`, `BorderRadius`, `Padding`, `Color`. |

## Behavior

- Opens at the cursor position, measured and clamped to the window with an 8px margin.
- Clicking an item runs its action and closes the menu.
- Clicking outside or pressing Escape closes it.
- The dismissal layer is modal: clicks under it do not reach the UI.
