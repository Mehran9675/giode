# Icons

Stateless icons in the popular formats. Size comes from `Width`/`Height` (default 24px),
color from `Color` (default black).

## Material set

```go
giode.IconMaterial("home", giode.Styles{Color: giode.HexColor("#fff"), Width: 20, Height: 20})
```

Unknown names panic. List the built-in names with `icon.Names()`. The set includes:
menu, close, check, plus, minus, chevron-down/up/left/right, arrow-up/down/left/right,
play, pause, stop, skip-next, skip-previous, refresh, search, home, settings, trash, edit,
star, heart, user, folder, file, info, warning, mail, link, copy, download, upload, eye,
logout, more-horiz, more-vert, external-link.

## SVG

```go
giode.IconSVG("M3 6h18v2H3V6zm0 5h18v2H3v-2z", giode.Styles{Width: 24, Height: 24})
```

Supports the commands M, L, H, V, C, S, Q, T, A, Z in absolute or relative form, scaled
from a 24x24 viewBox. `IconSVGCustom(pathData, viewBox, st)` scales from a custom viewBox.
Arc commands (`A`) degrade to straight lines.

## IconVG

Gio's native vector format:

```go
giode.IconIconVG(data, giode.Styles{Width: 24})
```

Invalid data panics.

## Image

```go
giode.IconImage(img, giode.Styles{Width: 24, Height: 24})
```
