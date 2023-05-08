package connection

// This package is used to describe connection as unicast, broadcast or multicast

import (
	"net"
	"strings"
)

type Mode int

const (
	Unicast Mode = iota
	Broadcast
	Multicast
)

var (
	mcAddr = []byte{239, 0, 0, 1}
)

func Type(addr string) Mode {
	addrSplit := strings.Split(addr, ":")
	ipString := addrSplit[0]
	ip := net.ParseIP(ipString)
	if ip[15] == 255 {
		return Broadcast
	} else if ip[12] == mcAddr[0] &&
		ip[13] == mcAddr[1] &&
		ip[14] == mcAddr[2] &&
		ip[15] == mcAddr[3] {
		return Multicast
	}
	return Unicast
}
