package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

func worker(ports chan int, wg *sync.WaitGroup) {
	for p := range ports {
		go func(port int) {
			defer wg.Done()
			address := fmt.Sprintf("scanme.nmap.org:%d", port)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			fmt.Printf("%d open\n", port)
			conn.Close()
		}(p)
	}
}

func sorted_worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	_, err := net.Dial("tcp", "scanme.namp.org:80")
	if err == nil {
		fmt.Println("connection successful")
	}
	var wg sync.WaitGroup

	for i := 1; i < 1024; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			address := fmt.Sprintf("scanme.nmap.org:%d", j)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				//port is closed or filtered
				return
			}
			conn.Close()
			fmt.Printf("%d open\n", j)
		}(i)
	}
	wg.Wait()

	//using work groups for more consistency with concurrency
	ports := make(chan int, 100) // create 100 work groups
	var wg2 sync.WaitGroup
	for i := 0; i < cap(ports); i++ {
		go worker(ports, &wg2)
	}
	for i := 1; i < 1024; i++ {
		wg2.Add(1)
		ports <- i
	}
	wg2.Wait()
	close(ports)

	//using work groups & sorting the results so we have a clearer result list & no longer require the wait sync group import
	ports2 := make(chan int, 100)
	results := make(chan int)
	var openports []int
	for i := 0; i < cap(ports2); i++ {
		go sorted_worker(ports2, results)
	}
	go func() {
		for i := 1; i <= 1024; i++ {
			ports2 <- i
		}
	}()
	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports2)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
