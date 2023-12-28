package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
)

func main() {
	var msg dns.Msg
	fqdn := dns.Fqdn("stacktitan.com")
	msg.SetQuestion(fqdn, dns.TypeA)
	r, err := dns.Exchange(&msg, "8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}

	if len(r.Answer) < 1 {
		fmt.Println("No Records for query found")
		return
	}

	for _, answer := range r.Answer {
		if a, ok := answer.(*dns.A); ok {
			fmt.Println(a.A)
		}
	}

}
