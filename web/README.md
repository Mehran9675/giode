# Documentation site

The website for the Giode documentation. It renders the Markdown files in
[`../docs`](../docs) in the browser, with no framework and no runtime dependencies.

## Run it

```sh
npm install
npm run dev       # http://localhost:5173
```

Production build:

```sh
npm run build     # writes dist/, including a copy of the docs
npm run preview   # serves dist/ locally
```

## How it works

- `index.html` is the page shell: sidebar, search and content area.
- `src/app.js` fetches `docs/*.md`, renders it with a small built-in Markdown renderer,
  and routes with the URL hash (`#/styles`, `#/components/menu`), so every page and heading
  has a link you can share. It also builds the sidebar, search and copy buttons for code
  blocks.
- `src/style.css` is the theme.
- `vite.config.js` serves `../docs` at `/docs/` during development and copies it into
  `dist/docs` when building.

## Deploying to GitHub Pages

`.github/workflows/deploy-pages.yml` builds the site and deploys it when a release is
published, or when the workflow is run by hand. To enable it once:

1. In the repository, open **Settings → Pages** and set **Source** to **GitHub Actions**.
2. Publish a release. The site appears at `https://<user>.github.io/giode/` a minute or two
   later.

## Adding a page

1. Create the Markdown file in `../docs`.
2. Add it to the `NAV` array in `src/app.js`. Pages missing there don't show in the sidebar
   or in search.
