package udp

import "net"

func newMCServer(srvAddr *net.UDPAddr) (*server, error) {
	c, err := net.ListenMulticastUDP("udp4", nil, srvAddr)
	if err != nil {
		return nil, err
	}
	return &server{conn: c}, nil
}
