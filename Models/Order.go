package Models

import (
	"TradeSimulator/Models/Enum"
	"time"
)

type Order struct {
	StockID   string         `json:"stockId"`   // 訂單ID
	OrderType Enum.OrderType `json:"orderType"` // 訂單類型 (0:Buy, 1:Sell)
	Price     float64        `json:"price"`     // 下單價格
	Quantity  int            `json:"quantity"`  // 下單數量
	Timestamp time.Time      `json:"-"`         // 下單時間
}
