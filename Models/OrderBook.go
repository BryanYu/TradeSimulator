package Models

import (
	"TradeSimulator/Models/Enum"
	"container/heap"
	"fmt"
	"time"
)

type OrderBook struct {
	buyOrders  PriorityQueue
	sellOrders PriorityQueue
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		buyOrders:  make(PriorityQueue, 0),
		sellOrders: make(PriorityQueue, 0),
	}
}

func (orderBook *OrderBook) AddOrder(order *Order) {
	switch order.OrderType {
	case Enum.Buy:
		heap.Push(&orderBook.buyOrders, order)
	case Enum.Sell:
		heap.Push(&orderBook.sellOrders, order)
	}
}

func (orderBook *OrderBook) MatchOrders() {
	for {
		fmt.Printf("WaitTrading...Time:%s \n", time.Now().Format("2006-04-02 15:04:05"))
		for orderBook.buyOrders.Len() > 0 && orderBook.sellOrders.Len() > 0 {
			fmt.Printf("TradeSimulator Start...Time:%s \n", time.Now().Format("2006-04-02 15:04:05"))
			buyOrder := orderBook.buyOrders[0]
			sellOrder := orderBook.sellOrders[0]

			if buyOrder.Price >= sellOrder.Price {
				quantity := min(buyOrder.Quantity, sellOrder.Quantity)
				buyOrder.Quantity -= quantity
				sellOrder.Quantity -= quantity
				fmt.Printf("Matched Order: Buy Order %d with Sell Order %d for %d units at price %.2f\n",
					buyOrder.ID, sellOrder.ID, quantity, sellOrder.Price)
				if buyOrder.Quantity == 0 {
					heap.Pop(&orderBook.buyOrders)
				}
				if sellOrder.Quantity == 0 {
					heap.Pop(&orderBook.sellOrders)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
