package textrender

import (
	"image"
	"image/png"
	"os"
)

func save(path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	png.Encode(f, img)
	f.Close()
}
