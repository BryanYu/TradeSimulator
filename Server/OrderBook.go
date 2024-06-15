package Server

import (
	"TradeSimulator/Models"
	"TradeSimulator/Models/Enum"
	"container/heap"
	"fmt"
	"time"
)

type OrderBook struct {
	buyOrders   Models.PriorityQueue
	sellOrders  Models.PriorityQueue
	sender      *MessageSender
	latestPrice *Models.LatestPrice
	tradeLogs   []Models.TradeLog
}

func NewOrderBook(sender *MessageSender) IOrderBook {
	return &OrderBook{
		buyOrders:  make(Models.PriorityQueue, 0),
		sellOrders: make(Models.PriorityQueue, 0),
		sender:     sender,
		tradeLogs:  make([]Models.TradeLog, 0),
		latestPrice: &Models.LatestPrice{
			StockID:            "StockId1",
			TradePrice:         100,
			Margin:             0,
			TotalTradeQuantity: 0,
		},
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

				tradeLog := Models.TradeLog{
					StockId:    "Stock1",
					BuyPrice:   buyOrder.Price,
					SellPrice:  sellOrder.Price,
					TradePrice: sellOrder.Price,
					Quantity:   quantity,
					TimeStamp:  time.Now().Unix(),
				}

				orderBook.tradeLogs = append(orderBook.tradeLogs, tradeLog)
				orderBook.setLatestPrice("Stock1", sellOrder.Price, quantity)

				// 推送最新成交價
				latestPrice := orderBook.GetLatestPrice()
				orderBook.sender.Send(Enum.Price, Enum.GetLatestPrice, latestPrice)

				// 推送交易log
				orderBook.sender.Send(Enum.TradeLog, Enum.GetTradeLog, tradeLog)

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

func (orderBook *OrderBook) setLatestPrice(stockId string, tradePrice float64, quantity int) {
	if orderBook.latestPrice == nil {
		orderBook.latestPrice = &Models.LatestPrice{
			StockID:            stockId,
			TradePrice:         tradePrice,
			Margin:             0,
			TotalTradeQuantity: quantity}
	} else {
		orderBook.latestPrice.TotalTradeQuantity += quantity
		orderBook.latestPrice.Margin = tradePrice - orderBook.latestPrice.TradePrice
		orderBook.latestPrice.TradePrice = tradePrice
	}
}

func (orderBook *OrderBook) GetLatestPrice() *Models.LatestPrice {
	return orderBook.latestPrice
}

func (orderBook *OrderBook) GetTradeLogs() []Models.TradeLog {
	return orderBook.tradeLogs
}
