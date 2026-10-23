package giode

import (
	"image"
	"testing"

	"gioui.org/app"
)

func TestSetSize_WidensBounds(t *testing.T) {
	a := &App{window: new(app.Window), cfg: Config{Resizable: true, MinSize: image.Pt(600, 400), MaxSize: image.Pt(1200, 800)}}
	a.SetSize(400, 840)
	if a.cfg.MinSize != image.Pt(400, 400) || a.cfg.MaxSize != image.Pt(1200, 840) {
		t.Errorf("bounds %v..%v, want 400x400..1200x840", a.cfg.MinSize, a.cfg.MaxSize)
	}

	locked := &App{window: new(app.Window), cfg: Config{Size: image.Pt(800, 600)}}
	locked.SetSize(400, 840)
	if locked.cfg.MinSize != image.Pt(400, 840) || locked.cfg.MaxSize != image.Pt(400, 840) {
		t.Errorf("fixed-size window bounds %v..%v, want locked at 400x840", locked.cfg.MinSize, locked.cfg.MaxSize)
	}
}
