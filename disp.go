package main

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Disp updates and displays a room's Img
func (room *Place) Disp() {
	room.DispBoost()
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

// DispBoost draws the boost to the room's target, if present.
func (room *Place) DispBoost() {
	if room.Booster.Present {
		room.Booster.Img.Clear()
		room.Booster.Img.Push(room.Booster.Posn)
		room.Booster.Img.Circle(10, 0)
		room.Booster.Img.Draw(room.Target)
	}
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
