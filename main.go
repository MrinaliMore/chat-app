package main

import (
	"chat-app/model"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//Creating map of connected clients
var clients = make(map[*websocket.Conn]bool)

//Creating broadcast connection channel
var broadcast = make(chan model.Message)

//Configure the upgrader
var upgrader = websocket.Upgrader{}

func main() {
	//Create simple file server
	fs := http.FileServer(http.Dir("/Users/mrinalimore/Go/src/chat-app/public"))
	http.Handle("/", fs)
	//Configure websocket routes
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Println("Http server started on localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Listen and Server error: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//Register our new clients
	clients[conn] = true

	for {
		var msg model.Message
		err := conn.ReadJSON(&msg)

		if err != nil {
			log.Printf("errors: %v", err)
			delete(clients, conn)
			break
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error : %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
