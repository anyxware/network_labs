package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"bufio"
)

func main() {
	fmt.Println("SERVER")
	fmt.Print("Input port: ")
	reader := bufio.NewReader(os.Stdin)
	port, err := reader.ReadString('\n')
	if err != nil {
		println("Failed to read from stdin:", err.Error())
		os.Exit(1)
	}
	port = port[:len(port)-1]

	listen, err := net.Listen("tcp", "localhost"+":"+port)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	for {
		// incoming request
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
	
		log.Printf("Got your message: %s", string(buffer))
	
		conn.Write([]byte("hello"))
		// close conn
	}
	//conn.Close()
}