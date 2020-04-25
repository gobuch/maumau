package main

import (
	"flag"
	"log"
)

func main() {
	var (
		flagPort  = flag.String("p", ":8080", "Port für den Server")
		flagDebug = flag.Bool("d", false, "Starte den Server im Debug-Modus")
	)
	flag.Parse()
	s := newServer(*flagPort)
	s.debug = *flagDebug
	log.Println("Starting server at ", *flagPort)
	err := s.run()
	if err != nil {
		log.Fatal("error main.s.run():", err)
	}
}
