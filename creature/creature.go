package creature

import (
	"math"
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
	monster.Vel = pixel.V(0, 0)

	monster.Img = imdraw.New(nil)
	monster.Img.Color = colornames.Darkolivegreen

	return monster
}

// Vector limitation. Takes in a pixel.Vec and a float64.
// If the magnitude of the given vector is greater than the limit,
// Then the magnitude is scaled down to the limit.
func limitVecMag(v pixel.Vec, lim float64) pixel.Vec {
	if v.Len() != 0 && v.Len() > lim {
		v = v.Scaled(lim / v.Len())
	}
	return v
}

func vecDist(v1, v2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
}

func magnitude(vec pixel.Vec) float64 {
	return vecDist(pixel.ZV, vec)
}

// Update is a method for a creature, taking in a room
// returning nothing, it alters the position and velocity of the creature
func (monster *Creature) Update(room boundry.Place, target pixel.Vec) {
	acc := target.Sub(monster.Posn) //limitVecMag(target.Sub(monster.Posn), vecDist(monster.Posn, target)/10)

	acc = acc.Sub(monster.Vel)
	acc = acc.Scaled(.095)
	monster.Vel = monster.Vel.Add(acc)
	monster.Vel = limitVecMag(monster.Vel, 6)
	// monster.Vel = monster.Vel.Scaled(1.3)

	center := monster.Posn
	radius := 20.0
	for _, obst := range room.Blocks {

		for vertexInd := 1; vertexInd < len(obst.Vertices); vertexInd++ {
			end1 := obst.Vertices[vertexInd]
			end2 := obst.Vertices[vertexInd-1]
			segVec := end2.Sub(end1)
			unitSegVec := segVec.Scaled(1 / magnitude(segVec))
			centerOffset := center.Sub(end1)
			projMag := centerOffset.Dot(unitSegVec)
			projVec := unitSegVec.Scaled(projMag)
			var closest pixel.Vec

			if projMag < 0 {
				closest = end1
			} else if projMag > magnitude(segVec) {
				closest = end2
			} else {
				closest = end1.Add(projVec)
			}

			dist := center.Sub(closest)

			if magnitude(dist) < (radius + 4) {
				offset := (dist.Scaled(1 / magnitude(dist))).Scaled((radius + 4) - magnitude(dist))
				offset = offset.Scaled(.30)
				monster.Vel = (monster.Vel.Add(offset)).Scaled(.88)
			}
		}

		end1 := obst.Vertices[0]
		end2 := obst.Vertices[len(obst.Vertices)-1]
		segVec := end2.Sub(end1)
		unitSegVec := segVec.Scaled(1 / magnitude(segVec))
		centerOffset := center.Sub(end1)
		projMag := centerOffset.Dot(unitSegVec)
		projVec := unitSegVec.Scaled(projMag)
		var closest pixel.Vec

		if projMag < 0 {
			closest = end1
		} else if projMag > magnitude(segVec) {
			closest = end2
		} else {
			closest = end1.Add(projVec)
		}

		dist := center.Sub(closest)

		if magnitude(dist) < (radius + 4) {
			offset := (dist.Scaled(1 / magnitude(dist))).Scaled((radius + 4) - magnitude(dist))
			offset = offset.Scaled(.30)
			monster.Vel = (monster.Vel.Add(offset)).Scaled(.88)
		}

	}

	monster.Posn = monster.Posn.Add(monster.Vel)

	// monster.Posn = monster.Posn.Add(monster.Vel)
	// if monster.Posn.X > (room.Rect.Max.X - 20) {
	// 	monster.Posn.X = (room.Rect.Max.X - 20)
	// 	monster.Vel.X *= -1
	// }
	// if monster.Posn.X < (room.Rect.Min.X + 20) {
	// 	monster.Posn.X = (room.Rect.Min.X + 20)
	// 	monster.Vel.X *= -1
	// }
	// if monster.Posn.Y > (room.Rect.Max.Y - 20) {
	// 	monster.Posn.Y = (room.Rect.Max.Y - 20)
	// 	monster.Vel.Y *= -1
	// }
	// if monster.Posn.Y < (room.Rect.Min.Y + 20) {
	// 	monster.Posn.Y = (room.Rect.Min.Y + 20)
	// 	monster.Vel.Y *= -1
	// }
}

// Disp draws a creature based on its Img
func (monster *Creature) Disp(win *pixelgl.Canvas) {
	monster.Img.Clear()
	monster.Img.Push(monster.Posn)
	monster.Img.Circle(20, 0)
	monster.Img.Draw(win)
}
