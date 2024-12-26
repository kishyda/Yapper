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
	AccountID   string `json:"accountID"`
	AccountName string `json:"accountName"`
	ChatID      string `json:"chatID"`
	ChatName    string `json:"chatName"`
	Message     string `json:"message"`
	Time        string `json:"time"`
	Close       bool   `json:"close"`
}

type Chat struct {
	ChatName string `json:"chatName"`
	ChatID   string `json:"chatID"`
	Users    string `json:"users"`
}

type Error struct {
	Message string `json:"message"`
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
	statement, err := tx.Prepare(`INSERT INTO ? (accountID, accountName, chatID, chatName, message, time) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Print("Error preparing statement\n")
		return
	}
	_, err = statement.Exec(message.ChatID, message.AccountID, message.AccountName, message.ChatID, message.ChatName, message.Message, message.Time)
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
	for _, user := range userArray {
		if connections[user] != nil {
			err := connections[user].WriteJSON(message)
			if err != nil {
				fmt.Print("Error writing message\n")
				return
			}
		}
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	connections[r.URL.Query().Get("accountID")] = conn
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
			w.Write([]byte("Connection closed"))
			return
		}
		writeMessageToDatabase(&messageObject)
		sendMessage(&messageObject)
	}
}

func createChat(w http.ResponseWriter, chatID string) {
	tx, err := database.Begin()
	if err != nil {
		fmt.Print("Error beginning transaction\n")
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		return
	}
	statement, err := tx.Prepare(`CREATE TABLE IF NOT EXISTS ? (AccountID TEXT, AccountName TEXT, ChatID TEXT, ChatName TEXT, Message TEXT, Time TEXT)`)
	if err != nil {
		fmt.Print("Error preparing statement\n")
		http.Error(w, "Error preparing statement", http.StatusInternalServerError)
		return
	}
	_, err = statement.Exec(chatID)
	if err != nil {
		fmt.Print("Error executing statement\n")
		http.Error(w, "Error executing statement", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		fmt.Print("Error committing transaction\n")
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Chat created"))
}

func getMessages(w http.ResponseWriter, chatID string) {
	messages := []Message{}
	rows, err := database.Query(`SELECT * FROM ?`, chatID)
	defer rows.Close()
	if err != nil {
		fmt.Print("Error querying database\n")
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		message := Message{}
		err = rows.Scan(&message.AccountID, &message.ChatID, &message.Message, &message.Time)
		messages = append(messages, message)
	}
	w.Header().Set("Content-Type", "application/json")
	jsonMessages, err := json.Marshal(messages)
	if err != nil {
		fmt.Print("Error marshalling messages\n")
		http.Error(w, "Error marshalling messages", http.StatusInternalServerError)
		return
	}
	w.Write(jsonMessages)
}

func addUser(w http.ResponseWriter, userID string, chatID string) {
	var users string
	database.QueryRow(`SELECT users FROM CHAT WHERE chatID = ?`, chatID).Scan(&users)
	var userArray []string
	json.Unmarshal([]byte(users), &userArray)
	userArray = append(userArray, userID)
	tx, err := database.Begin()
	if err != nil {
		fmt.Print("Error beginning transaction\n")
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		return
	}
	statement, err := tx.Prepare(`UPDATE CHAT SET users = ? WHERE chatID = ?`)
	if err != nil {
		fmt.Print("Error preparing statement\n")
		http.Error(w, "Error preparing statement", http.StatusInternalServerError)
		return
	}
	_, err = statement.Exec(userArray, chatID)
	if err != nil {
		fmt.Print("Error executing statement\n")
		http.Error(w, "Error executing statement", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		fmt.Print("Error committing transaction\n")
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User added to chat"))
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
		go createChat(w, r.URL.Query().Get("chatID"))
	})

	http.HandleFunc("/getMessages/", func(w http.ResponseWriter, r *http.Request) {
		go getMessages(w, r.URL.Query().Get("chatID"))
	})

	http.HandleFunc("/addUser/", func(w http.ResponseWriter, r *http.Request) {
		go addUser(w, r.URL.Query().Get("userID"), r.URL.Query().Get("chatID"))
	})

	http.ListenAndServe(":8080", nil)
}
