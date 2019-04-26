package main

import (
	"github.com/faiface/pixel"
)

// A Game represents the current game in progress.
type Game struct {
	Player Agent
	Room   Place

	Monsters []Creature
	Shots    []Shot

	DevMode bool
	Level   float64
	Spacing int
	Count   int
}

func (g Game) getRoom() Place {
	return g.Room
}

func makeGame(devMode bool) (g Game) {
	var p = MakeAgent(0, 0) // player

	// TODO: Make room bounds relative to Bounds
	room := MakePlace(pixel.R(-700, -600, 700, 600), 11)
	room.ToGrid(40)
	g.Player = p
	g.Room = room

	// g.Shots = make([]Shot, 0)
	// g.Monsters = make([]Creature, 0)

	g.DevMode = devMode
	g.Level = .02
	g.Spacing = 6
	g.Count = 88

	return g
}

func removeDead(monsters *[]Creature) {
	for monsterInd := 0; monsterInd < len(*monsters); monsterInd++ {
		if (*monsters)[monsterInd].Health <= 0 {
			(*monsters)[monsterInd] = (*monsters)[len(*monsters)-1]
			(*monsters) = (*monsters)[:len(*monsters)-1]
		}
	}
}

func (g *Game) updateShots() {
BulletLoop:
	for bulletInd := 0; bulletInd < len(g.Shots); bulletInd++ { //range g.Shots {
		g.Shots[bulletInd].update()

		for monsterInd := range g.Monsters {
			if vecDist(g.Shots[bulletInd].Posn2, g.Monsters[monsterInd].Posn) < 20 {
				g.Monsters[monsterInd].Health--
				g.Shots[bulletInd] = g.Shots[len(g.Shots)-1]
				g.Shots = g.Shots[:len(g.Shots)-1]
				bulletInd--
				continue BulletLoop
			}
		}

		// Does more checks than necesarry in the case that it does collide.
		if !(g.Room.Rect.Contains(g.Shots[bulletInd].Posn2)) || (collision(g.Room.Blocks, g.Shots[bulletInd].Posn2, 1) != pixel.ZV) {
			g.Shots[bulletInd] = g.Shots[len(g.Shots)-1]
			g.Shots = g.Shots[:len(g.Shots)-1]
			bulletInd--
			continue BulletLoop
		}

		if magnitude(g.Shots[bulletInd].Vel) < 2 {
			g.Shots[bulletInd] = g.Shots[len(g.Shots)-1]
			g.Shots = g.Shots[:len(g.Shots)-1]
			bulletInd--
			continue BulletLoop
		}

	}

	removeDead(&g.Monsters)
}

func removeAtInd(lst *[]interface{}, ind int) {
	if ind > len(*lst) {
		return
	}
	(*lst)[ind] = (*lst)[len(*lst)-1]
	*lst = (*lst)[:len(*lst)-1]
}

func (g *Game) updateGame() {

}
