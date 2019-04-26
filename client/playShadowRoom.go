package main

import (
	"fmt"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"golang.org/x/net/websocket"
)

func listen(receiver chan Game, ws *websocket.Conn) {
	var g Game
	err := websocket.JSON.Receive(ws, &g)
	for err == nil {
		receiver <- g
		err = websocket.JSON.Receive(ws, &g)
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
	game := makeGame(win, false)

	receiver := make(chan Game, 10)
	go listen(receiver, ws)
	tick := time.Tick(20 * time.Millisecond)

	for !win.Closed() {
		select {
		case newGameState := <-receiver:
			game.Player.Posn = newGameState.Player.Posn
			game.Player.Vel = newGameState.Player.Vel
			fmt.Println(newGameState)
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
