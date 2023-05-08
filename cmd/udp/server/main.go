package main

import (
	"fmt"
	"log"
	"net"
	"network/internal/udp/server"
)

func main() {
	fmt.Print(`
  _________                           ________                                _________                  .__              
 /   _____/__ ________   ___________  \______ \  __ ________   ___________   /   _____/ ______________  _|__| ____  ____  
 \_____  \|  |  \____ \_/ __ \_  __ \  |    |  \|  |  \____ \_/ __ \_  __ \  \_____  \_/ __ \_  __ \  \/ /  |/ ___\/ __ \ 
 /        \  |  /  |_> >  ___/|  | \/  |    '   \  |  /  |_> >  ___/|  | \/  /        \  ___/|  | \/\   /|  \  \__\  ___/
/_________/____/|   __/ \_____>__|    /_________/____/|   __/ \_____>__|    /_________/\_____>__|    \_/ |__|\_____>_____>
                |__|                                  |__|                                                             
`)

	var mode int
	fmt.Println("0 - unicast/broadcast")
	fmt.Println("1 - multicast")
	fmt.Print("Enter mode: ")
	fmt.Scan(&mode)

	var srvAddr string
	fmt.Print("Enter serving address: ")
	fmt.Scan(&srvAddr)
	s, err := net.ResolveUDPAddr("udp4", srvAddr)
	if err != nil {
		log.Fatal(err)
	}

	var l server.Server
	if mode == 0 {
		l, err = server.NewServer(s)
	} else {
		l, err = server.NewMCListener(s)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer l.Destroy()

	l.Serve()
}
