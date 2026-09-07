# Documentation site

This folder (`docs/`) holds the documentation source: one Markdown file per page, plus
`components/` with one file per component. The website itself lives in [`../web`](../web)
— a Vite app that renders these files; see `../web/README.md` for how to run it.

## Adding a page

1. Create `docs/<name>.md` (or `docs/<section>/<name>.md`).
2. Add an entry to the `NAV` array in `web/src/app.js`.

Markdown links between pages are relative to the current file (e.g. `[button](button.md)`
inside `components/README.md`); the site resolves them automatically.

## Supported Markdown

Headings (with anchor IDs), fenced code blocks with a language tag, tables, ordered and
unordered lists, horizontal rules, links, inline code, bold and italic.
