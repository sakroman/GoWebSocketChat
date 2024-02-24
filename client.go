package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/gorilla/websocket"
)

func Nicknames() string {
	name := namegenerator.NewNameGenerator(rand.Int63()).Generate()
	return name
}

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	Username   string

	broadcast chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		Username:   Nicknames(),
		broadcast:  make(chan []byte),
	}
}

func (c *Client) handleMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error reading message: %v", err)
			}
			break
		}
		var data Payload
		err = json.Unmarshal(payload, &data)
		if err != nil {
			fmt.Printf("error parsing payload: %v", err)
			continue
		}
		data.Username = c.Username

		data.TimeSent = time.Now().Format(time.TimeOnly)

		updatedPayload, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("error marshaling payload: %v", err)
			continue
		}

		for client := range c.manager.clients {
			client.broadcast <- []byte(updatedPayload)
		}

	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		msg := <-c.broadcast
		var data Payload
		err := json.Unmarshal(msg, &data)
		if err != nil {
			fmt.Printf("error parsing payload: %v", err)
			continue
		}

		chatMessage := data.ChatMessage
		fmt.Printf("Chat Message: %s\n", chatMessage)

		htmlContent := fmt.Sprintf("<div hx-swap-oob=\"beforeend:#chat-room\"><div><span class='fw-semibold'>%s (%s)</span><p>%s</p></div></div>", data.Username, data.TimeSent, chatMessage)

		if err := c.connection.WriteMessage(websocket.TextMessage, []byte(htmlContent)); err != nil {
			fmt.Println(err)
			return
		}
	}
}
