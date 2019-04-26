package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// A PowerUp represents any type that gives some sort of bonus to the player
type PowerUp interface {
	// Applies the PowerUp's effect to the player
	apply(player *Agent)
	// Alerts when the PowerUp effect ends
	isOver(player *Agent)
	// Applies the necessary effects for removal of the effect
	tearDown(player *Agent)
}

// Boost to give bonus to the player
type Boost struct {
	Posn    pixel.Vec
	Present bool

	Img *imdraw.IMDraw
}

func (room *Place) presentBoost() {
	if !room.Booster.Present {
		room.Booster.Posn = pixel.V(room.Rect.Center().X+(room.Rect.W()/2*(rand.Float64()-rand.Float64())/2), room.Rect.Center().Y+room.Rect.H()*(rand.Float64()-rand.Float64())/2)
		room.Booster.Present = true
	} else {
		room.Booster.Present = false
	}
}
