package Models

import (
	"TradeSimulator/Models/Enum"
	"time"
)

type Order struct {
	ID        int            // 訂單ID
	OrderType Enum.OrderType //訂單類型 (0:Buy, 1:Sell)
	Price     float64        // 下單價格
	Quantity  int            // 下單數量
	Timestamp time.Time      // 下單時間
}
