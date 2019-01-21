package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var env = flag.String("env", "prod", "Set environment")
var port = flag.String("port", ":8080", "http service address")
var key = flag.String("key", "server.key", "ssl key")
var crt = flag.String("crt", "server.crt", "ssl crt")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func checkOrigin(r *http.Request) bool {
	log.Println(r)
	// TODO: Check origin
	return true
}

func main() {
	flag.Parse()

	state := createState()

	go state.start()

	http.HandleFunc("/", index)
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		wsHandler(writer, request, state)
	})

	var err error
	if *env == "prod" {
		err = http.ListenAndServeTLS(*port, *crt, *key, nil)
	} else {
		err = http.ListenAndServe(*port, nil)
	}

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
