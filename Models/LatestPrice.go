package Models

type LatestPrice struct {
	StockID            string  `json:"stockId"`
	TradePrice         float64 `json:"tradePrice"`
	Margin             float64 `json:"margin"`
	TotalTradeQuantity int     `json:"totalTradeQuantity"`
}
