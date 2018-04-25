package main

import (
	"fmt"
	"math"
	"math/rand"
	"shadowRoom/boundry"
	"shadowRoom/creature"
	"shadowRoom/player"
	"sort"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

func vecDist(v1, v2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
}

func illuminate(room boundry.Place, cir player.Agent, point *imdraw.IMDraw, win *pixelgl.Window) {
	point.Clear()
	anglesToCheck := make([]float64, 0, 10)
	for _, block := range room.Blocks {
		for _, vertex := range block.Vertices {
			theta := math.Atan2((vertex.Y - cir.Posn.Y), (vertex.X - cir.Posn.X))
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
			theta := math.Atan2((vertex.Y - cir.Posn.Y), (vertex.X - cir.Posn.X))
			landed := boundry.Obstruct(cir.Posn, theta, room, block)
			if math.Abs(vecDist(landed, cir.Posn)-vecDist(vertex, cir.Posn)) > 1 {
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
			vec := boundry.Obstruct(cir.Posn, angle, room, block)
			shapePoints = append(shapePoints, vec)
			point.Push(vec)
			if !cir.Shade {
				point.Circle(4, 0)
			}
		}
		if cir.Shade {
			if cir.Fill {
				point.Polygon(0)
			} else {
				point.Polygon(1)
			}
		}
		for _, vertex := range shadedVertices {
			for vecInd := 1; vecInd < len(shapePoints); vecInd++ {
				point.Push(vertex)
				if !cir.Shade {
					point.Circle(4, 0)
				}
				point.Push(shapePoints[vecInd])
				if !cir.Shade {
					point.Circle(4, 0)
				}
				point.Push(shapePoints[vecInd-1])
				if !cir.Shade {
					point.Circle(4, 0)
				}
				if cir.Shade {
					if cir.Fill {
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
				if !cir.Shade {
					point.Circle(4, 0)
				}
				point.Push(shadedVertices[vertexInd])
				if !cir.Shade {
					point.Circle(4, 0)
				}
				point.Push(shapePoints[0])
				if !cir.Shade {
					point.Circle(4, 0)
				}
				if cir.Shade {
					if cir.Fill {
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
	var cir = player.MakeAgent(0, 0)

	// setting up the camera
	cam := player.MakeCamera(cir.Posn, win)

	// TODO: Make room bounds relative to Bounds
	room := boundry.MakePlace(pixel.R(-500, -350, 500, 350), 10)

	point := imdraw.New(nil)
	point.Color = colornames.Black

	blob := creature.MakeCreature(0, 100)

	// Main Draw Loop:
	for !win.Closed() {
		win.Clear(colornames.Black)

		cir.Update(win, room)

		cam.Attract(cir.Posn)
		cam.Matrix = pixel.IM.Moved(win.Bounds().Center().Sub(cam.Posn))

		win.SetMatrix(cam.Matrix)

		room.Disp(cir.Posn, win)

		blob.Update(room)
		blob.Disp(win)

		for _, bc := range room.Blocks {
			bc.Img.Draw(win)
		}

		illuminate(room, cir, point, win)

		cir.Disp(win)
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
