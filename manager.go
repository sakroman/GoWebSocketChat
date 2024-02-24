package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// var upgrader = websocket.Upgrader{
// 	// CheckOrigin: func(r *http.Request) bool {
// 	// 	return true
// 	// },
// }

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (m *Manager) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := NewClient(conn, m)
	m.addClient(client)

	go client.handleMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()

	defer m.Unlock()

	m.clients[client] = true
	m.broadcastParticipantList()
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
	m.broadcastParticipantList()
}

func (m *Manager) broadcastParticipantList() {
	// Construct the list of participant usernames
	var usernames []string
	for client := range m.clients {
		usernames = append(usernames, client.Username)
	}

	addParticipant := fmt.Sprintf("<div hx-swap-oob=\"innerHTML:#participants\"><div><span class='fw-semibold'>%s</span></div></div>", strings.Join(usernames, "<br>"))
	for client := range m.clients {
		err := client.connection.WriteMessage(websocket.TextMessage, []byte(addParticipant))
		if err != nil {
			fmt.Println("Error sending participant list to client:", err)
		}
	}
}
