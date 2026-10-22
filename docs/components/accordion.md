# Accordion

A list of sections that expand and collapse when their header is clicked.

```go
faq := giode.Accordion().
	Item("What is Giode?", func() giode.Element { return giode.Text(aboutText) }).
	Item("Is it free?", func() giode.Element { return giode.Text(licenseText) }).
	SingleOpen(true)
```

A section's content function only runs while the section is open.

| Method | Description |
| --- | --- |
| `Accordion(st...) *Accordion` | Creates an empty accordion. |
| `Item(label, content)` | Adds a section. Sections can be added at any time. |
| `SingleOpen(true)` | Opening a section closes the others. By default sections open independently. |
| `Expanded(i) bool`, `SetExpanded(i, v)` | Read or change a section's state. |

Styles: `Color` (headers and arrows), `BorderColor` (lines between sections). Opening and
closing aren't animated; wrap the content in a [Transition](transition.md) with `Collapse`
if you want that.
