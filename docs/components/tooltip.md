# Tooltip

Shows a short text below an element while the pointer is over it. The tooltip draws above
everything else.

```go
save := giode.Tooltip(giode.Button("Save"), "Save changes (Ctrl+S)")

// in the view, use the tooltip in place of the button:
giode.Row(giode.Styles{}, save)
```

Create it once, like other components, since it tracks hover over its trigger.

| Method | Description |
| --- | --- |
| `Tooltip(trigger, text, st...) *Tooltip` | Wraps `trigger`. |
| `SetContent(text)` | Changes the text. |
| `Hovered() bool` | Whether the trigger is hovered. |

Styles: `Background` (`#0f172a` by default), `Color` (`#f8fafc`), `FontSize` (13).
