package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"golang.org/x/net/websocket"
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

func makeGame(win *pixelgl.Window, devMode bool, ws *websocket.Conn) (g Game) {
	var p = MakeAgent(0, 0, win) // player

	// TODO: Make room bounds relative to Bounds
	g.Player = p
	g.Room = getInitRoomFromServer(ws)

	g.Shots = make([]Shot, 0)
	g.Monsters = make([]Creature, 0)

	g.DevMode = devMode
	g.Level = .02
	g.Spacing = 6
	g.Count = 88

	return g
}

func getInitRoomFromServer(ws *websocket.Conn) Place {
	websocket.JSON.Send(ws, "initState")
	var room Place
	websocket.JSON.Receive(ws, &room)
	blocks := make([]Obstacle, len(room.Blocks))
	for ind, block := range room.Blocks {
		blocks[ind] = MakeObstacle(block.Vertices, block.Center, block.Radius)
	}
	newRoom := MakePlace(room.Rect, len(blocks)-1, blocks...)
	return newRoom
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
