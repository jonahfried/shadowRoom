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
func run(win *pixelgl.Window, devMode, noSpawn bool) pixel.Vec {
	game := makeGame(win, devMode)

	point := imdraw.New(nil)
	point.Color = colornames.Black

	frameRate := time.Tick(time.Millisecond * 17)
	fiveSec := time.Tick(time.Second * 5)
	thirtySec := time.Tick(time.Second * 30)

	last := time.Now()
	frames := 0.0
	seconds := 0.0
	// Main Draw Loop:
	for !win.Closed() && game.Player.Health > 0 {
		if !win.Focused() {
			game.Paused = true
		}
		if game.Paused {
			if win.JustPressed(pixelgl.KeyP) {
				game.Paused = false
			}
			win.Update()
			continue
		}
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
			if !noSpawn {
				game.Monsters = append(game.Monsters, MakeCreature(&game.Room, &game.Player))
			}
		case <-thirtySec:
			game.PowerUps = append(game.PowerUps, MakeShotgun(game.Room.safeSpawnInRoom(10)))
		}

		PressHandler(win, &game)
		ReleaseHandler(win, &game)

		game.Player.playerKinamatics(&game.Room)

		game.Player.Cam.Attract(game.Player.Posn)
		game.Player.Cam.Matrix = pixel.IM.Moved(win.Bounds().Center().Sub(game.Player.Cam.Posn))
		win.SetMatrix(game.Player.Cam.Matrix)

		game.updateMonsters()
		game.updateShots()
		game.updatePowerUps()

		game.Disp(win)
		win.Update()
	}
	return game.Player.Posn
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
	noSpawn := flag.Bool("noSpawn", false, "stop the spawning of enemies")
	flag.Parse()

	win := getWindow()
	finalPosn := run(win, *devMode, *noSpawn)

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)
		endscreen := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		endTxt := text.New(finalPosn, endscreen) // Position ?
		endTxt.Color = colornames.Black
		fmt.Fprint(endTxt, "Game Over!\n(space to restart)")
		endTxt.Draw(win, pixel.IM)
		if win.JustPressed(pixelgl.KeySpace) {
			finalPosn = run(win, *devMode, *noSpawn)
		}
		win.Update()
	}

}

func getWindow() *pixelgl.Window {
	var win *pixelgl.Window
	cfg := pixelgl.WindowConfig{
		Title:     "Shadow Room!",
		Bounds:    pixel.R(0, 0, 1350, 725),
		VSync:     true,
		Resizable: true,
	}
	var err error
	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	// win.SetCursorVisible(false)
	return win
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	pixelgl.Run(starter)
}
