package main

import (
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
	Blocks             []Obstacle
	Vertices           []pixel.Vec
	GridRepresentation Grid
	Booster            Boost

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

	room.Booster.Present = false
	room.Booster.Posn = pixel.V(100, 100)
	room.Booster.Img = imdraw.New(nil)
	room.Booster.Img.Color = colornames.Royalblue
	// room.Booster.Img.Precision = 32

	room.Target = pixelgl.NewCanvas(rect)
	room.Img = imdraw.New(nil)
	return room
}

// Obstacle is a data structure defining a boundry inside a room
// it contains:
//	 a list of vertices defining its bounds
//	 and an *imdraw.IMdraw to store its look.
//	 also stores a Center and Radius set if generated randomly to avoid overlap
type Obstacle struct {
	Vertices []pixel.Vec
	Center   pixel.Vec
	Radius   float64

	Img *imdraw.IMDraw
}

// MakeObstacle takes in a list of vertices that define its boundries,
// and a bool representing whether or not it is a Room.
// returns an Obstacle
func MakeObstacle(vertices []pixel.Vec, center pixel.Vec, radius float64) (obst Obstacle) {
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
// and a Room to hold it, returning a randomly generated Obstacle
func MakeRandomBlock(radius, stdDev, vertScale float64, rect pixel.Rect, existingBlocks []Obstacle) Obstacle {
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
	return MakeObstacle(vertices, center, radius+stdDev)
}

func vecDist(v1, v2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
}

// Obstruct takes in a "viewers" position, and an angle of viewing.
// it also takes a room for its boundries
// and an Obstacle of interest to determine intersection
func Obstruct(posn pixel.Vec, angle float64, room Place, block Obstacle) (stopPoint pixel.Vec) {
	// TODO: Fix divide by zero error
	if (math.Cos(angle)) == 0 {
		panic("divide by zero")
	}
	slope := math.Sin(angle) / math.Cos(angle)
	yInt := posn.Y - (posn.X * slope)

	extension := 0.000000001

	var blocks [2]Obstacle

	blocks[0] = block
	blocks[1] = MakeObstacle(room.Vertices, pixel.ZV, 0)

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