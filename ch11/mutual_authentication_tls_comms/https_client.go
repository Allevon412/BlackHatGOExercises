package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var (
		err              error
		cert             tls.Certificate
		serverCert, body []byte
		pool             *x509.CertPool
		tlsConf          *tls.Config
		transport        *http.Transport
		client           *http.Client
		resp             *http.Response
	)

	if cert, err = tls.LoadX509KeyPair("clientCrt.pem", "clientKey.pem"); err != nil { // generated client cert using $ openssl req -nodes -x509 -newkey rsa:4096 -keyout clientKey.pem -out clientCrt.pem -days 365
		log.Fatalln(err)
	}

	if serverCert, err = ioutil.ReadFile("../server/serverCrt.pem"); err != nil { // generated server cert using $ openssl req -nodes -x509 -newkey rsa:4096 -keyout serverKey.pem -out serverCrt.pem -days 365
		log.Fatalln(err)
	}

	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(serverCert)

	tlsConf = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	tlsConf.BuildNameToCertificate()

	transport = &http.Transport{
		TLSClientConfig: tlsConf,
	}
	client = &http.Client{
		Transport: transport,
	}

	if resp, err = client.Get("https://server.blackhat-go.local:9443/hello"); err != nil {
		log.Fatalln(err)
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	fmt.Printf("Success: %s\n", body)
}
