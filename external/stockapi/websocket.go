package stockapi

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/gorilla/websocket"
)

//for serving to clients
type LiveData struct {
	Symbol             string  `json:"s"`
	CurrentPrice       float64 `json:"c"`
	OpenPrice          float64 `json:"o"`
	PercentChange      float64 `json:"p"`
	Difference         float64 `json:"d"`
	Raising            int     `json:"r"`
	previousClosePrice float64
}

//response from external API
type WebSocketResponse struct {
	Data *[]Data `json:"data,omitempty"`
}

type Data struct {
	S string  `json:"s,omitempty"` //symbol
	P float64 `json:"p,omitempty"` //last price
}

type ListeningClient struct {
	c  *websocket.Conn
	mu sync.Mutex
}

var (
	externalWebsocket   *websocket.Conn
	listeningClients    = make(map[string]*ListeningClient)
	loadForListeners    = make(map[string]*LiveData)
	quoteMap            = make(map[string]finnhub.Quote)
	liveDataPropagation func()
	symbols             = []string{"AAPL", "MSFT", "AMZN",
		"TSLA", "NVDA", "GOOG", "GOOGL", "FB", "NFLX",
		"CMCSA", "CSCO", "COST", "AVGO", "PEP", "PYPL", "INTC",
		"QCOM", "TXN", "INTU", "AMD", "TMUS", "HON", "AMAT", "SBUX",
		"CHTR", "MRNA", "AMGN", "ISRG", "ADP", "ADI", "LRCX", "MU",
		"GILD", "BKNG", "MDLZ", "CSX", "MRVL", "REGN", "FISV", "ASML",
		"JD", "KLAC", "NXPI", "ADSK", "LULU", "ILMN", "XLNX", "VRTX", "SNPS"}
)

func StartListening(token string) {
	if externalWebsocket != nil {
		externalWebsocket.Close()
	}

	if len(loadForListeners) < 1 {
		ready := make(chan int)
		for _, s := range symbols {
			go func(symbol string, ready chan int) {
				q, err := GetQuoteForSymbol(symbol)
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
	}

	externalWebsocket, _, err := websocket.DefaultDialer.Dial("wss://ws.finnhub.io?token="+token, nil)
	if err != nil {
		fmt.Printf("Failed to connect to external websocket %v", err)
		return
	}

	defer reattemptConnection(externalWebsocket, token)

	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})

		if err := externalWebsocket.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Printf("error when sending message to websocket %v\n", err)
			return
		}
	}

	data := &WebSocketResponse{}

	if liveDataPropagation == nil {
		liveDataPropagation = informListeners
		go liveDataPropagation()
	}

	for {
		err := externalWebsocket.ReadJSON(data)
		if err != nil {
			fmt.Printf("failed to read message from websocket %v\n", err)
			return
		}

		if data == nil {
			continue
		}

		{
			if l, ok := loadForListeners[(*data.Data)[0].S]; ok {
				l.CurrentPrice = (*data.Data)[0].P
				l.Difference = (*data.Data)[0].P - l.previousClosePrice
				l.PercentChange = (l.Difference / l.previousClosePrice) * 100

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
}

func reattemptConnection(w *websocket.Conn, token string) {
	w.Close()
	time.Sleep(1 * time.Minute)
	fmt.Printf("Reattempt to connect to external websocket...\n") //tends to happen when stock market is closed
	StartListening(token)
}

func AddWsListenerClient(id string, conn *websocket.Conn) {
	listeningClients[id] = &ListeningClient{
		c:  conn,
		mu: sync.Mutex{},
	}
	//sending initial load for the client
	err := conn.WriteJSON(loadForListeners)
	if err != nil {
		fmt.Printf("failed to send msg to the client %v %v\n", conn, err)
	}
}

func RemoveWsListenerClient(id string) {
	listeningClients[id].c.Close()
	delete(listeningClients, id)
}

func informListeners() {
	for {
		for _, l := range listeningClients {
			err := l.sendLoad(loadForListeners)
			if err != nil {
				fmt.Printf("failed to send msg to the client %v\n", err)
				continue
			}
		}

		time.Sleep(4 * time.Second)
	}
}
func (l *ListeningClient) sendLoad(v interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.c.WriteJSON(v)
}

func prepairLoad() {
	for _, s := range symbols {
		if l, ok := loadForListeners[s]; ok {
			l.OpenPrice = float64(*quoteMap[s].O)
			l.Symbol = s
		} else {
			loadForListeners[s] = &LiveData{
				OpenPrice:          float64(*quoteMap[s].O),
				previousClosePrice: float64(*quoteMap[s].Pc),
				Symbol:             s,
			}
		}
	}
}
