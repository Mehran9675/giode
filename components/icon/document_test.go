package icon

import "testing"

func TestLoadSVG_ViewBoxAndPaths(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
		<path d="M0 0 L24 0 L24 24 Z"/>
	</svg>`
	d, err := LoadSVG([]byte(svg))
	if err != nil {
		t.Fatalf("LoadSVG: %v", err)
	}
	if d.viewW != 24 || d.viewH != 24 {
		t.Errorf("viewBox size = %v,%v, want 24,24", d.viewW, d.viewH)
	}
	if d.viewX != 0 || d.viewY != 0 {
		t.Errorf("viewBox origin = %v,%v, want 0,0", d.viewX, d.viewY)
	}
	if len(d.paths) != 1 || d.paths[0] != "M0 0 L24 0 L24 24 Z" {
		t.Errorf("paths = %v, want one path", d.paths)
	}
}

func TestLoadSVG_OffsetViewBox(t *testing.T) {
	svg := `<svg viewBox="10 20 100 200"><path d="M0 0"/></svg>`
	d, err := LoadSVG([]byte(svg))
	if err != nil {
		t.Fatalf("LoadSVG: %v", err)
	}
	if d.viewX != 10 || d.viewY != 20 || d.viewW != 100 || d.viewH != 200 {
		t.Errorf("viewBox = %v,%v,%v,%v, want 10,20,100,200", d.viewX, d.viewY, d.viewW, d.viewH)
	}
}

func TestLoadSVG_NestedGroups(t *testing.T) {
	svg := `<svg viewBox="0 0 24 24">
		<g>
			<path d="M1 1"/>
			<g><path d="M2 2"/></g>
		</g>
		<path d="M3 3"/>
	</svg>`
	d, err := LoadSVG([]byte(svg))
	if err != nil {
		t.Fatalf("LoadSVG: %v", err)
	}
	if len(d.paths) != 3 {
		t.Fatalf("paths = %v, want 3 paths", d.paths)
	}
}

func TestLoadSVG_FallsBackToWidthHeight(t *testing.T) {
	svg := `<svg width="32px" height="32px"><path d="M0 0"/></svg>`
	d, err := LoadSVG([]byte(svg))
	if err != nil {
		t.Fatalf("LoadSVG: %v", err)
	}
	if d.viewW != 32 || d.viewH != 32 {
		t.Errorf("size = %v,%v, want 32,32", d.viewW, d.viewH)
	}
}

func TestLoadSVG_ErrorsWithoutSizeInfo(t *testing.T) {
	if _, err := LoadSVG([]byte(`<svg><path d="M0 0"/></svg>`)); err == nil {
		t.Error("LoadSVG with no viewBox or width/height: want error, got nil")
	}
}

func TestLoadSVG_ErrorsOnMalformedXML(t *testing.T) {
	if _, err := LoadSVG([]byte(`not xml`)); err == nil {
		t.Error("LoadSVG with malformed XML: want error, got nil")
	}
}

func TestMustLoadSVG_PanicsOnError(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustLoadSVG with malformed XML: want panic, got none")
		}
	}()
	MustLoadSVG([]byte(`not xml`))
}
