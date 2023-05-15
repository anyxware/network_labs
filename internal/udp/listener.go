package udp

import (
	"fmt"
	"log"
	"net"
)

type listener struct {
	msgCount int64
	conn     *net.UDPConn
}

func NewListener() (*listener, error) {
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

	var l *listener
	if serverMode == 0 {
		l, err = newListener(addr)
	} else {
		l, err = newMCListener(addr)
	}
	if err != nil {
		return nil, err
	}

	return l, nil
}

func newListener(addr *net.UDPAddr) (*listener, error) {
	c, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, err
	}
	return &listener{conn: c}, nil
}

func (l *listener) Listen() error {
	for {
		senderAddr := l.read()
		if !l.shouldRespond() {
			continue
		}
		l.write([]byte("lera"), senderAddr)
	}
}

func (l *listener) read() *net.UDPAddr {
	buffer := make([]byte, 1024)
	n, addr, err := l.conn.ReadFromUDP(buffer)
	if err != nil {
		log.Print(err)
	} else {
		log.Printf("[%d] %l", l.msgCount, string(buffer[:n]))
		l.msgCount++
	}
	return addr
}

func (l *listener) write(data []byte, addr *net.UDPAddr) {
	_, err := l.conn.WriteToUDP(data, addr)
	if err != nil {
		log.Print(err)
	}
}

func (l *listener) shouldRespond() bool {
	cType := connType(l.conn.LocalAddr().String())
	return cType == unicast
}

func (l *listener) Destroy() {
	l.conn.Close()
}
