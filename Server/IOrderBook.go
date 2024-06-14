package Server

import "TradeSimulator/Models"

type IOrderBook interface {
	AddOrder(order *Models.Order)
	MatchOrders()
}
