package Server

import "TradeSimulator/Models"

type IOrderBook interface {
	AddOrder(order *Models.Order)
	MatchOrders()
	GetLatestPrice() *Models.LatestPrice
	GetTradeLogs() []Models.TradeLog
}
