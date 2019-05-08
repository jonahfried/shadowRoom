package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// A Game represents the current game in progress.
type Game struct {
	Player Agent
	Room   Place

	Monsters []Creature
	Shots    []Shot
	PowerUps []PowerUp

	DevMode bool
	Level   float64
	Spacing float64
	Count   int
	Paused  bool
}

func (g Game) getRoom() Place {
	return g.Room
}

func makeGame(win *pixelgl.Window, devMode bool) (g Game) {
	var p = MakeAgent(0, 0, win) // player

	// TODO: Make room bounds relative to Bounds
	room := MakePlace(pixel.R(-700, -600, 700, 600), 11)
	room.Target.SetMatrix(pixel.IM.Moved(room.Target.Bounds().Center())) //Move this out of loop?
	room.ToGrid(40)
	g.Player = p
	g.Room = room

	g.Monsters = make([]Creature, 0)
	g.Shots = make([]Shot, 0)
	g.PowerUps = make([]PowerUp, 0)

	g.DevMode = devMode
	g.Level = .02
	g.Spacing = 6
	g.Count = 88
	g.Paused = false

	return g
}

func (g *Game) update() {
	g.updateMonsters()
	g.updateShots()
	g.updatePowerUps()
	if g.Player.TorchLevel > 2.5 {
		g.Player.TorchLevel *= 1 - 1e-4 //.9995
	}
}

func (g *Game) updatePowerUps() {
	for i := len(g.PowerUps) - 1; i >= 0; i-- {
		if g.PowerUps[i].shouldApply(&g.Player) {
			g.PowerUps[i].apply(&g.Player)
			removeAtInd(&g.PowerUps, i)
		}
	}
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
				g.Monsters[monsterInd].Img.Color = colornames.Red
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

		if g.Shots[bulletInd].Vel.Len() < 2 {
			g.Shots[bulletInd] = g.Shots[len(g.Shots)-1]
			g.Shots = g.Shots[:len(g.Shots)-1]
			bulletInd--
			continue BulletLoop
		}

	}

	removeDead(&g.Monsters)
}

func removeAtInd(lst *[]PowerUp, ind int) {
	if ind > len(*lst) {
		return
	}
	(*lst)[ind] = (*lst)[len(*lst)-1]
	*lst = (*lst)[:len(*lst)-1]
}
