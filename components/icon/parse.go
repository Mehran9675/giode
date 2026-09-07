package icon

import (
	"strconv"
	"strings"

	"gioui.org/f32"
	"gioui.org/op/clip"
)

// pathParser parses SVG path data into a clip.Path. State (current
// point, control points) persists across commands.
type pathParser struct {
	p         *clip.Path
	cur       f32.Point
	first     f32.Point
	cp        f32.Point
	subCp     f32.Point
	haveFirst bool
	prev      byte
}

// parsePath parses SVG path data into p. Unsupported arc commands
// degrade to straight lines.
func parsePath(p *clip.Path, d string) {
	pp := &pathParser{p: p}
	var (
		cmd    byte
		args   []float32
		startI int
	)
	flush := func() {
		if len(args) == 0 {
			return
		}
		pp.run(cmd, args)
		args = args[:0]
	}
	for i := 0; i < len(d); i++ {
		c := d[i]
		switch {
		case c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z':
			if len(args) > 0 {
				flush()
			}
			cmd = c
			startI = i + 1
		case c == ' ' || c == ',' || c == '\t' || c == '\n' || c == '\r':
			if startI < i {
				if v, ok := parseFloat(d[startI:i]); ok {
					args = append(args, v)
				}
			}
			startI = i + 1
		case c == '-':
			if startI < i {
				if v, ok := parseFloat(d[startI:i]); ok {
					args = append(args, v)
				}
			}
			startI = i
		}
	}
	if startI < len(d) {
		if v, ok := parseFloat(d[startI:]); ok {
			args = append(args, v)
		}
	}
	flush()
}

func parseFloat(s string) (float32, bool) {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 32)
	if err != nil {
		return 0, false
	}
	return float32(f), true
}

// run executes one command with its arguments.
func (pp *pathParser) run(cmd byte, args []float32) {
	rel := cmd >= 'a'
	cc := cmd
	if rel {
		cc = cmd - 'a' + 'A'
	}
	if cc == 'Z' {
		pp.close()
		return
	}
	n := map[byte]int{
		'M': 2, 'L': 2, 'H': 1, 'V': 1,
		'C': 6, 'S': 4, 'Q': 4, 'T': 2, 'A': 7,
	}[cc]
	if n == 0 || len(args) < n {
		return
	}
	point := func(off int) f32.Point {
		a := args[off : off+2]
		if rel {
			return f32.Pt(pp.cur.X+a[0], pp.cur.Y+a[1])
		}
		return f32.Pt(a[0], a[1])
	}
	scalar := func(v float32, hor bool) f32.Point {
		if rel {
			if hor {
				return f32.Pt(pp.cur.X+v, pp.cur.Y)
			}
			return f32.Pt(pp.cur.X, pp.cur.Y+v)
		}
		if hor {
			return f32.Pt(v, pp.cur.Y)
		}
		return f32.Pt(pp.cur.X, v)
	}
	for i := 0; i+n <= len(args); i += n {
		switch cc {
		case 'M':
			to := point(i)
			if pp.haveFirst && (i > 0 || pp.prev == 'M' || pp.prev == 'm') {
				pp.line(to)
			} else {
				pp.move(to)
			}
		case 'L':
			pp.line(point(i))
		case 'H':
			pp.line(scalar(args[i], true))
		case 'V':
			pp.line(scalar(args[i], false))
		case 'C':
			c1 := point(i)
			c2 := point(i + 2)
			to := point(i + 4)
			pp.p.CubeTo(c1, c2, to)
			pp.cur = to
			pp.cp = c2
			pp.subCp = c2
		case 'S':
			c1 := pp.cur
			if pp.prev == 'C' || pp.prev == 'c' || pp.prev == 'S' || pp.prev == 's' {
				c1 = f32.Pt(2*pp.cur.X-pp.cp.X, 2*pp.cur.Y-pp.cp.Y)
			}
			c2 := point(i)
			to := point(i + 2)
			pp.p.CubeTo(c1, c2, to)
			pp.cur = to
			pp.cp = c2
			pp.subCp = c2
		case 'Q':
			c1 := point(i)
			to := point(i + 2)
			pp.p.QuadTo(c1, to)
			pp.cur = to
			pp.cp = c1
			pp.subCp = c1
		case 'T':
			c1 := pp.cur
			if pp.prev == 'Q' || pp.prev == 'q' || pp.prev == 'T' || pp.prev == 't' {
				c1 = f32.Pt(2*pp.cur.X-pp.cp.X, 2*pp.cur.Y-pp.cp.Y)
			}
			to := point(i)
			pp.p.QuadTo(c1, to)
			pp.cur = to
			pp.cp = c1
			pp.subCp = c1
		case 'A':
			// rx ry rot largeArc sweep x y: approximate the arc with a
			// straight line.
			to := point(i + 5)
			pp.line(to)
		}
	}
	pp.prev = cc
}

func (pp *pathParser) move(to f32.Point) {
	pp.p.MoveTo(to)
	pp.cur = to
	pp.first = to
	pp.cp = to
	pp.subCp = to
	pp.haveFirst = true
}

func (pp *pathParser) line(to f32.Point) {
	pp.p.LineTo(to)
	pp.cur = to
	pp.cp = to
	pp.subCp = to
}

func (pp *pathParser) close() {
	if !pp.haveFirst {
		return
	}
	pp.p.Close()
	pp.cur = pp.first
	pp.cp = pp.subCp
}
