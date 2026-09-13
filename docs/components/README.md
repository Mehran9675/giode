# Components

Stateful components are created once and laid out every frame. They are built on the
stateless elements.

| Component | Constructor | State |
| --- | --- | --- |
| Button | [button.md](button.md) | clickable |
| Input (text field) | [input.md](input.md) | text |
| Input (checkbox) | [input.md](input.md) | bool |
| Input (slider) | [input.md](input.md) | value |
| Input (select) | [input.md](input.md) | selection |
| Tabs | [tabs.md](tabs.md) | selection |
| Scroll | [scroll.md](scroll.md) | offset |
| Drawer | [drawer.md](drawer.md) | external `*bool` |
| Dialog | [dialog.md](dialog.md) | external `*bool` |
| Menu | [menu.md](menu.md) | items, open |
| Progress | [progress.md](progress.md) | stateless |
| Spinner | [spinner.md](spinner.md) | stateless, animated |
| Image | [image.md](image.md) | stateless |
| Icons | [icons.md](icons.md) | stateless |
| Title | [title.md](title.md) | stateless |

Common patterns:

- **Constructors take styles last**: `giode.Button("Go", giode.Styles{...})`,
  `giode.Input.Text("Name", giode.Styles{...})`, `giode.Scroll(content, styles)`.
- **Callbacks**: handlers are plain Go closures and may capture anything; chainable
  setters return the component.
- **Styling**: every component accepts `Styles`; see the individual docs for which
  properties are honored.
