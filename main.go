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

	//order1 := &Models.Order{ID: "Stock1", OrderType: Enum.Buy, Price: 100.5, Quantity: 10, Timestamp: time.Now().Add(-time.Hour)}
	//order2 := &Models.Order{ID: "Stock1", OrderType: Enum.Buy, Price: 101.0, Quantity: 5, Timestamp: time.Now()}
	//order3 := &Models.Order{ID: "Stock1", OrderType: Enum.Sell, Price: 101.0, Quantity: 7, Timestamp: time.Now().Add(-time.Hour)}
	//order4 := &Models.Order{ID: "Stock1", OrderType: Enum.Sell, Price: 100.0, Quantity: 10, Timestamp: time.Now()}

	//orderBook.AddOrder(order1)
	//orderBook.AddOrder(order2)
	//orderBook.AddOrder(order3)
	//orderBook.AddOrder(order4)
	//
	go orderBook.MatchOrders()
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
