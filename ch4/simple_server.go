package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s\n", r.URL.Query().Get("name"))
}

type router struct{ Inner http.Handler }

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	switch req.URL.Path {
	case "/a":
		fmt.Fprintf(w, "Executing /a")
	case "/b":
		fmt.Fprintf(w, "Executing /b")
	case "/c":
		fmt.Fprintf(w, "Executing /c")
	case "/hello":
		r.Inner.ServeHTTP(w, req)
	default:
		http.Redirect(w, req, "https://google.com", 302)
		//http.Error(w, "404 Not Found", 404)
	}

}

func main() {
	var r router
	f := http.HandlerFunc(hello)
	r.Inner = f
	http.ListenAndServe(":8000", &r)
}
