package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"html/template"
	"net/http"
	"os"
	/* "strings" */
)

var con redis.Conn = nil

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/shorten/", shorten)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("listening...")

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	// TODO : Redirect to the URL form the request after trying to find it
	/* fmt.Println(strings.TrimPrefix(req.URL.Path, "/")) */
	// in Redis
	t, _ := template.ParseFiles("index.html")
	t.Execute(res, nil)
}

func shorten(res http.ResponseWriter, req *http.Request) {

}
