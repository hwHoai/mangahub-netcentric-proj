package types

import "encoding/json"

type TCPMessage struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
	Token   string          `json:"token"`
}