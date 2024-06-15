package Models

type TradeLog struct {
	StockId    string  `json:"stockId"`
	BuyPrice   float64 `json:"buyPrice"`
	SellPrice  float64 `json:"sellPrice"`
	TradePrice float64 `json:"tradePrice"`
	Quantity   int     `json:"quantity"`
	TimeStamp  int64   `json:"timeStamp"`
}
