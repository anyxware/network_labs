package main

import (
	"log"

	"network/internal/udp"
)

type Server interface {
	Serve() error
	Destroy()
}

func main() {
	proto := "pgm"

	var srv Server
	var err error
	switch proto {
	case "udp":
		srv, err = udp.NewServer()
	}
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Destroy()

	err = srv.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
