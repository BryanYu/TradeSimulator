package Server

import (
	"TradeSimulator/Models"
	"TradeSimulator/Models/Enum"
	"container/heap"
	"fmt"
	"sort"
	"time"
)

type OrderBook struct {
	buyOrders   Models.PriorityQueue
	sellOrders  Models.PriorityQueue
	sender      *MessageSender
	latestPrice *Models.LatestPrice
	tradeLogs   []Models.TradeLog
	openPrice   float64
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
		openPrice: 100,
	}
}

// AddOrder 下單
func (orderBook *OrderBook) AddOrder(order *Models.Order) {
	switch order.OrderType {
	case Enum.Buy:
		heap.Push(&orderBook.buyOrders, order)
	case Enum.Sell:
		heap.Push(&orderBook.sellOrders, order)
	}
	betterFivePrice := orderBook.GetBetterFivePrice()
	orderBook.sender.Send(Enum.BetterFivePrice, Enum.GetBetterFivePrice, betterFivePrice)
}

// MatchOrders 交易搓合
func (orderBook *OrderBook) MatchOrders() {
	fmt.Printf("WaitTrading...Time:%s \n", time.Now().Format("2006-04-02 15:04:05"))
	for {
		if !isOrdersExist(orderBook) {
			continue
		}
		for isOrderPriceMatch(orderBook.buyOrders[0], orderBook.sellOrders[0]) || isOrderHasMarketPrice(orderBook.buyOrders[0], orderBook.sellOrders[0]) {
			fmt.Printf("TradeSimulator Start...Time:%s \n", time.Now().Format("2006-04-02 15:04:05"))
			buyOrder := orderBook.buyOrders[0]
			sellOrder := orderBook.sellOrders[0]
			quantity := min(buyOrder.Quantity, sellOrder.Quantity)
			buyOrder.Quantity -= quantity
			sellOrder.Quantity -= quantity

			if buyOrder.Quantity == 0 {
				heap.Pop(&orderBook.buyOrders)
			}
			if sellOrder.Quantity == 0 {
				heap.Pop(&orderBook.sellOrders)
			}

			var tradePrice float64
			var buyPrice float64
			var sellPrice float64
			if buyOrder.IsMarketPrice && sellOrder.IsMarketPrice {
				latestPrice := orderBook.GetLatestPrice()
				tradePrice = latestPrice.TradePrice
				buyPrice = latestPrice.TradePrice
				sellPrice = latestPrice.TradePrice
			} else if buyOrder.IsMarketPrice {
				tradePrice = sellOrder.Price
				buyPrice = sellOrder.Price
				sellPrice = sellOrder.Price
			} else if sellOrder.IsMarketPrice {
				tradePrice = buyOrder.Price
				buyPrice = buyOrder.Price
				sellPrice = buyOrder.Price
			} else {
				tradePrice = sellOrder.Price
				buyPrice = buyOrder.Price
				sellPrice = sellOrder.Price
			}
			tradeLog := Models.TradeLog{
				StockId:    "Stock1",
				BuyPrice:   buyPrice,
				SellPrice:  sellPrice,
				TradePrice: tradePrice,
				Quantity:   quantity,
				TimeStamp:  time.Now().Unix(),
			}

			// 設定最新成交價
			orderBook.setLatestPrice("Stock1", tradePrice, quantity)

			// 交易紀錄
			orderBook.tradeLogs = append(orderBook.tradeLogs, tradeLog)
			// 推送單筆交易log
			orderBook.sender.Send(Enum.TradeLog, Enum.GetTradeLog, tradeLog)

			// 推送最佳買賣五檔報價
			betterFiveOrders := orderBook.GetBetterFivePrice()
			orderBook.sender.Send(Enum.BetterFivePrice, Enum.GetBetterFivePrice, betterFiveOrders)
			if !isOrdersExist(orderBook) {
				break
			}
		}
		// 推送最新成交價
		latestPrice := orderBook.GetLatestPrice()
		orderBook.sender.Send(Enum.Price, Enum.GetLatestPrice, latestPrice)
		time.Sleep(1 * time.Second)
	}

}

func isOrderHasMarketPrice(buyOrder *Models.Order, sellOrder *Models.Order) bool {
	return buyOrder.IsMarketPrice || sellOrder.IsMarketPrice
}

func isOrderPriceMatch(buyOrder *Models.Order, sellOrder *Models.Order) bool {
	return buyOrder.Price >= sellOrder.Price || buyOrder.IsMarketPrice || sellOrder.IsMarketPrice
}

func isOrdersExist(orderBook *OrderBook) bool {
	return orderBook.buyOrders.Len() > 0 && orderBook.sellOrders.Len() > 0
}

// setLatestPrice 設定最新價格資訊
func (orderBook *OrderBook) setLatestPrice(stockId string, tradePrice float64, quantity int) {
	if orderBook.latestPrice == nil {
		orderBook.latestPrice = &Models.LatestPrice{
			StockID:            stockId,
			TradePrice:         tradePrice,
			Margin:             0,
			TotalTradeQuantity: quantity}
	} else {
		orderBook.latestPrice.TotalTradeQuantity += quantity
		orderBook.latestPrice.Margin = tradePrice - orderBook.openPrice
		orderBook.latestPrice.TradePrice = tradePrice
	}
}

// GetLatestPrice 取得最新價格資訊
func (orderBook *OrderBook) GetLatestPrice() *Models.LatestPrice {
	return orderBook.latestPrice
}

// GetTradeLogs 取得交易紀錄
func (orderBook *OrderBook) GetTradeLogs() []Models.TradeLog {
	return orderBook.tradeLogs
}

// GetBetterFivePrice 取得買賣五檔最佳報價
func (orderBook *OrderBook) GetBetterFivePrice() Models.BetterFivePriceResponse {
	buyBetterOrders := getBetterFivePrice(orderBook.buyOrders, Enum.Buy)
	sellBetterOrders := getBetterFivePrice(orderBook.sellOrders, Enum.Sell)
	return Models.BetterFivePriceResponse{
		Buy:  buyBetterOrders,
		Sell: sellBetterOrders,
	}
}

func getBetterFivePrice(orderQueue Models.PriorityQueue, orderType Enum.OrderType) []Models.BetterFivePrice {
	priceGroups := make(map[float64]int)
	hasMarketPrice := make(map[bool]int)
	for _, order := range orderQueue {
		if order.IsMarketPrice {
			hasMarketPrice[order.IsMarketPrice] += order.Quantity
		} else {
			priceGroups[order.Price] += order.Quantity
		}
	}
	betterFivePrice := make([]Models.BetterFivePrice, 0)
	for price, totalQuantity := range priceGroups {
		betterFivePrice = append(betterFivePrice, Models.BetterFivePrice{
			Price:         price,
			TotalQuantity: totalQuantity,
			IsMarketPrice: false,
		})
	}

	if totalQuantity, find := hasMarketPrice[true]; find {
		betterFivePrice = append(betterFivePrice, Models.BetterFivePrice{
			Price:         0,
			TotalQuantity: totalQuantity,
			IsMarketPrice: true,
		})
	}

	if orderType == Enum.Buy {
		sort.Slice(betterFivePrice, func(i, j int) bool {
			if betterFivePrice[i].IsMarketPrice {
				return true
			} else if betterFivePrice[i].Price > betterFivePrice[j].Price {
				return true
			}
			return false
		})
	} else {
		sort.Slice(betterFivePrice, func(i, j int) bool {
			if betterFivePrice[i].IsMarketPrice {
				return true
			} else if betterFivePrice[i].Price < betterFivePrice[j].Price {
				return true
			}
			return false
		})
	}
	if len(betterFivePrice) >= 5 {
		return betterFivePrice[:5]
	} else {
		return betterFivePrice
	}
}
