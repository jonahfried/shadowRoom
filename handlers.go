package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func (game *Game) fire(win *pixelgl.Window) {
	mousePosn := game.Player.Cam.Matrix.Unproject(win.MousePosition())
	directionVec := mousePosn.Sub(game.Player.Posn)
	directionVec = directionVec.Scaled(1 / directionVec.Len())

	switch game.Player.GunType {
	case 1:
		var bullet Shot
		bullet.GunType = 1
		bullet.color = pixel.ToRGBA(colornames.Firebrick).Mul(pixel.Alpha(.7))
		bullet.Posn1 = game.Player.Posn
		bullet.Posn2 = game.Player.Posn.Add(directionVec.Scaled(10))
		bullet.Vel = directionVec.Scaled(14)
		bullet.Vel.Add(game.Player.Vel)
		game.Shots = append(game.Shots, bullet)
	case SHOTGUN:
		angle := math.Atan2(mousePosn.Y-game.Player.Posn.Y, mousePosn.X-game.Player.Posn.X)
		for shotCount := 0; shotCount < 5; shotCount++ {
			var bullet Shot
			bullet.GunType = 2
			bullet.color = pixel.ToRGBA(colornames.Firebrick).Mul(pixel.Alpha(.7))
			// offset := pixel.V(rand.Float64()*20, rand.Float64()*20)
			offset := (rand.Float64() - rand.Float64()) / 2.3
			// newDirection := directionVec.Add(offset)
			newAngle := angle + offset
			// newDirection = newDirection.Scaled(1 / magnitude(newDirection))
			bullet.Posn1 = game.Player.Posn
			newDirection := pixel.V(math.Cos(newAngle), math.Sin(newAngle))
			bullet.Posn2 = game.Player.Posn.Add(newDirection.Scaled(10))
			bullet.Vel = (bullet.Posn2.Sub(bullet.Posn1)).Scaled(2 + (rand.Float64()-rand.Float64())/2.3)
			game.Shots = append(game.Shots, bullet)
		}
		game.Player.Bullets[SHOTGUN]--
		if game.Player.Bullets[SHOTGUN] <= 0 {
			game.Player.GunType = 1
		}
	}
}

// PressHandler Handles key presses.
// Agent method taking in a window from which to accept inputs.
func PressHandler(win *pixelgl.Window, game *Game) {
	game.Player.Acc = pixel.ZV
	if win.Pressed(pixelgl.KeyA) {
		game.Player.Acc = game.Player.Acc.Sub(pixel.V(5, 0))
	}
	if win.Pressed(pixelgl.KeyD) {
		game.Player.Acc = game.Player.Acc.Add(pixel.V(5, 0))
	}
	if win.Pressed(pixelgl.KeyS) {
		game.Player.Acc = game.Player.Acc.Sub(pixel.V(0, 5))
	}
	if win.Pressed(pixelgl.KeyW) {
		game.Player.Acc = game.Player.Acc.Add(pixel.V(0, 5))
	}
	if win.Pressed(pixelgl.Key2) {
		if game.Player.Bullets[SHOTGUN] > 0 {
			game.Player.GunType = SHOTGUN
		}
	}
	if win.Pressed(pixelgl.Key1) {
		game.Player.GunType = 1
	}
	// if game.Player.devMode {
	// 	if win.JustPressed(pixelgl.KeyJ) {
	// 		game.Player.Posn.X--
	// 	}
	// 	if win.JustPressed(pixelgl.KeyL) {
	// 		game.Player.Posn.X++
	// 	}
	// 	if win.JustPressed(pixelgl.KeyK) {
	// 		game.Player.Posn.Y--
	// 	}
	// 	if win.JustPressed(pixelgl.KeyI) {
	// 		game.Player.Posn.Y++
	// 	}
	// 	if win.JustPressed(pixelgl.KeySpace) {
	// 		game.Player.Shade = !game.Player.Shade
	// 	}
	// 	if win.JustPressed(pixelgl.KeyF) {
	// 		game.Player.Fill = !game.Player.Fill
	// 	}
	// }
	if win.JustPressed(pixelgl.KeyP) {
		game.Paused = true
	}
	if win.JustPressed(pixelgl.MouseButton1) {
		game.fire(win)
	}
}

// ReleaseHandler Handles key releases.
// Agent method taking in a window from which to accept inputs.
func ReleaseHandler(win *pixelgl.Window, game *Game) {
	if win.JustReleased(pixelgl.KeyA) {
		game.Player.Acc = game.Player.Acc.Add(pixel.V(5, 0))
	}
	if win.JustReleased(pixelgl.KeyD) {
		game.Player.Acc = game.Player.Acc.Sub(pixel.V(5, 0))
	}
	if win.JustReleased(pixelgl.KeyS) {
		game.Player.Acc = game.Player.Acc.Add(pixel.V(0, 5))
	}
	if win.JustReleased(pixelgl.KeyW) {
		game.Player.Acc = game.Player.Acc.Sub(pixel.V(0, 5))
	}
}
