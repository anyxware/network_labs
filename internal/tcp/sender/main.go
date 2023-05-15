package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("CLIENT")
	fmt.Print("Input port: ")
	reader := bufio.NewReader(os.Stdin)
	port, err := reader.ReadString('\n')
	if err != nil {
		println("Failed to read from stdin:", err.Error())
		os.Exit(1)
	}
	port = port[:len(port)-1]
	
	tcpServer, err := net.ResolveTCPAddr("tcp", "localhost"+":"+port)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	for {
		fmt.Print("Input message: ")
		reader = bufio.NewReader(os.Stdin)
		_, err := reader.ReadString('\n')

		tmp := make([]byte,655360)
		_, err = conn.Write(tmp)
		if err != nil {
			println("Write data failed:", err.Error())
			os.Exit(1)
		}
	
		// buffer to get data
		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}

		println("Received message:", string(received))
	}

	//conn.Close()
}