package i18n

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// Intl returns the message for key in the app's locale with data
// filled into its {placeholders}. data is a map (map[string]any,
// map[string]string, ...), a struct (its fields by name, by json tag,
// and with a lower-case first letter), or name, value pairs:
//
//	i18n.Intl("greeting", map[string]any{"name": user.Name}) // "Hello, {name}!"
//	i18n.Intl("greeting", struct{ Name string }{user.Name})
//	i18n.Intl("greeting", "name", user.Name)
//
// When the message has plural forms, a "count" or "n" value picks the
// form by the language's rules:
//
//	i18n.Intl("unread", map[string]any{"count": 3}) // {"one": "{count} message", "other": "{count} messages"}
//
// A missing key returns the key itself, filled the same way.
func Intl(key string, data ...any) string {
	vals := values(data)
	m, ok := lookup(key)
	if !ok {
		return fill(key, []any{vals})
	}
	if m.forms == nil {
		return fill(m.text, []any{vals})
	}
	n, hasN := count(vals)
	if !hasN {
		return fill(m.forms["other"], []any{vals})
	}
	if _, ok := vals["n"]; !ok {
		vals["n"] = n
	}
	form := pluralForm(language(Locale()), n)
	if n == 0 && m.forms["zero"] != "" {
		form = "zero"
	}
	s, ok := m.forms[form]
	if !ok {
		s = m.forms["other"]
	}
	return fill(s, []any{vals})
}

// values flattens Intl's data into one map.
func values(data []any) map[string]any {
	vals := map[string]any{}
	if len(data) != 1 {
		for i := 0; i+1 < len(data); i += 2 {
			vals[fmt.Sprint(data[i])] = data[i+1]
		}
		return vals
	}
	v := reflect.ValueOf(data[0])
	for v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Map:
		for it := v.MapRange(); it.Next(); {
			vals[fmt.Sprint(it.Key().Interface())] = it.Value().Interface()
		}
	case reflect.Struct:
		t := v.Type()
		for i := range t.NumField() {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			x := v.Field(i).Interface()
			vals[f.Name] = x
			vals[strings.ToLower(f.Name[:1])+f.Name[1:]] = x
			if tag, _, _ := strings.Cut(f.Tag.Get("json"), ","); tag != "" && tag != "-" {
				vals[tag] = x
			}
		}
	}
	return vals
}

// count returns the integer "count" or "n" value.
func count(vals map[string]any) (int, bool) {
	for _, k := range []string{"count", "n"} {
		v := reflect.ValueOf(vals[k])
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int(v.Int()), true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int(v.Uint()), true
		case reflect.Float32, reflect.Float64:
			return int(v.Float()), true
		}
	}
	return 0, false
}

// LoadDir reads the message files (<locale>.json) in a folder on disk.
// A relative dir is looked for in the working directory, then next to
// the executable, then in a macOS app bundle's Resources folder.
func LoadDir(dir string) error {
	path, err := FindDir(dir)
	if err != nil {
		return err
	}
	return Load(os.DirFS(path), ".")
}

// FindDir resolves a translations folder as LoadDir does.
func FindDir(dir string) (string, error) {
	candidates := []string{dir}
	if !filepath.IsAbs(dir) {
		if exe, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exe)
			candidates = append(candidates,
				filepath.Join(exeDir, dir),
				filepath.Join(exeDir, "..", "Resources", dir))
		}
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c, nil
		}
	}
	return "", fmt.Errorf("i18n: no translations folder %q (looked in %s): %w", dir, strings.Join(candidates, ", "), fs.ErrNotExist)
}
