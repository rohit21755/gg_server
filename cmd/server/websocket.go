package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rohit21755/gg_server.git/ws"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// Allow frontend / mobile connections
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// serveWS upgrades the HTTP connection to a WebSocket connection
func serveWS(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}

	client := &ws.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	hub.Register <- client

	// Start goroutines
	go clientWriter(hub, client)
	go clientReader(hub, client)
}

// Reads messages FROM the client
func clientReader(hub *ws.Hub, client *ws.Client) {
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(512)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	client.Conn.SetPongHandler(func(appData string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println("WS Read error:", err)
			}
			break
		}

		// Echo or broadcast incoming messages
		hub.Broadcast <- message
	}
}

// Sends messages TO the client
func clientWriter(hub *ws.Hub, client *ws.Client) {
	ticker := time.NewTicker(45 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				// Hub closed the channel
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(msg)

			// Send queued messages in this frame to prevent blocking
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// Ping to keep connection alive
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
