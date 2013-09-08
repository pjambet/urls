package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

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
	urlParam := strings.TrimPrefix(req.URL.Path, "/")
	if urlParam != "" {
		conn, _ := getRedisConn()
		result, err := redis.String(conn.Do("GET", urlParam))
		if err != nil {
			fmt.Println("Doesn't not exist")
		}
		http.Redirect(res, req, result, 302)
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
	res.Header().Set("Content-Type", "application/json")
	fullURL := hash
	fmt.Fprint(res, Response{"success": true, "url": req.Host + "/" + fullURL})
	http.Redirect(res, req, "/", http.StatusFound)
}

func generateUniqueHash(url string) (string, error) {
	h := md5.New()
	io.WriteString(h, url)
	byteArray := h.Sum(nil)
	hash := hex.EncodeToString(byteArray)[0:6]
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
