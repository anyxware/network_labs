package main

import (
	"log"

	"network/internal/udp"
)

type Server interface {
	Serve()
	Destroy()
}

func main() {
	var srv Server
	srv, err := udp.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Destroy()
	srv.Serve()
}
