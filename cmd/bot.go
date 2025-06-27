package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender    string `json:"sender"`
	Content   string `json:"content"`
}

func main() {
	fmt.Println("🤖 DVCA Bot starting up...")

	// Connect to the WebSocket server
	wsURL := "ws://localhost:8090/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	// Send welcome message
	hello := Message{
		Sender:  "dvca-bot",
		Content: "Hello, humans. I am your insecure assistant. Type @ai help",
	}
	sendJSON(conn, hello)

	// Handle Ctrl+C to exit cleanly
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}

			var m Message
			if err := json.Unmarshal(msg, &m); err != nil {
				continue
			}

			handleMessage(conn, m)
		}
	}()

	<-interrupt
	fmt.Println("🛑 DVCA Bot shutting down.")
}

func handleMessage(conn *websocket.Conn, m Message) {
	if m.Sender == "dvca-bot" {
		return // Don't respond to itself
	}

	if m.Content == "@ai help" {
		sendJSON(conn, Message{
			Sender:  "dvca-bot",
			Content: "Commands: @ai joke, @ai hack, @ai status",
		})
	} else if m.Content == "@ai joke" {
		jokes := []string{
			"Why do hackers wear leather jackets? Because they crack systems!",
			"I told a firewall a joke. It blocked me.",
			"Knock knock. Who’s there? admin. admin who? Access denied.",
		}
		sendJSON(conn, Message{
			Sender:  "dvca-bot",
			Content: jokes[rand.Intn(len(jokes))],
		})
	} else if m.Content == "@ai hack" {
		sendJSON(conn, Message{
			Sender:  "dvca-bot",
			Content: "<i>Initiating simulated hack... done. You are root now (jk)</i>",
		})
	} else if m.Content == "@ai status" {
		sendJSON(conn, Message{
			Sender:  "dvca-bot",
			Content: "I'm alive and totally not vulnerable to anything. 💀",
		})
	}
}

func sendJSON(conn *websocket.Conn, m Message) {
	data, _ := json.Marshal(m)
	conn.WriteMessage(websocket.TextMessage, data)
}
