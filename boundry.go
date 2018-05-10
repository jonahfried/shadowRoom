package main

import (
	"container/heap"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Place is datastructure that holds a rectangle for its bounds
// as well as a slice of Blocks that it contains
// and a slice of its vertices
type Place struct {
	Rect               pixel.Rect
	Blocks             []Obsticle
	Vertices           []pixel.Vec
	GridRepresentation Grid
	Booster            Boost

	Img    *imdraw.IMDraw
	Target *pixelgl.Canvas
}

// Boost to give bonus to the player
type Boost struct {
	Posn    pixel.Vec
	Present bool

	Img *imdraw.IMDraw
}

// MakePlace turns specifications of a place, (rectangle and the Obsticles it contains)
// and returns a Place
func MakePlace(rect pixel.Rect, numBlocks int, blocks ...Obsticle) (room Place) {
	room.Rect = rect

	borderVertices := make([]pixel.Vec, 0, 4)
	borderVertices = append(borderVertices, rect.Min)
	borderVertices = append(borderVertices, pixel.V(rect.Min.X, rect.Max.Y))
	borderVertices = append(borderVertices, rect.Max)
	borderVertices = append(borderVertices, pixel.V(rect.Max.X, rect.Min.Y))

	room.Vertices = borderVertices

	if len(blocks) == 0 {
		blockList := make([]Obsticle, 0)
		for ind := 0; ind < numBlocks; ind++ {
			// TODO: randomize more of the blocks
			// -	Take the bounds and set variables for Radius/stdDev
			// - 	Then make sure the vertex won't land out of bounds
			// block := makeRandomBlock(50, 10, 6, pixel.V(400, 200))
			// blc := makeRandomBlock(50, 10, 6, pixel.V((rand.Float64()*400), (rand.Float64()*250)))
			blc := MakeRandomBlock(50, 10, 6, rect, blockList)
			blockList = append(blockList, blc)
		}
		room.Blocks = blockList
	} else {
		room.Blocks = blocks
	}

	room.Booster.Present = false
	room.Booster.Posn = pixel.V(100, 100)
	room.Booster.Img = imdraw.New(nil)
	room.Booster.Img.Color = colornames.Royalblue
	// room.Booster.Img.Precision = 32

	room.Target = pixelgl.NewCanvas(rect)
	room.Img = imdraw.New(nil)
	return room
}

func (room *Place) presentBoost() {
	if !room.Booster.Present {
		room.Booster.Posn = pixel.V(room.Rect.Center().X+(room.Rect.W()/2*(rand.Float64()-rand.Float64())/2), room.Rect.Center().Y+room.Rect.H()*(rand.Float64()-rand.Float64())/2)
		room.Booster.Present = true
	} else {
		room.Booster.Present = false
	}
}

// Disp updates and displays a room's Img
func (room *Place) Disp() {
	// room.Target.Clear(pixel.Alpha(0))
	// room.Target.Clear(colornames.Green)
	room.Target.SetComposeMethod(pixel.ComposeIn)
	room.Img.Clear()
	room.Img.Color = colornames.Burlywood
	room.Img.Push(room.Rect.Center().Sub(room.Rect.Size().Scaled(.5)))
	room.Img.Push(room.Rect.Center().Add(room.Rect.Size().Scaled(.5)))
	room.Img.Rectangle(3)
	room.Img.Draw(room.Target)

	// bcImg := imdraw.New(nil)
	// bcImg.Color = colornames.Black
	for _, bc := range room.Blocks {
		bc.Img.Draw(room.Target)
		// bcImg.Push(bc.Center)
		// bcImg.Circle(bc.Radius, 0)
	}

	room.Target.SetComposeMethod(pixel.ComposeOver)
	// bcImg.Draw(room.Target)
}

// Obsticle is a data structure defining a boundry inside a room
// it contains:
//	 a list of vertices defining its bounds
//	 and an *imdraw.IMdraw to store its look.
//	 also stores a Center and Radius set if generated randomly to avoid overlap
type Obsticle struct {
	Vertices []pixel.Vec
	Center   pixel.Vec
	Radius   float64

	Img *imdraw.IMDraw
}

// MakeObsticle takes in a list of vertices that define its boundries,
// and a bool representing whether or not it is a Room.
// returns an Obsticle
func MakeObsticle(vertices []pixel.Vec, center pixel.Vec, radius float64) (obst Obsticle) {
	obst.Vertices = vertices
	obst.Center = center
	obst.Radius = radius

	obst.Img = imdraw.New(nil)
	obst.Img.EndShape = 1

	obst.Img.Color = colornames.Burlywood
	for _, vert := range obst.Vertices {
		obst.Img.Push(vert)
	}
	obst.Img.Polygon(8)

	return obst
}

// MakeRandomBlock takes in a standard radius,
// a standard deviation from the radius, a number of vertices to create, (float64)
// and a Room to hold it, returning a randomly generated Obsticle
func MakeRandomBlock(radius, stdDev, vertScale float64, rect pixel.Rect, existingBlocks []Obsticle) Obsticle {
	centerX := rect.Center().X + (rand.Float64() * rect.W()) - (rect.W() / 2)
	centerY := rect.Center().Y + (rand.Float64() * rect.H()) - (rect.H() / 2)
	center := pixel.V(centerX, centerY)
TryLoop:
	for tries := 0; tries < 10; tries++ {
		for _, blc := range existingBlocks {
			if vecDist(blc.Center, center) < (blc.Radius + radius + stdDev) {
				centerX = rect.Center().X + (rand.Float64() * rect.W()) - (rect.W() / 2)
				centerY = rect.Center().Y + (rand.Float64() * rect.H()) - (rect.H() / 2)
				center = pixel.V(centerX, centerY)
				continue TryLoop
			}
		}
		break TryLoop
	}
	vertices := make([]pixel.Vec, 0, 3)
	// <= ?
	for angle := 2 * math.Pi / vertScale; angle < math.Pi*2; angle += 2 * math.Pi / vertScale {
		r := rand.NormFloat64()*stdDev + radius
		vertex := pixel.V(r*math.Cos(angle), (r * math.Sin(angle)))
		vertex = vertex.Add(center)

		vertex.X = math.Min(vertex.X, rect.Max.X)
		vertex.X = math.Max(vertex.X, rect.Min.X)

		vertex.Y = math.Min(vertex.Y, rect.Max.Y)
		vertex.Y = math.Max(vertex.Y, rect.Min.Y)

		vertices = append(vertices, vertex)
	}
	return MakeObsticle(vertices, center, radius+stdDev)
}

func vecDist(v1, v2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
}

// Obstruct takes in a "viewers" position, and an angle of viewing.
// it also takes a room for its boundries
// and an obsticle of interest to determine intersection
func Obstruct(posn pixel.Vec, angle float64, room Place, block Obsticle) (stopPoint pixel.Vec) {
	// TODO: Fix divide by zero error
	if (math.Cos(angle)) == 0 {
		panic("divide by zero")
	}
	slope := math.Sin(angle) / math.Cos(angle)
	yInt := posn.Y - (posn.X * slope)

	extension := 0.000000001

	var blocks [2]Obsticle

	blocks[0] = block
	blocks[1] = MakeObsticle(room.Vertices, pixel.ZV, 0)

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

// A Tile stores a rect to represent an area in a room
// and a bool to represent wheter it is enterable.
type Tile struct {
	Rect       pixel.Rect
	gScore     float64
	fScore     float64
	Enterable  bool
	xInd, yInd int

	foundFrom int
	reviewed  bool
}

func makeTile(rect pixel.Rect, xInd, yInd int) (tl Tile) {
	tl.Rect = rect
	tl.Enterable = true
	tl.gScore = math.MaxFloat64
	tl.fScore = math.MaxFloat64
	tl.xInd, tl.yInd = xInd, yInd
	tl.foundFrom = -1
	tl.reviewed = false
	return tl
}

// Grid is a structure that stores a rough tile representation of a room.
type Grid struct {
	GridMap     []Tile
	Room        *Place
	TilesPerRow int
	TileSize    pixel.Vec
}

func makeGrid(room *Place) (grid Grid) {
	grid.GridMap = make([]Tile, 0)
	grid.Room = room
	return grid
}

// ToGrid is a method for a Place, returning a rough tile-based
// representation of where in the room is enterable.
func (room *Place) ToGrid(grain int) {
	var grid Grid
	tilesPerRow := float64(grain)
	grid = makeGrid(room)
	grid.TilesPerRow = grain
	grid.TileSize = pixel.V(room.Rect.W()/tilesPerRow, room.Rect.H()/tilesPerRow)

	xInd := 0
	yInd := 0
	for j := room.Rect.Min.Y; j < room.Rect.Max.Y; j += grid.TileSize.Y {
		for i := room.Rect.Min.X; i < room.Rect.Max.X; i += grid.TileSize.X {
			rec := pixel.R(i, j, i+(room.Rect.W()/tilesPerRow), j+(room.Rect.H()/tilesPerRow))
			tl := makeTile(rec, xInd, yInd)
		VertexEval:
			for _, obst := range room.Blocks {
				if vecDist(obst.Center, rec.Center()) < obst.Radius+4 {
					tl.Enterable = false
					break VertexEval
				}

				// for _, vertex := range obst.Vertices {
				// 	if rec.Contains(vertex) {
				// 		tl.Enterable = false
				// 		break VertexEval
				// 	}
				// }
			}
			// if (start.X >= i && start.Y >= j) && (start.X < (i+room.Rect.W()/tilesPerRow) && start.Y < (j+room.Rect.H()/tilesPerRow)) {
			// 	grid.StartIndex = (yInd * grain) + xInd
			// }
			// if (goal.X >= i && goal.Y >= j) && (goal.X < (i+room.Rect.W()/tilesPerRow) && goal.Y < (j+room.Rect.H()/tilesPerRow)) {
			// 	grid.GoalIndex = (yInd * grain) + xInd
			// }
			grid.GridMap = append(grid.GridMap, tl)
			xInd++
		}
		xInd = 0
		yInd++
	}

	// img := imdraw.New(nil)
	// for _, tl := range grid.GridMap {
	// 	img.Color = colornames.Red
	// 	if tl.Enterable {
	// 		img.Color = colornames.Blue
	// 	}
	// 	img.Push(tl.Rect.Min)
	// 	img.Push(tl.Rect.Max)
	// 	img.Rectangle(0)
	// }
	// img.Draw(grid.Room.Target)

	// return grid
	room.GridRepresentation = grid
}

func traceBack(gridMap *[]Tile, traceInd int, goal pixel.Vec) pixel.Vec {
	for (*gridMap)[traceInd].foundFrom != -1 {
		next := (*gridMap)[traceInd].foundFrom
		if (*gridMap)[next].foundFrom == -1 {
			return (*gridMap)[traceInd].Rect.Center()
		}
		traceInd = next
	}
	return goal

}

func (grid *Grid) tileDist(ind1, ind2 int) float64 {
	v1 := grid.GridMap[ind1].Rect.Center()
	v2 := grid.GridMap[ind2].Rect.Center()
	return vecDist(v1, v2)

}

// AStar uses a* pathfinding to go between grid's start and goal posns
func AStar(grid Grid, startPosn, goal pixel.Vec, cir *Agent) pixel.Vec {

	tileGrid := make([]Tile, len(grid.GridMap))
	copy(tileGrid, grid.GridMap)

	min := grid.Room.Rect.Min

	shiftedStart := startPosn.Sub(min)
	shiftedGoal := goal.Sub(min)

	/* Will crash if monster is out of room bounds (if pushed from contact) */
	startXInd := int(math.Floor(shiftedStart.X / grid.TileSize.X))
	startYInd := int(math.Floor(shiftedStart.Y / grid.TileSize.Y))
	startIndex := (startYInd * grid.TilesPerRow) + startXInd
	tileGrid[startIndex].gScore = 0
	tileGrid[startIndex].fScore = 0

	// fmt.Println(min, startPosn, shiftedStart, startIndex, "\n")

	goalXInd := int(math.Floor(shiftedGoal.X / grid.TileSize.X))
	goalYInd := int(math.Floor(shiftedGoal.Y / grid.TileSize.Y))
	goalIndex := (goalYInd * grid.TilesPerRow) + goalXInd

	for _, blob := range cir.Monsters {
		if blob.Posn != startPosn {
			shiftedBlob := blob.Posn.Sub(min)
			blobXInd := int(math.Ceil(shiftedBlob.X / float64(grid.TilesPerRow)))
			blobYInd := int(math.Ceil(shiftedBlob.Y / float64(grid.TilesPerRow)))
			blobIndex := (blobYInd * grid.TilesPerRow) + blobXInd

			tileGrid[blobIndex].Enterable = false
		}
	}

	prior := make(PriorityQueue, 0)
	heap.Init(&prior)

	start := &Elem{
		Value:    startIndex,
		Priority: tileGrid[startIndex].fScore,
	}
	heap.Push(&prior, start)

	for prior.Len() > 0 {
		currentInd := heap.Pop(&prior).(int)
		if currentInd == goalIndex {
			return traceBack(&tileGrid, currentInd, goal)
		}

		tileGrid[currentInd].reviewed = true

		neighbors := make([]int, 0, 8)
		for xOffset := -1; xOffset <= 1; xOffset++ {
			for yOffset := -1; yOffset <= 1; yOffset++ {

				if !(xOffset == 0 && yOffset == 0) {

					neighbor := (currentInd + (yOffset * grid.TilesPerRow) + xOffset)
					if neighbor >= 0 && neighbor < len(grid.GridMap) {

						if tileGrid[neighbor].Enterable {
							neighbors = append(neighbors, neighbor)
						}

					}

				}

			}
		}

		// fmt.Println("neighbor length:", len(neighbors))

		for _, neighbor := range neighbors {
			// fmt.Print("checking neighbor: ")
			if !tileGrid[neighbor].reviewed {
				// fmt.Print("not reviewed: ")
				possibleGScore := tileGrid[currentInd].gScore + grid.tileDist(currentInd, neighbor)
				if possibleGScore < tileGrid[neighbor].gScore {
					// fmt.Print("changing foundfrom")
					tileGrid[neighbor].foundFrom = currentInd
					tileGrid[neighbor].gScore = possibleGScore
					tileGrid[neighbor].fScore = possibleGScore + grid.tileDist(neighbor, goalIndex)

					item := &Elem{
						Value:    neighbor,
						Priority: tileGrid[neighbor].fScore,
						// Index:    0,
					}
					heap.Push(&prior, item)
				}
			}
		}

	}
	// fmt.Println("Didn't Find Goal")
	return startPosn
}
