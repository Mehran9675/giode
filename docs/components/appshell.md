# AppShell

A whole-window layout: an optional header, left navbar, right aside and footer, with the main
content filling the rest.

```go
nav := giode.Section(func() giode.Element { // scrolls if the links don't fit
	return giode.Stack(giode.Styles{Gap: 4, Padding: giode.UniformInset(12)}, links...)
})

shell := giode.AppShell(giode.AppShellConfig{
	Header: giode.Header(56, giode.Text("My app", giode.Styles{Color: "white"})),
	Navbar: giode.Navbar(220, nav),
	Main:   giode.Main(page),
})

app.Run(func() giode.Element { return shell })
```

| Function | Placement | Size |
| --- | --- | --- |
| `Header(height, content, st...)` | Top, full width | `height` dp |
| `Navbar(width, content, st...)` | Left, between header and footer | `width` dp |
| `Aside(width, content, st...)` | Right, between header and footer | `width` dp |
| `Footer(height, content, st...)` | Bottom, full width | `height` dp |
| `Main(content, st...)` | Everything left over | |

Leave out any of the first four. `AppShell`'s own styles apply to the whole window, which it
fills by default.

`Section(content, st...)` is a scrolling area for a navbar or aside. It's a
[Scroll](scroll.md), so create it once rather than in the view.

Each slot's styles apply to its own area; `Background` defaults to a dark slate.
