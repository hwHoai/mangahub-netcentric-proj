package ws_utils_pool

import (

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	RoomID string
	Conn   *websocket.Conn
	Pool   ChatPool
}

type Message struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
	Sender  string `json:"sender"`
}

type ChatPool interface {
	Register(client *Client)
	Unregister(client *Client)
	Broadcast(msg Message)
}