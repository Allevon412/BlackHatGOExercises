package main

import (
	"fmt"
	"log"
	"os"
)

type CustomReader struct{}

func (CustomReader *CustomReader) Read(b []byte) (int, error) {
	fmt.Print("in > ")
	return os.Stdin.Read(b)
}

type CustomWriter struct{}

func (CustomWriter *CustomWriter) Write(b []byte) (int, error) {
	fmt.Print("Out > ")
	return os.Stdout.Write(b)
}

func main() {
	var (
		reader CustomReader
		writer CustomWriter
	)

	input := make([]byte, 4096)
	s, err := reader.Read(input)
	if err != nil {
		log.Fatalln("unable to read data")
	}
	fmt.Printf("read %d bytes from stdin\n", s)

	s, err = writer.Write(input)
	if err != nil {
		log.Fatalln("unable to write data")
	}
	fmt.Printf("Wrote %d bytes to stdout\n", s)
}
