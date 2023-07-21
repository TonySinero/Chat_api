package handler

import (
	"net/http"
	"strconv"

	"chat/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getMessagePostgres(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.Warnf("Handler getMessagePostgres (reading param): %s", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid id"})
		return
	}

	t, err := h.services.MessagePostgresService.GetMessage(id)
	if err != nil {
		logrus.Errorf("Handler getMessagePostgres (db get): %s", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "failed to get message"})
		return
	}

	ctx.JSON(http.StatusOK, t)
}

func (h *Handler) getMessagesPostgres(ctx *gin.Context) {
	var page int64 = 1
	var limit int64 = 10

	if ctx.Query("page") != "" {
		paramPage, err := strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil || paramPage < 0 {
			logrus.Warnf("No url request: %s", err)
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid url query"})
			return
		}
		page = paramPage
	}
	if ctx.Query("limit") != "" {
		paramLimit, err := strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil || paramLimit < 0 {
			logrus.Warnf("No url request: %s", err)
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid url query"})
			return
		}
		limit = paramLimit
	}

	todos, pages, err := h.services.MessagePostgresService.GetMessages(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.Header("pages", strconv.Itoa(pages))
	ctx.JSON(http.StatusOK, todos)
}

func (h *Handler) createMessagePostgres(ctx *gin.Context) {
	var input models.Message
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logrus.Warnf("Handler createMessagePostgres (binding JSON): %s", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid request"})
		return
	}

	id, err := h.services.MessagePostgresService.CreateMessage(&input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, id)
}

func (h *Handler) updateMessagePostgres(ctx *gin.Context) {
	var input models.Message
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		logrus.Warnf("Handler updateMessagePostgres (reading param): %s", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid request"})
		return
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logrus.Warnf("Handler updateMessagePostgres (binding JSON): %s", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid request"})
		return
	}
	input.ID = id
	err = h.services.MessagePostgresService.UpdateMessage(&input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
}

func (h *Handler) deleteMessagePostgres(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		logrus.Warnf("Handler deleteMessagePostgres (reading param): %s", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid id"})
		return
	}

	_, err = h.services.MessagePostgresService.DeleteMessageByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
