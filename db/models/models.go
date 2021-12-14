package models

type User struct {
	Email            string `json:"email"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	IsAdmin          bool   `json:"isAdmin"`
	DateCreated      string `json:"dateCreated"`
	DateUpdated      string `json:"dateUpdated"`
	DateLastAccessed string `json:"dateLastAccessed"`
	IsConfirmed      bool   `json:"isConfirmed"`
}

type Trades struct {
	Email         string          `json:"email"`
	Credits       float32         `json:"credits"`
	AccountValue  float32         `json:"accountValue"`  // credits + value of current held stocks
	BoughtStocks  []BoughtStock   `json:"boughtStocks"`  // buy long
	SoldStocks    []SoldStock     `json:"soldStocks"`    // sell long
	ShortStocks   []ShortStock    `json:"shortStocks"`   // sell short
	BoughtToCover []BoughtToCover `json:"boughtToCover"` // buy to cover - buy short
	HoldLong      []HoldLong      `json:"holdLong"`      // long stocks currently held by user
	HoldShort     []HoldShort     `json:"holdShort"`     // short stocks currently held by user
}

type Portfolio struct {
	Credits      float32       `json:"credits"`
	BoughtStocks []BoughtStock `json:"boughtStocks"` // buy long
	ShortStocks  []ShortStock  `json:"shortStocks"`  // sell short
}

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

type ListOfRegisteredUsers struct {
	Region string   `json:"region"` // future possibility to separate registered users by region
	Users  []string `json:"users"`
}

type HoldStocks struct {
	HoldLong  []HoldLong  `json:"holdLong"`  // long stocks currently held by user
	HoldShort []HoldShort `json:"holdShort"` // short stocks currently held by user
}

var (
	DefaultAmountOfCreditsForNewUsers = float32(100000)
)
