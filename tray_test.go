package giode

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// testPNG returns a small PNG.
func testPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	for i := range img.Pix {
		img.Pix[i] = 0xff
	}
	img.Set(3, 3, color.NRGBA{A: 0xff})
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestTray_ItemHelpers(t *testing.T) {
	ran := false
	items := []TrayItem{
		TrayAction("Open", func() { ran = true }),
		TraySeparator(),
		TrayCheckbox("Mute", true, nil),
		TraySubmenu("More", TrayAction("About", nil)),
	}
	if items[0].Label != "Open" || items[0].OnClick == nil {
		t.Error("TrayAction")
	}
	items[0].OnClick()
	if !ran {
		t.Error("action didn't run")
	}
	if !items[1].Separator || !items[2].Checkable || !items[2].Checked || len(items[3].Items) != 1 {
		t.Errorf("helpers built %+v", items)
	}
}
