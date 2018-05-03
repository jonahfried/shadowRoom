package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"shadowRoom/boundry"
	"shadowRoom/creature"
	"shadowRoom/player"
	"sort"
	"time"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
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

func timer(monsters *[]creature.Creature) {
	tick := time.Tick(time.Second * 5)
	for {
		*monsters = append(*monsters, creature.MakeCreature(0, 0))
		<-tick
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

// Acting main function
func run() {
	// setting up the window
	cfg := pixelgl.WindowConfig{
		Title:  "shadowRoom",
		Bounds: pixel.R(0, 0, 1350, 725),
		VSync:  true,
		// Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	lightingPic, err := loadPicture("drawing.png")
	if err != nil {
		panic(err)
	}

	lightingSprite := pixel.NewSprite(lightingPic, pixel.R(0, 0, 480, 480))

	// rooms := make(map[pixel.Vec]struct{ bottemLeft, topRight float64 })

	// setting up the Player
	var cir = player.MakeAgent(0, 0, win, lightingSprite)

	// setting up the camera
	// cam := player.MakeCamera(cir.Posn, win)

	// TODO: Make room bounds relative to Bounds
	room := boundry.MakePlace(pixel.R(-700, -600, 700, 600), 11)

	point := imdraw.New(nil)
	point.Color = colornames.Black

	cir.Monsters = append(cir.Monsters, creature.MakeCreature(0, 0))

	go timer(&cir.Monsters)

	last := time.Now()

	frames := 0.0
	seconds := 0.0

	// Main Draw Loop:
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		if frames > 20 {
			frames = 0
			seconds = 0
		}
		seconds += dt
		frames++
		last = time.Now()

		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(cir.Posn.Add(pixel.V(-20, 30)), basicAtlas)
		basicTxt.Color = colornames.Red

		fmt.Fprintln(basicTxt, "fps:", math.Round(frames/seconds))

		fmt.Println(1 / dt)

		win.Clear(colornames.Black)
		room.Disp()

		cir.Update(win, room)

		cir.Cam.Attract(cir.Posn)
		cir.Cam.Matrix = pixel.IM.Moved(win.Bounds().Center().Sub(cir.Cam.Posn))

		win.SetMatrix(cir.Cam.Matrix)
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
		// room.Disp()
		cir.Light(&room)

		// if rand.Float64() > .992 {
		// 	cir.Monsters = append(cir.Monsters, creature.MakeCreature(0, 0))
		// }

		for monsterInd := range cir.Monsters {
			grid := room.ToGrid(40, cir.Monsters[monsterInd].Posn, cir.Posn)
			monsterTarget := grid.AStar()
			cir.Monsters[monsterInd].Update(room, monsterTarget)
			room.Target.SetComposeMethod(pixel.ComposeIn)
			cir.Monsters[monsterInd].Disp(room.Target)
			room.Target.SetComposeMethod(pixel.ComposeOver)
		}

		room.Disp()
		room.Target.Draw(win, pixel.IM) //.Moved(win.Bounds().Center()))

		cir.DispShots(win)

		illuminate(room, cir, point, win)

		cir.Disp(win)
		basicTxt.Draw(win, pixel.IM)
		lightingSprite.Draw(win, pixel.IM)
		win.Update()
		// time.Sleep(1 / 2 * time.Second)
	}
	fmt.Print("\nWindow Closed: App Shutting Down \n")
	// fmt.Printf("Playerdata at close: %v \n \n", cir)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pixelgl.Run(run)
}
