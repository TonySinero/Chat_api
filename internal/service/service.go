package service

import (
	"chat/internal/models"
	"chat/internal/repository"
	"context"
	"errors"
)

type MessagePostgresService interface {
	GetMessage(id int) (*models.Message, error)
	GetMessages(page, limit int64) ([]models.Message, int, error)
	CreateMessage(message *models.Message) (int, error)
	UpdateMessage(message *models.Message) error
	DeleteMessageByID(id int) (int, error)
}

type Authorization interface {
	CreateUser(ctx context.Context, user *models.User) (*models.GenerateTokens, error)
	AuthUser(ctx context.Context, user *models.User) (tokens *models.GenerateTokens, err error)
	User(ctx context.Context, userID int) (*models.ResponseUser, error)
	Users(ctx context.Context, page, limit int64) ([]models.ResponseUser, error)
	UpdateUser(ctx context.Context, inputUser *models.ResponseUser) error
	DeleteUser(ctx context.Context, userID int) error
	RefreshToken(refreshToken string) (*models.GenerateTokens, error)
	CheckRole(neededRoles []string, givenRole string) error
	ParseToken(token string) (int, string, error)
	GenerateTokens(user *models.User) (*models.GenerateTokens, error)
	RestorePassword(ctx context.Context, restore *models.RestorePassword) error
}

type ChatService interface {
	CreateGroupChat(input models.GroupChatCreateInput) (*models.GroupChat, error)
	JoinGroupChat(userID, groupChatID string) error
	SendMessageToGroupChat(userID, groupChatID string, input models.MessageCreateInput) (*models.GroupMessage, error)
	SendPrivateMessage(senderID, recipientID string, input models.MessageCreateInput) (*models.GroupMessage, error)
	CanUserAccessGroupChat(userID, groupChatID string) bool
	CanUserSendMessage(senderID, recipientID string) bool
	AddUserToGroupChat(userID, groupChatID string) error
}

type Service struct {
	MessagePostgresService
	Authorization
	ChatService
}

var serviceFactories = map[string]func(*repository.Repository) interface{}{
	repository.PostgresDB: func(r *repository.Repository) interface{} {
		return &Service{
			MessagePostgresService: &PostgresService{repository: r},
			Authorization:          &AuthorizationService{repository: r},
			ChatService:            &ChatMessageService{repository: r},
		}
	},
}

func NewTodoService(dbType string, db *repository.Repository) (*Service, error) {
	serviceFactory, ok := serviceFactories[dbType]
	if !ok {
		return nil, errors.New("unsupported database type")
	}

	return serviceFactory(db).(*Service), nil
}
