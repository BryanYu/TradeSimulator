package Server

import (
	"TradeSimulator/Models"
	"TradeSimulator/Models/Enum"
	"container/heap"
	"fmt"
	"time"
)

type OrderBook struct {
	buyOrders  Models.PriorityQueue
	sellOrders Models.PriorityQueue
	sender     *TradeLogSender
}

func NewOrderBook(sender *TradeLogSender) IOrderBook {
	return &OrderBook{
		buyOrders:  make(Models.PriorityQueue, 0),
		sellOrders: make(Models.PriorityQueue, 0),
		sender:     sender,
	}
}

func (orderBook *OrderBook) AddOrder(order *Models.Order) {
	switch order.OrderType {
	case Enum.Buy:
		heap.Push(&orderBook.buyOrders, order)
	case Enum.Sell:
		heap.Push(&orderBook.sellOrders, order)
	}
}

func (orderBook *OrderBook) MatchOrders() {
	time.Sleep(10 * time.Second)
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
				orderBook.sender.Send(Enum.TradeLog, Models.TradeLog{
					StockId:    "Stock1",
					BuyPrice:   buyOrder.Price,
					SellPrice:  sellOrder.Price,
					TradePrice: buyOrder.Price,
					Quantity:   quantity,
					TimeStamp:  time.Now().Unix(),
				})
				if buyOrder.Quantity == 0 {
					heap.Pop(&orderBook.buyOrders)
				}
				if sellOrder.Quantity == 0 {
					heap.Pop(&orderBook.sellOrders)
				}
			}
			time.Sleep(10 * time.Second)
		}
		time.Sleep(10 * time.Second)
	}
}
