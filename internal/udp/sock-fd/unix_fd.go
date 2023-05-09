//go:build unix

package sockfd

import (
	"net"
)

func GetFd(conn *net.UDPConn) (int, error) {
	f, err := conn.File()
	if err != nil {
		return 0, err
	}
	fd := f.Fd()
	return int(fd), nil
}
