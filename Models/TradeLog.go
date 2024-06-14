package Models

type TradeLog struct {
	StockId    string
	BuyPrice   float64
	SellPrice  float64
	TradePrice float64
	Quantity   int
	TimeStamp  int64
}
