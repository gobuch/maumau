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
	router  *http.ServeMux
	clients []client
	server  *http.Server
	debug   bool
}

func newServer(addr string) *server {
	return &server{
		game:   newGame(),
		router: http.NewServeMux(),
		server: &http.Server{
			Addr:           addr,
			ReadTimeout:    3 * time.Second,
			WriteTimeout:   3 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
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
