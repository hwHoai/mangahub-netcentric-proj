package models

type UDPMessage struct {
	Action    string `json:"action"`     
	UserID    string `json:"user_id"`    
	Secret    string `json:"secret"`	  
	Content   string `json:"content"`    
}