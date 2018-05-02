package player

import (
	"math"
	"shadowRoom/boundry"

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

	Img *imdraw.IMDraw
}

// MakeAgent creates a new agent starting at a given (x, y) coordinate
func MakeAgent(x, y float64) (cir Agent) {
	cir.Posn = pixel.V(x, y)
	cir.Vel = pixel.ZV
	cir.Acc = pixel.ZV

	cir.Shade = true
	cir.Fill = true

	cir.Img = imdraw.New(nil)
	cir.Img.Color = colornames.Purple
	cir.Img.Push(pixel.V(500, 350))
	cir.Img.Circle(20, 0)
	return cir
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

// Light adds fading light (white circles) around an Agent's posn
func (cir *Agent) Light(room *boundry.Place) {
	img := imdraw.New(nil)
	col := (pixel.ToRGBA(colornames.Whitesmoke)).Mul(pixel.Alpha(.1))
	for fade := 80; fade > 1; fade-- {
		img.Color = col
		img.Push(cir.Posn)
		img.Circle(float64(fade*6), 0)
		// col = col.Mul(pixel.Alpha(1 / float64(fade)))
		col = (pixel.ToRGBA(colornames.Whitesmoke)).Mul(pixel.Alpha(1 / float64(fade)))
	}

	// room.Target.Clear(pixel.Alpha(0))
	// room.Target.SetComposeMethod()
	img.Draw(room.Target)
}

// Update is an agent method. Runs all necessary per-frame proccedures on agent.
// Takes in a pixelgl.Window from which to accept inputs.
func (cir *Agent) Update(win *pixelgl.Window, room boundry.Place) {
	cir.PressHandler(win)
	cir.ReleaseHandler(win)

	newPosn := cir.Posn.Add(cir.Vel)
	// for _, obst := range room.Blocks {
	// 	angle := math.Atan2(cir.Posn.Y-obst.Center.Y, cir.Posn.X-obst.Center.X)
	// 	lerpAmount := vecDist(newPosn, obst.Center) / 20
	// 	lerp := pixel.Lerp(obst.Center, newPosn, lerpAmount)
	// 	landed := boundry.Obstruct(lerp, angle, room, obst)
	// 	if !(vecDist(obst.Center, landed) > obst.Radius) {
	// 		cir.Vel = limitVecMag(cir.Vel.Add(limitVecMag(cir.Acc, 1.5)), 10).Scaled(.88)
	// 		if vecDist(cir.Vel, pixel.ZV) < .2 {
	// 			cir.Vel = pixel.ZV
	// 		}
	// 		return
	// 	}
	// }

	cir.Posn = newPosn
	cir.Vel = limitVecMag(cir.Vel.Add(limitVecMag(cir.Acc, 1.5)), 10).Scaled(.88)
	if vecDist(cir.Vel, pixel.ZV) < .2 {
		cir.Vel = pixel.ZV
	}

	cir.Posn.X = math.Max(cir.Posn.X-20, room.Rect.Min.X) + 20
	cir.Posn.X = math.Min(cir.Posn.X+20, room.Rect.Max.X) - 20

	cir.Posn.Y = math.Max(cir.Posn.Y-20, room.Rect.Min.Y) + 20
	cir.Posn.Y = math.Min(cir.Posn.Y+20, room.Rect.Max.Y) - 20
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
