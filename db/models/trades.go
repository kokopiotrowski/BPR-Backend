package models

type Trades struct {
	Email         string          `json:"email"`
	Credits       float32         `json:"credits"`
	BuyingPower   float32         `json:"buyingPower"`
	AccountValue  float32         `json:"accountValue"`  // credits + value of current held stocks
	BoughtStocks  []BoughtStock   `json:"boughtStocks"`  // buy long
	SoldStocks    []SoldStock     `json:"soldStocks"`    // sell long
	ShortStocks   []ShortStock    `json:"shortStocks"`   // sell short
	BoughtToCover []BoughtToCover `json:"boughtToCover"` // buy to cover - buy short
	HoldLong      []HoldLong      `json:"holdLong"`      // long stocks currently held by user
	HoldShort     []HoldShort     `json:"holdShort"`     // short stocks currently held by user
}
