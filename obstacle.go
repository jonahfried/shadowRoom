package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

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
	obst.Img.Polygon(ObstacleBorderWidth)

	return obst
}

// ObstacleBorderWidth is the amount of pixels thick the line is drawn between vertices.
const ObstacleBorderWidth = 8

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
