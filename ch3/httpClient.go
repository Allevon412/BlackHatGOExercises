package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Status struct {
	Message string
	Status  string
}

func DecodeJsonResp(response http.Response) bool {
	var status Status
	if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	log.Printf("%s -> %s\n", status.Status, status.Message)
	return true
}

func main() {
	res, err := http.Get("https://www.google.com/robots.txt")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(res.Status)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(string(body))

}
