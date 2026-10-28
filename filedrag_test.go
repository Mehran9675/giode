package giode

import (
	"reflect"
	"testing"

	"gioui.org/f32"
)

func TestParseURIList(t *testing.T) {
	list := "# a comment\r\nfile:///home/me/My%20Video.mkv\r\nfile://localhost/tmp/a.txt\r\nhttps://example.com/x\r\nfile://otherhost/share/b\r\n\r\n"
	got := parseURIList(list)
	want := []string{"/home/me/My Video.mkv", "/tmp/a.txt"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseURIList = %q, want %q", got, want)
	}
}

func TestFileDragState(t *testing.T) {
	a := &App{}
	if _, ok := a.FileDrag(); ok {
		t.Fatal("drag before any event")
	}
	a.drag.active, a.drag.pos = true, f32.Pt(10, 20)
	d, ok := a.FileDrag()
	if !ok || d.Position != f32.Pt(10, 20) || d.Kind != FileDragOver {
		t.Errorf("FileDrag = %+v, %v", d, ok)
	}
	for k, s := range map[FileDragKind]string{FileDragEnter: "enter", FileDragOver: "over", FileDragLeave: "leave", FileDragDrop: "drop", 0: "none"} {
		if k.String() != s {
			t.Errorf("%d.String() = %q, want %q", k, k.String(), s)
		}
	}
}
