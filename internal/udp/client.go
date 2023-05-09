package udp

import (
	"fmt"
	"log"
	"net"
	"syscall"

	sockfd "network/internal/udp/sock-fd"
)

type client struct {
	conn *net.UDPConn
}

func NewClient() (*client, error) {
	var srvAddr string
	fmt.Print("Enter receiver address: ")
	fmt.Scan(&srvAddr)

	addr, err := net.ResolveUDPAddr("udp4", srvAddr)
	if err != nil {
		log.Fatal(err)
	}
	
	c, err := newClient(nil, addr)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newClient(localAddr, remoteAddr *net.UDPAddr) (*client, error) {
	c, err := net.DialUDP("udp4", localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}
	return &client{conn: c}, nil
}

func (c *client) Send(msg []byte) {
	var ttl int
	fmt.Print("Enter ttl: ")
	fmt.Scan(&ttl)

	err := c.setTTL(ttl)
	if err != nil {
		log.Printf("failed set ttl for multicast: %s", err)
		return
	}

	c.write(msg)
	if !c.shouldWaitResponse() {
		return
	}
	c.read()
}

func (c *client) write(data []byte) {
	_, err := c.conn.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (c *client) read() {
	buffer := make([]byte, 1024)
	n, err := c.conn.Read(buffer)
	if err != nil {
		log.Print(err)
	} else {
		log.Print(string(buffer[:n]))
	}
}

func (c *client) setTTL(ttl int) error {
	fd, err := sockfd.GetFd(c.conn)
	if err != nil {
		return err
	}
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_MULTICAST_TTL, ttl)
}

func (c *client) shouldWaitResponse() bool {
	cType := connType(c.conn.RemoteAddr().String())
	return cType == unicast
}

func (c *client) Destroy() {
	c.conn.Close()
}
