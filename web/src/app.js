"use strict";

import "./style.css";

// ---------------------------------------------------------------------------
// Giode documentation site. Vanilla JS, no runtime dependencies.
// Markdown files are fetched from /docs/ and rendered client-side;
// navigation uses the URL hash so every page is linkable.
// ---------------------------------------------------------------------------

// --- Navigation -----------------------------------------------------------

const NAV = [
	{
		title: "Introduction",
		items: [
			{ path: "index", label: "Overview" },
			{ path: "getting-started", label: "Getting started" },
		],
	},
	{
		title: "Guides",
		items: [
			{ path: "window", label: "Window configuration" },
			{ path: "elements", label: "Elements" },
			{ path: "styles", label: "Style properties" },
			{ path: "routing", label: "Routing" },
			{ path: "extending", label: "Extending the library" },
			{ path: "architecture", label: "Architecture" },
			{ path: "low-level", label: "Low-level access" },
		],
	},
	{
		title: "Components",
		items: [
			{ path: "components", label: "Overview" },
			{ path: "components/button", label: "Button" },
			{ path: "components/checkbox", label: "Checkbox" },
			{ path: "components/dialog", label: "Dialog" },
			{ path: "components/drawer", label: "Drawer" },
			{ path: "components/dropdown", label: "Dropdown" },
			{ path: "components/icons", label: "Icons" },
			{ path: "components/image", label: "Image" },
			{ path: "components/input", label: "Input" },
			{ path: "components/menu", label: "Context menu" },
			{ path: "components/progress", label: "Progress" },
			{ path: "components/scroll", label: "Scroll" },
			{ path: "components/slider", label: "Slider" },
			{ path: "components/spinner", label: "Spinner" },
			{ path: "components/tabs", label: "Tabs" },
			{ path: "components/title", label: "Title" },
		],
	},
];

// --- Markdown renderer -----------------------------------------------------

function escapeHtml(s) {
	return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
}

// Inline formatting: code, links, bold, italic. The input is escaped
// first, so the tags added here are the only HTML in the output.
function inline(s) {
	s = escapeHtml(s);
	s = s.replace(/`([^`]+)`/g, "<code>$1</code>");
	s = s.replace(/\[([^\]]+)\]\(([^)]+)\)/g, (_, t, u) => `<a href="${u}">${t}</a>`);
	s = s.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
	s = s.replace(/(^|[^*])\*([^*]+)\*(?!\*)/g, "$1<em>$2</em>");
	return s;
}

function slug(s) {
	return s.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
}

function isBlockStart(l) {
	l = l.trim();
	return /^(#{1,6}\s|```)/.test(l) || /^([-*]|\d+\.)\s/.test(l) ||
		/^\s*(-{3,}|\*{3,})\s*$/.test(l);
}

function parseRow(line) {
	return line.split("|").map((c) => c.trim()).filter((c, i, a) => {
		if (c === "" && (i === 0 || i === a.length - 1)) return false;
		return true;
	});
}

function renderTable(header, rows) {
	const head = header.map((c) => `<th>${inline(c)}</th>`).join("");
	const body = rows.map((r) => {
		const pad = r.length < header.length ? new Array(header.length - r.length).fill("") : [];
		return "<tr>" + r.concat(pad).map((c) => `<td>${inline(c)}</td>`).join("") + "</tr>";
	}).join("");
	return `<div class="tablewrap"><table><thead><tr>${head}</tr></thead><tbody>${body}</tbody></table></div>`;
}

const codeStore = [];

function renderCodeBlock(code, lang) {
	const idx = codeStore.push(code) - 1;
	const label = lang ? `<span class="code-lang">${escapeHtml(lang)}</span>` : "";
	return `<div class="codeblock">${label}<button class="copy" data-code="${idx}">Copy</button>` +
		`<pre><code>${escapeHtml(code)}</code></pre></div>`;
}

function mdToHtml(text) {
	const lines = text.split(/\r?\n/);
	const out = [];
	let i = 0;
	while (i < lines.length) {
		const line = lines[i];
		const t = line.trim();

		// Fenced code block.
		if (/^```/.test(t)) {
			const lang = t.slice(3).trim();
			const buf = [];
			i++;
			while (i < lines.length && !/^```/.test(lines[i].trim())) {
				buf.push(lines[i]);
				i++;
			}
			i++;
			out.push(renderCodeBlock(buf.join("\n"), lang));
			continue;
		}

		// Heading.
		const h = line.match(/^(#{1,6})\s+(.*)$/);
		if (h) {
			const level = h[1].length;
			const text = h[2].trim();
			out.push(`<h${level} id="${slug(text)}">${inline(text)}</h${level}>`);
			i++;
			continue;
		}

		// Horizontal rule.
		if (/^\s*(-{3,}|\*{3,})\s*$/.test(line)) {
			out.push("<hr>");
			i++;
			continue;
		}

		// Table: the current line contains pipes and the next line is a
		// separator row.
		if (line.includes("|") && i + 1 < lines.length && /^\s*\|?[\s:|-]+\|[\s:|-]*$/.test(lines[i + 1].trim())) {
			const header = parseRow(line);
			i += 2;
			const rows = [];
			while (i < lines.length && lines[i].includes("|")) {
				rows.push(parseRow(lines[i]));
				i++;
			}
			out.push(renderTable(header, rows));
			continue;
		}

		// List.
		const li = line.match(/^\s*([-*]|\d+\.)\s+(.*)$/);
		if (li) {
			const ordered = /^\d+\./.test(li[1]);
			const items = [];
			while (i < lines.length) {
				const m = lines[i].match(/^\s*([-*]|\d+\.)\s+(.*)$/);
				if (!m) break;
				items.push(`<li>${inline(m[2].trim())}</li>`);
				i++;
			}
			out.push(`<${ordered ? "ol" : "ul"}>${items.join("")}</${ordered ? "ol" : "ul"}>`);
			continue;
		}

		// Paragraph.
		if (t !== "") {
			const buf = [];
			while (i < lines.length && lines[i].trim() !== "" && !isBlockStart(lines[i])) {
				buf.push(lines[i].trim());
				i++;
			}
			out.push(`<p>${inline(buf.join(" "))}</p>`);
			continue;
		}

		i++;
	}
	return out.join("\n");
}

// --- Routing ---------------------------------------------------------------

function currentRoute() {
	const h = decodeURIComponent(location.hash.replace(/^#\/?/, ""));
	const p = h.replace(/^\/+|\/+$/g, "").replace(/\.md$/, "");
	return p || "index";
}

// Resolve a markdown link relative to the current page into a route.
function routeForHref(href, current) {
	if (!href || /^https?:|^mailto:|^#/.test(href)) return null;
	let t = href.replace(/\.md$/, "").replace(/\/+$/, "");
	if (t.startsWith("../")) t = t.replace(/^(\.\.\/)+/, "");
	if (t.startsWith("./")) t = t.slice(2);
	if (t === "") return current;
	const dir = current.includes("/") ? current.slice(0, current.lastIndexOf("/")) : "";
	if (!t.startsWith("/")) t = dir ? dir + "/" + t : t;
	t = t.replace(/\/+$/, "");
	return t || "index";
}

// --- Loading ---------------------------------------------------------------

const docCache = new Map();

async function fetchDoc(path) {
	if (docCache.has(path)) return docCache.get(path);
	const candidates = ["/docs/" + path + ".md", "/docs/" + path + "/README.md"];
	let text = null;
	for (const c of candidates) {
		try {
			const r = await fetch(c);
			if (r.ok) {
				text = await r.text();
				break;
			}
		} catch (e) {
			// try the next candidate
		}
	}
	docCache.set(path, text);
	return text;
}

function pageTitle(text) {
	const m = text.match(/^#\s+(.*)$/m);
	return m ? m[1].trim() : "Documentation";
}

// --- Rendering -------------------------------------------------------------

const content = document.getElementById("content");

function wireLinks(root, route) {
	root.querySelectorAll("a").forEach((a) => {
		const target = routeForHref(a.getAttribute("href") || "", route);
		if (target !== null) {
			a.setAttribute("href", "#/" + target);
		} else {
			a.setAttribute("target", "_blank");
			a.setAttribute("rel", "noopener");
		}
	});
}

function wireCopy(root) {
	root.querySelectorAll(".copy").forEach((btn) => {
		btn.addEventListener("click", () => {
			const code = codeStore[Number(btn.dataset.code)] || "";
			navigator.clipboard.writeText(code).then(() => {
				const old = btn.textContent;
				btn.textContent = "Copied!";
				setTimeout(() => { btn.textContent = old; }, 1200);
			});
		});
	});
}

async function render() {
	const route = currentRoute();
	const text = await fetchDoc(route);
	if (text === null) {
		content.innerHTML = `<div class="missing"><h1>Page not found</h1>` +
			`<p>No documentation for <code>${escapeHtml(route)}</code>. <a href="#/">Back to the overview</a>.</p></div>`;
		document.title = "Not found — Giode";
	} else {
		content.innerHTML = mdToHtml(text);
		document.title = pageTitle(text) + " — Giode";
		wireLinks(content, route);
		wireCopy(content);
	}
	buildNav(route);
	closeMenu();
	window.scrollTo(0, 0);
}

function buildNav(route) {
	document.getElementById("nav").innerHTML = NAV.map((section) =>
		`<div class="nav-section">${section.title}</div>` +
		section.items.map((it) =>
			`<a class="nav-item${it.path === route ? " active" : ""}" href="#/${it.path}">${it.label}</a>`,
		).join(""),
	).join("");
}

// --- Search ----------------------------------------------------------------

const searchInput = document.getElementById("search");
const resultsEl = document.getElementById("results");
let searchTimer = null;
let searchPromise = null;

async function buildSearchIndex() {
	const jobs = NAV.flatMap((s) => s.items).map((it) =>
		fetchDoc(it.path).then((text) => ({ path: it.path, label: it.label, text })),
	);
	return Promise.all(jobs);
}

function snippet(text, q) {
	if (!text) return "";
	const idx = text.toLowerCase().indexOf(q.toLowerCase());
	if (idx < 0) return "";
	const s = Math.max(0, idx - 30);
	const e = Math.min(text.length, idx + 80);
	return (s > 0 ? "…" : "") + text.slice(s, e).replace(/\s+/g, " ") + (e < text.length ? "…" : "");
}

function showResults(rows, q) {
	if (rows.length === 0) {
		resultsEl.innerHTML = `<div class="result-empty">No results for “${escapeHtml(q)}”.</div>`;
	} else {
		resultsEl.innerHTML = rows.map((r) => {
			const snip = snippet(r.text, q.split(/\s+/)[0]);
			return `<a class="result" href="#/${r.path}">` +
				`<span class="result-title">${r.label}</span>` +
				(snip ? `<span class="result-snippet">${escapeHtml(snip)}</span>` : "") +
				"</a>";
		}).join("");
	}
	resultsEl.classList.add("open");
}

function hideResults() {
	resultsEl.classList.remove("open");
}

searchInput.addEventListener("input", () => {
	clearTimeout(searchTimer);
	const q = searchInput.value.trim();
	if (!q) {
		hideResults();
		return;
	}
	searchTimer = setTimeout(async () => {
		if (!searchPromise) searchPromise = buildSearchIndex();
		const docs = await searchPromise;
		const terms = q.toLowerCase().split(/\s+/);
		const scored = [];
		for (const d of docs) {
			const hay = (d.text || "").toLowerCase();
			let score = 0;
			let ok = true;
			for (const term of terms) {
				if (d.label.toLowerCase().includes(term)) {
					score += 3;
				} else if (hay.includes(term)) {
					score += 1;
				} else {
					ok = false;
					break;
				}
			}
			if (ok && score > 0) scored.push({ path: d.path, label: d.label, text: d.text, score });
		}
		scored.sort((a, b) => b.score - a.score);
		showResults(scored.slice(0, 8), q);
	}, 150);
});

searchInput.addEventListener("keydown", (e) => {
	if (e.key === "Escape") {
		searchInput.blur();
		hideResults();
	}
});

document.addEventListener("click", (e) => {
	if (!e.target.closest(".search")) hideResults();
});

document.addEventListener("keydown", (e) => {
	if (e.key === "/" && document.activeElement !== searchInput &&
		!["INPUT", "TEXTAREA"].includes(document.activeElement.tagName)) {
		e.preventDefault();
		searchInput.focus();
	}
});

// --- Mobile menu -----------------------------------------------------------

const menuBtn = document.getElementById("menu-btn");
const sidebar = document.getElementById("sidebar");
const overlay = document.getElementById("overlay");

function openMenu() {
	sidebar.classList.add("open");
	overlay.classList.add("open");
}

function closeMenu() {
	sidebar.classList.remove("open");
	overlay.classList.remove("open");
}

menuBtn.addEventListener("click", () => {
	sidebar.classList.contains("open") ? closeMenu() : openMenu();
});

overlay.addEventListener("click", closeMenu);

// --- Boot ------------------------------------------------------------------

window.addEventListener("hashchange", render);
render();
