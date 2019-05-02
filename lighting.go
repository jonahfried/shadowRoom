package main

import (
	"math"
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

func illuminate(room Place, cir Agent, point *imdraw.IMDraw) {
	point.Clear()
	var anglesToCheck []float64
	var shadedRoomCorners []pixel.Vec
	var obstructedPoints []pixel.Vec

	for _, block := range room.Blocks {
		getAnglesToCheck(&anglesToCheck, block, cir.Posn)
		getShadedRoomCorners(&shadedRoomCorners, room, block, cir.Posn)
		getObstructedPoints(&obstructedPoints, anglesToCheck, room, block, cir.Posn)

		shadePointsByViewMode(obstructedPoints, point, cir)
		shadeObstructedPointsToCorners(obstructedPoints, shadedRoomCorners, point, cir)
		shadeBetweenCorners(shadedRoomCorners, obstructedPoints[0], point, cir) // len(obstructedPoints) >= 1

	}
}

func getObstructedPoints(obstructedPoints *[]pixel.Vec, anglesToCheck []float64, room Place, block Obstacle, posn pixel.Vec) {
	*obstructedPoints = make([]pixel.Vec, 0)
	for _, angle := range anglesToCheck {
		vec := Obstruct(posn, angle, room, block)
		*obstructedPoints = append(*obstructedPoints, vec)
	}
}

func getShadedRoomCorners(shadedRoomCorners *[]pixel.Vec, room Place, block Obstacle, posn pixel.Vec) {
	*shadedRoomCorners = make([]pixel.Vec, 0) // Can just set len to 0 (instead of allocating each time)?
	for _, vertex := range room.Vertices {
		theta := math.Atan2((vertex.Y - posn.Y), (vertex.X - posn.X))
		landed := Obstruct(posn, theta, room, block)
		if math.Abs(vecDist(landed, posn)-vecDist(vertex, posn)) > 1 { // Magic Number
			*shadedRoomCorners = append(*shadedRoomCorners, vertex)
		}
	}
}

func shadeBetweenCorners(shadedRoomCorners []pixel.Vec, obstructedPoint pixel.Vec, point *imdraw.IMDraw, cir Agent) {
	if len(shadedRoomCorners) > 1 {
		for vertexInd := 1; vertexInd < len(shadedRoomCorners); vertexInd++ {
			vecs := [3]pixel.Vec{shadedRoomCorners[vertexInd-1], shadedRoomCorners[vertexInd], obstructedPoint}
			shadePointsByViewMode(vecs[0:3], point, cir)
		}
	}
}

func getAnglesToCheck(anglesToCheck *[]float64, block Obstacle, posn pixel.Vec) {
	*anglesToCheck = make([]float64, 0, 10) // Can just set len to 0 (instead of allocating each time)?
	for _, vertex := range block.Vertices {
		theta := math.Atan2((vertex.Y - posn.Y), (vertex.X - posn.X))
		*anglesToCheck = append(*anglesToCheck, theta)
	}
	length := len(*anglesToCheck)
	for k := 0; k < length; k++ {
		for offset := -.000001; offset <= .000001; offset += .000001 { // Magic Number
			*anglesToCheck = append(*anglesToCheck, (*anglesToCheck)[k]+offset)
		}
	}
	sort.Float64s(*anglesToCheck)
}

func shadeObstructedPointsToCorners(obstructedPoints, shadedRoomCorners []pixel.Vec, point *imdraw.IMDraw, cir Agent) {
	for _, vertex := range shadedRoomCorners {
		for vecInd := 1; vecInd < len(obstructedPoints); vecInd++ {
			vecs := [3]pixel.Vec{vertex, obstructedPoints[vecInd], obstructedPoints[vecInd-1]}
			shadePointsByViewMode(vecs[0:3], point, cir)
		}
	}
}

func shadePointsByViewMode(vecs []pixel.Vec, point *imdraw.IMDraw, cir Agent) {
	for _, vec := range vecs {
		point.Push(vec)
	}
	point.Polygon(0)
}

// playerTorch adds fading light (white circles) around an Agent's posn
func (p *Agent) playerTorch(level float64, count, spacing int, room *Place) {
	room.Target.SetComposeMethod(pixel.ComposeOver)
	img := imdraw.New(nil)
	img.Precision = 32
	col := (pixel.ToRGBA(colornames.Cornsilk)).Mul(pixel.Alpha(level)) // (pixel.ToRGBA(colornames.Whitesmoke)).Mul(pixel.Alpha(cir.Level))
	for fade := 1; fade < count; fade++ {
		img.Color = col
		img.Push(p.Posn)
		img.Circle(float64(fade*spacing), 0)
	}
	img.Draw(room.Target)
	room.Target.SetComposeMethod(pixel.ComposeIn)

}
