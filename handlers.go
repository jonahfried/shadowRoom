package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func (cir *Agent) fire(win *pixelgl.Window) {
	mousePosn := cir.Cam.Matrix.Unproject(win.MousePosition())
	directionVec := mousePosn.Sub(cir.Posn)
	directionVec = directionVec.Scaled(1 / magnitude(directionVec))

	switch cir.GunType {
	case 1:
		var bullet Shot
		bullet.GunType = 1
		bullet.color = pixel.ToRGBA(colornames.Firebrick).Mul(pixel.Alpha(.7))
		bullet.Posn1 = cir.Posn
		bullet.Posn2 = cir.Posn.Add(directionVec.Scaled(10))
		// bullet.Vel.Add(cir.Vel)
		bullet.Vel = directionVec.Scaled(14)
		cir.Shots = append(cir.Shots, bullet)
	case 2:
		angle := math.Atan2(mousePosn.Y-cir.Posn.Y, mousePosn.X-cir.Posn.X)
		for shotCount := 0; shotCount < 5; shotCount++ {
			var bullet Shot
			bullet.GunType = 2
			bullet.color = pixel.ToRGBA(colornames.Firebrick).Mul(pixel.Alpha(.7))
			// offset := pixel.V(rand.Float64()*20, rand.Float64()*20)
			offset := (rand.Float64() - rand.Float64()) / 2.3
			// newDirection := directionVec.Add(offset)
			newAngle := angle + offset
			// newDirection = newDirection.Scaled(1 / magnitude(newDirection))
			bullet.Posn1 = cir.Posn
			newDirection := pixel.V(math.Cos(newAngle), math.Sin(newAngle))
			bullet.Posn2 = cir.Posn.Add(newDirection.Scaled(10))
			bullet.Vel = (bullet.Posn2.Sub(bullet.Posn1)).Scaled(2 + (rand.Float64()-rand.Float64())/2.3)
			cir.Shots = append(cir.Shots, bullet)
		}
		cir.Bullets--
		if cir.Bullets <= 0 {
			cir.GunType = 1
		}
	}
}

// PressHandler Handles key presses.
// Agent method taking in a window from which to accept inputs.
func (cir *Agent) PressHandler(win *pixelgl.Window) {
	cir.Acc = pixel.ZV
	if win.Pressed(pixelgl.KeyA) {
		cir.Acc = cir.Acc.Sub(pixel.V(5, 0))
	}
	if win.Pressed(pixelgl.KeyD) {
		cir.Acc = cir.Acc.Add(pixel.V(5, 0))
	}
	if win.Pressed(pixelgl.KeyS) {
		cir.Acc = cir.Acc.Sub(pixel.V(0, 5))
	}
	if win.Pressed(pixelgl.KeyW) {
		cir.Acc = cir.Acc.Add(pixel.V(0, 5))
	}

	if cir.devMode {
		if win.JustPressed(pixelgl.KeyJ) {
			cir.Posn.X--
		}
		if win.JustPressed(pixelgl.KeyL) {
			cir.Posn.X++
		}
		if win.JustPressed(pixelgl.KeyK) {
			cir.Posn.Y--
		}
		if win.JustPressed(pixelgl.KeyI) {
			cir.Posn.Y++
		}
		if win.JustPressed(pixelgl.KeySpace) {
			cir.Shade = !cir.Shade
		}
		if win.JustPressed(pixelgl.KeyF) {
			cir.Fill = !cir.Fill
		}
		if win.JustPressed(pixelgl.KeyUp) {
			cir.Level += .001
		}
		if win.JustPressed(pixelgl.KeyDown) {
			cir.Level -= .001
		}
		if win.JustPressed(pixelgl.KeyRight) {
			cir.Spacing++
		}
		if win.JustPressed(pixelgl.KeyLeft) {
			cir.Spacing--
		}
		if win.Pressed(pixelgl.KeyComma) {
			cir.Count--
		}
		if win.Pressed(pixelgl.KeyPeriod) {
			cir.Count++
		}
	}

	if win.JustPressed(pixelgl.MouseButton1) {
		cir.fire(win)
	}
}

// ReleaseHandler Handles key releases.
// Agent method taking in a window from which to accept inputs.
func (cir *Agent) ReleaseHandler(win *pixelgl.Window) {
	if win.JustReleased(pixelgl.KeyA) {
		cir.Acc = cir.Acc.Add(pixel.V(5, 0))
	}
	if win.JustReleased(pixelgl.KeyD) {
		cir.Acc = cir.Acc.Sub(pixel.V(5, 0))
	}
	if win.JustReleased(pixelgl.KeyS) {
		cir.Acc = cir.Acc.Add(pixel.V(0, 5))
	}
	if win.JustReleased(pixelgl.KeyW) {
		cir.Acc = cir.Acc.Sub(pixel.V(0, 5))
	}
}
