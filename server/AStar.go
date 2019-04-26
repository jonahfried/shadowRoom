package main

import (
	"container/heap"
	"math"

	"github.com/faiface/pixel"
)

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
	room.gridRepresentation = grid
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
func (game *Game) AStar(startPosn, goal pixel.Vec) pixel.Vec {
	grid := game.Room.gridRepresentation

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

	for _, blob := range game.Monsters {
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
