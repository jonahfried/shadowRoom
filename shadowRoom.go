package main

import (
	"fmt"
	_ "image/png"
	"math"
	"math/rand"
	"sort"
	"time"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

func illuminate(room Place, cir Agent, point *imdraw.IMDraw, win *pixelgl.Window) {
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
			landed := Obstruct(cir.Posn, theta, room, block)
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
			vec := Obstruct(cir.Posn, angle, room, block)
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
func run(window ...*pixelgl.Window) {
	var win *pixelgl.Window
	if len(window) == 0 {
		// setting up the window
		cfg := pixelgl.WindowConfig{
			Title:  "shadowRoom",
			Bounds: pixel.R(0, 0, 1350, 725),
			VSync:  true,
			// Resizable: true,
		}
		var err error
		win, err = pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}
	} else {
		win = window[0]
	}

	// rooms := make(map[pixel.Vec]struct{ bottemLeft, topRight float64 })

	// setting up the Player
	var cir = MakeAgent(0, 0, win)

	// setting up the camera
	// cam := player.MakeCamera(cir.Posn, win)

	// TODO: Make room bounds relative to Bounds
	room := MakePlace(pixel.R(-700, -600, 700, 600), 11)
	room.ToGrid(40)

	point := imdraw.New(nil)
	point.Color = colornames.Black

	frameRate := time.Tick(time.Millisecond * 17)
	fiveSec := time.Tick(time.Second * 5)
	thirtySec := time.Tick(time.Second * 30)

	last := time.Now()

	frames := 0.0
	seconds := 0.0

	// Main Draw Loop:
	for !win.Closed() && cir.Health > 0 {
		dt := time.Since(last).Seconds()
		if frames > 200 {
			frames = 0
			seconds = 0
		}

		select {
		case <-frameRate:

		case <-fiveSec:
			cir.Monsters = append(cir.Monsters, MakeCreature(&room, &cir))
		case <-thirtySec:
			room.presentBoost()

		}

		seconds += dt
		frames++
		last = time.Now()

		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(cir.Posn.Add(pixel.V(-20, 30)), basicAtlas)
		basicTxt.Color = colornames.Red

		fmt.Fprintln(basicTxt, "fps:", math.Round(frames/seconds))

		// fmt.Println(1 / dt)

		win.Clear(colornames.Black)
		room.Disp()

		cir.Update(win, &room)

		cir.Cam.Attract(cir.Posn)
		cir.Cam.Matrix = pixel.IM.Moved(win.Bounds().Center().Sub(cir.Cam.Posn))

		win.SetMatrix(cir.Cam.Matrix)

		//Move this out of loop?
		room.Target.SetMatrix(pixel.IM.Moved(room.Target.Bounds().Center()))

		// TESTING PRIORITY QUEUE
		/*
			prior := make(priorityqueue.PriorityQueue, 0)
			heap.Init(&prior)
			for testCount := 0; testCount < 10; testCount++ {
				randomInt := rand.Intn(20)
				item := &priorityqueue.Elem{
					Value:    randomInt,
					Priority: float64(randomInt),
				}
				heap.Push(&prior, item)
			}
			for prior.Len() > 1 {
				fmt.Println(heap.Pop(&prior))
			}

			fmt.Print("\n")
		*/

		room.Target.Clear(pixel.Alpha(0))
		room.Disp()
		cir.Light(&room)

		room.Target.SetComposeMethod(pixel.ComposeIn)
		if room.Booster.Present {
			room.Booster.Img.Clear()
			room.Booster.Img.Push(room.Booster.Posn)
			room.Booster.Img.Circle(10, 0)
			room.Booster.Img.Draw(room.Target)
		}

		for monsterInd := range cir.Monsters {
			monsterTarget := AStar(room.GridRepresentation, cir.Monsters[monsterInd].Posn, cir.Posn, &cir)
			cir.Monsters[monsterInd].Update(room, &cir, monsterTarget, cir.Monsters)
			cir.Monsters[monsterInd].Disp(room.Target)
		}

		cir.DispShots(room.Target)

		room.Target.SetComposeMethod(pixel.ComposeOver)
		room.Disp()

		room.Target.Draw(win, pixel.IM) //.Moved(win.Bounds().Center()))
		basicTxt.Draw(win, pixel.IM)
		cir.Disp(win)
		illuminate(room, cir, point, win)

		win.Update()
		// time.Sleep(1 / 2 * time.Second)
	}
	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)
		endscreen := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		endTxt := text.New(cir.Posn.Add(pixel.V(-20, 30)), endscreen)
		endTxt.Color = colornames.Black
		fmt.Fprint(endTxt, "Game Over!\n(space to restart)")
		endTxt.Draw(win, pixel.IM)
		if win.JustPressed(pixelgl.KeySpace) {
			run(win)
		}
		win.Update()
	}
}

func starter() {
	run()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pixelgl.Run(starter)
}
