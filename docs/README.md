# Documentation site

This folder (`docs/`) holds the documentation source: one Markdown file per page, with
`components/` and `cli/` for those sections. The website lives in [`../web`](../web), a
Vite app that renders these files next to its home and examples pages; see
`../web/README.md` for how to run it.

## Adding a page

1. Create `docs/<name>.md` (or `docs/<section>/<name>.md`).
2. Add an entry to `NAV` in `web/src/nav.js`.

Markdown links between pages are relative to the current file (`[button](button.md)`
inside `components/README.md`); the site resolves them, including `#section` anchors.

## Supported Markdown

Headings (with anchor IDs and a table of contents from `##` and `###`), fenced code with
highlighting for Go, shell, TOML, JSON and YAML, tables, lists, block quotes (shown as
callouts; one mentioning "alpha" gets the alpha style), horizontal rules, links, inline
code, bold and italic.

Consecutive code blocks with a `tab` label become tabs, and the one matching the reader's
system opens first:

    ```sh tab="macOS / Linux"
    curl -fsSL https://example.com/install.sh | sh
    ```

    ```powershell tab="Windows"
    irm https://example.com/install.ps1 | iex
    ```
