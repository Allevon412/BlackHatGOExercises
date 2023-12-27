package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"html/template"
	"net/http"
	"os"
)

type auth struct {
	username string
	password string
}

type trivial struct {
}

var x_template = `
<html>
	<body>
		Hello {{.}}
	</body>
</html>`

func (t *trivial) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("Executing trivial middleware")
	next(w, r)
}

func (a *auth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	RequestUsername := r.URL.Query().Get("username")
	RequestPassword := r.URL.Query().Get("password")
	if RequestUsername != a.username || RequestPassword != a.password {
		http.Error(w, "Unauthorized", 401)
		return
	}
	ctx := context.WithValue(r.Context(), "username", RequestUsername)
	r = r.WithContext(ctx)
	next(w, r)
}

func hello(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	fmt.Fprintf(w, "hi %s", username)
}

func main() {
	t, err := template.New("hello").Parse(x_template)
	if err != nil {
		panic(err)
	}
	t.Execute(os.Stdout, "<script>alert('world')</script>")
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET")
	n := negroni.Classic()
	n.Use(&auth{username: "admin", password: "password"})
	n.UseHandler(r)
	http.ListenAndServe(":8000", n)
}
