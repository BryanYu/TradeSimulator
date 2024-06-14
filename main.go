package main

import (
	"TradeSimulator/Server"
	"log"
	"net/http"
)

func main() {
	socketServer := Server.NewSocketServer()
	socketServer.InitialServer()
	socketServer.RegisterEvent("/")
	server := socketServer.GetServer()
	go server.Serve()
	defer server.Close()
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
