package mcp

import (
	"encoding/json"
	"reflect"
	"strings"
)

var rawMessageType = reflect.TypeOf(json.RawMessage(nil))

// SchemaFor derives a JSON Schema from a Go type, the way Func does
// for tool inputs. Structs become objects (see Func for the tags
// honored), slices arrays and maps objects; interface and
// json.RawMessage values accept anything.
func SchemaFor(t reflect.Type) map[string]any {
	s := schemaFor(t, 0)
	if s["type"] != "object" {
		// Tool inputs must be objects.
		return map[string]any{"type": "object", "properties": map[string]any{}}
	}
	return s
}

func schemaFor(t reflect.Type, depth int) map[string]any {
	if t == nil || depth > 8 {
		return map[string]any{}
	}
	if t == rawMessageType {
		return map[string]any{}
	}
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return map[string]any{"type": "string"}
	case reflect.Bool:
		return map[string]any{"type": "boolean"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]any{"type": "integer"}
	case reflect.Float32, reflect.Float64:
		return map[string]any{"type": "number"}
	case reflect.Slice, reflect.Array:
		return map[string]any{"type": "array", "items": schemaFor(t.Elem(), depth+1)}
	case reflect.Map:
		return map[string]any{"type": "object", "additionalProperties": schemaFor(t.Elem(), depth+1)}
	case reflect.Struct:
		props := map[string]any{}
		var required []string
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			name, opts, _ := strings.Cut(f.Tag.Get("json"), ",")
			if name == "-" {
				continue
			}
			if name == "" {
				name = f.Name
			}
			ps := schemaFor(f.Type, depth+1)
			if d := f.Tag.Get("description"); d != "" {
				ps["description"] = d
			}
			if e := f.Tag.Get("enum"); e != "" {
				ps["enum"] = strings.Split(e, ",")
			}
			props[name] = ps
			if !strings.Contains(opts, "omitempty") && f.Type.Kind() != reflect.Pointer {
				required = append(required, name)
			}
		}
		s := map[string]any{"type": "object", "properties": props}
		if len(required) > 0 {
			s["required"] = required
		}
		return s
	}
	return map[string]any{}
}
