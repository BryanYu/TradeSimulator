package main

import (
	"TradeSimulator/Server"
	"log"
	"net/http"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func main() {
	socketServer := Server.NewSocketServer()
	tradeLogSender := Server.NewTradeLogSender(socketServer)
	orderBook := Server.NewOrderBook(tradeLogSender)

	socketServer.InitialServer(orderBook)
	socketServer.RegisterEvent("/")
	server := socketServer.GetServer()

	go server.Serve()
	defer server.Close()

	go orderBook.MatchOrders()
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
