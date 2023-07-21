package service

import (
	"chat/internal/models"
	"chat/internal/repository"
	"fmt"
)

type PostgresService struct {
	repository *repository.Repository
}

func (t *PostgresService) GetMessage(id int) (*models.Message, error) {
	user, err := t.repository.AppMessagePostgres.GetMessageByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (t *PostgresService) GetMessages(page, limit int64) ([]models.Message, int, error) {
	users, pages, err := t.repository.AppMessagePostgres.GetMessages(page, limit)
	if err != nil {
		return nil, 0, err
	}
	return users, pages, nil
}

func (t *PostgresService) CreateMessage(message *models.Message) (int, error) {
	id, err := t.repository.AppMessagePostgres.CreateMessage(message)
	if err != nil {
		return 0, fmt.Errorf("something went wrong when creating a user:%w", err)
	}
	return id, nil
}

func (t *PostgresService) UpdateMessage(message *models.Message) error {
	err := t.repository.AppMessagePostgres.UpdateMessage(message)
	if err != nil {
		return err
	}
	return nil
}

func (t *PostgresService) DeleteMessageByID(id int) (int, error) {
	Id, err := t.repository.AppMessagePostgres.DeleteMessageByID(id)
	if err != nil {
		return 0, err
	}
	return Id, nil
}
