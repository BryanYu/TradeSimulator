package Enum

type SocketEvent string

const (
	Join           SocketEvent = "join"           //加入報價頻道
	Order          SocketEvent = "order"          //下單
	GetLatestPrice SocketEvent = "getLatestPrice" //取得最新報價
	GetTradeLogs   SocketEvent = "getTradeLogs"   //取得所有交易紀錄
	GetTradeLog    SocketEvent = "getTradeLog"    //取得當次成交交易紀錄
)
