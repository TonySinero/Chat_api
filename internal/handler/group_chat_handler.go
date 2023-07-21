package handler

import (
	"chat/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createGroupChat(ctx *gin.Context) {
	var input models.GroupChatCreateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	chat, err := h.services.ChatService.CreateGroupChat(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group chat"})
		return
	}

	ctx.JSON(http.StatusCreated, chat)
}

func (h *Handler) joinGroupChat(ctx *gin.Context) {
	groupChatID := ctx.Param("id")
	userID := ctx.MustGet("id").(string)

	if !h.services.ChatService.CanUserAccessGroupChat(userID, groupChatID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.services.ChatService.AddUserToGroupChat(userID, groupChatID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join group chat"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Joined group chat successfully"})
}

func (h *Handler) sendMessageToGroupChat(ctx *gin.Context) {
	groupChatID := ctx.Param("id")
	userID := ctx.MustGet("id").(string)

	if !h.services.ChatService.CanUserAccessGroupChat(userID, groupChatID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var input models.MessageCreateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	message, err := h.services.ChatService.SendMessageToGroupChat(userID, groupChatID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	ctx.JSON(http.StatusCreated, message)
}

func (h *Handler) sendPrivateMessage(ctx *gin.Context) {
	recipientID := ctx.Param("id")
	senderID := ctx.MustGet("id").(string)

	if !h.services.ChatService.CanUserSendMessage(senderID, recipientID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var input models.MessageCreateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	message, err := h.services.ChatService.SendPrivateMessage(senderID, recipientID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send private message"})
		return
	}

	ctx.JSON(http.StatusCreated, message)
}
