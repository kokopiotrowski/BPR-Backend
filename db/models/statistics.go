package models

type Statistics struct {
	Email             string                    `json:"email"`
	Credits           []CreditsStatus           `json:"credits"`
	AccountValue      []AccountValueStatus      `json:"accountValue"`
	OwnedStocksAmount []OwnedStocksAmountStatus `json:"ownedStocksAmount"`
	HoldStocks        []HoldStocksStatus        `json:"holdStocks"`
	TotalGain         []TotalGain               `json:"totalGain"`
}

type CreditsStatus struct {
	Date    string  `json:"d"`
	Credits float32 `json:"c"`
}

type OwnedStocksAmountStatus struct {
	Date   string `json:"d"`
	Amount int64  `json:"a"`
}

type AccountValueStatus struct {
	Date        string  `json:"d"`
	AcountValue float32 `json:"value"`
}

type HoldStocksStatus struct {
	Date   string      `json:"d"`
	HLong  []HoldLong  `json:"holdLong"`
	HShort []HoldShort `json:"holdShort"`
}

type TotalGain struct {
	Date string  `json:"d"`
	Gain float32 `json:"gain"`
}
