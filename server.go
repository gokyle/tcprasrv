package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

var addrRegex = regexp.MustCompile("^(.+):\\d{1,5}$")

func remoteAddr(conn net.Conn) string {
	addr := conn.RemoteAddr().String()
	return addrRegex.ReplaceAllString(addr, "$1")
}

func serve(address string) {
	srv, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("failed to set up server:", err.Error())
	}

	log.Println("listening on", address)
	for {
		conn, err := srv.Accept()
		if err != nil {
			log.Println("error accepting connection:", err.Error())
			continue
		}
		go func(conn net.Conn) {
			conn.Write([]byte(remoteAddr(conn)))
			conn.Close()
		}(conn)
	}
}

func main() {
	var port string

	if port = os.Getenv("PORT"); port == "" {
		port = "4321"
	}
	go serve(":" + port)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
	<-sigc
	log.Println("shutting down.")
}
