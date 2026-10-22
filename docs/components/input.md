# Input

The form controls are grouped under `giode.Input`:

| Constructor | Control |
| --- | --- |
| `giode.Input.Text(placeholder, st...)` | Text field |
| `giode.Input.Checkbox(label, st...)` | Checkbox, with an optional label (`""` for none) |
| `giode.Input.Slider(st...)` | Horizontal slider from 0 to 1 |
| `giode.Input.Select(st, options...)` | Dropdown list |
| `giode.Input.AutoComplete(placeholder, options, st...)` | Text field with suggestions |

Like every component, create them once and lay them out every frame.

## Text field

```go
email := giode.Input.Text("you@example.com").
	Label("Email").
	OnChange(func(text string) { validate(text) }).
	Submit(func(text string) { save(text) })
```

| Method | Description |
| --- | --- |
| `Text() string` | Current text. |
| `SetText(s)` | Replaces the text without calling `OnChange`. |
| `ChangeText(s)` | Replaces the text and calls `OnChange`, as typing would. |
| `OnChange(fn)` | Called after every edit. |
| `Submit(fn)` | Called when Enter is pressed (single-line fields). |
| `SubmitText()` | Calls the `Submit` handler with the current text. |
| `Password(true)` | Shows bullets instead of the text. |
| `Multiline(true)` | Allows several lines; Enter adds a line instead of submitting. |
| `ReadOnly(true)` | Text can be selected but not edited. |
| `MaxLength(n)` | Limits the text to `n` characters. |
| `Label(s)` | Name for automation; the placeholder is used otherwise. |
| `Focus()`, `Focused()` | Request keyboard focus; report whether it has it. |

Styles: `Color` (text, `theme.text` by default), `Background`, `BorderColor`
(`theme.border`), `BorderWidth` (1), `BorderRadius` (6), `Padding` (12×8), `FontSize`. The
field fills the available width unless `Width` is set.

## Checkbox

```go
notify := giode.Input.Checkbox("Email me about updates").
	OnChange(func(checked bool) { prefs.Notify = checked })
notify.SetChecked(true)
```

| Method | Description |
| --- | --- |
| `Checked() bool` | Current state. |
| `SetChecked(v)` | Sets the state without calling `OnChange`. |
| `Toggle()` | Flips the state and calls `OnChange`, as a click would. |
| `OnChange(fn)` | Called when the state changes. |
| `Icon(src, st...)` | Replaces the check mark with any icon source; `Icon(nil)` restores it. |

Styles: `Color` (label), `Background` (fill when checked), `BorderColor` (border when
unchecked), `BorderWidth`, `BorderRadius`, `Width`/`Height` (box size, 18dp by default).

## Slider

```go
volume := giode.Input.Slider().
	OnChange(func(v float32) { player.SetVolume(v) })
```

| Method | Description |
| --- | --- |
| `Value() float32` | Current value, 0 to 1. |
| `SetValue(v)` | Sets the value without calling `OnChange`. |
| `Change(v)` | Sets the value and calls `OnChange`. |
| `OnChange(fn)` | Called while the thumb is dragged. |
| `Dragging() bool` | Whether the user is dragging the thumb, so code syncing it to an outside value can wait. |

Styles: `Color` (filled part), `Background` (track), `BorderColor` (thumb). The slider fills
the available width.

## Select

```go
size := giode.Input.Select(giode.Styles{Width: 160}, "Small", "Medium", "Large").
	OnChange(func(i int) { order.Size = i })
```

| Method | Description |
| --- | --- |
| `Selected() int`, `Value() string` | Selected index and text. |
| `SetSelected(i)` | Selects without calling `OnChange`. |
| `Choose(i)` | Selects and calls `OnChange`, as picking from the list would. |
| `Open()`, `Close()`, `Click()` | Open or close the list from code. |
| `Options() []string` | The options. |
| `OnChange(fn)` | Called when the user picks an option. |
| `HoverColor(c)` | Background of the option under the pointer. |
| `Opened() bool` | Whether the list is open. |

The list opens just below the field (above it near the bottom of the window) and draws over
everything else. Clicking anywhere outside it closes it without changing the selection.

Keyboard: Tab to the field, then Enter, Space or an arrow key opens the list. In the open
list, Up and Down move the highlight, Enter or Space picks it, and Escape closes the list.

Styles: `Background` (field and list), `Color` (text), `BorderColor`, `BorderRadius`,
`Width` (field width, in dp).

## AutoComplete

```go
city := giode.Input.AutoComplete("City", []string{"Paris", "London", "Tokyo"}).
	OnSelect(func(value string) { trip.City = value })
```

| Method | Description |
| --- | --- |
| `Text() string` | Current text. |
| `SetText(s)` | Replaces the text without opening the suggestions, for loading a saved value. |
| `SetOptions(options)` | Replaces the suggestion list. |
| `OnSelect(fn)` | Called when a suggestion is clicked. |

Suggestions are the options containing the typed text (ignoring case), shown below the
field. They're picked by clicking; there's no keyboard selection yet. Clicking elsewhere
hides them until the text changes, and the click still goes through, so clicking back into
the field places the caret. Styles are the same as for Select.
