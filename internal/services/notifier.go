package services

import "github.com/rohit21755/gg_server.git/ws"

var Hub *ws.Hub

func InitNotifier(h *ws.Hub) {
	Hub = h
}

func NotifyAll(msg []byte) {
	if Hub != nil {
		Hub.Broadcast <- msg
	}
}
