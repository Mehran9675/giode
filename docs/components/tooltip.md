# Tooltip

Wraps any element as a trigger: hovering it shows a small panel just below, rendered
above everything else regardless of nesting (see [extending.md](../extending.md#rendering-above-everything-portals)).

```go
trigger := giode.Button("Hover me", giode.Styles{Background: "gray"})
tip := giode.Tooltip(trigger, "A helpful hint")

// in the view, lay out tip in trigger's place:
giode.Stack(giode.Styles{}, tip)
```

Create it once with `New`/`giode.Tooltip` - it wraps a stateful trigger, so it must be
laid out every frame like any other component, not rebuilt inside the view function.

## API

| Member | Description |
| --- | --- |
| `Tooltip(trigger Element, content string, st ...Styles) *Tooltip` | Wraps trigger. |
| `SetContent(content string)` | Replaces the tooltip text. |
| `Hovered() bool` | Whether the trigger was hovered on the last frame. |

## Styles

| Property | Role |
| --- | --- |
| `Background` | Panel fill; defaults to a dark slate. |
| `Color` | Text color; defaults to near-white. |
| `FontSize` | Text size; defaults to 13. |
