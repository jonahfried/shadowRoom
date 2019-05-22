package main

import "github.com/faiface/pixel"

// Returns the magnitude of the distance between v1 and v2
func vecDist(v1, v2 pixel.Vec) float64 {
	return v1.Sub(v2).Len()
}

// Vector limitation. Takes in a pixel.Vec and a float64.
// If the magnitude of the given vector is greater than the limit,
// Then the magnitude is scaled down to the limit.
func limitVecMag(v pixel.Vec, lim float64) pixel.Vec {
	if v.Len() != 0 && v.Len() > lim {
		v = v.Scaled(lim / v.Len())
	}
	return v
}

// Returns the point on a line segment closest to a given position. (Either one of the two ends or a point between them)
func closestPointOnSegment(end1, end2, posn pixel.Vec) (closest pixel.Vec) {
	segVec := end2.Sub(end1)
	unitSegVec := segVec.Scaled(1 / segVec.Len())
	posnOffset := posn.Sub(end1)
	projMag := posnOffset.Dot(unitSegVec)
	projVec := unitSegVec.Scaled(projMag)

	if projMag < 0 {
		closest = end1
	} else if projMag > segVec.Len() {
		closest = end2
	} else {
		closest = end1.Add(projVec)
	}
	return closest
}
