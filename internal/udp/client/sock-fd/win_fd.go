//go:build windows

package sockfd

// Here is a hack to extract windows socket handle

import (
	"golang.org/x/sys/windows"
	"net"
	"sync"
	"syscall"
	"unsafe"
)

type operation struct {
	o          syscall.Overlapped
	runtimeCtx uintptr
	mode       int32
	errno      int32
	qty        uint32
	fd         *fd
	buf        syscall.WSABuf
	msg        windows.WSAMsg
	sa         syscall.Sockaddr
	rsa        *syscall.RawSockaddrAny
	rsan       int32
	handle     syscall.Handle
	flags      uint32
	bufs       []syscall.WSABuf
}

type fdMutex struct {
	state uint64
	rsema uint32
	wsema uint32
}

type pollDesc struct {
	runtimeCtx uintptr
}

type fileKind byte

type fd struct {
	fdmu           fdMutex
	Sysfd          syscall.Handle
	rop            operation
	wop            operation
	pd             pollDesc
	l              sync.Mutex
	lastbits       []byte
	readuint16     []uint16
	readbyte       []byte
	readbyteOffset int
	csema          uint32
	skipSyncNotif  bool
	IsStream       bool
	ZeroReadIsEOF  bool
	isFile         bool
	kind           fileKind
}

type netFD struct {
	pfd         fd
	family      int
	sotype      int
	isConnected bool
	net         string
	laddr       net.Addr
	raddr       net.Addr
}

type udpConn struct {
	fd *netFD
}

func GetFd1(conn *net.UDPConn) syscall.Handle {
	c := (*udpConn)(unsafe.Pointer(conn))
	h := c.fd.pfd.Sysfd
	return h
}

func GetFd(conn *net.UDPConn) (syscall.Handle, error) {
	connPtr := (*uintptr)(unsafe.Pointer(conn))
	fdPtr := unsafe.Pointer(*connPtr)
	handlePtr := (*syscall.Handle)(unsafe.Add(fdPtr, 16))
	return *handlePtr, nil
}
