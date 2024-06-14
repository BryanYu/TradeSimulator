package Response

type GetLatestFivePriceResponse struct {
	BuyOrders  []GetLatestFivePriceData
	SellOrders []GetLatestFivePriceData
}

type GetLatestFivePriceData struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
