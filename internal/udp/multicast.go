package udp

import "net"

func newMCListener(srvAddr *net.UDPAddr) (*listener, error) {
	c, err := net.ListenMulticastUDP("udp4", nil, srvAddr)
	if err != nil {
		return nil, err
	}
	return &listener{conn: c}, nil
}
