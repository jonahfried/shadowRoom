package player

import (
	"fmt"
	"math"
	"shadowRoom/boundry"
	"shadowRoom/creature"

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

	Level   float64
	Spacing int
	Count   int

	Cam Camera

	Monsters []creature.Creature
	Shots    []Shot
	ShotsImg *imdraw.IMDraw

	Fade   pixel.Picture
	Sprite *pixel.Sprite
	Img    *imdraw.IMDraw
}

// MakeAgent creates a new agent starting at a given (x, y) coordinate
func MakeAgent(x, y float64, win *pixelgl.Window, sprite *pixel.Sprite) (cir Agent) {
	cir.Posn = pixel.V(x, y)
	cir.Vel = pixel.ZV
	cir.Acc = pixel.ZV

	cir.Shade = true
	cir.Fill = true

	cir.Cam = MakeCamera(cir.Posn, win)

	cir.Level = 0.02
	cir.Spacing = 6
	cir.Count = 88

	cir.Monsters = make([]creature.Creature, 0)
	cir.Shots = make([]Shot, 0)
	cir.ShotsImg = imdraw.New(nil)
	cir.ShotsImg.Color = pixel.ToRGBA(colornames.Firebrick).Mul(pixel.Alpha(.7))

	cir.Fade = sprite.Picture()
	cir.Sprite = sprite

	cir.Img = imdraw.New(nil)
	cir.Img.Color = colornames.Purple
	cir.Img.Push(pixel.V(500, 350))
	cir.Img.Circle(20, 0)
	return cir
}

// A Shot keeps track of the players attacks
type Shot struct {
	Posn1, Posn2, Vel pixel.Vec
}

func (cir Agent) fire(win *pixelgl.Window) (bullet Shot) {
	mousePosn := cir.Cam.Matrix.Unproject(win.MousePosition())
	directionVec := mousePosn.Sub(cir.Posn)
	directionVec = directionVec.Scaled(1 / magnitude(directionVec))
	bullet.Posn1 = cir.Posn
	bullet.Posn2 = cir.Posn.Add(directionVec.Scaled(10))
	bullet.Vel = directionVec.Scaled(14)
	return bullet
}

// DispShots displays shots
func (cir *Agent) DispShots(win *pixelgl.Canvas) {
	cir.ShotsImg.Clear()
	for _, bullet := range cir.Shots {
		cir.ShotsImg.Push(bullet.Posn1)
		cir.ShotsImg.Push(bullet.Posn2)
		cir.ShotsImg.Line(4)
	}
	cir.ShotsImg.Draw(win)
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

// PressHandler Handles key presses.
// Agent method taking in a window from which to accept inputs.
func (cir *Agent) PressHandler(win *pixelgl.Window) {
	cir.Acc = pixel.ZV
	if win.Pressed(pixelgl.KeyA) {
		cir.Acc = cir.Acc.Sub(pixel.V(5, 0))
	}
	if win.Pressed(pixelgl.KeyD) {
		cir.Acc = cir.Acc.Add(pixel.V(5, 0))
	}
	if win.Pressed(pixelgl.KeyS) {
		cir.Acc = cir.Acc.Sub(pixel.V(0, 5))
	}
	if win.Pressed(pixelgl.KeyW) {
		cir.Acc = cir.Acc.Add(pixel.V(0, 5))
	}
	if win.JustPressed(pixelgl.KeyJ) {
		cir.Posn.X--
	}
	if win.JustPressed(pixelgl.KeyL) {
		cir.Posn.X++
	}
	if win.JustPressed(pixelgl.KeyK) {
		cir.Posn.Y--
	}
	if win.JustPressed(pixelgl.KeyI) {
		cir.Posn.Y++
	}
	if win.JustPressed(pixelgl.KeySpace) {
		cir.Shade = !cir.Shade
	}
	if win.JustPressed(pixelgl.KeyF) {
		cir.Fill = !cir.Fill
	}
	if win.JustPressed(pixelgl.KeyUp) {
		cir.Level += .001
	}
	if win.JustPressed(pixelgl.KeyDown) {
		cir.Level -= .001
	}
	if win.JustPressed(pixelgl.KeyRight) {
		cir.Spacing++
	}
	if win.JustPressed(pixelgl.KeyLeft) {
		cir.Spacing--
	}
	if win.Pressed(pixelgl.KeyComma) {
		cir.Count--
	}
	if win.Pressed(pixelgl.KeyPeriod) {
		cir.Count++
	}

	if win.JustPressed(pixelgl.MouseButton1) {
		cir.Shots = append(cir.Shots, cir.fire(win))
	}
}

// ReleaseHandler Handles key releases.
// Agent method taking in a window from which to accept inputs.
func (cir *Agent) ReleaseHandler(win *pixelgl.Window) {
	if win.JustReleased(pixelgl.KeyA) {
		cir.Acc = cir.Acc.Add(pixel.V(5, 0))
	}
	if win.JustReleased(pixelgl.KeyD) {
		cir.Acc = cir.Acc.Sub(pixel.V(5, 0))
	}
	if win.JustReleased(pixelgl.KeyS) {
		cir.Acc = cir.Acc.Add(pixel.V(0, 5))
	}
	if win.JustReleased(pixelgl.KeyW) {
		cir.Acc = cir.Acc.Sub(pixel.V(0, 5))
	}
}

// Disp Handles display of an agent. Clears, pushes, adds shape, and draws.
func (cir *Agent) Disp(win *pixelgl.Window) {
	cir.Img.Clear()
	cir.Img.Push(cir.Posn)
	cir.Img.Circle(20, 0)
	cir.Img.Draw(win)
}

func magnitude(vec pixel.Vec) float64 {
	return vecDist(pixel.ZV, vec)
}

// Light adds fading light (white circles) around an Agent's posn
// func (cir *Agent) Light(room *boundry.Place) {
// 	room.Target.SetComposeMethod(pixel.ComposePlus)
// 	cir.Sprite.DrawColorMask(room.Target, pixel.IM.Moved(cir.Posn), pixel.Alpha(cir.Level))
// 	// cir.Sprite.Draw(room.Target, pixel.IM.Moved(cir.Posn))
// 	// room.Target.SetComposeMethod(pixel.ComposeIn)
// }

// Light adds fading light (white circles) around an Agent's posn
func (cir *Agent) Light(room *boundry.Place) {
	img := imdraw.New(nil)
	img.Precision = 32
	col := (pixel.ToRGBA(colornames.Whitesmoke)).Mul(pixel.Alpha(cir.Level))
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

}

// // Light adds fading light (white circles) around an Agent's posn
// func (cir *Agent) Light(room *boundry.Place) {
// 	img := imdraw.New(nil)
// 	col := (pixel.ToRGBA(colornames.Whitesmoke)).Mul(pixel.Alpha(.1))
// 	for fade := 80; fade > 1; fade-- {
// 		img.Color = col
// 		img.Push(cir.Posn)
// 		img.Circle(float64(fade*6), 0)
// 		// col = col.Mul(pixel.Alpha(1 / float64(fade)))
// 		col = (pixel.ToRGBA(colornames.Whitesmoke)).Mul(pixel.Alpha(.8 / float64(fade)))
// 	}

// 	// room.Target.Clear(pixel.Alpha(0))
// 	// room.Target.SetComposeMethod()
// 	img.Draw(room.Target)

// }

// Update is an agent method. Runs all necessary per-frame proccedures on agent.
// Takes in a pixelgl.Window from which to accept inputs.
func (cir *Agent) Update(win *pixelgl.Window, room boundry.Place) {
	cir.PressHandler(win)
	cir.ReleaseHandler(win)

	// if vecDist(cir.Vel, pixel.ZV) < .2 {
	// 	cir.Vel = pixel.ZV
	// }

	center := cir.Posn
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
				cir.Vel = (cir.Vel.Add(offset)).Scaled(.88)
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
			cir.Vel = (cir.Vel.Add(offset)).Scaled(.88)
		}

	}
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

	// Update Shots
	for bulletInd := len(cir.Shots) - 1; bulletInd >= 0; bulletInd-- { //range cir.Shots {
		cir.Shots[bulletInd].Posn1 = cir.Shots[bulletInd].Posn1.Add(cir.Shots[bulletInd].Vel)
		cir.Shots[bulletInd].Posn2 = cir.Shots[bulletInd].Posn2.Add(cir.Shots[bulletInd].Vel)

		// midPoint := (cir.Shots[bulletInd].Posn1.Add(cir.Shots[bulletInd].Posn2)).Scaled(1 / 2)
		for monsterInd := range cir.Monsters {
			if vecDist(cir.Shots[bulletInd].Posn2, cir.Monsters[monsterInd].Posn) < 20 {
				fmt.Println("shot")
				cir.Monsters[monsterInd].Health--
				// BAD
				// if len(cir.Shots) > 0 {
				// 	cir.Shots

				// }
			}
		}
	}
	filter(&cir.Monsters)
}

// // Update Shots
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

func filter(monsters *[]creature.Creature) {
	for monsterInd := 0; monsterInd < len(*monsters); monsterInd++ {
		if (*monsters)[monsterInd].Health <= 0 {
			(*monsters)[monsterInd] = (*monsters)[len(*monsters)-1]
			(*monsters) = (*monsters)[:len(*monsters)-1]
		}
	}
}

func vecDist(v1, v2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
}

// Camera is a struct storing a Matrix to set the window
type Camera struct {
	Posn, vel          pixel.Vec
	maxForce, maxSpeed float64
	Matrix             pixel.Matrix
}

// MakeCamera takes in a starting position and window and returns a Camera
func MakeCamera(posn pixel.Vec, win *pixelgl.Window) (cam Camera) {
	cam.Posn = posn
	cam.vel = pixel.ZV
	cam.maxSpeed = 10

	cam.Matrix = pixel.IM.Moved(win.Bounds().Center().Sub(posn))

	return cam
}

// Attract updates a Camera's velocity and position to follow a given point
func (cam *Camera) Attract(target pixel.Vec) {
	acc := limitVecMag(target.Sub(cam.Posn), vecDist(cam.Posn, target)/10)
	// scale := cam.maxForce / 8
	// acc = acc.Scaled(scale)
	acc = acc.Sub(cam.vel)

	cam.vel = cam.vel.Add(acc)
	cam.vel = limitVecMag(cam.vel, cam.maxSpeed)

	cam.Posn = cam.Posn.Add(cam.vel)

}
