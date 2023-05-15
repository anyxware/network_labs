package main

import (
	"fmt"
	"log"

	"network/internal/udp"
)

type Sender interface {
	Send(msg []byte)
	Destroy()
}

func main() {
	var sender Sender
	sender, err := udp.NewSender()
	if err != nil {
		log.Fatal(err)
	}
	defer sender.Destroy()

	for {
		var msg string
		fmt.Print("Enter message: ")
		fmt.Scan(&msg)
		sender.Send([]byte(msg))
	}
}
