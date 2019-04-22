package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Camera is a struct storing a Matrix to set the window
type Camera struct {
	Posn, vel          pixel.Vec
	maxForce, maxSpeed float64
	Matrix             pixel.Matrix
}

// MakeCamera takes in a starting position and window and returns a Camera
func MakeCamera(posn pixel.Vec, win *pixelgl.Window) (cam Camera) {
	cam.Posn = posn
	cam.vel = pixel.ZV
	cam.maxSpeed = 10

	cam.Matrix = pixel.IM.Moved(win.Bounds().Center().Sub(posn))

	return cam
}

// Attract updates a Camera's velocity and position to follow a given point
func (cam *Camera) Attract(target pixel.Vec) {
	acc := limitVecMag(target.Sub(cam.Posn), vecDist(cam.Posn, target)/10)
	// scale := cam.maxForce / 8
	// acc = acc.Scaled(scale)
	acc = acc.Sub(cam.vel)

	cam.vel = cam.vel.Add(acc)
	cam.vel = limitVecMag(cam.vel, cam.maxSpeed)

	cam.Posn = cam.Posn.Add(cam.vel)

}
