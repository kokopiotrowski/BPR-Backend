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
	BoughtStocks  []BoughtStock   `json:"boughtStocks"`  // buy long
	SoldStocks    []SoldStock     `json:"soldStocks"`    // sell long
	ShortStocks   []ShortStock    `json:"shortStocks"`   // sell short
	BoughtToCover []BoughtToCover `json:"boughtToCover"` // buy to cover - buy short
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
	Price  float32 `json:"p"` //price for one stock
	Date   int64   `json:"d"`
}

type ShortStock struct {
	Symbol string  `json:"s"`
	Amount int64   `json:"am"`
	Price  float32 `json:"p"` //price for one stock
	Date   int64   `json:"d"`
}

type BoughtToCover struct {
	Symbol string  `json:"s"`
	Amount int64   `json:"am"`
	Price  float32 `json:"p"` //price for one stock
	Date   int64   `json:"d"`
}

type Statistics struct {
	Email             string                    `json:"email"`
	Credits           []CreditsStatus           `json:"creditsStatus"`
	OwnedStocksAmount []OwnedStocksAmountStatus `json:"ownedStocksAmount"`
}

type CreditsStatus struct {
	Date    string  `json:"d"`
	Credits float64 `json:"c"`
}

type OwnedStocksAmountStatus struct {
	Date   string `json:"d"`
	Amount string `json:"a"`
}

var (
	DefaultAmountOfCreditsForNewUsers = float32(100000)
)
