package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/faiface/pixel"

	"golang.org/x/net/websocket"
)

type hub struct{ clients chan *client }

// WorldState stores the position of the circle to disp
type WorldState struct {
	Posn pixel.Vec
}

func (h *hub) handleWebsocketConnection(ws *websocket.Conn) {
	newClient := client{ws, make(chan string)}
	fmt.Println("New Client", newClient)
	h.clients <- &newClient
	fmt.Println("Added to hub", newClient)

	var m string
	err := websocket.JSON.Receive(ws, &m)
	for err == nil {
		fmt.Println(m)
		newClient.msg <- m
		err = websocket.JSON.Receive(ws, &m)
	}
	fmt.Println(m)
	log.Print(err)
}

func (h *hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var s websocket.Server
	s.Handler = func(ws *websocket.Conn) {
		h.handleWebsocketConnection(ws)
	}
	s.ServeHTTP(w, r)
}

func (h *hub) matchPlayers() {
	for {
		p1 := <-h.clients
		server := Server{p1, makeGame(false)}
		go server.listen()
	}
}

func main() {
	h := hub{make(chan *client)}
	go h.matchPlayers()

	fmt.Println("Beginning Server")
	go http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "index.html") }))
	http.ListenAndServe(":1234", &h)
}
