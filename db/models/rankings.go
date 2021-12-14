package models

type Rankings struct {
	Region             string             `json:"region"`
	Date               string             `json:"date"`
	AccountValueRating []UserAccountValue `json:"accountValue"`
	MostPopularStock   []MostPopularStock `json:"stockAmount"`
}

type UserAccountValue struct {
	Username     string  `json:"username"`
	AccountValue float32 `json:"accountValue"`
}

type UserCredits struct {
	Username string  `json:"email"`
	Credits  float32 `json:"credits"`
}

type MostPopularStock struct {
	Symbol string `json:"symbol"`
	Amount int64  `json:"amount"`
}
