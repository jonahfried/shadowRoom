package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// A Game represents the current game in progress.
type Game struct {
	Player *Agent
	Room   *Place

	DevMode bool
	Level   float64
	Spacing int
	Count   int
}

func (g Game) getRoom() Place {
	return *g.Room
}

func makeGame(win *pixelgl.Window, devMode bool) (g Game) {
	var p = MakeAgent(0, 0, win, devMode) // player

	// TODO: Make room bounds relative to Bounds
	room := MakePlace(pixel.R(-700, -600, 700, 600), 11)
	room.Target.SetMatrix(pixel.IM.Moved(room.Target.Bounds().Center())) //Move this out of loop?
	room.ToGrid(40)
	g.Player = &p
	g.Room = &room

	g.DevMode = devMode
	g.Level = .02
	g.Spacing = 6
	g.Count = 88

	return g
}
