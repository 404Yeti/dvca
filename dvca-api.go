package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	ID      int    `json:"id"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

var db *sql.DB

func main() {
	var err error

	// Open the database (creates it if it doesn't exist)
	db, err = sql.Open("sqlite3", "./dvca.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sender TEXT,
			content TEXT
		)
	`)
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}

	http.HandleFunc("/api/send", sendMessage)
	http.HandleFunc("/api/messages", getMessages)

	fmt.Println("📡 DVCA API server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// POST /api/send
func sendMessage(w http.ResponseWriter, r *http.Request) {
	var m Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	stmt, err := db.Prepare("INSERT INTO messages(sender, content) VALUES(?, ?)")
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	_, err = stmt.Exec(m.Sender, m.Content)
	if err != nil {
		http.Error(w, "Insert failed", 500)
		return
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, "Message stored")
}

// GET /api/messages?user=xxx
func getMessages(w http.ResponseWriter, r *http.Request) {
	// ❌ DELIBERATE SQLi RISK
	user := r.URL.Query().Get("user")
	query := fmt.Sprintf("SELECT id, sender, content FROM messages WHERE sender = '%s'", user)

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Query error", 500)
		return
	}
	defer rows.Close()

	var results []Message
	for rows.Next() {
		var m Message
		rows.Scan(&m.ID, &m.Sender, &m.Content)
		results = append(results, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
