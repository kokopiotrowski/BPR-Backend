package models


type BoughtStock struct {
	Symbol string  `json:"s"`
	Amount int64   `json:"am"`
	Price  float32 `json:"p"` //price for one stock
	Date   int64   `json:"d"`
}

type SoldStock struct {
	Symbol string  `json:"s"`
	Amount int64   `json:"am"`
	Price  float32 `json:"p"` // buy price for one stock
	Date   int64   `json:"d"`
}

type ShortStock struct {
	Symbol string  `json:"s"`
	Amount int64   `json:"am"`
	Price  float32 `json:"p"` // buy price for one stock
	Date   int64   `json:"d"`
}

type BoughtToCover struct {
	Symbol string  `json:"s"`
	Amount int64   `json:"am"`
	Price  float32 `json:"p"` // buy price for one stock
	Date   int64   `json:"d"`
}

type HoldLong struct {
	Symbol       string  `json:"s"`
	Amount       int64   `json:"am"`
	Price        float32 `json:"p"` // buy price for one stock
	CurrentPrice float32 `json:"c"`
	Gain         float32 `json:"g"`
	GainP        float32 `json:"gp"`
}

type HoldShort struct {
	Symbol       string  `json:"s"`
	Amount       int64   `json:"am"`
	Price        float32 `json:"p"` // buy price for one stock
	CurrentPrice float32 `json:"c"`
	Gain         float32 `json:"g"`
	GainP        float32 `json:"gp"`
}

type HoldStocks struct {
	HoldLong  []HoldLong  `json:"holdLong"`  // long stocks currently held by user
	HoldShort []HoldShort `json:"holdShort"` // short stocks currently held by user
}


