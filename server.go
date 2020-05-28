package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
	utils "gitlab.com/stockboi/marco-polo/lib"
)

func main() {
	go playWithHashMaps()
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket/", server)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Serving at localhost:3000...")
	log.Println("Check Home Page at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func playWithHashMaps() {
	outer := make(map[string]*utils.HashMap)
	outer["api1"] = utils.NewHashMapDefault()
	outer["api2"] = utils.NewHashMapDefault()

	outer["api1"].Insert(`{"id":"10.0.0.1", "lat":12.4, "long":13.5, "tags":{"a":"b"}}`)
	outer["api2"].Insert(`{"id":"10.0.0.2", "lat":12.4, "long":13.5, "tags":{"a":"b"}}`)

	fmt.Println(outer["api1"].FlushMap())
	time.Sleep(10 * time.Second)
	fmt.Println(outer["api2"].FlushMap())
}
