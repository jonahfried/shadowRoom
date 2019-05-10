package main

import (
	"math"

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
	Radius    float64
	Health    int

	Img *imdraw.IMDraw
}

// MakeCreature takes a starting x and y (float64), and returns a *Creature
func MakeCreature(room *Place, cir *Agent) (monster Creature) {
	monster.Radius = 20
	monster.Posn = room.safeSpawnInRoom(monster.Radius)
	monster.Vel = pixel.ZV
	monster.Health = 5

	monster.Img = imdraw.New(nil)
	monster.Img.Color = colornames.Darkolivegreen

	return monster
}

// Update is a method for a creature, taking in a room
// returning nothing, it alters the position and velocity of the creature
func (monster *Creature) Update(room Place, cir *Agent, target pixel.Vec, monsters *[]Creature) {
	acc := target.Sub(monster.Posn) //limitVecMag(target.Sub(monster.Posn), vecDist(monster.Posn, target)/10)

	acc = acc.Sub(monster.Vel)
	acc = acc.Scaled(.095)
	monster.Vel = monster.Vel.Add(acc)
	monster.Vel = limitVecMag(monster.Vel, 6)
	// monster.Vel = monster.Vel.Scaled(1.3)

	offset := collision(room.Blocks, monster.Posn, monster.Radius, &monster.Vel)
	monster.Vel = monster.Vel.Add(offset.Scaled(.2))

	monster.playerCollision(cir)
	monster.creatureCollision(monsters)

	monster.Posn = monster.Posn.Add(monster.Vel)
	monster.Posn.X = math.Min(monster.Posn.X, room.Rect.Max.X-20)
	monster.Posn.X = math.Max(monster.Posn.X, room.Rect.Min.X+20)

	monster.Posn.Y = math.Min(monster.Posn.Y, room.Rect.Max.Y-20)
	monster.Posn.Y = math.Max(monster.Posn.Y, room.Rect.Min.Y+20)
}

func (monster *Creature) playerCollision(cir *Agent) {
	minRadius := monster.Radius + cir.Radius
	playerDist := cir.Posn.Sub(monster.Posn)
	if playerDist.Len() < minRadius && playerDist.Len() > 0 {
		changeBy := (minRadius - playerDist.Len())
		monster.Vel = monster.Vel.Sub(playerDist.Unit().Scaled(changeBy))
		cir.Vel = cir.Vel.Add(playerDist.Unit().Scaled(changeBy))
		cir.Img.Color = colornames.Red
		cir.Health--
	}
}

func (monster *Creature) creatureCollision(monsters *[]Creature) {
	for blobInd := range *monsters {
		minRadius := monster.Radius + (*monsters)[blobInd].Radius
		dist := (*monsters)[blobInd].Posn.Sub(monster.Posn)
		if dist.Len() < minRadius && dist.Len() > 0 {
			changeBy := (minRadius - dist.Len()) / 2
			monster.Vel = monster.Vel.Sub(dist.Scaled(changeBy / dist.Len()))
			(*monsters)[blobInd].Vel = (*monsters)[blobInd].Vel.Add(dist.Scaled(changeBy / dist.Len()))
			// monsters[blobInd].Posn = monsters[blobInd].Posn.Sub(dist)
		}
	}
}

func (game *Game) updateMonsters() {
	for monsterInd := range game.Monsters {
		monsterTarget := game.AStar(game.Monsters[monsterInd].Posn, game.Player.Posn)
		game.Monsters[monsterInd].Update(game.Room, &game.Player, monsterTarget, &game.Monsters)
	}
}
