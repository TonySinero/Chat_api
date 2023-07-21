package repository

import (
	"chat/internal/models"
	"database/sql"

	"github.com/pkg/errors"
)

type ChatPostgres struct {
	db *sql.DB
}

func NewGroupChatRepository(db *sql.DB) *ChatPostgres {
	return &ChatPostgres{db: db}
}

func (r *ChatPostgres) CreateGroupChat(groupChat *models.GroupChat) (*models.GroupChat, error) {
	query := `
		INSERT INTO group_chats (name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, created_at, updated_at
	`
	var createdGroupChat models.GroupChat
	err := r.db.QueryRow(query, groupChat.Name, groupChat.Description, groupChat.CreatedAt, groupChat.UpdatedAt).Scan(
		&createdGroupChat.ID,
		&createdGroupChat.Name,
		&createdGroupChat.Description,
		&createdGroupChat.CreatedAt,
		&createdGroupChat.UpdatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create group chat")
	}

	return &createdGroupChat, nil
}

func (r *ChatPostgres) AddMember(groupChatID, userID string) error {
	query := `
		INSERT INTO group_chat_members (group_chat_id, user_id)
		VALUES ($1, $2)
	`
	_, err := r.db.Exec(query, groupChatID, userID)
	if err != nil {
		return errors.Wrap(err, "failed to add member to group chat")
	}

	return nil
}

func (r *ChatPostgres) CreateMessage(message *models.GroupMessage) (*models.GroupMessage, error) {
	query := `
		INSERT INTO group_chat_messages (sender_id, group_chat_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, sender_id, group_chat_id, content, created_at
	`
	var createdMessage models.GroupMessage
	err := r.db.QueryRow(query, message.SenderID, message.GroupChatID, message.Content, message.CreatedAt).Scan(
		&createdMessage.ID,
		&createdMessage.SenderID,
		&createdMessage.GroupChatID,
		&createdMessage.Content,
		&createdMessage.CreatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create group message")
	}

	return &createdMessage, nil
}

func (r *ChatPostgres) IsMember(groupChatID, userID string) bool {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM group_chat_members
			WHERE group_chat_id = $1 AND user_id = $2
		)
	`
	var isMember bool
	err := r.db.QueryRow(query, groupChatID, userID).Scan(&isMember)
	if err != nil {
		return false
	}

	return isMember
}

func (r *ChatPostgres) CanSendMessage(senderID, recipientID string) bool {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM group_chat_members m1
			JOIN group_chat_members m2 ON m1.group_chat_id = m2.group_chat_id
			WHERE m1.user_id = $1 AND m2.user_id = $2
		)
	`
	var canSendMessage bool
	err := r.db.QueryRow(query, senderID, recipientID).Scan(&canSendMessage)
	if err != nil {
		return false
	}

	return canSendMessage
}
