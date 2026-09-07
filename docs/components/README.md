# Components

Stateful components are created once and laid out every frame. They are built on the
stateless elements.

| Component | Docs | State |
| --- | --- | --- |
| Button | [button.md](button.md) | clickable |
| Checkbox | [checkbox.md](checkbox.md) | bool |
| Dialog | [dialog.md](dialog.md) | open/closed |
| Drawer | [drawer.md](drawer.md) | open/closed, animated |
| Dropdown | [dropdown.md](dropdown.md) | selection, open |
| Icons | [icons.md](icons.md) | stateless |
| Image | [image.md](image.md) | stateless |
| Input | [input.md](input.md) | text |
| Menu | [menu.md](menu.md) | items, open |
| Progress | [progress.md](progress.md) | stateless |
| Scroll | [scroll.md](scroll.md) | offset |
| Slider | [slider.md](slider.md) | value |
| Spinner | [spinner.md](spinner.md) | stateless, animated |
| Tabs | [tabs.md](tabs.md) | selection |
| Title | [title.md](title.md) | stateless |

Common patterns:

- **Chaining**: constructors and setters return the component, so handles configure in one
  expression.
- **Callbacks**: handlers are plain Go closures and may capture anything.
- **Styling**: every component takes `Styles(...)`; see the individual docs for which
  properties are honored.
