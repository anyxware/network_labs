package main

import (
	"fmt"
	"log"
	"net"
	"network/internal/udp/client"
)

func readAddr() *net.UDPAddr {
	var srvAddr string
	fmt.Print("Enter receiver address: ")
	fmt.Scan(&srvAddr)

	addr, err := net.ResolveUDPAddr("udp4", srvAddr)
	if err != nil {
		log.Fatal(err)
	}

	return addr
}

func main() {

	addr := readAddr()

	s, err := client.NewClient(nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Destroy()

	for {
		var msg string
		fmt.Print("Enter message: ")
		fmt.Scan(&msg)

		if msg == "new" {
			addr = readAddr()
			s.ChangeReceiver(addr)
			continue
		}

		s.Send([]byte(msg))
	}
}
