package models

type Portfolio struct {
	Credits      float32       `json:"credits"`
	BoughtStocks []BoughtStock `json:"boughtStocks"` // buy long
	ShortStocks  []ShortStock  `json:"shortStocks"`  // sell short
}
