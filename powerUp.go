package main

import (
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
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

// Shotgun implements PowerUp and provides the player with 10 uses of the SHOTGUN weapon
type Shotgun struct {
	Posn   pixel.Vec
	Color  color.RGBA
	Radius float64
}

// ShotgunWeight is the weight for determining spawn prob
const ShotgunWeight = 1

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

// Torch increases the player light source
type Torch struct {
	Posn   pixel.Vec
	Color  color.RGBA
	Radius float64
}

// TorchWeight is the weight for determining spawn prob
const TorchWeight = 1

// MakeTorch returns a Torch
func MakeTorch(posn pixel.Vec) (t Torch) {
	t.Posn = posn
	t.Color = colornames.Red
	t.Radius = 10
	return t
}

func (t Torch) shouldApply(player *Agent) bool {
	return vecDist(t.Posn, player.Posn) <= player.Radius+t.Radius
}
func (t Torch) apply(player *Agent) {
	if player.TorchLevel <= 10 {
		// player.TorchLevel += 3
		player.Torches += 1
	}
}
func (t Torch) isOver(player *Agent)   {}
func (t Torch) tearDown(player *Agent) {}
func (t Torch) color() color.RGBA {
	return t.Color
}
func (t Torch) position() pixel.Vec {
	return t.Posn
}
