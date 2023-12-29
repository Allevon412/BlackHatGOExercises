package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func parse(filename string) (map[string]string, error) {
	records := make(map[string]string)
	hFile, err := os.Open(filename)
	if err != nil {
		return records, err
	}
	defer hFile.Close()
	scanner := bufio.NewScanner(hFile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			return records, fmt.Errorf("%s is not a valid line", line)
		}
		records[parts[0]] = parts[1]
	}
	return records, scanner.Err()
}

func main() {
	var RecordLock sync.RWMutex

	records, err := parse("C:\\Users\\Brendan Ortiz\\Documents\\GOProjcets\\BHGO\\ch5\\DNS_Proxy\\proxy.config")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", records)

	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		if len(req.Question) < 1 {
			m := new(dns.Msg)
			m.SetRcode(req, dns.RcodeServerFailure)
			w.WriteMsg(m)
			return
		}

		name := req.Question[0].Name
		parts := strings.Split(name, ".")
		if len(parts) > 1 {
			name = strings.Join(parts[:len(parts)-1], ".")
		}

		RecordLock.RLock()
		match, ok := records[name]
		RecordLock.RUnlock()

		if !ok {
			m := new(dns.Msg)
			m.SetRcode(req, dns.RcodeServerFailure)
			w.WriteMsg(m)
			return
		}
		resp, err := dns.Exchange(req, match)
		if err != nil {
			m := new(dns.Msg)
			m.SetRcode(req, dns.RcodeServerFailure)
			w.WriteMsg(m)
			return
		}

		if err := w.WriteMsg(resp); err != nil {
			m := new(dns.Msg)
			m.SetRcode(req, dns.RcodeServerFailure)
			w.WriteMsg(m)
			return
		}
	})

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.Signal(10))
		for sig := range sigs {
			switch sig {
			case syscall.Signal(10):
				log.Println("SIGUSR1: reloading records")
				RecordLock.Lock()
				parse("C:\\Users\\Brendan Ortiz\\Documents\\GOProjcets\\BHGO\\ch5\\DNS_Proxy\\proxy.config")
				RecordLock.Unlock()
			}
		}
	}()

	log.Fatal(dns.ListenAndServe(":53", "udp", nil))
}
