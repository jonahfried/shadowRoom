package main

import (
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"golang.org/x/net/websocket"
)

func listen(receiver chan Agent, ws *websocket.Conn) {
	var p Agent
	err := websocket.JSON.Receive(ws, &p)
	for err == nil {
		receiver <- p
		err = websocket.JSON.Receive(ws, &p)
	}
}

func initializeConn() *websocket.Conn {
	origin := "http://localhost/"
	url := "ws://localhost:1234/"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	return ws
}

func run() {
	win := getWindow()
	ws := initializeConn()
	game := makeGame(win, false, ws)

	receiver := make(chan Agent, 20)
	go listen(receiver, ws)
	tick := time.Tick(75 * time.Millisecond)

	for !win.Closed() {
		select {
		case newPlayer := <-receiver:
			game.Player.Posn = newPlayer.Posn
			game.Player.Vel = newPlayer.Vel
		case <-tick:
		}
		PressHandler(win, ws)
		// ReleaseHandler(win, ws)
		game.Player.playerKinamatics(&game.Room)
		game.Player.Cam.Attract(game.Player.Posn)
		game.Player.Cam.Matrix = pixel.IM.Moved(win.Bounds().Center().Sub(game.Player.Cam.Posn))
		win.SetMatrix(game.Player.Cam.Matrix)

		win.Clear(colornames.Black)
		game.Room.Disp()
		game.Room.Target.Clear(pixel.Alpha(0))
		game.Player.playerTorch(game.Level, game.Count, game.Spacing, &game.Room)
		game.Room.Disp()
		game.DispShots(game.Room.Target)
		game.Room.Target.Draw(win, pixel.IM)
		game.Player.Disp(win)
		illuminate(game.Room, game.Player, win)

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
	pixelgl.Run(run)
}
