package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// A Game represents the current game in progress.
type Game struct {
	Player Agent
	Room   *Place
}

func (g Game) getRoom() Place {
	return *g.Room
}

func makeGame(win *pixelgl.Window, devMode bool) (g Game) {
	var cir = MakeAgent(0, 0, win, devMode) // player

	// TODO: Make room bounds relative to Bounds
	room := MakePlace(pixel.R(-700, -600, 700, 600), 11)
	room.Target.SetMatrix(pixel.IM.Moved(room.Target.Bounds().Center())) //Move this out of loop?
	room.ToGrid(40)
	g.Player = cir
	g.Room = &room
	return g
}
