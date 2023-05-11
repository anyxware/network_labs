package udp

import (
	"fmt"
	"log"
	"net"
)

type server struct {
	msgCount int64
	conn     *net.UDPConn
}

func NewServer() (*server, error) {
	var serverMode int
	fmt.Println("0 - unicast/broadcast")
	fmt.Println("1 - multicast")
	fmt.Print("Enter mode: ")
	fmt.Scan(&serverMode)

	var srvAddr string
	fmt.Print("Enter serving address: ")
	fmt.Scan(&srvAddr)
	addr, err := net.ResolveUDPAddr("udp4", srvAddr)
	if err != nil {
		return nil, err
	}

	var srv *server
	if serverMode == 0 {
		srv, err = newServer(addr)
	} else {
		srv, err = newMCServer(addr)
	}
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func newServer(srvAddr *net.UDPAddr) (*server, error) {
	c, err := net.ListenUDP("udp4", srvAddr)
	if err != nil {
		return nil, err
	}
	return &server{conn: c}, nil
}

func (s *server) Serve() error {
	for {
		senderAddr := s.read()
		if !s.shouldRespond() {
			continue
		}
		s.write([]byte("lera"), senderAddr)
	}
}

func (s *server) read() *net.UDPAddr {
	buffer := make([]byte, 1024)
	n, addr, err := s.conn.ReadFromUDP(buffer)
	if err != nil {
		log.Print(err)
	} else {
		log.Printf("[%d] %s", s.msgCount, string(buffer[:n]))
		s.msgCount++
	}
	return addr
}

func (s *server) write(data []byte, addr *net.UDPAddr) {
	_, err := s.conn.WriteToUDP(data, addr)
	if err != nil {
		log.Print(err)
	}
}

func (s *server) shouldRespond() bool {
	cType := connType(s.conn.LocalAddr().String())
	return cType == unicast
}

func (s *server) Destroy() {
	s.conn.Close()
}
