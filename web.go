package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	/* REDISTOGO_URL=redis://redistogo:843a6dc681ee128046391c888529f6f1@koi.redistogo.com:9934/ */
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
	urlParam := strings.TrimPrefix(req.URL.Path, "/")
	if urlParam != "" {
		fmt.Println(urlParam)
		conn, _ := getRedisConn()
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
	conn, _ := getRedisConn()
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

func getRedisConn() (redis.Conn, error) {
	connectionString := os.Getenv("REDISTOGO_URL")
	var conn redis.Conn
	var err error
	if connectionString != "" {
		conn, err = redisurl.ConnectToURL(connectionString)
	} else {
		conn, err = redis.Dial("tcp", ":6379")
	}
	if err != nil {
		fmt.Println("Failed to connect to redis")
	}
	fmt.Println("Connected to redis")
	return conn, err
}
