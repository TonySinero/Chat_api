package handler

import (
	"chat/internal/repository"
	"chat/internal/server"
	"chat/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services        *service.Service
	WebSocketServer *server.WebSocketServer
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services:        services,
		WebSocketServer: server.NewWebSocketServer(),
	}
}

func (h *Handler) InitRoutes(dbType string) *gin.Engine {
	r := gin.Default()
	r.Use(h.CorsMiddleware)

	auth := r.Group("/auth")
	{
		auth.POST("/user", h.createUser)
		auth.POST("/login", h.authUser)
		auth.POST("/restore", h.restorePassword)
		auth.POST("/refresh", h.RefreshToken)
	}

	switch dbType {
	case repository.PostgresDB:
		r.Use(h.parseAuthHeader, h.checkRole)
		h.initPostgresRoutes(r)
	}

	r.GET("/ws", func(c *gin.Context) {
		h.WebSocketServer.HandleConnections(c.Writer, c.Request)
	})

	return r
}

func (h *Handler) initPostgresRoutes(r *gin.Engine) {
	r.PUT("/postgres/user/:id", h.updateUser)
	r.DELETE("/postgres/user/:id", h.deleteUser)
	r.GET("/postgres/users", h.getUsers)
	r.GET("/postgres/user/:id", h.getUser)

	r.GET("/postgres/messages", h.getMessagesPostgres)
	r.GET("/postgres/message/:id", h.getMessagePostgres)
	r.POST("/postgres/message", h.createMessagePostgres)
	r.PUT("/postgres/message/:id", h.updateMessagePostgres)
	r.DELETE("/postgres/message/:id", h.deleteMessagePostgres)

	r.POST("/postgres/group-chat", h.createGroupChat)
	r.POST("/postgres/group-chat/:id/join", h.joinGroupChat)
	r.POST("/postgres/group-chat/:id/message", h.sendMessageToGroupChat)
	r.POST("/postgres/user/:id/message", h.sendPrivateMessage)
}
