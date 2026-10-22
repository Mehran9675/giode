# Context menu

A menu that opens at the pointer or below a trigger, stays inside the window, and supports
icons, separators, disabled items and submenus.

```go
edit := giode.Menu().
	Item("Cut", cut).
	Item("Copy", copy).
	IconItem(trashIcon, "Delete", remove).
	Separator().
	DisabledItem("Paste special")
```

## Right-click anywhere

```go
app.SetContextMenu(edit)
```

A right-click anywhere in the window opens the menu at the pointer. The app draws it on top;
there's nothing to place in the view.

## Dropdown from a button

`Trigger` wraps any element so clicking it opens the menu below it, like a "File" menu:

```go
file := giode.Menu().Item("New", newDoc).Item("Open…", openDoc)
fileButton := file.Trigger(giode.Button("File").Variant(giode.VariantGhost))

// in the view: place the trigger, and lay the menu out anywhere
giode.Row(giode.Styles{}, fileButton, file)
```

## Submenus

```go
share := giode.Menu().Item("Email", byEmail).Item("Link", copyLink)
file := giode.Menu().Item("New", newDoc).SubMenu("Share", share)
```

A submenu opens to the side when its item is hovered, and stays open while the pointer moves
into it. Don't register a submenu with `SetContextMenu`, `Trigger` or `OpenAt` itself.

## Opening it yourself

```go
menu.OpenAt(f32.Pt(x, y)) // window coordinates, in pixels
menu.Toggle(f32.Pt(x, y))
menu.Close()
menu.Opened()
```

Lay the menu out somewhere in the view when you open it this way; it draws above everything
regardless of where.

## Methods

| Method | Description |
| --- | --- |
| `Menu(st...) *Menu` | Creates an empty menu. |
| `Item(label, fn)` | Adds an item. |
| `IconItem(icon, label, fn)` | Adds an item with an icon: any icon source, such as a loaded SVG or `IconPath`. |
| `DisabledItem(label)` | Adds a grayed-out item. |
| `Separator()` | Adds a line. |
| `SubMenu(label, menu)` | Adds an item that opens `menu` on hover. |
| `Trigger(el) Element` | Wraps `el` so clicking it opens the menu below it. |

Styles: `Background`, `Color` (item text), `BorderRadius`, `Padding`.

## Behavior

- Clicking an item runs its function and closes the menu.
- Clicking outside or pressing Escape closes it, including any open submenus. While it's
  open, clicks outside don't reach the rest of the UI.
- Near the bottom edge it opens above the pointer instead of below.
