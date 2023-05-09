package udp

import (
	"net"
	"strings"
)

type mode int

const (
	unicast mode = iota
	broadcast
	multicast
)

var (
	mcAddr = []byte{239, 0, 0, 1}
)

func connType(addr string) mode {
	addrSplit := strings.Split(addr, ":")
	ipString := addrSplit[0]
	ip := net.ParseIP(ipString)
	if ip[15] == 255 {
		return broadcast
	} else if ip[12] == mcAddr[0] &&
		ip[13] == mcAddr[1] &&
		ip[14] == mcAddr[2] &&
		ip[15] == mcAddr[3] {
		return multicast
	}
	return unicast
}
