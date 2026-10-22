# Documentation source

This folder holds the documentation: one Markdown file per page, with component pages in
`components/`. The website in [`../web`](../web) renders these files; see its README to run
it locally.

## Adding a page

1. Create `docs/<name>.md` (or `docs/<section>/<name>.md`).
2. Add it to the `NAV` array in `web/src/app.js`.

Links between pages are relative to the current file, for example `[Button](button.md)`
from inside `components/`.

## Supported Markdown

Headings (each gets an anchor), fenced code blocks with a language, tables, ordered and
unordered lists, horizontal rules, links, inline code, bold and italic.
