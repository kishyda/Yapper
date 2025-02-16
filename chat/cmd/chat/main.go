package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/rs/cors"
)

type jsonOject struct {
	ChatID string `json:"chatID"`
	UserID string `json:"userID"`
}

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

type Success struct {
	Message string `json:"message"`
}

type Connection struct {
	Conn  *websocket.Conn
	mutex *sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var c = cors.New(cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	AllowedHeaders:   []string{"*"},
	AllowCredentials: true,
})

var connections = make(map[string]*Connection)

var version string
var database *sql.DB
var err error

func (c *Connection) WriteJSON(v interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Conn.WriteJSON(v)
}

func writeMessageToDatabase(conn *Connection, message *Message) {
	tx, err := database.Begin()
	if err != nil {
		fmt.Printf("Error beginning transaction: %s\n", err)
		conn.WriteJSON(Error{Message: "writeMessageToDatabase Error: Error beginning transaction"})
		return
	}

	statement, err := tx.Prepare(fmt.Sprintf(`INSERT INTO %s (accountID, accountName, chatID, chatName, message, time) VALUES (?, ?, ?, ?, ?, ?)`, message.ChatID))
	if err != nil {
		fmt.Printf("Error preparing statement: %s\n", err)
		conn.WriteJSON(Error{Message: "writeMessageToDatabase Error: Error preparing statement"})
		return
	}

	_, err = statement.Exec(message.AccountID, message.AccountName, message.ChatID, message.ChatName, message.Message, message.Time)
	if err != nil {
		fmt.Printf("Error executing statement: %s\n", err)
		conn.WriteJSON(Error{Message: "writeMessageToDatabase Error: Error executing statement"})
		return
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("Error committing transaction: %s\n", err)
		conn.WriteJSON(Error{Message: "writeMessageToDatabase Error: Error committing transaction"})
		return
	}

	fmt.Print("Message written to database\n")
	conn.WriteJSON(Success{Message: "Message written to database"})
}

func sendMessage(conn *Connection, message *Message) {
	var users string
	err := database.QueryRow(`SELECT users FROM CHAT WHERE chatName = ?`, message.ChatName).Scan(&users)
	if err != nil {
		fmt.Print("Error querying database ", err, "\n")
		conn.WriteJSON(Error{Message: "Error querying database"})
		return
	}

	var userArray []string
	err = json.Unmarshal([]byte(users), &userArray)
	if err != nil {
		fmt.Print("Error unmarshalling users ", err, "\n")
	}

	for _, user := range userArray {
		if connections[user] != nil {
			fmt.Print("Sending message to ", user, "\n")
			err := connections[user].Conn.WriteJSON(message)
			if err != nil {
				fmt.Print("Error writing message\n")
				conn.WriteJSON(Error{Message: "Error writing message"})
				return
			}
		}
	}
	conn.WriteJSON(Success{Message: "Message sent"})
}

func handleConnection(w http.ResponseWriter, r *http.Request, userID string) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("Error upgrading connection\n")
		return
	}
	if c == nil {
		fmt.Print("Connection is nil\n")
		return
	}
	conn := &Connection{Conn: c, mutex: &sync.Mutex{}}
	conn.WriteJSON(`'{"message": "Connection established"}'`)

	connections[userID] = conn

	go func() {
		defer conn.Conn.Close()
		defer delete(connections, userID)
		message := Message{}
		for {
			err := conn.Conn.ReadJSON(&message)
			if err != nil {
				if reflect.TypeOf(err) == reflect.TypeOf(&websocket.CloseError{}) {
					fmt.Print("websocket closed")
				} else {
					fmt.Printf("Unexpected close error: %v\n", err)
				}
				return
			}
			if message.Close == true {
				fmt.Print("Closing websocket\n")
				return
			}
			go sendMessage(conn, &message)
			go writeMessageToDatabase(conn, &message)
		}
	}()
}

func createChat(w http.ResponseWriter, r *http.Request) {
	tx, err := database.Begin()
	if err != nil {
		fmt.Print("Error beginning transaction\n")
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		return
	}

	body := jsonOject{}
	err = json.NewDecoder(r.Body).Decode(&body)
	chatID := body.ChatID

	statement, err := tx.Prepare(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (AccountID TEXT, AccountName TEXT, ChatID TEXT, ChatName TEXT, Message TEXT, Time TEXT)`, chatID))
	if err != nil {
		fmt.Print("Error preparing statement", err, "\n")
		http.Error(w, "Error preparing statement", http.StatusInternalServerError)
		return
	}
	_, err = statement.Exec()
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

func getMessages(w http.ResponseWriter, r *http.Request) {
	body := jsonOject{}
	err = json.NewDecoder(r.Body).Decode(&body)
	chatID := body.ChatID
	rows, err := database.Query(`SELECT * FROM ?`, chatID)
	defer rows.Close()
	if err != nil {
		fmt.Print("Error querying database\n")
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	messages := []Message{}
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

func getChats(w http.ResponseWriter, r *http.Request) {
	body := jsonOject{}
	err = json.NewDecoder(r.Body).Decode(&body)
	userID := body.UserID

	rows, err := database.Query(`SELECT * FROM CHAT WHERE users LIKE ?`, "%"+userID+"%")
	defer rows.Close()
	if err != nil {
		fmt.Print("Error querying database\n")
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	chats := []Chat{}
	for rows.Next() {
		chat := Chat{}
		err = rows.Scan(&chat.ChatID, &chat.ChatName, &chat.Users)
		chats = append(chats, chat)
	}
	w.Header().Set("Content-Type", "application/json")
	jsonChats, err := json.Marshal(chats)
	if err != nil {
		fmt.Print("Error marshalling chats\n")
		http.Error(w, "Error marshalling chats", http.StatusInternalServerError)
		return
	}
	w.Write(jsonChats)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var users string
	body := jsonOject{}
	err = json.NewDecoder(r.Body).Decode(&body)
	userID := body.UserID
	chatID := body.ChatID
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
	database, err = sql.Open("sqlite3", "./db/database.sqlite")
	if err != nil {
		fmt.Println("Error opening database\n", err)
		return
	}
	database.QueryRow(`SELECT sqlite_version()`).Scan(&version)

	http.Handle("/ws/", c.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleConnection(w, r, r.URL.Query().Get("userID"))
	})))

	http.Handle("/createChat", c.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			go createChat(w, r)
		}
	})))

	http.Handle("/getMessages", c.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			go getMessages(w, r)
		}
	})))

	http.Handle("/getChats", c.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			go getChats(w, r)
		}
	})))

	http.Handle("/addUser", c.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			go addUser(w, r)
		}
	})))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print(err)
	}
}
