package main

import (
	"TradeSimulator/Models"
	"TradeSimulator/Models/Enum"
	"fmt"
	"time"
)

func main() {
	orderBook := Models.NewOrderBook()

	order1 := &Models.Order{ID: 1, OrderType: Enum.Buy, Price: 100.5, Quantity: 10, Timestamp: time.Now()}
	order2 := &Models.Order{ID: 2, OrderType: Enum.Sell, Price: 100.0, Quantity: 5, Timestamp: time.Now()}
	order3 := &Models.Order{ID: 3, OrderType: Enum.Buy, Price: 101.0, Quantity: 7, Timestamp: time.Now()}
	order4 := &Models.Order{ID: 4, OrderType: Enum.Sell, Price: 101.0, Quantity: 10, Timestamp: time.Now()}
	order5 := &Models.Order{ID: 5, OrderType: Enum.Buy, Price: 102.0, Quantity: 15, Timestamp: time.Now()}

	orderBook.AddOrder(order1)
	orderBook.AddOrder(order2)
	orderBook.AddOrder(order3)
	orderBook.AddOrder(order4)
	orderBook.AddOrder(order5)

	go orderBook.MatchOrders()
	for {
		var input string
		fmt.Scanln(&input)
	}
}
