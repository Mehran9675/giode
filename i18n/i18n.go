// Package i18n translates an app's text and tells its layout which way
// to read: messages by locale with plurals and placeholders, the
// system's locale, and right-to-left languages.
//
//	//go:embed locales
//	var locales embed.FS
//
//	i18n.Load(locales, "locales") // locales/en.json, locales/fa.json, ...
//	i18n.SetLocale(i18n.SystemLocale())
//
//	giode.Text(i18n.T("greeting", "name", user.Name))  // "Hello, {name}!"
//	giode.Text(i18n.N("unread", count))                // {"one": "1 message", "other": "{n} messages"}
//
// A message file is a JSON object of keys to text, or to plural forms
// (zero, one, two, few, many, other) chosen by the language's plural
// rules. Lookups fall back from "fa-IR" to "fa" to the fallback locale
// (English by default) to the key itself.
//
// In a right-to-left locale (Arabic, Persian, Hebrew, Urdu, ...) giode
// lays boxes out mirrored, as CSS does for direction: rtl, and aligns
// text from the right.
package i18n

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"sync"
	"sync/atomic"

	"gioui.org/io/system"
)

var (
	mu       sync.RWMutex
	catalogs = map[string]map[string]message{}
	locale   = "en"
	fallback = "en"
	version  atomic.Uint64
	hooks    []func()
)

// message is a text or plural forms.
type message struct {
	text  string
	forms map[string]string
}

func (m *message) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &m.text); err == nil {
		return nil
	}
	return json.Unmarshal(b, &m.forms)
}

// Load reads the message files in dir of fsys: <locale>.json each.
func Load(fsys fs.FS, dir string) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || path.Ext(e.Name()) != ".json" {
			continue
		}
		b, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		var msgs map[string]message
		if err := json.Unmarshal(b, &msgs); err != nil {
			return fmt.Errorf("i18n: %s: %w", e.Name(), err)
		}
		add(strings.TrimSuffix(e.Name(), ".json"), msgs)
	}
	return nil
}

// Add adds messages for a locale: text by key, or plural forms as a
// map[string]string.
func Add(loc string, messages map[string]any) {
	msgs := make(map[string]message, len(messages))
	for k, v := range messages {
		switch v := v.(type) {
		case string:
			msgs[k] = message{text: v}
		case map[string]string:
			msgs[k] = message{forms: v}
		}
	}
	add(loc, msgs)
}

func add(loc string, msgs map[string]message) {
	loc = normalize(loc)
	mu.Lock()
	c := catalogs[loc]
	if c == nil {
		c = map[string]message{}
		catalogs[loc] = c
	}
	for k, v := range msgs {
		c[k] = v
	}
	mu.Unlock()
	changed()
}

// SetLocale makes loc ("fa", "pt-BR") the app's locale.
func SetLocale(loc string) {
	mu.Lock()
	locale = normalize(loc)
	mu.Unlock()
	changed()
}

// SetFallback sets the locale used for keys the current one lacks.
func SetFallback(loc string) {
	mu.Lock()
	fallback = normalize(loc)
	mu.Unlock()
	changed()
}

// Locale returns the app's locale.
func Locale() string {
	mu.RLock()
	defer mu.RUnlock()
	return locale
}

// Locales returns the locales with messages.
func Locales() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(catalogs))
	for l := range catalogs {
		out = append(out, l)
	}
	return out
}

// OnChange registers fn to run when the locale or messages change;
// giode uses it to redraw.
func OnChange(fn func()) {
	mu.Lock()
	hooks = append(hooks, fn)
	mu.Unlock()
}

// Version counts locale and message changes.
func Version() uint64 { return version.Load() }

func changed() {
	version.Add(1)
	mu.RLock()
	hs := append([]func(){}, hooks...)
	mu.RUnlock()
	for _, h := range hs {
		h()
	}
}

// rtlLanguages read right to left.
var rtlLanguages = map[string]bool{"ar": true, "fa": true, "he": true, "iw": true, "ur": true, "ps": true, "sd": true,
	"ug": true, "yi": true, "dv": true, "ckb": true, "ku-arab": true}

// RTL reports whether the locale reads right to left.
func RTL() bool {
	l := Locale()
	return rtlLanguages[l] || rtlLanguages[language(l)]
}

// Direction returns the locale's text direction.
func Direction() system.TextDirection {
	if RTL() {
		return system.RTL
	}
	return system.LTR
}

// SystemLocale returns the user's preferred locale from the platform,
// "en" when it can't tell.
func SystemLocale() string {
	if l := systemLocale(); l != "" {
		return normalize(l)
	}
	return "en"
}

// T returns the message for key in the locale, with placeholders
// filled: {name} from args given as name, value pairs or a single
// map[string]any.
func T(key string, args ...any) string {
	m, ok := lookup(key)
	if !ok {
		return fill(key, args)
	}
	if m.forms != nil {
		return fill(m.forms["other"], args)
	}
	return fill(m.text, args)
}

// N returns the plural form of key for n, the language's rules
// choosing among its forms; {n} is n.
func N(key string, n int, args ...any) string {
	m, ok := lookup(key)
	if !ok {
		return fill(key, append(args, "n", n))
	}
	args = append(args, "n", n)
	if m.forms == nil {
		return fill(m.text, args)
	}
	form := pluralForm(language(Locale()), n)
	if n == 0 && m.forms["zero"] != "" {
		form = "zero"
	}
	s, ok := m.forms[form]
	if !ok {
		s = m.forms["other"]
	}
	return fill(s, args)
}

// lookup finds key in the locale, its language, then the fallback.
func lookup(key string) (message, bool) {
	mu.RLock()
	defer mu.RUnlock()
	for _, l := range []string{locale, language(locale), fallback, language(fallback)} {
		if m, ok := catalogs[l][key]; ok {
			return m, true
		}
	}
	return message{}, false
}

// fill replaces {name} placeholders.
func fill(s string, args []any) string {
	if len(args) == 0 || !strings.Contains(s, "{") {
		return s
	}
	vals := map[string]string{}
	if len(args) == 1 {
		if m, ok := args[0].(map[string]any); ok {
			for k, v := range m {
				vals[k] = fmt.Sprint(v)
			}
		}
	}
	for i := 0; i+1 < len(args); i += 2 {
		vals[fmt.Sprint(args[i])] = fmt.Sprint(args[i+1])
	}
	var b strings.Builder
	for {
		i := strings.IndexByte(s, '{')
		if i < 0 {
			b.WriteString(s)
			return b.String()
		}
		j := strings.IndexByte(s[i:], '}')
		if j < 0 {
			b.WriteString(s)
			return b.String()
		}
		name := s[i+1 : i+j]
		b.WriteString(s[:i])
		if v, ok := vals[name]; ok {
			b.WriteString(v)
		} else {
			b.WriteString(s[i : i+j+1])
		}
		s = s[i+j+1:]
	}
}

// normalize turns pt_BR.UTF-8 into pt-br.
func normalize(l string) string {
	l = strings.TrimSpace(l)
	if i := strings.IndexAny(l, ".@"); i >= 0 {
		l = l[:i]
	}
	return strings.ToLower(strings.ReplaceAll(l, "_", "-"))
}

// language is a locale's language: "pt" for "pt-br".
func language(l string) string {
	if i := strings.IndexByte(l, '-'); i >= 0 {
		return l[:i]
	}
	return l
}

// pluralForm applies the CLDR cardinal plural rules of common languages
// to an integer.
func pluralForm(lang string, n int) string {
	if n < 0 {
		n = -n
	}
	n10, n100 := n%10, n%100
	switch lang {
	case "ja", "zh", "ko", "vi", "th", "id", "ms", "fa", "tr", "ka", "km", "lo", "my":
		// Persian and Turkish use "one" for 1 in CLDR, but translations
		// commonly don't; "other" with {n} reads right either way.
		if (lang == "fa" || lang == "tr") && n == 1 {
			return "one"
		}
		return "other"
	case "fr", "pt", "hi", "bn", "hy":
		if n == 0 || n == 1 {
			return "one"
		}
		return "other"
	case "ru", "uk", "be", "sr", "hr", "bs":
		switch {
		case n10 == 1 && n100 != 11:
			return "one"
		case n10 >= 2 && n10 <= 4 && (n100 < 12 || n100 > 14):
			return "few"
		}
		return "many"
	case "pl":
		switch {
		case n == 1:
			return "one"
		case n10 >= 2 && n10 <= 4 && (n100 < 12 || n100 > 14):
			return "few"
		}
		return "many"
	case "cs", "sk":
		switch {
		case n == 1:
			return "one"
		case n >= 2 && n <= 4:
			return "few"
		}
		return "other"
	case "ar":
		switch {
		case n == 0:
			return "zero"
		case n == 1:
			return "one"
		case n == 2:
			return "two"
		case n100 >= 3 && n100 <= 10:
			return "few"
		case n100 >= 11 && n100 <= 99:
			return "many"
		}
		return "other"
	case "he", "iw":
		switch {
		case n == 1:
			return "one"
		case n == 2:
			return "two"
		}
		return "other"
	case "lt":
		switch {
		case n10 == 1 && (n100 < 11 || n100 > 19):
			return "one"
		case n10 >= 2 && (n100 < 11 || n100 > 19):
			return "few"
		}
		return "other"
	case "ro":
		switch {
		case n == 1:
			return "one"
		case n == 0 || (n100 >= 2 && n100 <= 19):
			return "few"
		}
		return "other"
	}
	if n == 1 {
		return "one"
	}
	return "other"
}
