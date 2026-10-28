package giode

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseContextAction(t *testing.T) {
	id, paths, ok := parseContextAction([]string{"-v", "--giode-action=ask", "/a b/c.txt", "/d"})
	if !ok || id != "ask" || !reflect.DeepEqual(paths, []string{"/a b/c.txt", "/d"}) {
		t.Errorf("got %q %q %v", id, paths, ok)
	}
	if _, _, ok := parseContextAction([]string{"file.txt"}); ok {
		t.Error("plain arguments parsed as an action")
	}
}

func TestContextActions_BuildAndConfig(t *testing.T) {
	b, _ := json.Marshal([]ContextAction{{ID: "ask", Label: "Ask (toml)"}, {ID: "open", Label: "Open"}})
	old := buildContextActions
	buildContextActions = base64.URLEncoding.EncodeToString(b)
	defer func() { buildContextActions = old }()
	a := &App{cfg: Config{ContextActions: []ContextAction{{ID: "ask", Label: "Ask (Go)"}, {ID: "sum", Label: "Summarize", Folders: true}}}}
	got := a.ContextActions()
	want := []ContextAction{{ID: "ask", Label: "Ask (Go)"}, {ID: "open", Label: "Open"}, {ID: "sum", Label: "Summarize", Folders: true}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
}
