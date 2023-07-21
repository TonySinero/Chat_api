package repository

import (
	"chat/internal/models"
	"context"
	"database/sql"
	"errors"
)

type AppMessagePostgres interface {
	GetMessageByID(id int) (*models.Message, error)
	GetMessages(page, limit int64) ([]models.Message, int, error)
	CreateMessage(message *models.Message) (int, error)
	UpdateMessage(message *models.Message) error
	DeleteMessageByID(id int) (int, error)
}

type AuthorizationApp interface {
	CreateUser(ctx context.Context, user *models.User) error
	CheckByEmail(ctx context.Context, restore *models.RestorePassword) error
	UserById(ctx context.Context, userID int) (*models.ResponseUser, error)
	UserByPhone(ctx context.Context, user *models.User) (*models.User, error)
	Users(ctx context.Context, page, limit int64) ([]models.ResponseUser, error)
	UpdateUser(ctx context.Context, inputUser *models.ResponseUser) error
	DeleteUser(ctx context.Context, userID int) error
	UserRoleById(userId int) (*models.User, error)
	RestorePassword(ctx context.Context, restore *models.RestorePassword) error
}

type GroupChatRepository interface {
	CreateGroupChat(groupChat *models.GroupChat) (*models.GroupChat, error)
	AddMember(groupChatID, userID string) error
	CreateMessage(message *models.GroupMessage) (*models.GroupMessage, error)
	IsMember(groupChatID, userID string) bool
	CanSendMessage(senderID, recipientID string) bool
}

const (
	PostgresDB string = "postgres"
)

type Repository struct {
	AppMessagePostgres
	AuthorizationApp
	GroupChatRepository
}

func NewRepository(dbType string, db interface{}) (*Repository, error) {
	switch dbType {
	case "postgres":
		PostgresDB, ok := db.(*sql.DB)
		if !ok {
			return nil, errors.New("invalid database postgres connection")
		}
		return &Repository{
			AppMessagePostgres:  NewMessagePostgres(PostgresDB),
			AuthorizationApp:    NewAuthRepository(PostgresDB),
			GroupChatRepository: NewGroupChatRepository(PostgresDB),
		}, nil
	default:
		return nil, errors.New("unsupported database type")
	}
}
