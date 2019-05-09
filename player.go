package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Agent keeps a posn, vel, acc, and a *pixel.IMDraw.
// Used to keep all information for the movable image together.
type Agent struct {
	Posn, Vel, Acc   pixel.Vec
	maxAcc, maxSpeed float64
	Radius           float64

	Cam Camera

	Health     int
	TorchLevel float64
	Torches    int

	GunType int
	Bullets map[int]int

	Img *imdraw.IMDraw
}

// MakeAgent creates a new agent starting at a given (x, y) coordinate
func MakeAgent(x, y float64, win *pixelgl.Window) (cir Agent) {
	cir.Posn = pixel.V(x, y)
	cir.Vel, cir.Acc = pixel.ZV, pixel.ZV
	cir.maxAcc = 2
	cir.maxSpeed = 8
	cir.Radius = 20

	cir.Cam = MakeCamera(cir.Posn, win)

	cir.Health = 100
	cir.TorchLevel = 3
	cir.Torches = 1

	cir.GunType = 1
	cir.Bullets = make(map[int]int)
	cir.Img = imdraw.New(nil)
	cir.Img.Color = colornames.Purple

	return cir
}

// A Shot keeps track of the players attacks
type Shot struct {
	Posn1, Posn2, Vel pixel.Vec
	GunType           int
	color             pixel.RGBA
}

func (bullet *Shot) update() {
	bullet.Posn1 = bullet.Posn1.Add(bullet.Vel)
	bullet.Posn2 = bullet.Posn2.Add(bullet.Vel)

	if bullet.GunType == 2 {
		bullet.Vel = bullet.Vel.Scaled(.94)
		bullet.color = (bullet.color).Add(pixel.ToRGBA(colornames.Darkorange).Mul(pixel.Alpha(.03))).Add(pixel.Alpha(.02)) //  Mul(pixel.Alpha(1))
	}
}

// Checks to see if a given position should receive a collision force from a list of obstacles
func collision(blocks []Obstacle, posn pixel.Vec, radius float64, vel *pixel.Vec) (force pixel.Vec) {
	posn = posn.Add(*vel)
	minDist := radius + ObstacleBorderWidth/2
	for _, obst := range blocks {
		obst.Vertices = append(obst.Vertices, obst.Vertices[0]) // adding the first element to the end to complete the shape
		for vertexInd := 1; vertexInd < len(obst.Vertices); vertexInd++ {
			closest := closestPointOnSegment(obst.Vertices[vertexInd], obst.Vertices[vertexInd-1], posn)
			dist := posn.Sub(closest)
			if dist.Len() < (minDist) {
				dir := dist.Unit()
				*vel = vel.Sub(dir.Scaled(vel.Dot(dir)))
				dir = dir.Scaled((minDist - dist.Len()))
				force = force.Add(dir)
			}
		}

	}
	return force
}

// Returns the point on a line segment closest to a given position. (Either one of the two ends or a point between them)
func closestPointOnSegment(end1, end2, posn pixel.Vec) (cloestest pixel.Vec) {
	segVec := end2.Sub(end1)
	unitSegVec := segVec.Scaled(1 / segVec.Len())
	posnOffset := posn.Sub(end1)
	projMag := posnOffset.Dot(unitSegVec)
	projVec := unitSegVec.Scaled(projMag)
	var closest pixel.Vec

	if projMag < 0 {
		closest = end1
	} else if projMag > segVec.Len() {
		closest = end2
	} else {
		closest = end1.Add(projVec)
	}
	return closest
}

const friction = .8

// playerKinamatics runs necessary per-frame movements on agent.
func (cir *Agent) playerKinamatics(room *Place) {
	cir.Acc = limitVecMag(cir.Acc, cir.maxAcc)
	cir.Vel = cir.Vel.Add(cir.Acc)
	cir.Vel = limitVecMag(cir.Vel, cir.maxSpeed)

	offset := collision(room.Blocks, cir.Posn, cir.Radius, &cir.Vel)
	cir.Vel = cir.Vel.Add(offset.Scaled(.2))

	cir.Posn = cir.Posn.Add(cir.Vel)
	cir.Vel = limitVecMag(cir.Vel.Scaled(friction), 10)

	cir.Posn.X = math.Max(cir.Posn.X-20, room.Rect.Min.X) + 20
	cir.Posn.X = math.Min(cir.Posn.X+20, room.Rect.Max.X) - 20

	cir.Posn.Y = math.Max(cir.Posn.Y-20, room.Rect.Min.Y) + 20
	cir.Posn.Y = math.Min(cir.Posn.Y+20, room.Rect.Max.Y) - 20
}
