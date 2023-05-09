package main

import (
	"fmt"
	"log"

	"network/internal/udp"
)

type Client interface {
	Send(msg []byte)
	Destroy()
}

func main() {
	var c Client
	c, err := udp.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Destroy()

	for {
		var msg string
		fmt.Print("Enter message: ")
		fmt.Scan(&msg)
		c.Send([]byte(msg))
	}
}
