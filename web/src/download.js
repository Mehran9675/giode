// Downloads for the examples pages: the visitor's system, and the file of
// a repository's latest GitHub release that suits it.

// detectOS returns windows, macos, linux, android or ios, and whether the
// CPU looks like ARM.
export function detectOS() {
	const ua = navigator.userAgent || "";
	const platform = (navigator.userAgentData && navigator.userAgentData.platform) || navigator.platform || "";
	const arm = /arm|aarch64/i.test(ua) || /arm/i.test(platform);
	if (/android/i.test(ua)) return { os: "android", arm };
	// iPadOS reports itself as a Mac; touch gives it away.
	if (/iphone|ipad|ipod/i.test(ua) || (/mac/i.test(platform) && navigator.maxTouchPoints > 1)) return { os: "ios", arm: true };
	if (/win/i.test(platform) || /windows/i.test(ua)) return { os: "windows", arm };
	if (/mac/i.test(platform) || /mac os/i.test(ua)) return { os: "macos", arm };
	if (/linux|x11|cros/i.test(platform + ua)) return { os: "linux", arm };
	return { os: "", arm };
}

// Files that are never the app itself.
const SKIP = /(\.sha256|\.sha512|\.sig|\.asc|checksums?|sha256sums|\.txt|\.json|\.blockmap|symbols|\.aab)$|native-debug-symbols|anyprobe/i;

// kinds maps a platform to its file types, best first, and words that
// mark an archive as being for it.
const KINDS = {
	windows: { exts: [/-setup\.exe$/i, /\.msi$/i, /\.msix(bundle)?$/i, /\.exe$/i], words: /windows|win64|win32|\bwin\b/i },
	macos: { exts: [/\.dmg$/i, /\.pkg$/i], words: /mac|darwin|osx|apple/i },
	linux: { exts: [/\.AppImage$/i, /\.flatpak$/i, /\.deb$/i, /\.rpm$/i], words: /linux/i },
	android: { exts: [/\.apk$/i], words: /android/i },
	ios: { exts: [/\.ipa$/i], words: /ios|iphone/i },
};

const ARCHIVE = /\.(zip|tar\.gz|tgz|tar\.xz|7z)$/i;

// classify returns the platform an asset is for, and its rank among that
// platform's files (lower is better), or null.
export function classify(name) {
	if (SKIP.test(name)) return null;
	for (const [os, k] of Object.entries(KINDS)) {
		const i = k.exts.findIndex((re) => re.test(name));
		if (i >= 0) {
			// A .exe named for another system (an installer for a
			// tool) still counts as Windows.
			return { os, rank: i };
		}
	}
	if (ARCHIVE.test(name)) {
		for (const [os, k] of Object.entries(KINDS)) {
			if (k.words.test(name)) return { os, rank: 10 };
		}
	}
	return null;
}

function archScore(name, arm) {
	const isArm = /arm64|aarch64|arm\b/i.test(name);
	const isX86 = /amd64|x86_64|x64|intel/i.test(name);
	const universal = /universal/i.test(name);
	if (universal) return 0;
	if (arm) return isArm ? 0 : isX86 ? 2 : 1;
	return isX86 ? 0 : isArm ? 2 : 1;
}

// pick sorts a release's assets by platform, best file first for each.
export function pick(assets, arm) {
	const by = {};
	for (const a of assets) {
		const c = classify(a.name);
		if (!c) continue;
		(by[c.os] = by[c.os] || []).push({ ...a, rank: c.rank, arch: archScore(a.name, arm) });
	}
	for (const os of Object.keys(by)) {
		// A file for another CPU comes last, whatever its type.
		const wrong = (a) => (a.arch === 2 ? 1 : 0);
		by[os].sort((x, y) => wrong(x) - wrong(y) || x.rank - y.rank || x.arch - y.arch || x.name.localeCompare(y.name));
	}
	return by;
}

// latestRelease fetches a repository's newest release (prereleases
// included, drafts not), cached for the session.
export async function latestRelease(repo) {
	const key = "release:" + repo;
	try {
		const cached = sessionStorage.getItem(key);
		if (cached) return JSON.parse(cached);
	} catch (e) {
		// storage blocked: fetch every time
	}
	const r = await fetch(`https://api.github.com/repos/${repo}/releases?per_page=10`, {
		headers: { Accept: "application/vnd.github+json" },
	});
	if (!r.ok) throw new Error(r.status === 403 ? "GitHub's rate limit was reached; try again in a few minutes" : `GitHub answered ${r.status}`);
	const releases = (await r.json()).filter((rel) => !rel.draft);
	const rel = releases.find((x) => x.assets && x.assets.length) || releases[0] || null;
	const out = rel && {
		tag: rel.tag_name,
		name: rel.name || rel.tag_name,
		url: rel.html_url,
		date: rel.published_at,
		prerelease: rel.prerelease,
		assets: (rel.assets || []).map((a) => ({ name: a.name, url: a.browser_download_url, size: a.size })),
	};
	try {
		sessionStorage.setItem(key, JSON.stringify(out));
	} catch (e) {
		// ignore
	}
	return out;
}

export function formatSize(n) {
	if (!n) return "";
	const units = ["B", "KB", "MB", "GB"];
	let i = 0;
	while (n >= 1024 && i < units.length - 1) {
		n /= 1024;
		i++;
	}
	return `${n < 10 && i > 0 ? n.toFixed(1) : Math.round(n)} ${units[i]}`;
}

// fileKind describes an asset for people: "Installer (.exe)".
export function fileKind(name) {
	const n = name.toLowerCase();
	if (n.endsWith("-setup.exe")) return "Installer";
	if (n.endsWith(".msix") || n.endsWith(".msixbundle")) return "MSIX package";
	if (n.endsWith(".msi")) return "Installer";
	if (n.endsWith(".exe")) return "Program";
	if (n.endsWith(".dmg")) return "Disk image";
	if (n.endsWith(".pkg")) return "Installer";
	if (n.endsWith(".appimage")) return "AppImage";
	if (n.endsWith(".flatpak")) return "Flatpak";
	if (n.endsWith(".deb")) return "Debian / Ubuntu package";
	if (n.endsWith(".rpm")) return "Fedora / openSUSE package";
	if (n.endsWith(".apk")) return "Android app";
	if (n.endsWith(".ipa")) return "iOS app";
	return "Archive";
}
