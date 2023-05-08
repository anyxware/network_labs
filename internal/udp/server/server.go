package server

// Simple udp listener implementation

import (
	"log"
	"net"
	"network/internal/udp/connection"
)

type Server struct {
	msgCount int64
	conn     *net.UDPConn
}

func NewServer(srvAddr *net.UDPAddr) (*Server, error) {
	c, err := net.ListenUDP("udp4", srvAddr)
	if err != nil {
		return nil, err
	}
	return &Server{conn: c}, nil
}

func (s *Server) Serve() {
	for {
		senderAddr := s.read()
		if !s.shouldRespond() {
			continue
		}
		s.write([]byte("lera"), senderAddr)
	}
}

func (s *Server) read() *net.UDPAddr {
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

func (s *Server) write(data []byte, addr *net.UDPAddr) {
	_, err := s.conn.WriteToUDP(data, addr)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) shouldRespond() bool {
	connType := connection.Type(s.conn.LocalAddr().String())
	return connType == connection.Unicast
}

func (s *Server) Destroy() {
	s.conn.Close()
}
