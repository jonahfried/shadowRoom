package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Disp updates and displays a room's Img
func (room *Place) Disp() {
	room.Img.Clear()
	room.Img.Color = colornames.Burlywood
	room.Img.Push(room.Rect.Center().Sub(room.Rect.Size().Scaled(.5)))
	room.Img.Push(room.Rect.Center().Add(room.Rect.Size().Scaled(.5)))
	room.Img.Rectangle(3)
	room.Img.Draw(room.Target)

	for _, bc := range room.Blocks {
		bc.Img.Draw(room.Target)
	}

}

func dispPowerUps(powerUps []PowerUp, canvas *pixelgl.Canvas) {
	img := imdraw.New(nil)
	for _, powerUp := range powerUps {
		img.Color = powerUp.color()
		img.Push(powerUp.position())
		img.Circle(10, 0) // 10 assumed powerup radius
	}
	img.Draw(canvas)
}

// Disp Handles display of an agent. Clears, pushes, adds shape, and draws.
func (cir *Agent) Disp(win *pixelgl.Window) {
	cir.Img.Clear()
	cir.Img.Push(cir.Posn)
	cir.Img.Circle(20, 0)
	cir.Img.Draw(win)
	cir.Img.Color = colornames.Purple
}

// DispShots displays shots
func (game *Game) DispShots(canv *pixelgl.Canvas) {
	img := imdraw.New(nil)
	for _, bullet := range game.Shots {
		img.Color = bullet.color
		img.Push(bullet.Posn1)
		img.Push(bullet.Posn2)
		img.Line(4)
	}
	img.Draw(canv)
}

// Disp draws a creature based on its Img
func (monster *Creature) Disp(win *pixelgl.Canvas) {
	monster.Img.Clear()
	monster.Img.Push(monster.Posn)
	monster.Img.Circle(20, 0)
	monster.Img.Draw(win)
	monster.Img.Color = colornames.Darkolivegreen
}

// Disp draws all necesary pieces of the Game
func (game *Game) Disp(win *pixelgl.Window) {
	win.Clear(colornames.Black)
	game.Room.Target.Clear(pixel.Alpha(0))
	game.Player.playerTorch(game.Level, game.Count, game.Spacing, &game.Room)
	game.Room.Disp()
	for monsterInd := range game.Monsters {
		game.Monsters[monsterInd].Disp(game.Room.Target)
	}
	game.DispShots(game.Room.Target)
	dispPowerUps(game.PowerUps, game.Room.Target)
	game.Room.Target.Draw(win, pixel.IM) //.Moved(win.Bounds().Center()))
	game.Player.Disp(win)
	illuminate(game.Room, game.Player, win)
	// if game.DevMode {
	// 	fpsDisp(frames/seconds, game.Player.Posn, win)
	// }
}
