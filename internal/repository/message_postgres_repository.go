package repository

import (
	"chat/internal/models"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

type MessagePostgres struct {
	db *sql.DB
}

func NewMessagePostgres(db *sql.DB) *MessagePostgres {
	return &MessagePostgres{db: db}
}

func (u *MessagePostgres) GetMessageByID(id int) (*models.Message, error) {
	var message models.Message
	result := u.db.QueryRow("SELECT id, username, content FROM messages WHERE id = $1", id)
	if err := result.Scan(&message.ID, &message.Username, &message.Content); err != nil {
		logrus.Errorf("GetMessageByID: error while scanning for message:%s", err)
		return nil, fmt.Errorf("GetMessageByID: repository error:%w", err)
	}
	return &message, nil
}

func (u *MessagePostgres) GetMessages(page, limit int64) ([]models.Message, int, error) {
	transaction, err := u.db.Begin()
	if err != nil {
		logrus.Errorf("GetMessages: can not starts transaction:%s", err)
		return nil, 0, fmt.Errorf("GetMessages: can not starts transaction:%w", err)
	}
	var Messages []models.Message
	var query string
	var pages int
	var rows *sql.Rows
	if page == 0 || limit == 0 {
		query = "SELECT id, username, content FROM messages ORDER BY id"
		rows, err = transaction.Query(query)
		if err != nil {
			logrus.Errorf("GetMessages: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("GetMessages:repository error:%w", err)
		}
		pages = 1
	} else {
		query = "SELECT id, username, content FROM messages ORDER BY id LIMIT $1 OFFSET $2"
		rows, err = transaction.Query(query, limit, (page-1)*limit)
		if err != nil {
			logrus.Errorf("GetMessages: can not executes a query:%s", err)
			return nil, 0, fmt.Errorf("GetMessages:repository error:%w", err)
		}
	}
	for rows.Next() {
		var Message models.Message
		if err := rows.Scan(&Message.ID, &Message.Username, &Message.Content); err != nil {
			logrus.Errorf("Error while scanning for message:%s", err)
			return nil, 0, fmt.Errorf("GetMessages:repository error:%w", err)
		}
		Messages = append(Messages, Message)
	}
	if pages == 0 {
		query = "SELECT CEILING(COUNT(id)/$1::float) FROM messages"
		row := transaction.QueryRow(query, limit)
		if err := row.Scan(&pages); err != nil {
			logrus.Errorf("Error while scanning for pages:%s", err)
		}
	}
	return Messages, pages, transaction.Commit()
}

func (u *MessagePostgres) CreateMessage(message *models.Message) (int, error) {
	var id int
	row := u.db.QueryRow("INSERT INTO messages (username, content) VALUES ($1, $2) RETURNING id", message.Username, message.Content)
	if err := row.Scan(&id); err != nil {
		logrus.Errorf("CreateMessage: error while scanning for messages:%s", err)
		return 0, fmt.Errorf("CreateMessage: error while scanning for messages:%w", err)
	}
	return id, nil
}

func (u *MessagePostgres) UpdateMessage(message *models.Message) error {
	_, err := u.db.Exec("UPDATE messages SET username = $1, content = $2 WHERE id = $3", message.Username, message.Content, message.ID)
	if err != nil {
		logrus.Errorf("UpdateMessage: error while updating messages:%s", err)
		return fmt.Errorf("UpdateMessage: error while updating messages:%w", err)
	}
	return nil
}

func (u *MessagePostgres) DeleteMessageByID(id int) (int, error) {
	var messageId int
	row := u.db.QueryRow("DELETE FROM messages WHERE id=$1 RETURNING id", id)
	if err := row.Scan(&messageId); err != nil {
		logrus.Errorf("DeleteMessageByID: error while scanning for messages:%s", err)
		return 0, fmt.Errorf("DeleteMessageByID: error while scanning for messages:%w", err)
	}
	return messageId, nil
}
