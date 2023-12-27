package main

import (
	"fmt"
	"log"
	"os"
	"shodan/shodan"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("usage: [APIKEY] [SEARCH_TERM]")
	}
	//apiKey := os.Getenv("SHODAN_API_KEY")
	apiKey := os.Args[1]
	s := shodan.New(apiKey)
	info, err := s.APIInfo()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("Query Credits: %d\nScanCredits: %d\n\n", info.QueryCredits, info.ScanCredits)

	hostSearch, err := s.HostSearch(os.Args[2])
	if err != nil {
		log.Panicln(err)
	}
	for _, host := range hostSearch.Matches {
		fmt.Printf("%18s%18d\n", host.IPString, host.Port)
	}
}
