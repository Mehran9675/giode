//go:build darwin && !nometal

package giode

// metalOnly reports a build whose Gio can't render with OpenGL: macOS
// and iOS without the nometal tag.
const metalOnly = true
