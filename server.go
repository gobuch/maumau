package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type server struct {
	game    *Game
	port    string
	router  *http.ServeMux
	clients []*client
	msg     chan []byte
	server  *http.Server
	debug   bool
}

func newServer() *server {
	return &server{
		game:   newGame(),
		port:   *flagPort,
		router: http.NewServeMux(),
		server: &http.Server{
			Addr:           *flagPort,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		debug: *flagDebug,
	}
}

func (s *server) run() error {
	s.routes()
	s.server.Handler = s.router
	return s.server.ListenAndServe()
}

// sendState ranges over all clients creates the hand
// for each player, encodes it to JSON and sends the
// result into the message channel
func (s *server) sendState() {
	for _, c := range s.clients {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetIndent("", "  ")
		s.game.state()
		ps, _ := s.playerState(c.playerID)
		err := enc.Encode(ps)
		if err != nil {
			log.Println("sendState encoding error:", err)
		}
		b := buf.Bytes()
		c.messages <- b
	}
}
