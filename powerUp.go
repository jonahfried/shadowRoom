package main

import (
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// A PowerUp represents any type that gives some sort of bonus to the player
type PowerUp interface {
	// Returns whether the PowerUp should be applied to the player
	shouldApply(player *Agent) bool
	// Applies the PowerUp's effect to the player
	apply(player *Agent)
	// Alerts when the PowerUp effect ends
	isOver(player *Agent)
	// Applies the necessary effects for removal of the effect
	tearDown(player *Agent)
	position() pixel.Vec
	color() color.RGBA
}

// Boost to give bonus to the player
type Boost struct {
	Posn    pixel.Vec
	Present bool

	Img *imdraw.IMDraw
}

// func (room *Place) presentBoost() {
// 	if !room.Booster.Present {
// 		room.Booster.Posn = room.safeSpawnInRoom(5) //pixel.V(room.Rect.Center().X+(room.Rect.W()/2*(rand.Float64()-rand.Float64())/2), room.Rect.Center().Y+room.Rect.H()*(rand.Float64()-rand.Float64())/2)
// 		room.Booster.Present = true
// 	} else {
// 		room.Booster.Present = false
// 	}
// }

// Shotgun implements PowerUp and provides the player with 10 uses of the SHOTGUN weapon
type Shotgun struct {
	Posn   pixel.Vec
	Color  color.RGBA
	Radius float64
}

// MakeShotgun returns a Shotgun
func MakeShotgun(posn pixel.Vec) (shot Shotgun) {
	shot.Posn = posn
	shot.Color = colornames.Blue
	shot.Radius = 10
	return shot
}

// SHOTGUN is the number associated with the shotgun
const SHOTGUN = 2

func (shot Shotgun) shouldApply(player *Agent) bool {
	return vecDist(shot.Posn, player.Posn) <= player.Radius+shot.Radius
}
func (shot Shotgun) apply(player *Agent) {
	player.Bullets[SHOTGUN] += 10
}
func (shot Shotgun) isOver(player *Agent)   {}
func (shot Shotgun) tearDown(player *Agent) {}
func (shot Shotgun) color() color.RGBA {
	return shot.Color
}
func (shot Shotgun) position() pixel.Vec {
	return shot.Posn
}
