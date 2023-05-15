package main

import (
	"log"

	"network/internal/udp"
)

type Listener interface {
	Listen() error
	Destroy()
}

func main() {
	proto := "pgm"

	var listener Listener
	var err error
	switch proto {
	case "udp":
		listener, err = udp.NewListener()
	}
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Destroy()

	err = listener.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
