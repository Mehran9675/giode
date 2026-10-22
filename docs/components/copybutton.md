# CopyButton

A small button that copies text to the clipboard, then shows a check mark for a moment.

```go
copyKey := giode.CopyButton(func() string { return apiKey.Text() }).
	Labels("Copy key", "Copied").
	OnCopy(func(string) { toasts.Show(giode.ToastInfo, "Key copied") })
```

The function is called when the button is clicked, so it always copies the current value.

| Method | Description |
| --- | --- |
| `CopyButton(text, st...)` | Creates the button. |
| `Labels(label, copied)` | Text before and after copying ("Copy" and "Copied" by default). Pass empty strings for an icon-only button. |
| `Feedback(d)` | How long the check mark shows, 1.5s by default. |
| `OnCopy(fn)` | Called with the copied text. |
| `Copied() bool` | Whether the check mark is showing. |
| `Click()` | Copies from code. |

Styles: `Color` (icon and text, `theme.text` by default), `Background`, `Padding`,
`BorderRadius`, `FontSize`.

To copy without a button, use `app.SetClipboard(text)`; see
[Window configuration](../window.md#clipboard).
