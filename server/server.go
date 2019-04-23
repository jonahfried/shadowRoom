package main

import (
	"fmt"
	"log"
	"time"

	"github.com/faiface/pixel"
	"golang.org/x/net/websocket"
)

// Move is an alias for string
type Move string

type client struct {
	ws  *websocket.Conn
	msg chan Move
}

// Server holds the worldstate and the player connections
type Server struct {
	p1   *client
	game Game
}

// Sends worldstate through websock conn to the client
func (serv *Server) broadcast() {
	// fmt.Println("Broadcasting the WorldState", serv.world)
	err := websocket.JSON.Send(serv.p1.ws, serv.game)
	if err != nil {
		log.Fatal("Error while broadcasting the WS:", err)
	}
}

// Listens through the connections for commands from the clients
// And updates the game
func (serv *Server) listen() {
	fmt.Println("LISTENING")
	tick := time.Tick(15 * time.Millisecond)
	for {
		select {
		case msg := <-serv.p1.msg:
			fmt.Println(msg)
			serv.broadcast()
		case <-tick:

			game.Player.playerKinamatics(&game.Room)

			game.Player.Cam.Attract(game.Player.Posn)
			game.Player.Cam.Matrix = pixel.IM.Moved(win.Bounds().Center().Sub(game.Player.Cam.Posn))

			win.SetMatrix(game.Player.Cam.Matrix)

			game.Room.Target.Clear(pixel.Alpha(0))
			game.Player.playerTorch(game.Level, game.Count, game.Spacing, &game.Room)
			game.Room.Disp()

			game.updateMonsters()

			game.updateShots()
			game.DispShots(game.Room.Target)
			game.Room.Target.Draw(win, pixel.IM) //.Moved(win.Bounds().Center()))

			if devMode {
				fpsDisp(frames/seconds, game.Player.Posn, win)
			}

			game.Player.Disp(win)
			illuminate(game.Room, game.Player, win)

		}
	}
}
