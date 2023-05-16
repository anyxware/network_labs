package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/ishidawataru/sctp"
)

func serveClient(conn net.Conn, bufsize int) error {
	for {
		buf := make([]byte, bufsize+128) // add overhead of SCTPSndRcvInfoWrappedConn
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read failed: %v", err)
			return err
		}
		log.Printf("read: %d", n)
		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("write failed: %v", err)
			return err
		}
		log.Printf("write: %d", n)
	}
}

func main() {
	var server = flag.Bool("server", false, "")
	var ip = flag.String("ip", "0.0.0.0", "")
	var port = flag.Int("port", 0, "")
	var lport = flag.Int("lport", 0, "")
	var bufsize = flag.Int("bufsize", 256, "")

	flag.Parse()

	ips := []net.IPAddr{}

	for _, i := range strings.Split(*ip, ",") {
		if a, err := net.ResolveIPAddr("ip", i); err == nil {
			log.Printf("Resolved address '%s' to %s", i, a)
			ips = append(ips, *a)
		} else {
			log.Printf("Error resolving address '%s': %v", i, err)
		}
	}

	addr := &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    *port,
	}
	log.Printf("raw addr: %+v\n", addr.ToRawSockAddrBuf())

	if *server {
		ln, err := sctp.ListenSCTP("sctp", addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("Listen on %s", ln.Addr())

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatalf("failed to accept: %v", err)
			}
			log.Printf("Accepted Connection from RemoteAddr: %s", conn.RemoteAddr())
			wconn := sctp.NewSCTPSndRcvInfoWrappedConn(conn.(*sctp.SCTPConn))

			go serveClient(wconn, *bufsize)
		}

	} else {
		var laddr *sctp.SCTPAddr
		if *lport != 0 {
			laddr = &sctp.SCTPAddr{
				Port: *lport,
			}
		}
		conn, err := sctp.DialSCTP("sctp", laddr, addr)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		log.Printf("Dail LocalAddr: %s; RemoteAddr: %s", conn.LocalAddr(), conn.RemoteAddr())

		ppid := 0
		for {
			info := &sctp.SndRcvInfo{
				Stream: uint16(ppid),
				PPID:   uint32(ppid),
			}
			ppid += 1
			conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
			buf := make([]byte, *bufsize)
			n, err := rand.Read(buf)
			if n != *bufsize {
				log.Fatalf("failed to generate random string len: %d", *bufsize)
			}
			n, err = conn.SCTPWrite(buf, info)
			if err != nil {
				log.Fatalf("failed to write: %v", err)
			}
			log.Printf("write: len %d", n)
			n, info, err = conn.SCTPRead(buf)
			if err != nil {
				log.Fatalf("failed to read: %v", err)
			}
			log.Printf("read: len %d, info: %+v", n, info)
			time.Sleep(time.Second)
		}
	}
}
