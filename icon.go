package giode

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// decodeIcon parses PNG or JPEG icon bytes into an image.
func decodeIcon(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

// setIconBytes stores the window icon. The icon is applied to the
// window once its platform handle becomes available.
func (a *App) setIconBytes(data []byte) {
	if len(data) == 0 {
		return
	}
	img, err := decodeIcon(data)
	if err != nil {
		a.logf(LogLevelWarning, "giode: invalid window icon: %v", err)
		os.Stderr.WriteString("giode: invalid window icon: " + err.Error() + "\n")
		return
	}
	a.icon = img
}
