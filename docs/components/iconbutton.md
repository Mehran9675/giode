# IconButton

A button that shows only an icon.

```go
del := giode.IconButton(giode.Icon("trash", giode.Styles{Width: 18}), giode.Styles{Color: "theme.danger"}).
	Variant(giode.IconButtonGhost).
	Label("Delete").
	OnClick(func() { remove(item) })
```

Give it a `Label`: the button has no visible text, so the label is what screen readers and
automation (the [MCP server](../mcp.md)) use to identify it.

| Method | Description |
| --- | --- |
| `IconButton(icon, st...)` | Creates the button around any element, usually an icon. |
| `OnClick(fn)` | Click handler. |
| `Icon(el)` | Replaces the icon. |
| `Label(s)` | Name for automation. |
| `Variant(v)` | `IconButtonContained` (default), `IconButtonOutlined` or `IconButtonGhost`. |
| `Styles(st)` | Replaces the styles. |
| `Click()` | Clicks it from code. |

## Variants

- **Contained** fills the button with `Background` (`theme.primary` by default).
- **Outlined** is transparent with a border in `BorderColor` (or `Background`), tinted on
  hover.
- **Ghost** shows only the icon until hovered, then tints with `Color` (or `Background`).

An icon that doesn't set its own color takes the button's `Color`. Padding is 8dp on every
side by default.
