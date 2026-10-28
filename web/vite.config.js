import { defineConfig } from "vite";
import fs from "node:fs";
import path from "node:path";

const docsDir = path.resolve(import.meta.dirname, "../docs");

// Serves ../docs over /docs/ in dev and copies the markdown into the
// build output, so the site stays dependency-free at runtime.
function docsPlugin() {
	return {
		name: "giode-docs",
		configureServer(server) {
			server.middlewares.use("/docs/", (req, res, next) => {
				// Connect strips the mount prefix, so req.url is
				// already relative to ../docs.
				const p = path.join(docsDir, decodeURIComponent(req.url));
				if (!p.startsWith(docsDir)) {
					res.statusCode = 403;
					res.end();
					return;
				}
				if (!p.endsWith(".md")) {
					next();
					return;
				}
				fs.readFile(p, (err, data) => {
					if (err) {
						res.statusCode = 404;
						res.end("not found");
						return;
					}
					res.setHeader("Content-Type", "text/markdown; charset=utf-8");
					res.end(data);
				});
			});
		},
		closeBundle() {
			const out = path.resolve(import.meta.dirname, "dist/docs");
			fs.rmSync(out, { recursive: true, force: true });
			fs.mkdirSync(out, { recursive: true });
			copyMd(docsDir, out);
		},
	};
}

// snippetsPlugin serves "virtual:snippets": the samples in
// snippets/home/main.go, by the name in the "// NAME" comment before
// each, so the home page shows code that compiles.
function snippetsPlugin() {
	const id = "virtual:snippets";
	const file = path.resolve(import.meta.dirname, "snippets/home/main.go");
	return {
		name: "giode-snippets",
		resolveId(src) {
			return src === id ? "\0" + id : null;
		},
		load(src) {
			if (src !== "\0" + id) return null;
			this.addWatchFile(file);
			const text = fs.readFileSync(file, "utf8");
			const out = {};
			const parts = text.split(/^\/\/ ([A-Z]+)\n/m);
			for (let i = 1; i < parts.length; i += 2) {
				out[parts[i].toLowerCase()] = parts[i + 1].replace(/\nvar _ = [\s\S]*$/, "").replace(/\nfunc fetchUser[\s\S]*$/, "").trim();
			}
			return "export default " + JSON.stringify(out) + ";";
		},
	};
}

function copyMd(src, dst) {
	for (const e of fs.readdirSync(src, { withFileTypes: true })) {
		const s = path.join(src, e.name);
		const d = path.join(dst, e.name);
		if (e.isDirectory()) {
			fs.mkdirSync(d, { recursive: true });
			copyMd(s, d);
		} else if (e.name.endsWith(".md")) {
			fs.copyFileSync(s, d);
		}
	}
}

export default defineConfig({
	// Relative asset URLs so the build works when served from a
	// subpath, as on GitHub Pages.
	base: "./",
	plugins: [docsPlugin(), snippetsPlugin()],
});
