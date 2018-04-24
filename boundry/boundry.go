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
