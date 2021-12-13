package api

import (
	"fmt"
	"log"
	"net/http"
	"stockx-backend/external/stockapi"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	id := xid.New()

	fmt.Printf("New client %v connected to websocket\n", id.String())

	err = ws.WriteMessage(1, []byte("Connected to live stock data"))
	if err != nil {
		log.Println(err)
	}

	stockapi.AddWsListenerClient(id.String(), ws)

	go reader(id.String(), ws)
}

func reader(id string, conn *websocket.Conn) {
	defer closeConnection(id)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func closeConnection(id string) {
	stockapi.RemoveWsListenerClient(id)
	fmt.Printf("Websocket connection closed with client: %v", id)
}
