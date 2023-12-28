package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"os"
	"text/tabwriter"
)

type result struct {
	IPAddress string
	Hostname  string
}

type empty struct{}

func main() {
	var (
		flDomain      = flag.String("d", "", "The domain to perform guessing against")
		flWordlist    = flag.String("w", "", "The wordlist to use for guessing")
		flWorkerCount = flag.Int("c", 100, "The amount of workers to use.")
		flServerAddr  = flag.String("s", "8.8.8.8:53", "The DNS server to use")
	)
	flag.Parse()

	if *flDomain == "" || *flWordlist == "" {
		fmt.Println("-domain & -wordlist are required.")
		os.Exit(1)
	}
	fmt.Println(*flWorkerCount, *flServerAddr)

	var results []result
	fqdns := make(chan string, *flWorkerCount)
	gather := make(chan []result)
	tracker := make(chan empty)

	hFile, err := os.Open(*flWordlist)
	if err != nil {
		panic(err)
	}
	defer hFile.Close()
	scanner := bufio.NewScanner(hFile)

	//start workers
	for i := 0; i < *flWorkerCount; i++ {
		go worker(tracker, fqdns, gather, *flServerAddr)
	}

	// do work
	for scanner.Scan() {
		fqdns <- fmt.Sprintf("%s.%s", scanner.Text(), *flDomain)
	}

	// gather data
	go func() {
		for r := range gather {
			results = append(results, r...)
		}
		var e empty
		tracker <- e
	}()

	// close worker channels
	close(fqdns)
	for i := 0; i < *flWorkerCount; i++ {
		<-tracker
	}
	close(gather)
	<-tracker

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\n", r.Hostname, r.IPAddress)
	}
	w.Flush()
}

func lookupCNAME(fqdn, serverAddr string) ([]string, error) {
	var m dns.Msg
	var fqdns []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	reply, err := dns.Exchange(&m, serverAddr)
	if err != nil {
		return fqdns, err
	}
	if len(reply.Answer) < 1 {
		return fqdns, errors.New("No Answer")
	}
	for _, answer := range reply.Answer {
		if c, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, c.Target)
		}
	}
	return fqdns, nil
}

func lookupA(fqdn, serverAddr string) ([]string, error) {
	var m dns.Msg
	var ips []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	reply, err := dns.Exchange(&m, serverAddr)

	if err != nil {
		return ips, err
	}

	if len(reply.Answer) < 1 {
		return ips, errors.New("No Answer")
	}

	for _, answer := range reply.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips, nil
}

func lookup(fqdn, serverAddr string) []result {
	var results []result
	var fqdn_copy = fqdn // don't modify the original
	for {
		cnames, err := lookupCNAME(fqdn_copy, serverAddr)
		if err == nil && len(cnames) > 0 {
			fqdn_copy = cnames[0]
			continue
		}
		ips, err := lookupA(fqdn_copy, serverAddr)
		if err != nil {
			break
		}
		for _, ip := range ips {
			results = append(results, result{IPAddress: ip, Hostname: fqdn})
		}
		break
	}
	return results
}

func worker(tracker chan empty, fqdns chan string, gather chan []result, serverAddr string) {
	for fqdn := range fqdns {
		results := lookup(fqdn, serverAddr)
		if len(results) > 0 {
			gather <- results
		}
	}
	var e empty
	tracker <- e
}
