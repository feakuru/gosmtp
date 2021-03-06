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
	"bytes"
	"github.com/feakuru/gosmtp/confreaders"
	"github.com/feakuru/gosmtp/cmddispatch"
	"github.com/feakuru/gosmtp/mailstorage"
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

	mailstorage.Init()
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
	var currentCommand cmddispatch.StoredCommand
	var msg []byte
	reqLen, err := conn.Read(buf)
	dataRead := false
	data := []byte("")
	for ; err == nil; reqLen, err = conn.Read(buf) {
		log.Printf("Got request of len %d bytes from %s: %q", reqLen, conn.RemoteAddr(), buf[:reqLen])
		if dataRead {
			if bytes.Equal(buf[:reqLen], []byte(".\r\n")) {
				log.Println("\nsender: ", currentCommand.StrdSender, "\nrcpts: ", currentCommand.StrdRcpts, "\ndata: ", string(data))
				mailstorage.PutEmailToDB(currentCommand.StrdSender, currentCommand.StrdRcpts, string(data))
				dataRead = false
				if _, errConn := conn.Write([]byte("200 OK\r\n")); errConn != nil {
					return fmt.Errorf("can't write to connection: %s", errConn)
				}
			}
			for i := 0; i < reqLen; i++ {
				data = append(data, buf[i])
			}
		} else {
			msg, currentCommand = dispatchResponse(buf[:reqLen], currentCommand)
			if strings.Split(string(msg), " ")[0] == "354" {
				dataRead = true
			}
			if _, errConn := conn.Write(msg); errConn != nil {
				return fmt.Errorf("can't write to connection: %s", errConn)
			}
			log.Println("\nsender: ", currentCommand.StrdSender, "\nrcpts: ", currentCommand.StrdRcpts, "\ndata: ", string(data))
		}
		buf = make([]byte, 1024)
	}
	if err != nil {
		return fmt.Errorf("can't read data from connection: %s", err)
	}
	return nil
}

func dispatchResponse(command []byte, prevCommand cmddispatch.StoredCommand) ([]byte, cmddispatch.StoredCommand) {
	cmdBytes := bytes.Split(bytes.Trim(command, "\r\n"), []byte(":"))
	cmd := cmdBytes[0]
	var arg []byte
	if len(cmdBytes) > 1 {
		arg = cmdBytes[1]
	} else {
		arg = []byte("")
	}

	var currentCommand cmddispatch.StoredCommand
	var msg string

	currentCommand, msg = cmddispatch.Command(cmd, arg, prevCommand)
	return []byte(msg), currentCommand
}
