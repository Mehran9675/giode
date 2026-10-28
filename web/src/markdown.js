// A small Markdown renderer for the docs: headings (with anchors),
// paragraphs, lists, tables, block quotes (as callouts), fenced code with
// syntax highlighting and tab groups, and inline code, links, bold and
// italic. Input is escaped first, so the tags added here are the only
// HTML in the output.

export function escapeHtml(s) {
	return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
}

export function slug(s) {
	return s.toLowerCase().replace(/<[^>]+>/g, "").replace(/[`*]/g, "").replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
}

function inline(s) {
	s = escapeHtml(s);
	const codes = [];
	s = s.replace(/`([^`]+)`/g, (_, c) => {
		codes.push(c);
		return `\u0000${codes.length - 1}\u0000`;
	});
	s = s.replace(/\[([^\]]+)\]\(([^)\s]+)\)/g, (_, t, u) => `<a href="${u}">${t}</a>`);
	s = s.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
	s = s.replace(/(^|[^*\w])\*([^*\s][^*]*)\*(?!\*)/g, "$1<em>$2</em>");
	s = s.replace(/\u0000(\d+)\u0000/g, (_, i) => `<code>${codes[Number(i)]}</code>`);
	return s;
}

// --- Syntax highlighting -----------------------------------------------------

const GO_KEYWORDS = new Set(("break case chan const continue default defer else fallthrough for func go goto if " +
	"import interface map package range return select struct switch type var").split(" "));
const GO_BUILTINS = new Set(("append bool byte cap close complex copy delete error false float32 float64 imag int " +
	"int8 int16 int32 int64 iota len make max min new nil panic print println real recover rune string true uint " +
	"uint8 uint16 uint32 uint64 uintptr any").split(" "));
const SH_KEYWORDS = new Set("if then else fi for do done case esac in export echo cd sudo".split(" "));

// commandPosition reports whether the word at i starts a shell command:
// first on its line, or after a pipe, &&, ; or $(.
function commandPosition(code, i) {
	let j = i - 1;
	while (j >= 0 && (code[j] === " " || code[j] === "\t")) j--;
	return j < 0 || code[j] === "\n" || code[j] === "|" || code[j] === "&" || code[j] === ";" || code[j] === "(";
}

function span(cls, text) {
	return cls ? `<span class="tk-${cls}">${escapeHtml(text)}</span>` : escapeHtml(text);
}

// highlight tokenizes code with a few regular expressions per language;
// good enough for docs, and free of dependencies.
export function highlight(code, lang) {
	lang = (lang || "").toLowerCase();
	let rules;
	if (lang === "go") {
		rules = [
			[/^\/\/[^\n]*/, "comment"],
			[/^\/\*[\s\S]*?\*\//, "comment"],
			[/^"(?:\\.|[^"\\\n])*"/, "string"],
			[/^`[^`]*`/, "string"],
			[/^'(?:\\.|[^'\\\n])+'/, "string"],
			[/^\d[\d_.xXa-fA-F]*/, "number"],
			[/^[A-Za-z_]\w*(?=\s*\()/, (w) => GO_KEYWORDS.has(w) ? "keyword" : GO_BUILTINS.has(w) ? "builtin" : "func"],
			[/^[A-Za-z_]\w*/, (w) => GO_KEYWORDS.has(w) ? "keyword" : GO_BUILTINS.has(w) ? "builtin" : /^[A-Z]/.test(w) ? "type" : ""],
		];
	} else if (["sh", "bash", "shell", "powershell", "ps1", "console"].includes(lang)) {
		rules = [
			[/^#[^\n]*/, "comment"],
			[/^"(?:\\.|[^"\\])*"/, "string"],
			[/^'[^']*'/, "string"],
			[/^\$\(?[A-Za-z_][\w]*\)?/, "builtin"],
			[/^--?[A-Za-z][\w-]*/, "flag"],
			[/^[A-Za-z_][\w.-]*/, (w, at, code) => SH_KEYWORDS.has(w) ? "keyword" : commandPosition(code, at) ? "func" : ""],
		];
	} else if (lang === "toml" || lang === "ini") {
		rules = [
			[/^#[^\n]*/, "comment"],
			[/^\[\[?[^\]\n]+\]\]?/, "type"],
			[/^"(?:\\.|[^"\\])*"/, "string"],
			[/^'[^']*'/, "string"],
			[/^(true|false)\b/, "builtin"],
			[/^-?\d[\d_.]*/, "number"],
			[/^[A-Za-z_][\w.-]*(?=\s*=)/, "func"],
		];
	} else if (["json", "yaml", "yml", "xml", "html"].includes(lang)) {
		rules = [
			[/^#[^\n]*/, lang === "yaml" || lang === "yml" ? "comment" : ""],
			[/^"(?:\\.|[^"\\])*"/, "string"],
			[/^-?\d[\d.]*/, "number"],
			[/^(true|false|null)\b/, "builtin"],
			[/^<\/?[\w:-]+/, "keyword"],
		];
	} else {
		return escapeHtml(code);
	}
	let out = "";
	let i = 0;
	let plain = "";
	const flush = () => {
		out += escapeHtml(plain);
		plain = "";
	};
	outer: while (i < code.length) {
		const rest = code.slice(i);
		for (const [re, cls] of rules) {
			// A rule may only match at a word boundary, so identifiers
			// inside other words aren't split.
			if (/\w/.test(code[i - 1] || "") && /^\w/.test(rest)) break;
			const m = rest.match(re);
			if (m && m[0].length > 0) {
				flush();
				const c = typeof cls === "function" ? cls(m[0], i, code) : cls;
				out += span(c, m[0]);
				i += m[0].length;
				continue outer;
			}
		}
		plain += code[i];
		i++;
	}
	flush();
	return out;
}

// --- Blocks ---------------------------------------------------------------

let codeStore = [];

// codeFor returns the source of a rendered code block, for its copy
// button.
export function codeFor(i) {
	return codeStore[i] || "";
}

function codeBlock(code, lang, title) {
	const idx = codeStore.push(code) - 1;
	const label = title || lang;
	return `<div class="codeblock">` +
		`<div class="code-head"><span class="code-lang">${escapeHtml(label || "")}</span>` +
		`<button class="copy" data-code="${idx}" aria-label="Copy code">Copy</button></div>` +
		`<pre><code class="lang-${escapeHtml(lang || "text")}">${highlight(code, lang)}</code></pre></div>`;
}

let tabGroup = 0;

function tabs(blocks) {
	const g = tabGroup++;
	const head = blocks.map((b, i) =>
		`<button class="tab${i === 0 ? " active" : ""}" data-group="${g}" data-tab="${i}">${escapeHtml(b.tab)}</button>`).join("");
	const body = blocks.map((b, i) =>
		`<div class="tab-panel${i === 0 ? " active" : ""}" data-group="${g}" data-tab="${i}" data-label="${escapeHtml(b.tab)}">` +
		codeBlock(b.code, b.lang) + `</div>`).join("");
	return `<div class="tabs"><div class="tab-head">${head}</div>${body}</div>`;
}

function isBlockStart(l) {
	const t = l.trim();
	return /^(#{1,6}\s|```|>)/.test(t) || /^([-*]|\d+\.)\s/.test(t) || /^(-{3,}|\*{3,})$/.test(t);
}

function parseRow(line) {
	const cells = line.trim().replace(/^\|/, "").replace(/\|$/, "").split(/(?<!\\)\|/);
	return cells.map((c) => c.trim().replace(/\\\|/g, "|"));
}

function table(header, rows) {
	const head = header.map((c) => `<th>${inline(c)}</th>`).join("");
	const body = rows.map((r) => {
		const cells = r.concat(new Array(Math.max(0, header.length - r.length)).fill(""));
		return "<tr>" + cells.map((c) => `<td>${inline(c)}</td>`).join("") + "</tr>";
	}).join("");
	return `<div class="tablewrap"><table><thead><tr>${head}</tr></thead><tbody>${body}</tbody></table></div>`;
}

// parseFence reads a fence's info string: ```go, ```sh tab="macOS".
function parseFence(info) {
	const lang = (info.match(/^([\w+-]+)/) || [])[1] || "";
	const tab = (info.match(/tab="([^"]+)"/) || [])[1] || "";
	const title = (info.match(/title="([^"]+)"/) || [])[1] || "";
	return { lang, tab, title };
}

// render turns Markdown into HTML, and returns the h2/h3 headings for a
// table of contents.
export function render(text) {
	codeStore = [];
	return blocks(text);
}

function blocks(text) {
	const lines = text.split(/\r?\n/);
	const out = [];
	const toc = [];
	let pendingTabs = [];
	const flushTabs = () => {
		if (pendingTabs.length === 1) {
			out.push(codeBlock(pendingTabs[0].code, pendingTabs[0].lang, pendingTabs[0].tab));
		} else if (pendingTabs.length > 1) {
			out.push(tabs(pendingTabs));
		}
		pendingTabs = [];
	};
	let i = 0;
	while (i < lines.length) {
		const line = lines[i];
		const t = line.trim();

		if (/^```/.test(t)) {
			const f = parseFence(t.slice(3).trim());
			const buf = [];
			i++;
			while (i < lines.length && !/^```\s*$/.test(lines[i].trim())) {
				buf.push(lines[i]);
				i++;
			}
			i++;
			const code = buf.join("\n");
			if (f.tab) {
				pendingTabs.push({ code, lang: f.lang, tab: f.tab });
				continue;
			}
			flushTabs();
			out.push(codeBlock(code, f.lang, f.title));
			continue;
		}
		if (t === "") {
			i++;
			continue;
		}
		flushTabs();

		const h = line.match(/^(#{1,6})\s+(.*)$/);
		if (h) {
			const level = h[1].length;
			const raw = h[2].trim();
			const id = slug(raw);
			if (level === 2 || level === 3) toc.push({ level, id, text: raw.replace(/`/g, "") });
			const anchor = level > 1 ? `<a class="anchor" href="#${id}" aria-label="Link to this section">#</a>` : "";
			out.push(`<h${level} id="${id}">${inline(raw)}${anchor}</h${level}>`);
			i++;
			continue;
		}
		if (/^(-{3,}|\*{3,})$/.test(t)) {
			out.push("<hr>");
			i++;
			continue;
		}
		if (t.startsWith(">")) {
			const buf = [];
			while (i < lines.length && lines[i].trim().startsWith(">")) {
				buf.push(lines[i].trim().replace(/^>\s?/, ""));
				i++;
			}
			const first = buf[0] || "";
			const kind = /^\*\*(note|warning|tip|alpha)\b/i.test(first) ? first.match(/^\*\*(\w+)/)[1].toLowerCase()
				: /\balpha\b/i.test(first) ? "alpha" : "note";
			out.push(`<blockquote class="callout callout-${kind}">${blocks(buf.join("\n")).html}</blockquote>`);
			continue;
		}
		if (line.includes("|") && i + 1 < lines.length && /^\s*\|?\s*:?-{2,}/.test(lines[i + 1])) {
			const header = parseRow(line);
			i += 2;
			const rows = [];
			while (i < lines.length && lines[i].includes("|") && lines[i].trim() !== "") {
				rows.push(parseRow(lines[i]));
				i++;
			}
			out.push(table(header, rows));
			continue;
		}
		const li = line.match(/^(\s*)([-*]|\d+\.)\s+(.*)$/);
		if (li) {
			const ordered = /^\d+\./.test(li[2]);
			const items = [];
			while (i < lines.length) {
				const m = lines[i].match(/^\s*([-*]|\d+\.)\s+(.*)$/);
				if (m) {
					items.push(m[2].trim());
					i++;
				} else if (lines[i].trim() !== "" && /^\s{2,}/.test(lines[i]) && items.length) {
					items[items.length - 1] += " " + lines[i].trim();
					i++;
				} else {
					break;
				}
			}
			const tag = ordered ? "ol" : "ul";
			out.push(`<${tag}>${items.map((it) => `<li>${inline(it)}</li>`).join("")}</${tag}>`);
			continue;
		}
		const buf = [];
		while (i < lines.length && lines[i].trim() !== "" && !isBlockStart(lines[i])) {
			buf.push(lines[i].trim());
			i++;
		}
		out.push(`<p>${inline(buf.join(" "))}</p>`);
	}
	flushTabs();
	return { html: out.join("\n"), toc };
}
