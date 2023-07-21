package service

import (
	"chat/internal/models"
	"chat/internal/repository"
	"time"

	"github.com/pkg/errors"
)

type ChatMessageService struct {
	repository *repository.Repository
}

func (s *ChatMessageService) CreateGroupChat(input models.GroupChatCreateInput) (*models.GroupChat, error) {
	groupChat := &models.GroupChat{
		Name:        input.Name,
		Description: input.Description,
		Members:     input.Members,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdGroupChat, err := s.repository.GroupChatRepository.CreateGroupChat(groupChat)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create group chat")
	}

	return createdGroupChat, nil
}

func (s *ChatMessageService) JoinGroupChat(userID, groupChatID string) error {
	err := s.repository.GroupChatRepository.AddMember(groupChatID, userID)
	if err != nil {
		return errors.Wrap(err, "failed to join group chat")
	}

	return nil
}

func (s *ChatMessageService) SendMessageToGroupChat(userID, groupChatID string, input models.MessageCreateInput) (*models.GroupMessage, error) {
	message := &models.GroupMessage{
		Content:     input.Content,
		SenderID:    userID,
		GroupChatID: groupChatID,
		CreatedAt:   time.Now(),
	}

	createdMessage, err := s.repository.GroupChatRepository.CreateMessage(message)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send message to group chat")
	}

	return createdMessage, nil
}

func (s *ChatMessageService) SendPrivateMessage(senderID, recipientID string, input models.MessageCreateInput) (*models.GroupMessage, error) {
	message := &models.GroupMessage{
		Content:     input.Content,
		SenderID:    senderID,
		RecipientID: recipientID,
		CreatedAt:   time.Now(),
	}

	createdMessage, err := s.repository.GroupChatRepository.CreateMessage(message)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send private message")
	}

	return createdMessage, nil
}

func (s *ChatMessageService) CanUserAccessGroupChat(userID, groupChatID string) bool {
	return s.repository.GroupChatRepository.IsMember(groupChatID, userID)
}

func (s *ChatMessageService) CanUserSendMessage(senderID, recipientID string) bool {
	return s.repository.GroupChatRepository.CanSendMessage(senderID, recipientID)
}

func (s *ChatMessageService) AddUserToGroupChat(userID, groupChatID string) error {
	err := s.repository.GroupChatRepository.AddMember(groupChatID, userID)
	if err != nil {
		return errors.Wrap(err, "failed to add user to group chat")
	}

	return nil
}
