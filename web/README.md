# Documentation site (Vite)

The front-end for the Giode documentation. It renders the Markdown in [`../docs`](../docs)
— no runtime dependencies, no framework.

## Run it

```sh
npm install
npm run dev       # http://localhost:5173
```

Production build:

```sh
npm run build     # outputs dist/, with the docs copied in
npm run preview   # serves dist/ locally
```

## How it works

- `index.html` — the page shell (sidebar, search, content).
- `src/app.js` — fetches `/docs/*.md`, renders them with a small built-in Markdown
  renderer, and routes via the URL hash (`#/styles`, `#/components/menu`).
- `src/style.css` — the site theme.
- `vite.config.js` — serves `../docs` at `/docs/` in dev and copies the Markdown into
  `dist/docs` on build (both handled by a tiny inline plugin).

Every page has a stable URL: share links directly. The sidebar, mobile menu, code-block
copy buttons and full-text search are included.

## Adding a page

1. Create the Markdown file in `../docs`.
2. Add an entry to the `NAV` array in `src/app.js`.
