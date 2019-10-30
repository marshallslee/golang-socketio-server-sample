package main

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"os"
)

type person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello from HTTP request.")
}

func name(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if err := r.ParseForm(); err != nil {
			log.Println(err.Error())
		}
		name := r.Form.Get("name")

		fmt.Fprintf(w, "My name is %s\n", name)
		log.Printf("My name is %s\n", name)
	}
}

func main() {
	// Start logging
	logPath := "/socketiosample/output.log"
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	log.Println("Start logging..")
	defer f.Close()
	log.SetOutput(f)

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/socketio", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("Connected:", s.ID())
		return nil
	})

	server.OnEvent("/socketio", "name_event", func(s socketio.Conn, p person) {
		log.Println("First name:", p.FirstName)
		log.Println("Last name:", p.LastName)
	})

	server.OnError("/socketio", func(e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/socketio", func(s socketio.Conn, msg string) {
		log.Println("closed", msg)
	})

	go server.Serve()
	defer server.Close()

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/person/{name}", name)
	http.Handle("/", server)

	log.Println("Serving at localhost:12379...")
	log.Fatal(http.ListenAndServe(":12379", nil))
}
