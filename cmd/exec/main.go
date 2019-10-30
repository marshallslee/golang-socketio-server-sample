package main

import (
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"os"
)

var (
	router       = gin.Default()
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

func hello(c *gin.Context) {
	log.Println("Hello from HTTP request.")
}

func test(c *gin.Context) {
	log.Println("Testing..")
}

func person(c *gin.Context) {
	name := c.Param("name")
	log.Printf("My name is %s\n", name)
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
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
}

func main() {
	defer logFile.Close()

	router.GET("/hello", hello)
	router.GET("/test", test)
	router.GET("/person/:name", person)
	router.GET("/", gin.WrapF(socketHandler))
	router.Run(":12379")

	log.Println("Serving at localhost:12379...")
}
