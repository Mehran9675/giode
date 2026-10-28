# Website (Vite)

Giode's website: the home page, the examples pages and the documentation, which renders
the Markdown in [`../docs`](../docs). Plain JavaScript, no framework, no runtime
dependencies.

## Run it

```sh
npm install
npm run dev       # http://localhost:5173
npm run build     # dist/, with the docs copied in
npm run preview   # serves dist/
```

## Layout

| File | What it is |
| --- | --- |
| `index.html` | The shell: top bar, footer, theme bootstrap. |
| `src/main.js` | The hash router (`#/`, `#/docs/<page>`, `#/examples/<slug>`), theme toggle, copy buttons and tabs. Old links like `#/styles` redirect to `#/docs/styles`. |
| `src/pages/home.js` | The home page. Its code samples come from `snippets/home/main.go`. |
| `src/pages/docs.js` | Docs: navigation, search, table of contents, previous/next, edit links. |
| `src/pages/examples.js` | The examples list and each app's page with its downloads. |
| `src/nav.js` | The docs' navigation. |
| `src/projects.js` | The apps on the examples pages. |
| `src/download.js` | Detects the visitor's system and picks the right file from the app's latest GitHub release. |
| `src/markdown.js` | The Markdown renderer and syntax highlighter. |
| `src/style.css` | The design: tokens for dark and light, then each page. |
| `public/` | The favicon and the examples' icons and screenshots. |
| `snippets/` | Go programs holding the code the site shows (home page, recipes), so `go vet ./...` keeps it compiling. |

## Adding an example app

Add an entry to `src/projects.js` and put a 256px square `icon.png` (and screenshots, as
`.webp`) in `public/examples/<slug>/`. Downloads need nothing more: the app's page reads
the repository's latest release and offers the file for the visitor's system, matching
file names (`.dmg`, `-setup.exe`, `.msix`, `.AppImage`, `.deb`, `.rpm`, `.apk`, and
archives named for a system).

## Deploy to GitHub Pages

`.github/workflows/deploy-pages.yml` builds `web/` and deploys it when a release is
published (or on a manual run). Enable it once in **Settings → Pages → Source: GitHub
Actions**; the site appears at `https://<user>.github.io/giode/`. The build uses relative
URLs, so it works from the repository's subpath.
