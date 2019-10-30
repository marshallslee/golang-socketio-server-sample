package main

import (
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	logFile      *os.File
	socketServer *socketio.Server
)

func init() {
	// Start logging
	logPath := "/socketiosample/output.log"
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	log.Println("Start logging..")
	log.SetOutput(logFile)

	socketServer, err = socketio.NewServer(nil)
	if err != nil {
		log.Println(err.Error())
	}
}

type PersonInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello from HTTP request.")
}

func test(w http.ResponseWriter, r *http.Request) {
	log.Println("Testing..")
}

func person(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	name := param["name"]
	log.Printf("My name is %s\n", name)
}

func main() {
	defer logFile.Close()

	socketServer.OnConnect("/socketio", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("Connected:", s.ID())
		return nil
	})

	socketServer.OnEvent("/socketio", "name_event", func(s socketio.Conn, p PersonInfo) {
		log.Println("First name:", p.FirstName)
		log.Println("Last name:", p.LastName)
	})

	socketServer.OnError("/socketio", func(e error) {
		log.Println("meet error:", e)
	})

	socketServer.OnDisconnect("/socketio", func(s socketio.Conn, msg string) {
		log.Println("closed", msg)
	})

	go socketServer.Serve()
	defer socketServer.Close()

	router := mux.NewRouter()
	router.HandleFunc("/hello", hello)
	router.HandleFunc("/test", test)
	router.HandleFunc("/person/{name}", person)
	router.Handle("/", socketServer)

	log.Println("Serving at localhost:12379...")
	log.Fatal(http.ListenAndServe(":12379", router))
}
