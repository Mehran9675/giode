# Pagination

Previous and next buttons with page numbers in between.

```go
pager := giode.Pagination(20). // 20 pages
	OnChange(func(page int) { loadPage(page) })
```

| Method | Description |
| --- | --- |
| `Pagination(total, st...) *Pagination` | Creates the control, starting at page 1. |
| `Current() int` | Current page, starting at 1. |
| `SetCurrent(page)` | Changes the page without calling `OnChange`. Out-of-range values are clamped. |
| `SetTotal(total)` | Changes the page count. |
| `OnChange(fn)` | Called when the user picks a page. |

With up to 7 pages, every number shows. With more, it shows the first and last pages, the
pages next to the current one, and "…" for the gaps.

The styles apply to every button. The current page uses the contained
[Button](button.md) variant and the others the outlined one.
