package handler

import (
	"chat/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CorsMiddleware(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "*")
	ctx.Header("Access-Control-Allow-Headers", "*")
	ctx.Header("Content-Type", "application/json")

	if ctx.Request.Method != "OPTIONS" {
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusOK)
	}
}

func (h *Handler) parseAuthHeader(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: "empty auth header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: "invalid header"})
		return
	}

	if len(headerParts[1]) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: "token is empty"})
		return
	}

	id, role, err := h.services.Authorization.ParseToken(headerParts[1])

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.Set("role", role)
	ctx.Set("id", id)
}

func (h *Handler) checkRole(ctx *gin.Context) {
	necessaryRole := []string{string(models.RoleUser), string(models.RoleAdmin)}
	if err := h.services.Authorization.CheckRole(necessaryRole, ctx.GetString("role")); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{Message: "not enough rights"})
		return
	}
}
