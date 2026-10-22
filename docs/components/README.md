# Components

Components keep state between frames. Create each one once, outside the view function, and
lay it out every frame. For components that belong to items of a changing list, use
[`Keyed`](../state.md#per-item-state-for-dynamic-lists).

Some entries below are stateless elements documented alongside the components (Progress,
Spinner, Image, Title, Table, icons); build those in the view like any element.

## Actions

| Component | |
| --- | --- |
| [Button](button.md) | Text button with variants, icons, disabled and loading states |
| [IconButton](iconbutton.md) | Button showing only an icon |
| [CopyButton](copybutton.md) | Copies text to the clipboard and confirms it |
| [Pressable](pressable.md) | Makes any content clickable |
| [Context menu](menu.md) | Right-click menus, dropdown menus and submenus |

## Input

| Component | |
| --- | --- |
| [Input](input.md) | Text field, checkbox, slider, select and autocomplete |

## Layout and navigation

| Component | |
| --- | --- |
| [Scroll](scroll.md) | Scrolling container |
| [VirtualList](virtuallist.md) | List that builds only the visible rows |
| [TreeView](treeview.md) | Expandable tree with keyboard navigation |
| [Tabs](tabs.md) | Tab bar with one view per tab |
| [Accordion](accordion.md) | Collapsible sections |
| [Pagination](pagination.md) | Page number controls |
| [AppShell](appshell.md) | Window layout with header, sidebars and footer |
| [Table](table.md) | Rows with aligned columns |

## Overlays and feedback

| Component | |
| --- | --- |
| [Dialog](dialog.md) | Modal dialog |
| [Drawer](drawer.md) | Panel sliding in from the side |
| [Tooltip](tooltip.md) | Hint shown on hover |
| [Toaster](toaster.md) | Notifications that disappear on their own |
| [Transition](transition.md) | Animates a component in and out |
| [Progress](progress.md) | Progress bar |
| [Spinner](spinner.md) | Loading indicator |

## Window

| Component | |
| --- | --- |
| [TitleBar](titlebar.md) | Custom title bar for frameless windows |

## Content

| Component | |
| --- | --- |
| [Icons](icons.md) | Loading and drawing your own icons |
| [Image](image.md) | Images with fit modes |
| [Title](title.md) | Headings H1 to H6 |

## Common patterns

- Styles come last and are optional: `giode.Button("Go", giode.Styles{...})`.
- Setters return the component, so they chain: `giode.Button("Go").Icon(i).OnClick(f)`.
- Handlers are ordinary Go functions and run on the [UI goroutine](../threading.md).
- Unstyled components use the theme colors, so they follow light and dark mode.
