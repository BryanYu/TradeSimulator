package Models

type BetterFivePriceResponse struct {
	Buy  []BetterFivePrice `json:"buy"`
	Sell []BetterFivePrice `json:"sell"`
}

type BetterFivePrice struct {
	Price         float64 `json:"price"`
	TotalQuantity int     `json:"totalQuantity"`
}
