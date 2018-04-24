package main

import (
	"fmt"
	"math"
	"math/rand"
	"shadowRoom/boundry"
	"shadowRoom/creature"
	"sort"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// Type agent. Keeps a posn, vel, acc, and a *pixel.IMDraw.
// Used to keep all information for the movable image together.
type agent struct {
	posn pixel.Vec
	vel  pixel.Vec
	acc  pixel.Vec

	shade bool
	fill  bool

	img *imdraw.IMDraw
}

// Creates a new agent starting at a given (x, y) coordinate
func makeAgent(x, y float64) (cir agent) {
	cir.posn = pixel.V(x, y)
	cir.vel = pixel.ZV
	cir.acc = pixel.ZV

	cir.shade = true
	cir.fill = true

	cir.img = imdraw.New(nil)
	cir.img.Color = colornames.Purple
	cir.img.Push(pixel.V(500, 350))
	cir.img.Circle(20, 0)
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

// Handles key presses. Agent method taking in a window from which to accept inputs.
func (cir *agent) pressHandler(win *pixelgl.Window) {
	cir.acc = pixel.ZV
	if win.Pressed(pixelgl.KeyA) {
		cir.acc = cir.acc.Sub(pixel.V(5, 0))
	}
	if win.Pressed(pixelgl.KeyD) {
		cir.acc = cir.acc.Add(pixel.V(5, 0))
	}
	if win.Pressed(pixelgl.KeyS) {
		cir.acc = cir.acc.Sub(pixel.V(0, 5))
	}
	if win.Pressed(pixelgl.KeyW) {
		cir.acc = cir.acc.Add(pixel.V(0, 5))
	}
	if win.JustPressed(pixelgl.KeyJ) {
		cir.posn.X--
	}
	if win.JustPressed(pixelgl.KeyL) {
		cir.posn.X++
	}
	if win.JustPressed(pixelgl.KeyK) {
		cir.posn.Y--
	}
	if win.JustPressed(pixelgl.KeyI) {
		cir.posn.Y++
	}
	if win.JustPressed(pixelgl.KeySpace) {
		cir.shade = !cir.shade
	}
	if win.JustPressed(pixelgl.KeyF) {
		cir.fill = !cir.fill
	}
}

// Handles key releases. Agent method taking in a window from which to accept inputs.
func (cir *agent) releaseHandler(win *pixelgl.Window) {
	if win.JustReleased(pixelgl.KeyA) {
		cir.acc = cir.acc.Add(pixel.V(5, 0))
	}
	if win.JustReleased(pixelgl.KeyD) {
		cir.acc = cir.acc.Sub(pixel.V(5, 0))
	}
	if win.JustReleased(pixelgl.KeyS) {
		cir.acc = cir.acc.Add(pixel.V(0, 5))
	}
	if win.JustReleased(pixelgl.KeyW) {
		cir.acc = cir.acc.Sub(pixel.V(0, 5))
	}
}

// Handles display of an agent. Clears, pushes, adds shape, and draws.
func (cir *agent) disp(win *pixelgl.Window) {
	cir.img.Clear()
	cir.img.Push(cir.posn)
	cir.img.Circle(20, 0)
	cir.img.Draw(win)
}

// agent method. Runs all necessary per-frame proccedures on agent.
// Takes in a pixelgl.Window from which to accept inputs.
func (cir *agent) update(win *pixelgl.Window, room boundry.Place) {
	cir.pressHandler(win)
	cir.releaseHandler(win)

	cir.posn = cir.posn.Add(cir.vel)
	cir.vel = limitVecMag(cir.vel.Add(limitVecMag(cir.acc, 1.5)), 10).Scaled(.88)
	if vecDist(cir.vel, pixel.ZV) < .2 {
		cir.vel = pixel.ZV
	}

	// if !r.rect.Contains(cir.posn) {
	// if cir.posn.X < r.rect.Center().X-100 { //-(r.rect.W()/2) {
	// 	cir.posn.X = r.rect.Center().X - 100
	// }
	// cir.posn.X = math.Max(cir.posn.X, r.rect.Center().X-r.rect.W())
	cir.posn.X = math.Max(cir.posn.X-20, room.Rect.Min.X) + 20
	cir.posn.X = math.Min(cir.posn.X+20, room.Rect.Max.X) - 20

	cir.posn.Y = math.Max(cir.posn.Y-20, room.Rect.Min.Y) + 20
	cir.posn.Y = math.Min(cir.posn.Y+20, room.Rect.Max.Y) - 20
	// }
}

func vecDist(v1, v2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
}

func obstruct(posn pixel.Vec, angle float64, room boundry.Place, block boundry.Obsticle) (stopPoint pixel.Vec) {
	// TODO: Fix divide by zero error
	if (math.Cos(angle)) == 0 {
		panic("divide by zero")
	}
	slope := math.Sin(angle) / math.Cos(angle)
	yInt := posn.Y - (posn.X * slope)

	extension := 0.000000001

	var blocks [2]boundry.Obsticle

	blocks[0] = block
	blocks[1] = boundry.MakeObsticle(room.Vertices, true)

	stopPoint = pixel.V(math.MaxFloat64, (math.MaxFloat64*slope)+yInt)

	for _, block := range blocks {
		for ind := 0; ind < len(block.Vertices)-1; ind++ {
			deltaY := (block.Vertices[ind].X - block.Vertices[ind+1].X)
			if deltaY != 0 {
				edgeSlope := (block.Vertices[ind].Y - block.Vertices[ind+1].Y) / deltaY
				edgeYInt := block.Vertices[ind].Y - (block.Vertices[ind].X * edgeSlope)
				xInterception := (edgeYInt - yInt) / (slope - edgeSlope)
				interception := pixel.V(xInterception, (slope*xInterception)+yInt)
				if vecDist(interception, posn) < vecDist(stopPoint, posn) && xInterception <= math.Max(block.Vertices[ind].X, block.Vertices[ind+1].X)+extension && xInterception >= math.Min(block.Vertices[ind].X, block.Vertices[ind+1].X)-extension && (xInterception-posn.X)*math.Cos(angle) > 0 {
					stopPoint = interception
				}
			} else {
				yInterception := (slope * block.Vertices[ind].X) + yInt
				interception := pixel.V(block.Vertices[ind].X, yInterception)

				if vecDist(interception, posn) < vecDist(stopPoint, posn) && yInterception <= math.Max(block.Vertices[ind].Y, block.Vertices[ind+1].Y)+extension && yInterception >= math.Min(block.Vertices[ind].Y, block.Vertices[ind+1].Y)-extension && (interception.X-posn.X)*math.Cos(angle) > 0 {
					stopPoint = interception
				}

			}

		}
		deltaY := (block.Vertices[len(block.Vertices)-1].X - block.Vertices[0].X)
		if deltaY != 0 {
			edgeSlope := (block.Vertices[len(block.Vertices)-1].Y - block.Vertices[0].Y) / deltaY
			edgeYInt := block.Vertices[len(block.Vertices)-1].Y - (block.Vertices[len(block.Vertices)-1].X * edgeSlope)
			xInterception := (edgeYInt - yInt) / (slope - edgeSlope)
			interception := pixel.V(xInterception, (slope*xInterception)+yInt)
			if vecDist(interception, posn) < vecDist(stopPoint, posn) && xInterception <= math.Max(block.Vertices[len(block.Vertices)-1].X, block.Vertices[0].X) && xInterception >= math.Min(block.Vertices[len(block.Vertices)-1].X, block.Vertices[0].X) && (xInterception-posn.X)*math.Cos(angle) > 0 {
				stopPoint = interception
			}
		} else {
			yInterception := (slope * block.Vertices[len(block.Vertices)-1].X) + yInt
			interception := pixel.V(block.Vertices[len(block.Vertices)-1].X, yInterception)

			if vecDist(interception, posn) < vecDist(stopPoint, posn) && yInterception <= math.Max(block.Vertices[len(block.Vertices)-1].Y, block.Vertices[0].Y) && yInterception >= math.Min(block.Vertices[len(block.Vertices)-1].Y, block.Vertices[0].Y) && (interception.X-posn.X)*math.Cos(angle) > 0 {
				stopPoint = interception
			}

		}

	}

	return stopPoint

}

type camera struct {
	posn, vel          pixel.Vec
	maxForce, maxSpeed float64
	matrix             pixel.Matrix
}

func makeCamera(posn pixel.Vec, win *pixelgl.Window) (cam camera) {
	cam.posn = posn
	cam.vel = pixel.ZV
	cam.maxSpeed = 10

	cam.matrix = pixel.IM.Moved(win.Bounds().Center().Sub(posn))

	return cam
}

func (cam *camera) attract(target pixel.Vec) {
	acc := limitVecMag(target.Sub(cam.posn), vecDist(cam.posn, target)/10)
	// scale := cam.maxForce / 8
	// acc = acc.Scaled(scale)
	acc = acc.Sub(cam.vel)

	cam.vel = cam.vel.Add(acc)
	cam.vel = limitVecMag(cam.vel, cam.maxSpeed)

	cam.posn = cam.posn.Add(cam.vel)

}

// Acting main function
func run() {
	// setting up the window
	cfg := pixelgl.WindowConfig{
		Title:  "shadowRoom",
		Bounds: pixel.R(0, 0, 1350, 725),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// rooms := make(map[pixel.Vec]struct{ bottemLeft, topRight float64 })

	// setting up the Player
	var cir = makeAgent(0, 0)

	// setting up the camera
	cam := makeCamera(cir.posn, win)

	blockList := make([]boundry.Obsticle, 0)
	for ind := 0; ind < 1; ind++ {
		// TODO: randomize more of the blocks
		// -	Take the bounds and set variables for radius/stdDev
		// - 	Then make sure the vertex won't land out of bounds
		// block := makeRandomBlock(50, 10, 6, pixel.V(400, 200))
		// blc := makeRandomBlock(50, 10, 6, pixel.V((rand.Float64()*400), (rand.Float64()*250)))
		blc := boundry.MakeRandomBlock(50, 10, 6, pixel.V((rand.Float64()*800)-400, (rand.Float64()*500)-250))
		blockList = append(blockList, blc)
	}

	// TODO: Make room bounds relative to Bounds
	room := boundry.MakePlace(pixel.R(-500, -350, 500, 350), blockList)

	point := imdraw.New(nil)
	point.Color = colornames.Black

	blob := creature.MakeCreature(0, 100)

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)

		cir.update(win, room)

		cam.attract(cir.posn)
		cam.matrix = pixel.IM.Moved(win.Bounds().Center().Sub(cam.posn))

		win.SetMatrix(cam.matrix)

		room.Img.Draw(win)

		blob.Update(room)
		blob.Disp(win)

		point.Clear()
		anglesToCheck := make([]float64, 0, 10)
		for _, block := range room.Blocks {
			if !block.IsRoom {
				for _, vertex := range block.Vertices {
					theta := math.Atan2((vertex.Y - cir.posn.Y), (vertex.X - cir.posn.X))
					anglesToCheck = append(anglesToCheck, theta)
				}
				length := len(anglesToCheck)
				for k := 0; k < length; k++ {
					for offset := -.000001; offset <= .000001; offset += .000001 {
						anglesToCheck = append(anglesToCheck, anglesToCheck[k]+offset)
					}
				}
				shadedVertices := make([]pixel.Vec, 0)
				for _, vertex := range room.Vertices {
					theta := math.Atan2((vertex.Y - cir.posn.Y), (vertex.X - cir.posn.X))
					landed := obstruct(cir.posn, theta, room, block)
					if math.Abs(vecDist(landed, cir.posn)-vecDist(vertex, cir.posn)) > 1 {
						// point.Push(vertex)
						shadedVertices = append(shadedVertices, vertex)
						// if !cir.shade {
						// 	point.Circle(4, 0)
						// }
					}
				}
				shapePoints := make([]pixel.Vec, 0)
				sort.Float64s(anglesToCheck)
				for _, angle := range anglesToCheck {
					vec := obstruct(cir.posn, angle, room, block)
					shapePoints = append(shapePoints, vec)
					point.Push(vec)
					if !cir.shade {
						point.Circle(4, 0)
					}
				}
				if cir.shade {
					if cir.fill {
						point.Polygon(0)
					} else {
						point.Polygon(1)
					}
				}
				for _, vertex := range shadedVertices {
					for vecInd := 1; vecInd < len(shapePoints); vecInd++ {
						point.Push(vertex)
						if !cir.shade {
							point.Circle(4, 0)
						}
						point.Push(shapePoints[vecInd])
						if !cir.shade {
							point.Circle(4, 0)
						}
						point.Push(shapePoints[vecInd-1])
						if !cir.shade {
							point.Circle(4, 0)
						}
						if cir.shade {
							if cir.fill {
								point.Polygon(0)
							} else {
								point.Polygon(1)
							}
						}
					}
				}

				if len(shadedVertices) > 1 {
					for vertexInd := 1; vertexInd < len(shadedVertices); vertexInd++ {
						point.Push(shadedVertices[vertexInd-1])
						if !cir.shade {
							point.Circle(4, 0)
						}
						point.Push(shadedVertices[vertexInd])
						if !cir.shade {
							point.Circle(4, 0)
						}
						point.Push(shapePoints[0])
						if !cir.shade {
							point.Circle(4, 0)
						}
						if cir.shade {
							if cir.fill {
								point.Polygon(0)
							} else {
								point.Polygon(1)
							}
						}

					}
				}

				anglesToCheck = make([]float64, 0, 10)
				shapePoints = make([]pixel.Vec, 0)
				point.Draw(win)
			}
		}

		for _, bc := range room.Blocks {
			if !bc.IsRoom {
				bc.Img.Draw(win)
			}
		}

		cir.disp(win)
		win.Update()
		// time.Sleep(1 / 2 * time.Second)
	}
	fmt.Print("\nWindow Closed: App Shutting Down \n")
	fmt.Printf("Playerdata at close: %v \n \n", cir)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pixelgl.Run(run)
}
