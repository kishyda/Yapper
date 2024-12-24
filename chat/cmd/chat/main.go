package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default
		return true
	},
}

func writeToDatabase(message string) {
	// Write to database
}

func main() {
	// connections := make(map[string]*websocket.Conn)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("got a connection")
		go func() {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Print("Error reading body")
				return
			}
			defer r.Body.Close()
			jsonObject := make(map[string]string)
			json.Unmarshal([]byte(body), &jsonObject)
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				fmt.Print("Error upgrading connection")
				return
			}
			defer conn.Close()
			for {
				_, p, _ := conn.ReadMessage()
				if string(p) == "quit" {
					conn.Close()
					break
				}
				// conn.WriteMessage(messageType, p)
				writeToDatabase(string(p))
			}
		}()
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("got a connection\n")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Print("Error upgrading connection\n")
			return
		}
		defer conn.Close()
		conn.WriteMessage(websocket.TextMessage, []byte("Hello"))
	})
	http.ListenAndServe(":8080", nil)
}
