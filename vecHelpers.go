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
