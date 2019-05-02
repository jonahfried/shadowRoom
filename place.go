package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// Place is datastructure that holds a rectangle for its bounds
// as well as a slice of Blocks that it contains
// and a slice of its vertices
type Place struct {
	Rect               pixel.Rect
	Blocks             []Obstacle
	Vertices           []pixel.Vec
	GridRepresentation Grid

	Img    *imdraw.IMDraw
	Target *pixelgl.Canvas
}

// MakePlace turns specifications of a place, (rectangle and the Obstacles it contains)
// and returns a Place
func MakePlace(rect pixel.Rect, numBlocks int, blocks ...Obstacle) (room Place) {
	room.Rect = rect

	borderVertices := make([]pixel.Vec, 0, 4)
	borderVertices = append(borderVertices, rect.Min)
	borderVertices = append(borderVertices, pixel.V(rect.Min.X, rect.Max.Y))
	borderVertices = append(borderVertices, rect.Max)
	borderVertices = append(borderVertices, pixel.V(rect.Max.X, rect.Min.Y))

	room.Vertices = borderVertices

	if len(blocks) == 0 {
		blockList := make([]Obstacle, 0)
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

	// room.Booster.Present = false
	// room.Booster.Posn = pixel.V(100, 100)
	// room.Booster.Img = imdraw.New(nil)
	// room.Booster.Img.Color = colornames.Royalblue
	// room.Booster.Img.Precision = 32

	room.Target = pixelgl.NewCanvas(rect)
	room.Img = imdraw.New(nil)
	return room
}

// Returns a random point within radius dist of the room's walls
func (room Place) randomPlaceInRoom(radius float64) pixel.Vec {
	x := rand.Float64()*(room.Rect.W()-2*radius) + room.Rect.Min.X + radius
	y := rand.Float64()*(room.Rect.H()-2*radius) + room.Rect.Min.Y + radius
	return pixel.V(x, y)
}

// safeSpawnInRoom returns a pixel.Vec at least `radius` dist from obstacles.
// If cannot find one within 20 attempts, gives up and just gives the next random location.
func (room Place) safeSpawnInRoom(radius float64) pixel.Vec {
	attempt := room.randomPlaceInRoom(radius)
TryLoop:
	for i := 0; i < 20; i++ {
		for _, block := range room.Blocks {
			if vecDist(block.Center, attempt) < block.Radius+radius {
				attempt = room.randomPlaceInRoom(radius)
				continue TryLoop
			}
		}
		break TryLoop
	}
	return attempt
}
