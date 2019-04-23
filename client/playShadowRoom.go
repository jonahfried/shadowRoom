package main

import (
	"fmt"
	"log"
	"time"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"golang.org/x/net/websocket"
)

func listen(receiver chan pixel.Vec, ws *websocket.Conn) {
	var v pixel.Vec
	err := websocket.JSON.Receive(ws, &v)
	// fmt.Println("this is the new vec: ", v)
	for err == nil {
		receiver <- v
		err = websocket.JSON.Receive(ws, &v)
	}
}

func run() {
	win := getWindow()
	origin := "http://localhost/"
	url := "ws://localhost:1234/"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	v := pixel.V(300, 50)
	im := imdraw.New(nil)
	im.Color = colornames.Purple

	receiver := make(chan pixel.Vec)
	go listen(receiver, ws)
	tick := time.Tick(20 * time.Millisecond)

	for !win.Closed() {
		// fmt.Println(v)
		select {
		case v = <-receiver:
			// fmt.Println("received", v)
		case <-tick:
		}

		win.Clear(colornames.Whitesmoke)
		im.Clear()

		im.Push(v)
		im.Circle(20, 0)
		im.Draw(win)

		if win.JustPressed(pixelgl.KeySpace) {
			websocket.JSON.Send(ws, "SENDIND TO SERVER")
			fmt.Println("SENDING TO SERVER")
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
	pixelgl.Run(run)
}
