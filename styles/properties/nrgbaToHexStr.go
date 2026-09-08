package properties

import (
	"fmt"
	"image/color"
)

func CalcColorReverse(c color.NRGBA) string {
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}
