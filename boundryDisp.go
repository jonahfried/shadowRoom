package main

import (
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
