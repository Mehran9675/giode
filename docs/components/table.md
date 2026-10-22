# Table

Rows of cells whose columns line up: each column is as wide as its widest cell.

```go
giode.Table(giode.Styles{},
	giode.TableHeader(giode.Styles{Background: "theme.surface"},
		giode.TableCell(giode.Styles{}, giode.Text("Name")),
		giode.TableCell(giode.Styles{}, giode.Text("Role")),
	),
	giode.TableRow(giode.Styles{},
		giode.TableCell(giode.Styles{}, giode.Text("Ada")),
		giode.TableCell(giode.Styles{}, giode.Text("Engineer")),
	),
)
```

Tables are stateless: build them in the view every frame.

| Function | Description |
| --- | --- |
| `Table(st, rows...)` | Stacks the rows and lines up their columns. |
| `TableHeader(st, cells...)` | A row in bold. |
| `TableRow(st, cells...)` | A row. |
| `TableCell(st, content)` | One cell. Padding is 12×8 by default. |

Elements other than rows, such as a `Divider` between groups, can be placed among the rows
and are left as they are.

There's no sorting, selection or virtualization. Sort your data before building the rows,
and use a [Pressable](pressable.md) or `Box.OnClick` for clickable rows. For thousands of
rows, use [VirtualList](virtuallist.md).
