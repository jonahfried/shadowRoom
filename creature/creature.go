package creature

import (
	"shadowRoom/boundry"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Creature is a struct containing information on its
// - current position
// - current velocity
// - visual representation (*imdraw.IMDraw)
type Creature struct {
	Posn, Vel pixel.Vec

	Img *imdraw.IMDraw
}

// MakeCreature takes a starting x and y (float64), and returns a creature
func MakeCreature(x, y float64) (monster Creature) {
	monster.Posn = pixel.V(x, y)
	monster.Vel = pixel.V(3, 4)

	monster.Img = imdraw.New(nil)
	monster.Img.Color = colornames.Darkseagreen

	return monster
}

// Update is a method for a creature, taking in a room
// returning nothing, it alters the position and velocity of the creature
func (monster *Creature) Update(room boundry.Place) {
	monster.Posn = monster.Posn.Add(monster.Vel)

	if monster.Posn.X > (room.Rect.Max.X - 20) {
		monster.Posn.X = (room.Rect.Max.X - 20)
		monster.Vel.X *= -1
	}
	if monster.Posn.X < (room.Rect.Min.X + 20) {
		monster.Posn.X = (room.Rect.Min.X + 20)
		monster.Vel.X *= -1
	}
	if monster.Posn.Y > (room.Rect.Max.Y - 20) {
		monster.Posn.Y = (room.Rect.Max.Y - 20)
		monster.Vel.Y *= -1
	}
	if monster.Posn.Y < (room.Rect.Min.Y + 20) {
		monster.Posn.Y = (room.Rect.Min.Y + 20)
		monster.Vel.Y *= -1
	}
}

// Disp draws a creature based on its Img
func (monster *Creature) Disp(win *pixelgl.Window) {
	monster.Img.Clear()
	monster.Img.Push(monster.Posn)
	monster.Img.Circle(20, 0)
	monster.Img.Draw(win)
}
