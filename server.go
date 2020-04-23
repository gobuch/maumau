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
	}
}

func (s *server) run() error {
	s.routes()
	s.server.Handler = s.router
	return s.server.ListenAndServe()
}

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
