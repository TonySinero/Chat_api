package models

import "time"

type Message struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

type GroupMessage struct {
	ID          int       `json:"id"`
	SenderID    string    `json:"senderId"`
	RecipientID string    `json:"recipientId,omitempty"`
	GroupChatID string    `json:"groupChatId,omitempty"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
}
type ErrorResponse struct {
	Message string `json:"message"`
}

type GenerateTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponseAuth struct {
	Message       string `json:"message"`
	ResponseError string `json:"response_error"`
	Status        string `json:"status"`
}

type User struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}
type ResponseUser struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

type Post struct {
	Email    string
	Password string
}

type RestorePassword struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"password"`
}

type GroupChatCreateInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Members     []string `json:"members" binding:"required"`
}

type MessageCreateInput struct {
	Content     string `json:"content" binding:"required"`
	SenderID    string `json:"senderId" binding:"required"`
	GroupChatID string `json:"groupChatId" binding:"required"`
}

type GroupChat struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Members     []string  `json:"members"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
