package main

import (
	"fmt"
	"log"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	zctx, _ := zmq.NewContext()

	s, _ := zctx.NewSocket(zmq.PUB)
	s.Bind("pgm://224.0.0.1:5555")

	log.Print("start publishing")
	i := 0
	for {
		time.Sleep(5 * time.Second)
		s.Send("dev.to", zmq.SNDMORE)
		s.Send("lera", 0)
		fmt.Println(i)
		i++
	}
}
