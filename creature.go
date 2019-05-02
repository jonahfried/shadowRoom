package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

// Creature is a struct containing information on its
// - current position
// - current velocity
// - visual representation (*imdraw.IMDraw)
type Creature struct {
	Posn, Vel pixel.Vec
	Health    int

	Img *imdraw.IMDraw
}

// MakeCreature takes a starting x and y (float64), and returns a *Creature
func MakeCreature(room *Place, cir *Agent) (monster Creature) {
	xPosn := room.Rect.Center().X + (rand.Float64()-rand.Float64())*room.Rect.W()/2
	yPosn := room.Rect.Center().Y + (rand.Float64()-rand.Float64())*room.Rect.H()/2
	posn := pixel.V(xPosn, yPosn)
TryLoop:
	for tries := 0; tries < 10; tries++ {
		if vecDist(posn, cir.Posn) < 80 {
			xPosn = room.Rect.Center().X + (rand.Float64()-rand.Float64())*room.Rect.W()/2
			yPosn = room.Rect.Center().Y + (rand.Float64()-rand.Float64())*room.Rect.H()/2
			posn = pixel.V(xPosn, yPosn)
			continue TryLoop
		}
		for _, obst := range room.Blocks {
			if vecDist(posn, obst.Center) < (obst.Radius + 20) {
				xPosn = room.Rect.Center().X + (rand.Float64()-rand.Float64())*room.Rect.W()/2
				yPosn = room.Rect.Center().Y + (rand.Float64()-rand.Float64())*room.Rect.H()/2
				posn = pixel.V(xPosn, yPosn)
				continue TryLoop
			}
			break TryLoop
		}
	}
	monster.Posn = posn
	monster.Vel = pixel.V(0, 0)

	monster.Health = 5
	monster.Img = imdraw.New(nil)
	monster.Img.Color = colornames.Darkolivegreen

	return monster
}

// Update is a method for a creature, taking in a room
// returning nothing, it alters the position and velocity of the creature
func (monster *Creature) Update(room Place, cir *Agent, target pixel.Vec, monsters []Creature) {
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
			unitSegVec := segVec.Scaled(1 / segVec.Len())
			centerOffset := center.Sub(end1)
			projMag := centerOffset.Dot(unitSegVec)
			projVec := unitSegVec.Scaled(projMag)
			var closest pixel.Vec

			if projMag < 0 {
				closest = end1
			} else if projMag > segVec.Len() {
				closest = end2
			} else {
				closest = end1.Add(projVec)
			}

			dist := center.Sub(closest)

			if dist.Len() < (radius + 4) {
				offset := (dist.Scaled(1 / dist.Len())).Scaled((radius + 4) - dist.Len())
				offset = offset.Scaled(.30)
				monster.Vel = (monster.Vel.Add(offset)).Scaled(.88)
			}
		}

		end1 := obst.Vertices[0]
		end2 := obst.Vertices[len(obst.Vertices)-1]
		segVec := end2.Sub(end1)
		unitSegVec := segVec.Scaled(1 / segVec.Len())
		centerOffset := center.Sub(end1)
		projMag := centerOffset.Dot(unitSegVec)
		projVec := unitSegVec.Scaled(projMag)
		var closest pixel.Vec

		if projMag < 0 {
			closest = end1
		} else if projMag > segVec.Len() {
			closest = end2
		} else {
			closest = end1.Add(projVec)
		}

		dist := center.Sub(closest)

		if dist.Len() < (radius + 4) {
			offset := (dist.Scaled(1 / dist.Len())).Scaled((radius + 4) - dist.Len())
			offset = offset.Scaled(.30)
			monster.Vel = (monster.Vel.Add(offset)).Scaled(.88)
		}

	}

	playerDist := cir.Posn.Sub(monster.Posn)
	if playerDist.Len() < 40 && playerDist.Len() > 0 {
		changeBy := (40 - playerDist.Len()) / 2
		monster.Vel = monster.Vel.Sub(playerDist.Scaled(2 * changeBy / playerDist.Len()))
		cir.Vel = cir.Vel.Add(playerDist.Scaled(changeBy / playerDist.Len()))
		cir.Img.Color = colornames.Red
		cir.Health--
		// monsters[blobInd].Posn = monsters[blobInd].Posn.Sub(dist)
	}

	for blobInd := range monsters {
		dist := monsters[blobInd].Posn.Sub(monster.Posn)
		if dist.Len() < 40 && dist.Len() > 0 {
			changeBy := (40 - dist.Len()) / 2
			monster.Vel = monster.Vel.Sub(dist.Scaled(changeBy / dist.Len()))
			monsters[blobInd].Vel = monsters[blobInd].Vel.Add(dist.Scaled(changeBy / dist.Len()))
			// monsters[blobInd].Posn = monsters[blobInd].Posn.Sub(dist)
		}
	}

	monster.Posn = monster.Posn.Add(monster.Vel)
	monster.Posn.X = math.Min(monster.Posn.X, room.Rect.Max.X-20)
	monster.Posn.X = math.Max(monster.Posn.X, room.Rect.Min.X+20)

	monster.Posn.Y = math.Min(monster.Posn.Y, room.Rect.Max.Y-20)
	monster.Posn.Y = math.Max(monster.Posn.Y, room.Rect.Min.Y+20)
}

func (game *Game) updateMonsters() {
	for monsterInd := range game.Monsters {
		monsterTarget := game.AStar(game.Monsters[monsterInd].Posn, game.Player.Posn)
		game.Monsters[monsterInd].Update(game.Room, &game.Player, monsterTarget, game.Monsters)
	}
}

// // a runner is a kind of monster that "leaps" at a target after some time "tracking" them
// type runner struct {
// 	Posn, Vel                  pixel.Vec
// 	Health                     int
// 	Angle, TrackTime, AngleVel float64

// 	Img *imdraw.IMDraw
// }

// // Update updates a runner type monster
// func (monster *runner) Update(target pixel.Vec, dt float64) {
// 	if monster.Len()Vel) < .1 {
// 		monster.Vel = pixel.ZV
// 	}
// 	if monster.TrackTime < 10 {
// 		monster.TrackTime += dt
// 		targetAngle := math.Atan2(monster.Posn.Y-target.Y, monster.Posn.X-target.X)
// 		goalChange := monster.Angle - targetAngle
// 		monster.AngleVel += goalChange * .3
// 	}
// }

// // Disp displays a runner
// func (monster runner) Disp(target pixel.Target) {
// 	dirVec := pixel.V(math.Cos(monster.Angle), math.Sin(monster.Angle))
// 	monster.Img.Clear()
// 	monster.Img.Push(monster.Posn.Add(dirVec))

// 	dirVec =
// 		monster.Img.Push(monster.Posn.Add(dirVec))
// 	monster.Img.Push(monster.Posn.Add(dirVec))
// }
