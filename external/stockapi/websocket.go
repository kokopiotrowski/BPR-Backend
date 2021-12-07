package stockapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/gorilla/websocket"
)

//for serving to clients
type LiveData struct {
	Symbol        string  `json:"s"`
	CurrentPrice  float64 `json:"c"`
	OpenPrice     float64 `json:"o"`
	PercentChange float64 `json:"p"`
	Difference    float64 `json:"d"`
	Raising       int     `json:"r"`
}

//response from external API
type WebSocketResponse struct {
	Data *[]Data `json:"data,omitempty"`
}

type Data struct {
	S string  `json:"s,omitempty"` //symbol
	P float64 `json:"p,omitempty"` //last price
}

var (
	listeningClients = make(map[string]*websocket.Conn)
	loadForListeners = make(map[string]*LiveData)
	quoteMap         = make(map[string]finnhub.Quote)
	symbols          = []string{"AAPL", "MSFT", "AMZN",
		"TSLA", "NVDA", "GOOG", "GOOGL", "FB", "NFLX",
		"CMCSA", "CSCO", "COST", "AVGO", "PEP", "PYPL", "INTC",
		"QCOM", "TXN", "INTU", "AMD", "TMUS", "HON", "AMAT", "SBUX",
		"CHTR", "MRNA", "AMGN", "ISRG", "ADP", "ADI", "LRCX", "MU",
		"GILD", "BKNG", "MDLZ", "CSX", "MRVL", "REGN", "FISV", "ASML",
		"JD", "KLAC", "NXPI", "ADSK", "LULU", "ILMN", "XLNX", "VRTX", "SNPS"}
)

func StartListening(token string) {
	ready := make(chan int)

	for _, s := range symbols {
		go func(symbol string, ready chan int) {
			q, err := getQuoteForSymbol(symbol)
			if err != nil {
				ready <- 1
			}

			quoteMap[symbol] = q
			ready <- 0
		}(s, ready)
	}

	//waiting for all requests to be ready
	for _, s := range symbols {
		if <-ready != 0 {
			fmt.Printf("retrieving stock data for websocket failed %v", s)
			return
		}
	}

	prepairLoad()

	w, _, err := websocket.DefaultDialer.Dial("wss://ws.finnhub.io?token="+token, nil)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})

		if err := w.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Printf("error when sending message to websocket %v\n", err)
			return
		}
	}

	var res WebSocketResponse

	go startInformingListeners()

	for {
		err := w.ReadJSON(&res)
		if err != nil {
			panic(err)
		}

		if l, ok := loadForListeners[(*res.Data)[0].S]; ok {
			l.CurrentPrice = (*res.Data)[0].P
			l.Difference = (*res.Data)[0].P - l.OpenPrice
			l.PercentChange = (l.Difference / l.OpenPrice) * 100

			if l.PercentChange > 0 {
				l.Raising = 1
			} else if l.PercentChange < 0 {
				l.Raising = -1
			} else {
				l.Raising = 0
			}
		}
	}
}

func AddWsListenerClient(id string, conn *websocket.Conn) {
	listeningClients[id] = conn
	//sending initial load for the client
	err := conn.WriteJSON(loadForListeners)
	if err != nil {
		fmt.Printf("failed to send msg to the client %v %v\n", conn, err)
	}
}

func RemoveWsListenerClient(id string) {
	delete(listeningClients, id)
}

func startInformingListeners() {
	for {
		for _, c := range listeningClients {
			err := c.WriteJSON(loadForListeners)
			if err != nil {
				fmt.Printf("failed to send msg to the client %v %v\n", c, err)
				continue
			}
		}

		time.Sleep(4 * time.Second)
	}
}

func prepairLoad() {
	for _, s := range symbols {
		if l, ok := loadForListeners[s]; ok {
			l.OpenPrice = float64(*quoteMap[s].O)
			l.Symbol = s
		} else {
			loadForListeners[s] = &LiveData{
				OpenPrice: float64(*quoteMap[s].O),
				Symbol:    s,
			}
		}
	}
}
