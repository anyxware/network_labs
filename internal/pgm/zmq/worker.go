package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	zctx, _ := zmq.NewContext()

	s, _ := zctx.NewSocket(zmq.SUB)
	s.Connect("pgm://224.0.0.1:5555")
	s.SetSubscribe("dev.to")

	for {
		_, err := s.Recv(0)
		if err != nil {
			panic(err)
		}

		if msg, err := s.Recv(0); err != nil {
			panic(err)
		} else {
			fmt.Println("Received message by PGM:", msg)
		}
	}
}
