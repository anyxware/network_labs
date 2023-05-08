package server

import "net"

func NewMCListener(srvAddr *net.UDPAddr) (*Server, error) {
	c, err := net.ListenMulticastUDP("udp4", nil, srvAddr)
	if err != nil {
		return nil, err
	}
	return &Server{conn: c}, nil
}
