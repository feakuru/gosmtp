package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	log.Fatal(listenTCP())
}

func listenTCP() error {
	var listenAddr string
	flag.StringVar(&listenAddr, "listen-addr", "0.0.0.0:25", "address to listen")
	flag.Parse()

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("error listening on %q: %s", listenAddr, err)
	}
	defer func() {
		closeErr := l.Close()
		if closeErr != nil {
			log.Printf("can't close listen socket: %s", err)
		}
	}()

	fmt.Printf("Listening on %q\n", listenAddr)

	for {
		conn, err := l.Accept()
		defer conn.Close()
		log.Printf("Connection from %s...", conn.RemoteAddr())

		if err != nil {
			fmt.Printf("Error accepting connection %q: %s", listenAddr, err)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.Printf("...successful (%s)", conn.RemoteAddr())
		}

		if err := handleRequest(conn); err != nil {
			log.Printf("Stopped handling requests from %s: %s", conn.RemoteAddr(), err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
	}
}

func handleRequest(conn net.Conn) error {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	for ; err == nil; reqLen, err = conn.Read(buf) {
		log.Printf("Got request of len %d bytes from %s: %q", reqLen, conn.RemoteAddr(), buf[:reqLen])
		if _, errConn := conn.Write(dispatchResponse(buf[:reqLen])); errConn != nil {
			return fmt.Errorf("can't write to connection: %s", errConn)
		}
	}
	if err != nil {
		return fmt.Errorf("can't read data from connection: %s", err)
	}
	return nil
}

func dispatchResponse(command []byte) []byte {
	return []byte("Message received.\n")
}
