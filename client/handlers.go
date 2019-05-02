package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"golang.org/x/net/websocket"
)

func (game *Game) fire(win *pixelgl.Window) {
	mousePosn := game.Player.Cam.Matrix.Unproject(win.MousePosition())
	directionVec := mousePosn.Sub(game.Player.Posn)
	directionVec = directionVec.Scaled(1 / magnitude(directionVec))

	switch game.Player.GunType {
	case 1:
		var bullet Shot
		bullet.GunType = 1
		bullet.color = pixel.ToRGBA(colornames.Firebrick).Mul(pixel.Alpha(.7))
		bullet.Posn1 = game.Player.Posn
		bullet.Posn2 = game.Player.Posn.Add(directionVec.Scaled(10))
		// bullet.Vel.Add(game.Player.Vel)
		bullet.Vel = directionVec.Scaled(14)
		game.Shots = append(game.Shots, bullet)
	case 2:
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
		game.Player.Bullets--
		if game.Player.Bullets <= 0 {
			game.Player.GunType = 1
		}
	}
}

// PressHandler Handles key presses.
// Agent method taking in a window from which to accept inputs.
func PressHandler(win *pixelgl.Window, ws *websocket.Conn) {
	if win.JustPressed(pixelgl.KeyA) {
		msg := "a"
		websocket.JSON.Send(ws, msg)
	}
	if win.JustPressed(pixelgl.KeyD) {
		msg := "d"
		websocket.JSON.Send(ws, msg)
	}
	if win.JustPressed(pixelgl.KeyS) {
		websocket.JSON.Send(ws, "s")
	}
	if win.JustPressed(pixelgl.KeyW) {
		msg := "w"
		websocket.JSON.Send(ws, msg)
	}

	// if win.JustPressed(pixelgl.MouseButton1) {
	// 	game.fire(win)
	// }
}

// ReleaseHandler Handles key releases.
// Agent method taking in a window from which to accept inputs.
func ReleaseHandler(win *pixelgl.Window, ws *websocket.Conn) {
	if win.Pressed(pixelgl.KeyA) {
		websocket.Message.Send(ws, "a")
	}
	if win.Pressed(pixelgl.KeyD) {
		websocket.Message.Send(ws, "d")
	}
	if win.Pressed(pixelgl.KeyS) {
		websocket.Message.Send(ws, "s")
	}
	if win.Pressed(pixelgl.KeyW) {
		websocket.Message.Send(ws, "w")
	}
}