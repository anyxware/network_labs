package udp

import (
	"fmt"
	"log"
	"net"
	"syscall"

	sockfd "network/internal/udp/sock-fd"
)

type sender struct {
	conn *net.UDPConn
}

func NewSender() (*sender, error) {
	var srvAddr string
	fmt.Print("Enter receiver address: ")
	fmt.Scan(&srvAddr)

	addr, err := net.ResolveUDPAddr("udp4", srvAddr)
	if err != nil {
		log.Fatal(err)
	}
	
	c, err := newSender(nil, addr)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newSender(localAddr, remoteAddr *net.UDPAddr) (*sender, error) {
	c, err := net.DialUDP("udp4", localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}
	return &sender{conn: c}, nil
}

func (s *sender) Send(msg []byte) {
	var ttl int
	fmt.Print("Enter ttl: ")
	fmt.Scan(&ttl)

	err := s.setTTL(ttl)
	if err != nil {
		log.Printf("failed set ttl for multicast: %s", err)
		return
	}

	s.write(msg)
	if !s.shouldWaitResponse() {
		return
	}
	s.read()
}

func (s *sender) write(data []byte) {
	_, err := s.conn.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *sender) read() {
	buffer := make([]byte, 1024)
	n, err := s.conn.Read(buffer)
	if err != nil {
		log.Print(err)
	} else {
		log.Print(string(buffer[:n]))
	}
}

func (s *sender) setTTL(ttl int) error {
	fd, err := sockfd.GetFd(s.conn)
	if err != nil {
		return err
	}
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_MULTICAST_TTL, ttl)
}

func (s *sender) shouldWaitResponse() bool {
	cType := connType(s.conn.RemoteAddr().String())
	return cType == unicast
}

func (s *sender) Destroy() {
	s.conn.Close()
}
