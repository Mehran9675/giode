package elements

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

// debugf receives layout diagnostics; nil discards them.
var debugf atomic.Pointer[func(format string, v ...any)]

// SetDebugLogger sets the function layout diagnostics go to, such as
// the warning about FlexGrow in an unbounded container; nil turns them
// off. The app sets it from Config.Logger at LogLevelDebug; component
// authors don't need to.
func SetDebugLogger(fn func(format string, v ...any)) {
	if fn == nil {
		debugf.Store(nil)
		return
	}
	debugf.Store(&fn)
}

var (
	warnedMu sync.Mutex
	warned   = map[string]bool{}
)

// warnOnce logs a diagnostic the first time a given message comes up.
func warnOnce(format string, v ...any) {
	fn := debugf.Load()
	if fn == nil {
		return
	}
	msg := fmt.Sprintf(format, v...)
	warnedMu.Lock()
	seen := warned[msg]
	warned[msg] = true
	warnedMu.Unlock()
	if !seen {
		(*fn)("%s", msg)
	}
}

// describe names an element for diagnostics: its kind and the first
// text inside it, which is usually enough to find it in the code.
func describe(el Element) string {
	kind := fmt.Sprintf("%T", el)
	switch e := el.(type) {
	case *BoxElement:
		kind = "Box"
		if s := firstText(e, 4); s != "" {
			return kind + " containing " + strconv.Quote(s)
		}
	case *TextElement:
		return "Text " + strconv.Quote(e.str)
	case *Pressable:
		kind = "Pressable"
		if e.label != "" {
			return kind + " " + strconv.Quote(e.label)
		}
	}
	return kind
}

// firstText returns the text of the first Text element within depth
// levels of b.
func firstText(b *BoxElement, depth int) string {
	if depth == 0 {
		return ""
	}
	for _, c := range b.children {
		switch e := c.(type) {
		case *TextElement:
			return e.str
		case *BoxElement:
			if s := firstText(e, depth-1); s != "" {
				return s
			}
		}
	}
	return ""
}
