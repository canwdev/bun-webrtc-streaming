package main

import (
	"crypto/rand"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

//go:embed public/*
var publicFS embed.FS

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type client struct {
	role string
	conn *websocket.Conn
	mu   sync.Mutex
}

var (
	clients   = map[string]*client{}
	clientsMu sync.RWMutex
)

func genID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	id := genID()
	c := &client{conn: conn}

	clientsMu.Lock()
	clients[id] = c
	clientsMu.Unlock()

	defer func() {
		conn.Close()
		clientsMu.Lock()
		delete(clients, id)
		clientsMu.Unlock()
	}()

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		if role, ok := msg["role"].(string); ok {
			c.role = role
		}

		msg["from"] = id

		if to, ok := msg["to"].(string); ok {
			clientsMu.RLock()
			target, exists := clients[to]
			clientsMu.RUnlock()
			if exists {
				payload, _ := json.Marshal(msg)
				target.mu.Lock()
				target.conn.WriteMessage(websocket.TextMessage, payload)
				target.mu.Unlock()
			}
		} else {
			payload, _ := json.Marshal(msg)
			clientsMu.RLock()
			for otherID, other := range clients {
				if otherID != id {
					other.mu.Lock()
					other.conn.WriteMessage(websocket.TextMessage, payload)
					other.mu.Unlock()
				}
			}
			clientsMu.RUnlock()
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3661"
	}

	publicSub, _ := fs.Sub(publicFS, "public")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if websocket.IsWebSocketUpgrade(r) {
			handleWS(w, r)
			return
		}

		path := r.URL.Path
		var file string
		switch path {
		case "/", "/index.html":
			file = "index.html"
		case "/broadcaster":
			file = "broadcaster.html"
		case "/viewer":
			file = "viewer.html"
		default:
			http.NotFound(w, r)
			return
		}

		data, err := fs.ReadFile(publicSub, file)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	addr := ":" + port
	fmt.Printf("Server: http://localhost:%s\n", port)
	fmt.Printf("  Broadcaster: http://localhost:%s/broadcaster\n", port)
	fmt.Printf("  Viewer:      http://localhost:%s/viewer\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
