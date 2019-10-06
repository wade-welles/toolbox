package poly

import "github.com/richardwilkes/toolbox/xmath/geom"

type intersection struct {
	edge0 *edgeNode
	edge1 *edgeNode
	point geom.Point
	next  *intersection
}

func (inter *intersection) process(op clipOp, pt geom.Point, outPoly *polygonNode) *polygonNode {
	e0 := inter.edge0
	e1 := inter.edge1

	// Only generate output for contributing intersections
	if (e0.bundleAbove[clipping] || e0.bundleAbove[subject]) && (e1.bundleAbove[clipping] || e1.bundleAbove[subject]) {
		n0 := e0.outAbove
		n1 := e1.outAbove
		iPt := inter.point
		iPt.Y += pt.Y
		inClip := (e0.bundleAbove[clipping] && !e0.clipSide) ||
			(e1.bundleAbove[clipping] && e1.clipSide) ||
			(!e0.bundleAbove[clipping] && !e1.bundleAbove[clipping] && e0.clipSide && e1.clipSide)
		inSubj := (e0.bundleAbove[subject] && !e0.subjectSide) ||
			(e1.bundleAbove[subject] && e1.subjectSide) ||
			(!e0.bundleAbove[subject] && !e1.bundleAbove[subject] && e0.subjectSide && e1.subjectSide)

		// Determine quadrant occupancies
		var br, bl, tr, tl bool
		e0InClip := inClip != e0.bundleAbove[clipping]
		e1InClip := inClip != e1.bundleAbove[clipping]
		e0InSubj := inSubj != e0.bundleAbove[subject]
		e1InSubj := inSubj != e1.bundleAbove[subject]
		e10InClip := e1InClip != e0.bundleAbove[clipping]
		e10InSubj := e1InSubj != e0.bundleAbove[subject]
		switch op {
		case subtractOp, intersectOp:
			tr = inClip && inSubj
			tl = e1InClip && e1InSubj
			br = e0InClip && e0InSubj
			bl = e10InClip && e10InSubj
		case xorOp:
			tr = inClip != inSubj
			tl = e1InClip != e1InSubj
			br = e0InClip != e0InSubj
			bl = e10InClip != e10InSubj
		case unionOp:
			tr = inClip || inSubj
			tl = e1InClip || e1InSubj
			br = e0InClip || e0InSubj
			bl = e10InClip || e10InSubj
		}
		switch calcVertexType(tr, tl, br, bl) {
		case externalMinimum:
			outPoly = e0.addLocalMin(outPoly, iPt)
			e1.outAbove = e0.outAbove
		case externalRightIntermediate:
			if n0 != nil {
				n0.addRight(iPt)
				e1.outAbove = n0
				e0.outAbove = nil
			}
		case externalLeftIntermediate:
			if n1 != nil {
				n1.addLeft(iPt)
				e0.outAbove = n1
				e1.outAbove = nil
			}
		case externalMaximum:
			if n0 != nil && n1 != nil {
				n0.addLeft(iPt)
				n0.mergeRight(n1, outPoly)
				e0.outAbove = nil
				e1.outAbove = nil
			}
		case internalMinimum:
			outPoly = e0.addLocalMin(outPoly, iPt)
			e1.outAbove = e0.outAbove
		case internalLeftIntermediate:
			if n0 != nil {
				n0.addLeft(iPt)
				e1.outAbove = n0
				e0.outAbove = nil
			}
		case internalRightIntermediate:
			if n1 != nil {
				n1.addRight(iPt)
				e0.outAbove = n1
				e1.outAbove = nil
			}
		case internalMaximum:
			if n0 != nil && n1 != nil {
				n0.addRight(iPt)
				n0.mergeLeft(n1, outPoly)
				e0.outAbove = nil
				e1.outAbove = nil
			}
		case internalMaximumAndMinimum:
			if n0 != nil && n1 != nil {
				n0.addRight(iPt)
				n0.mergeLeft(n1, outPoly)
				outPoly = e0.addLocalMin(outPoly, iPt)
				e1.outAbove = e0.outAbove
			}
		case externalMaximumAndMinimum:
			if n0 != nil && n1 != nil {
				n0.addLeft(iPt)
				n0.mergeRight(n1, outPoly)
				outPoly = e0.addLocalMin(outPoly, iPt)
				e1.outAbove = e0.outAbove
			}
		default:
		}
	}

	// Swap bundle sides in response to edge crossing
	if e0.bundleAbove[clipping] {
		e1.clipSide = !e1.clipSide
	}
	if e1.bundleAbove[clipping] {
		e0.clipSide = !e0.clipSide
	}
	if e0.bundleAbove[subject] {
		e1.subjectSide = !e1.subjectSide
	}
	if e1.bundleAbove[subject] {
		e0.subjectSide = !e0.subjectSide
	}

	return outPoly
}
