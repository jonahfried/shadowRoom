package boundry

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

// Place is datastructure that holds a rectangle for its bounds
// as well as a slice of Blocks that it contains
// and a slice of its vertices
type Place struct {
	Rect     pixel.Rect
	Blocks   []Obsticle
	Vertices []pixel.Vec

	Img *imdraw.IMDraw
}

// MakePlace turns specifications of a place, (rectangle and the Obsticles it contains)
// and returns a Place
func MakePlace(rect pixel.Rect, blocks []Obsticle) (room Place) {
	room.Rect = rect

	borderVertices := make([]pixel.Vec, 0, 4)
	borderVertices = append(borderVertices, rect.Min)
	borderVertices = append(borderVertices, pixel.V(rect.Min.X, rect.Max.Y))
	borderVertices = append(borderVertices, rect.Max)
	borderVertices = append(borderVertices, pixel.V(rect.Max.X, rect.Min.Y))

	room.Vertices = borderVertices

	border := MakeObsticle(borderVertices, true)
	room.Blocks = append(blocks, border)

	room.Img = imdraw.New(nil)
	room.Img.Color = colornames.Black
	room.Img.Push(room.Rect.Center().Sub(room.Rect.Size().Scaled(.5)))
	room.Img.Push(room.Rect.Center().Add(room.Rect.Size().Scaled(.5)))
	room.Img.Rectangle(2)
	return room
}

// Obsticle is a data structure defining a boundry inside a room
// it contains:
//	 a list of vertices defining its bounds
//	 and an *imdraw.IMdraw to store its look.
type Obsticle struct {
	Vertices []pixel.Vec
	Img      *imdraw.IMDraw
	IsRoom   bool
}

// MakeObsticle takes in a list of vertices that define its boundries,
// and a bool representing whether or not it is a Room.
// returns an Obsticle
func MakeObsticle(vertices []pixel.Vec, isRoom bool) (obst Obsticle) {
	obst.Vertices = vertices
	obst.IsRoom = isRoom
	obst.Img = imdraw.New(nil)

	obst.Img.Color = colornames.Darkgrey
	for _, vert := range obst.Vertices {
		obst.Img.Push(vert)
	}
	obst.Img.Polygon(0)

	obst.Img.Color = colornames.Slategrey
	for _, vert := range obst.Vertices {
		obst.Img.Push(vert)
	}
	obst.Img.Polygon(1)

	return obst
}

// MakeRandomBlock takes in a standard radius,
// a standard deviation from the radius, a number of vertices to create, (float64)
// and a center (pixel.Vec), returning a randomly generated Obsticle
func MakeRandomBlock(radius, stdDev, vertScale float64, center pixel.Vec) Obsticle {
	vertices := make([]pixel.Vec, 0, 3)
	for angle := 2 * math.Pi / vertScale; angle < math.Pi*2; angle += 2 * math.Pi / vertScale {
		r := rand.NormFloat64()*stdDev + radius
		vertex := pixel.V(r*math.Cos(angle), (r*math.Sin(angle) + radius))
		vertex = vertex.Add(center)
		vertices = append(vertices, vertex)
	}
	return MakeObsticle(vertices, false)
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
	blocks[1] = MakeObsticle(room.Vertices, true)

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
