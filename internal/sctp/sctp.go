package main

import (
	"flag"
	"log"
	"net"
	"strings"
	"time"

	"github.com/ishidawataru/sctp"
)

func serveClient(conn *sctp.SCTPConn) error {
	buf := make([]byte, 1024) // add overhead of SCTPSndRcvInfoWrappedConn
	for {
		n, info, err := conn.SCTPRead(buf)
		if err != nil {
			log.Printf("read failed: %v", err)
			return err
		}
		log.Printf("received: %s, %+v", buf[:n], info)
		_, err = conn.SCTPWrite([]byte("vanya"), info)
		if err != nil {
			log.Printf("write failed: %v", err)
			return err
		}
	}
}

func main() {

	var server = flag.Bool("server", false, "")
	var ip = flag.String("ip", "0.0.0.0", "")
	var port = flag.Int("port", 0, "")

	flag.Parse()

	var ips []net.IPAddr

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
			conn, err := ln.AcceptSCTP()
			if err != nil {
				log.Fatalf("failed to accept: %v", err)
			}
			log.Printf("Accepted Connection from RemoteAddr: %s", conn.RemoteAddr())
			//wconn := sctp.NewSCTPSndRcvInfoWrappedConn(conn.(*sctp.SCTPConn))

			go serveClient(conn)
		}

	} else {
		var laddr *sctp.SCTPAddr
		conn, err := sctp.DialSCTP("sctp", laddr, addr)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}
		//wconn := sctp.NewSCTPSndRcvInfoWrappedConn(conn)

		log.Printf("Dial LocalAddr: %s; RemoteAddr: %s", conn.LocalAddr(), conn.RemoteAddr())

		ppid := 0
		buf := make([]byte, 1024)
		for {
			info := &sctp.SndRcvInfo{
				Stream: uint16(ppid),
				PPID:   uint32(ppid),
			}
			info = info
			ppid += 1
			conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)

			_, err = conn.Write([]byte("lera"))
			if err != nil {
				log.Fatalf("failed to write: %v", err)
			}
			n, info, err := conn.SCTPRead(buf)
			if err != nil {
				log.Fatalf("failed to read: %v", err)
			}
			log.Printf("received: %s, %+v", buf[:n], info)
			time.Sleep(5 * time.Second)
		}
	}
}
