package client

// Simple udp sender implementation

import (
	"fmt"
	"log"
	"net"
	connfd "network/internal/udp/client/sock-fd"
	"network/internal/udp/connection"
	"syscall"
)

type Client struct {
	conn *net.UDPConn
}

func NewClient(localAddr, remoteAddr *net.UDPAddr) (*Client, error) {
	c, err := net.DialUDP("udp4", localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: c}, nil
}

func (c *Client) Send(msg []byte) {
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

func (c *Client) write(data []byte) {
	_, err := c.conn.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (c *Client) read() {
	buffer := make([]byte, 1024)
	n, err := c.conn.Read(buffer)
	if err != nil {
		log.Print(err)
	} else {
		log.Print(string(buffer[:n]))
	}
}

func (c *Client) setTTL(ttl int) error {
	fd, err := connfd.GetFd(c.conn)
	if err != nil {
		return err
	}
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_MULTICAST_TTL, ttl)
}

func (c *Client) shouldWaitResponse() bool {
	connType := connection.Type(c.conn.RemoteAddr().String())
	return connType == connection.Unicast
}

func (c *Client) Destroy() {
	c.conn.Close()
}

func (c *Client) ChangeReceiver(remoteAddr *net.UDPAddr) {
	newConn, err := net.DialUDP("udp4", nil, remoteAddr)
	if err != nil {
		log.Print(err)
	}
	c.conn = newConn
}
