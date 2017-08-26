package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"errors"
	"os"
	"time"
	"strings"
	"github.com/feakuru/gosmtp/confreaders"
	"github.com/feakuru/gosmtp/workers"
)

func main() {
	log.Fatal(listenTCP())
}

func listenTCP() error {
	var listenAddr string
	var workersNum string
	var confName string
	var conf map[string]string
	flag.StringVar(&listenAddr, "listen-addr", "0.0.0.0:25", "address to listen")
	flag.StringVar(&workersNum, "workers-number", "4", "workers quantity")
	flag.StringVar(&confName, "conf", "conf.json", "configuration file name")
	flag.Parse()

	if _, err := os.Stat(confName); os.IsNotExist(err) {
		log.Printf("no config found: \"%s\", continuing without it", confName)
	} else {
		if strings.Split(confName, ".")[1] == "json" {
			conf = confreaders.ReadJSONConfig(confName)
		} else if strings.Split(confName, ".")[1] == "yaml" {
			conf = confreaders.ReadYAMLConfig(confName)
		} else {
			return errors.New("can't read your config")
		}
		if val, ok := conf["listen-addr"]; ok {
			listenAddr = val;
		}
		if val, ok := conf["workers-number"]; ok {
			workersNum = val;
		}
	}


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

	workers.WorkerPool(4, func () {
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
	})
	return nil
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
