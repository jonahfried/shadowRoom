package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

// Move is an alias for string
type Move string

type client struct {
	ws  *websocket.Conn
	msg chan string
}

// Server holds the worldstate and the player connections
type Server struct {
	p1    *client
	World Game
}

// Sends worldstate through websock conn to the client
func (serv *Server) broadcast() {
	// fmt.Println("Broadcasting the WorldState", serv.World.Room)
	err := websocket.JSON.Send(serv.p1.ws, serv.World)
	if err != nil {
		log.Fatal("Error while broadcasting the WS:", err)
	}
}

func (serv *Server) broadcastRoom() {
	err := websocket.JSON.Send(serv.p1.ws, serv.World.Room)
	if err != nil {
		log.Fatal("Error while broadcasting the WS:", err)
	}
}

// Listens through the connections for commands from the clients
// And updates the game
func (serv *Server) listen() {
	fmt.Println("LISTENING")
	tick := time.Tick(100 * time.Millisecond)
	for {
		select {
		case msg := <-serv.p1.msg:
			fmt.Println(msg)
			serv.World.Player.moveHandler(msg)
			serv.broadcast()
			fmt.Println("Broadcasted")
		case <-tick:
			serv.World.Player.playerKinamatics(&serv.World.Room)

			// serv.game.updateMonsters()

			// serv.game.updateShots()

			serv.broadcast()
		}
	}
}
