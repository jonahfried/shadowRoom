package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Agent keeps a posn, vel, acc, and a *pixel.IMDraw.
// Used to keep all information for the movable image together.
type Agent struct {
	Posn, Vel, Acc pixel.Vec

	Shade bool
	Fill  bool

	// Dev Tools
	Level   float64
	Spacing int
	Count   int
	devMode bool

	Cam Camera

	Health int

	Monsters []Creature
	Shots    []Shot
	ShotsImg *imdraw.IMDraw
	GunType  int
	Bullets  int

	Img *imdraw.IMDraw
}

// MakeAgent creates a new agent starting at a given (x, y) coordinate
func MakeAgent(x, y float64, win *pixelgl.Window, devMode bool) (cir Agent) {
	cir.Posn = pixel.V(x, y)
	cir.Vel = pixel.ZV
	cir.Acc = pixel.ZV

	cir.Shade = true
	cir.Fill = true

	cir.Cam = MakeCamera(cir.Posn, win)

	cir.Level = 0.02
	cir.Spacing = 6
	cir.Count = 88
	cir.devMode = devMode

	cir.Health = 100

	cir.Monsters = make([]Creature, 0)
	cir.Shots = make([]Shot, 0)
	cir.ShotsImg = imdraw.New(nil)
	cir.GunType = 1
	cir.Bullets = 0

	cir.Img = imdraw.New(nil)
	cir.Img.Color = colornames.Purple
	cir.Img.Push(pixel.V(500, 350))
	cir.Img.Circle(20, 0)
	return cir
}

// A Shot keeps track of the players attacks
type Shot struct {
	Posn1, Posn2, Vel pixel.Vec
	GunType           int
	color             pixel.RGBA
}

// DispShots displays shots
func (cir *Agent) DispShots(win *pixelgl.Canvas) {
	cir.ShotsImg.Clear()
	for _, bullet := range cir.Shots {
		cir.ShotsImg.Color = bullet.color
		cir.ShotsImg.Push(bullet.Posn1)
		cir.ShotsImg.Push(bullet.Posn2)
		cir.ShotsImg.Line(4)
	}
	cir.ShotsImg.Draw(win)
}

// Disp Handles display of an agent. Clears, pushes, adds shape, and draws.
func (cir *Agent) Disp(win *pixelgl.Window) {
	cir.Img.Clear()
	cir.Img.Push(cir.Posn)
	cir.Img.Circle(20, 0)
	cir.Img.Draw(win)
	cir.Img.Color = colornames.Purple
}

// Light adds fading light (white circles) around an Agent's posn
// func (cir *Agent) Light(room *boundry.Place) {
// 	room.Target.SetComposeMethod(pixel.ComposePlus)
// 	cir.Sprite.DrawColorMask(room.Target, pixel.IM.Moved(cir.Posn), pixel.Alpha(cir.Level))
// 	// cir.Sprite.Draw(room.Target, pixel.IM.Moved(cir.Posn))
// 	// room.Target.SetComposeMethod(pixel.ComposeIn)
// }

// Light adds fading light (white circles) around an Agent's posn
func (cir *Agent) Light(room *Place) {
	room.Target.SetComposeMethod(pixel.ComposeOver)
	img := imdraw.New(nil)
	img.Precision = 32
	// col := (pixel.ToRGBA(colornames.Whitesmoke)).Mul(pixel.Alpha(cir.Level))
	col := (pixel.ToRGBA(colornames.Cornsilk)).Mul(pixel.Alpha(cir.Level))
	for fade := 1; fade < cir.Count; fade++ {
		img.Color = col
		img.Push(cir.Posn)
		img.Circle(float64(fade*cir.Spacing), 0)
		// col = col.Mul(pixel.Alpha(1 / float64(fade)))
		// col = (pixel.ToRGBA(colornames.Whitesmoke)).Mul(pixel.Alpha(.8 / (float64(fade))))
		// col = col.Mul(pixel.Alpha(.95))

	}

	// room.Target.Clear(pixel.Alpha(0))
	// room.Target.SetComposeMethod()
	img.Draw(room.Target)
	room.Target.SetComposeMethod(pixel.ComposeIn)

}

// Checks to see if a given position should receive a collision force from a list of obstacles
func collision(blocks []Obstacle, posn pixel.Vec, radius float64) (force pixel.Vec) {
	for _, obst := range blocks {
		vertices := make([]pixel.Vec, 0, 10)
		vertices = append(vertices, obst.Vertices...)
		vertices = append(vertices, obst.Vertices[0]) // adding the first element to the end to complete the shape
		for vertexInd := 1; vertexInd < len(vertices); vertexInd++ {
			closest := closestPointOnSegment(vertices[vertexInd], vertices[vertexInd-1], posn)
			dist := posn.Sub(closest)
			if dist.Len() < (radius + 4) {
				offset := (dist.Scaled(1 / dist.Len())).Scaled((radius + 4) - dist.Len())
				offset = offset.Scaled(.30)
				force = force.Add((offset).Scaled(.88))
			}
		}

	}
	return force
}

// Returns the point on a line segment closest to a given position. (Either one of the two ends or a point between them)
func closestPointOnSegment(end1, end2, posn pixel.Vec) (cloestest pixel.Vec) {
	segVec := end2.Sub(end1)
	unitSegVec := segVec.Scaled(1 / magnitude(segVec))
	posnOffset := posn.Sub(end1)
	projMag := posnOffset.Dot(unitSegVec)
	projVec := unitSegVec.Scaled(projMag)
	var closest pixel.Vec

	if projMag < 0 {
		closest = end1
	} else if projMag > magnitude(segVec) {
		closest = end2
	} else {
		closest = end1.Add(projVec)
	}
	return closest
}

// Update is an agent method. Runs all necessary per-frame proccedures on agent.
// Takes in a pixelgl.Window from which to accept inputs.
func (cir *Agent) Update(win *pixelgl.Window, room *Place) {
	cir.PressHandler(win)
	cir.ReleaseHandler(win)

	// if vecDist(cir.Vel, pixel.ZV) < .2 {
	// 	cir.Vel = pixel.ZV
	// }

	offset := collision(room.Blocks, cir.Posn, 20)
	cir.Vel = cir.Vel.Add(offset)

	cir.Vel = limitVecMag(cir.Vel.Add(limitVecMag(cir.Acc, 1.5)), 10).Scaled(.88)

	cir.Posn = cir.Posn.Add(cir.Vel)
	// cir.Vel = limitVecMag(cir.Vel.Add(limitVecMag(cir.Acc, 1.5)), 10).Scaled(.88)
	// if vecDist(cir.Vel, pixel.ZV) < .2 {
	// 	cir.Vel = pixel.ZV
	// }

	cir.Posn.X = math.Max(cir.Posn.X-20, room.Rect.Min.X) + 20
	cir.Posn.X = math.Min(cir.Posn.X+20, room.Rect.Max.X) - 20

	cir.Posn.Y = math.Max(cir.Posn.Y-20, room.Rect.Min.Y) + 20
	cir.Posn.Y = math.Min(cir.Posn.Y+20, room.Rect.Max.Y) - 20

	if vecDist(room.Booster.Posn, cir.Posn) < 30 && room.Booster.Present {
		room.Booster.Present = false
		cir.Bullets += 10
		cir.GunType = 2
	}

	// Update Shots
BulletLoop:
	for bulletInd := 0; bulletInd < len(cir.Shots); bulletInd++ { //range cir.Shots {
		cir.Shots[bulletInd].Posn1 = cir.Shots[bulletInd].Posn1.Add(cir.Shots[bulletInd].Vel)
		cir.Shots[bulletInd].Posn2 = cir.Shots[bulletInd].Posn2.Add(cir.Shots[bulletInd].Vel)

		if cir.Shots[bulletInd].GunType == 2 {
			cir.Shots[bulletInd].Vel = cir.Shots[bulletInd].Vel.Scaled(.94)
			cir.Shots[bulletInd].color = (cir.Shots[bulletInd].color).Add(pixel.ToRGBA(colornames.Darkorange).Mul(pixel.Alpha(.03))).Add(pixel.Alpha(.02)) //  Mul(pixel.Alpha(1))
		}

		// midPoint := (cir.Shots[bulletInd].Posn1.Add(cir.Shots[bulletInd].Posn2)).Scaled(1 / 2)
		for monsterInd := range cir.Monsters {
			if vecDist(cir.Shots[bulletInd].Posn2, cir.Monsters[monsterInd].Posn) < 20 {
				// fmt.Println("shot")
				cir.Monsters[monsterInd].Health--
				cir.Monsters[monsterInd].Img.Color = colornames.Red
				// cir.Monsters[monsterInd].Vel = cir.Monsters[monsterInd].Vel.Add((cir.Shots[bulletInd].Vel).Scaled(1))
				cir.Shots[bulletInd] = cir.Shots[len(cir.Shots)-1]
				cir.Shots = cir.Shots[:len(cir.Shots)-1]
				bulletInd--
				continue BulletLoop
				// BAD
				// if len(cir.Shots) > 0 {
				// 	cir.Shots

				// }
			}
		}
		// Does more checks than necesarry in the case that it does collide.
		if !(room.Rect.Contains(cir.Shots[bulletInd].Posn2)) || (collision(room.Blocks, cir.Shots[bulletInd].Posn2, 1) != pixel.ZV) {
			cir.Shots[bulletInd] = cir.Shots[len(cir.Shots)-1]
			cir.Shots = cir.Shots[:len(cir.Shots)-1]
			bulletInd--
			continue BulletLoop
		}
		if magnitude(cir.Shots[bulletInd].Vel) < 2 {
			cir.Shots[bulletInd] = cir.Shots[len(cir.Shots)-1]
			cir.Shots = cir.Shots[:len(cir.Shots)-1]
			bulletInd--
			continue BulletLoop
		}
	}
	filter(&cir.Monsters)
}

// // TO PULL OUT (Update Shots):
// for bulletInd := range cir.Shots {
// 	cir.Shots[bulletInd].Posn1 = cir.Shots[bulletInd].Posn1.Add(cir.Shots[bulletInd].Vel)
// 	cir.Shots[bulletInd].Posn2 = cir.Shots[bulletInd].Posn2.Add(cir.Shots[bulletInd].Vel)

// 	// midPoint := (cir.Shots[bulletInd].Posn1.Add(cir.Shots[bulletInd].Posn2)).Scaled(1 / 2)
// 	for _, monster := range cir.Monsters {
// 		if vecDist(cir.Shots[bulletInd].Posn2, monster.Posn) < 20 {
// 			fmt.Println("shot")
// 			monster.Health--
// 			// BAD
// 			// if len(cir.Shots) > 0 {
// 			// 	cir.Shots[bulletInd] = cir.Shots[len(cir.Shots)-1]
// 			// }
// 		}
// 	}
// }
// cir.Monsters = filter(cir.Monsters)

func filter(monsters *[]Creature) {
	for monsterInd := 0; monsterInd < len(*monsters); monsterInd++ {
		if (*monsters)[monsterInd].Health <= 0 {
			(*monsters)[monsterInd] = (*monsters)[len(*monsters)-1]
			(*monsters) = (*monsters)[:len(*monsters)-1]
		}
	}
}
