package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func writeToDatabase(message string) {
	// Write to database
}

func main() {
	connIDs := make(map[string]string)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			jsonObject := make(map[string]string)
			json.Unmarshal([]byte(r.Body.Close().Error()), &jsonObject)
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			for {
				_, p, _ := conn.ReadMessage()
				// conn.WriteMessage(messageType, p)
				writeToDatabase(string(p))
			}
		}()
	})
	http.ListenAndServe(":8080", nil)
}
