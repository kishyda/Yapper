package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Message struct {
	AccountID string `json:"accountID"`
	ChatID    string `json:"chatID"`
	Message   string `json:"message"`
	Time      string `json:"time"`
	Close     bool   `json:"close"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connections = make(map[string]*websocket.Conn)

var version string
var database *sql.DB
var err error

func writeMessageToDatabase(message *Message) {
	tx, err := database.Begin()
	if err != nil {
		fmt.Print("Error beginning transaction\n")
		return
	}
	statement, err := tx.Prepare(`INSERT INTO messages (accountID, chatID, message, time) VALUES (?, ?, ?, ?)`)
	if err != nil {
		fmt.Print("Error preparing statement\n")
		return
	}
	_, err = statement.Exec(message.AccountID, message.ChatID, message.Message, message.Time)
	if err != nil {
		fmt.Print("Error executing statement\n")
		return
	}

	err = tx.Commit()
	if err != nil {
		fmt.Print("Error committing transaction\n")
		return
	}
	fmt.Print("Message written to database\n")
}

func sendMessage(message *Message) {
	var users string
	err := database.QueryRow(`SELECT users FROM CHAT WHERE chatID = ?`, message.ChatID).Scan(&users)
	if err != nil {
		fmt.Print("Error querying database\n")
		return
	}
	var userArray []string
	json.Unmarshal([]byte(users), &userArray)
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("Error upgrading connection\n")
		return
	}
	for {
		messageObject := Message{}
		_, p, _ := conn.ReadMessage()
		err := json.Unmarshal([]byte(p), &messageObject)
		if err != nil {
			fmt.Print("Error unmarshalling message\n")
			return
		}
		if messageObject.Close != false {
			return
		}
		writeMessageToDatabase(&messageObject)
		sendMessage(&messageObject)
	}
}

func createChat(w http.ResponseWriter, r *http.Request) {
	// Create a chat
}

func main() {

	database, err = sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		log.Println("Error opening database\n", err)
		return
	}
	database.QueryRow(`SELECT sqlite_version()`).Scan(&version)

	// Make sure the initial connection has data regarding the username, and user id
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		go handleConnection(w, r)
	})

	http.HandleFunc(`/createChat/`, func(w http.ResponseWriter, r *http.Request) {
		go createChat(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
