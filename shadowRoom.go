package main

import (
	"flag"
	"fmt"
	_ "image/png"
	"math"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

// Acting main function
func run(win *pixelgl.Window, devMode, noSpawns bool) pixel.Vec {
	var cir = MakeAgent(0, 0, win, devMode) // player

	// TODO: Make room bounds relative to Bounds
	room := MakePlace(pixel.R(-700, -600, 700, 600), 11)
	room.Target.SetMatrix(pixel.IM.Moved(room.Target.Bounds().Center())) //Move this out of loop?
	room.ToGrid(40)

	point := imdraw.New(nil)
	point.Color = colornames.Black

	frameRate := time.Tick(time.Millisecond * 17)
	fiveSec := time.Tick(time.Second * 5)
	thirtySec := time.Tick(time.Second * 30)

	last := time.Now()
	frames := 0.0
	seconds := 0.0
	// Main Draw Loop:
	for !win.Closed() && cir.Health > 0 {
		dt := time.Since(last).Seconds()
		last = time.Now()
		seconds += dt
		if frames > 1000 {
			frames = 0
			seconds = 0
		}
		frames++

		select {
		case <-frameRate:

		case <-fiveSec:
			if !noSpawns {
				cir.Monsters = append(cir.Monsters, MakeCreature(&room, &cir))
			}
		case <-thirtySec:
			room.presentBoost()

		}

		win.Clear(colornames.Black)
		room.Disp()

		cir.Update(win, &room)

		cir.Cam.Attract(cir.Posn)
		cir.Cam.Matrix = pixel.IM.Moved(win.Bounds().Center().Sub(cir.Cam.Posn))

		win.SetMatrix(cir.Cam.Matrix)

		room.Target.Clear(pixel.Alpha(0))
		cir.Light(&room)
		room.Disp()

		updateMonsters(&cir.Monsters, &room, &cir)

		cir.DispShots(room.Target)
		room.Target.Draw(win, pixel.IM) //.Moved(win.Bounds().Center()))

		if devMode {
			fpsDisp(frames/seconds, cir.Posn, win)
		}

		cir.Disp(win)
		illuminate(room, cir, point)
		point.Draw(win)

		win.Update()
		// time.Sleep(1 / 2 * time.Second)
	}
	return cir.Posn
}

func fpsDisp(fps float64, posn pixel.Vec, win *pixelgl.Window) {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(posn.Add(pixel.V(-20, 30)), basicAtlas)
	basicTxt.Color = colornames.Red
	fmt.Fprintln(basicTxt, "fps:", math.Round(fps))
	basicTxt.Draw(win, pixel.IM)
}

func starter() {
	devMode := flag.Bool("dev", false, "runs with access to dev buttons")
	noSpawns := flag.Bool("noSpawns", false, "stop the spawning of enemies")
	flag.Parse()

	win := getWindow()
	finalPosn := run(win, *devMode, *noSpawns)

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)
		endscreen := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		endTxt := text.New(finalPosn, endscreen) // Position ?
		endTxt.Color = colornames.Black
		fmt.Fprint(endTxt, "Game Over!\n(space to restart)")
		endTxt.Draw(win, pixel.IM)
		if win.JustPressed(pixelgl.KeySpace) {
			finalPosn = run(win, *devMode, *noSpawns)
		}
		win.Update()
	}

}

func getWindow() *pixelgl.Window {
	var win *pixelgl.Window
	cfg := pixelgl.WindowConfig{
		Title:     "shadowRoom",
		Bounds:    pixel.R(0, 0, 1350, 725),
		VSync:     true,
		Resizable: true,
	}
	var err error
	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	return win
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	pixelgl.Run(starter)
}
