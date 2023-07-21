package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	conn *websocket.Conn
}

type WebSocketServer struct {
	upgrader  *websocket.Upgrader
	clients   map[*Client]bool
	broadcast chan []byte
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
	}
}

func (s *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("Failed to upgrade to WebSocket:", err)
		return
	}

	client := &Client{conn: conn}
	s.clients[client] = true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logrus.Error("Error reading message:", err)
			delete(s.clients, client)
			break
		}

		s.broadcast <- message
	}

	conn.Close()
}

func (s *WebSocketServer) Run(port string, router *gin.Engine) error {
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
		ErrorLog:     log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}

	go s.handleMessages()

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error running server: %s", err)
	}
	return nil
}

func (s *WebSocketServer) Shutdown() error {
	close(s.broadcast)
	for client := range s.clients {
		client.conn.Close()
		delete(s.clients, client)
	}
	return nil
}

func (s *WebSocketServer) handleMessages() {
	for message := range s.broadcast {
		for client := range s.clients {
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logrus.Error("Failed to send message to client:", err)
				delete(s.clients, client)
			}
		}
	}
}
