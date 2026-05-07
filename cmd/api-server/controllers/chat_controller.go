package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mangahub/proto/message"
)

type ChatController struct {
	messageClient message.GRPCMessageServiceClient
}

func NewChatController(messageClient message.GRPCMessageServiceClient) *ChatController {
	return &ChatController{
		messageClient: messageClient,
	}
}

// GET /api/v1/mangas/:id/messages
func (cc *ChatController) GetChatHistory(c *gin.Context) {
	roomID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := cc.messageClient.GetChatHistory(c.Request.Context(), &message.GetChatHistoryRequest{
		RoomId: roomID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Chat history retrieved",
		"data":    resp.Messages,
	})
}
