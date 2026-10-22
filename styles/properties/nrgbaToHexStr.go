package properties

import "image/color"

// CalcColorReverse formats c as a "#RRGGBBAA" Color.
func CalcColorReverse(c color.NRGBA) string {
	const digits = "0123456789ABCDEF"
	b := [9]byte{'#'}
	for i, v := range [4]uint8{c.R, c.G, c.B, c.A} {
		b[1+2*i] = digits[v>>4]
		b[2+2*i] = digits[v&0x0f]
	}
	return string(b[:])
}
