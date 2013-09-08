package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	conn, e := redis.Dial("tcp", ":6379")
	if e != nil {
		fmt.Println("Failed to connect to redis")
	}
	fmt.Println("Connected to redis")
	conn.Do("SET", "foo", "bar")
	conn.Do("SET", "baz", "qux")
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
	urlParam := strings.TrimPrefix(req.URL.Path, "/")
	if urlParam != "" {
		fmt.Println(urlParam)
		conn, e := redis.Dial("tcp", ":6379")
		if e != nil {
			fmt.Println("Failed to connect to redis")
		}
		result, err := redis.String(conn.Do("GET", urlParam))
		if err != nil {
			fmt.Println("Doesn't not exist")
		}
		fmt.Printf("%#v\n", result)
		http.Redirect(res, req, result, 302)
		fmt.Println("Let's redirect")
	} else {
		t, _ := template.ParseFiles("index.html")
		t.Execute(res, nil)
	}
}

func shorten(res http.ResponseWriter, req *http.Request) {
	conn, e := redis.Dial("tcp", ":6379")
	if e != nil {
		fmt.Println("Failed to connect to redis")
	}
	url := req.FormValue("url")
	hash, _ := generateUniqueHash(url)
	conn.Do("SET", hash, url)
	http.Redirect(res, req, "/", http.StatusFound)
}

func generateUniqueHash(url string) (string, error) {
	h := md5.New()
	io.WriteString(h, url)
	byteArray := h.Sum(nil)
	fmt.Println("%x", byteArray)
	hash := hex.EncodeToString(byteArray)[0:6]
	fmt.Println(hash)
	return hash, nil
}
