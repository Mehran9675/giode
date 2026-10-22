package elements

import "image"

// Position corrections.
//
// Some things are recorded with their window position while a child is
// being laid out: popup anchors (RegisterPortal) and semantic nodes
// (RecordSemantic). AbsolutePosition only knows offsets pushed before
// the child is laid out, but containers such as a Row, a Stack or a
// Scroll decide where a child goes after laying it out. Those
// containers mark what was recorded while laying each child out and,
// once the child's offset is known, shift it there.

// PositionMark marks how many position-dependent records existed at
// some point of the layout.
type PositionMark struct {
	semantic, portal, probe int
}

// probes are positions components keep for later use (a menu
// trigger's anchor), corrected like portal anchors.
var probes []*image.Point

// TrackPosition stores p, a window position computed from
// AbsolutePosition, in *dst, and keeps *dst corrected for the rest of
// the frame as enclosing containers place the caller. Use it when a
// component computes a position now but uses it elsewhere (a menu
// trigger whose menu is laid out in another part of the tree).
func TrackPosition(dst *image.Point, p image.Point) {
	*dst = p
	probes = append(probes, dst)
}

// MarkPositions returns a mark of the records made so far this frame.
func MarkPositions() PositionMark {
	m := PositionMark{portal: len(portalQueue), probe: len(probes)}
	if semantics != nil {
		m.semantic = len(*semantics)
	}
	return m
}

// ShiftPositions moves everything recorded since mark by off. Call it
// from a container that places a child with op.Offset after laying it
// out, with a mark taken just before the child's Layout:
//
//	mark := elements.MarkPositions()
//	dims := child.Layout(gtx)          // records at the container's origin
//	off := placeChild(dims)
//	elements.ShiftPositions(mark, off) // now at the child's real position
//
// Containers that push the offset with PushOffset before laying the
// child out don't need it.
func ShiftPositions(mark PositionMark, off image.Point) {
	shiftRange(mark, MarkPositions(), off)
}

// ShiftPositionsRange moves the records made between two marks by off,
// for a container that records other things after the child (a
// scrollbar, say) before it knows the child's offset.
func ShiftPositionsRange(from, to PositionMark, off image.Point) {
	shiftRange(from, to, off)
}

// shiftRange moves the records between two marks by off.
func shiftRange(from, to PositionMark, off image.Point) {
	if off == (image.Point{}) {
		return
	}
	if semantics != nil {
		nodes := *semantics
		for i := from.semantic; i < to.semantic && i < len(nodes); i++ {
			nodes[i].Bounds = nodes[i].Bounds.Add(off)
		}
	}
	for i := from.portal; i < to.portal && i < len(portalQueue); i++ {
		if !portalQueue[i].fixed {
			portalQueue[i].anchor = portalQueue[i].anchor.Add(off)
		}
	}
	for i := from.probe; i < to.probe && i < len(probes); i++ {
		*probes[i] = probes[i].Add(off)
	}
}
