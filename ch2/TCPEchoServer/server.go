package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func echo2(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	s, err := reader.ReadString(0x0d)
	if err != nil {
		log.Fatalln("Unable to read data")
	}
	log.Printf("Read %d bytes: %s", len(s), s)

	log.Printf("Writing Data")
	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString(s); err != nil {
		log.Fatalln("unable to write data")
	}
	writer.Flush()
}

func echo3(conn net.Conn) {
	defer conn.Close()

	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatalln("unable to read/write data")
	}
}

func echo(conn net.Conn) {
	defer conn.Close()

	b := make([]byte, 512)
	for {
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("Client Disconnected")
			break
		}
		if err != nil {
			log.Println("Unexpected Error")
			break
		}
		log.Printf("Received %d bytes: %s\n", size, string(b))

		log.Println("writing Data")
		if _, err := conn.Write(b[0:size]); err != nil {
			log.Fatalln("unable to write data")
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":20080")

	if err != nil {
		log.Fatalln("Unable to bing to port")
	}

	log.Println("Listening on 0.0.0.0:20080")
	for {
		conn, err := listener.Accept()
		log.Println("Received Connection")
		if err != nil {
			log.Fatalln("unable to accept connection")
		}
		go echo2(conn)
	}

}
