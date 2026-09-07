# Dropdown

Select one option from a list.

```go
size := giode.NewSelect("Small", "Medium", "Large").
	OnChange(func(i int) { ... })

// in the view:
giode.Box(giode.Styles{}, size)
```

## API

| Member | Description |
| --- | --- |
| `NewSelect(options ...string) *Select` | Creates the dropdown. |
| `Selected() int` | Selected index. |
| `Value() string` | Selected option text. |
| `SetSelected(i int)` | Changes the selection without invoking OnChange. |
| `OnChange(fn func(index int)) *Select` | Change handler. |
| `Styles(st Styles) *Select` | Styles. |
| `Opened() bool` | Whether the panel is open. |

## Styles

| Property | Role |
| --- | --- |
| `Background` | Field and panel fill. |
| `Color` | Text. |
| `BorderColor` | Field border. |
| `BorderRadius` | Corners. |
| `Width` | Fixed field width. |

The panel opens below the field, clamped to the bounds of the select's parent container.
Note: as in all immediate-mode toolkits, siblings laid out after the select paint above an
open panel.
