package chats

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	GetAllChats(*[]models.ChatPreviews, uuid.UUID) *apperror.AppError
	GetMessagesInChat(*[]models.Messages, uuid.UUID, uuid.UUID, *models.PaginatedQuery) *apperror.AppError
	SaveMessages(*models.Messages) *apperror.AppError
	ReadMessages(uuid.UUID, uuid.UUID) *apperror.AppError
}

type serviceImpl struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(logger *zap.Logger, repo Repository) Service {
	return &serviceImpl{
		repo,
		logger,
	}
}

func (s *serviceImpl) GetAllChats(chats *[]models.ChatPreviews, userId uuid.UUID) *apperror.AppError {
	err := s.repo.GetAllChats(chats, userId)
	if err != nil {
		s.logger.Error("Could not get all chats", zap.Error(err), zap.String("userId", userId.String()))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all chats")
	}

	return nil
}

func (s *serviceImpl) GetMessagesInChat(msgs *[]models.Messages, sendUserId uuid.UUID, recvUserId uuid.UUID, paginated *models.PaginatedQuery) *apperror.AppError {
	err := s.repo.GetMessagesInChat(msgs, sendUserId, recvUserId, paginated)
	if err != nil {
		s.logger.Error("Could not get messages in chat",
			zap.Error(err),
			zap.String("senderUserId", sendUserId.String()),
			zap.String("receiveruserId", recvUserId.String()))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get messages in chat")
	}

	return nil
}

func (s *serviceImpl) SaveMessages(msg *models.Messages) *apperror.AppError {
	err := s.repo.SaveMessages(msg)
	if err != nil {
		s.logger.Error("Could not save message",
			zap.Error(err))
		return apperror.
			New(apperror.InternalServerError).
			Describe("error while sending message")
	}

	return nil
}

func (s *serviceImpl) ReadMessages(sendUserId uuid.UUID, recvUserId uuid.UUID) *apperror.AppError {
	err := s.repo.ReadMessages(sendUserId, recvUserId)
	if err != nil {
		s.logger.Error("Could not update mesages read status",
			zap.Error(err),
			zap.String("senderUserId", sendUserId.String()),
			zap.String("receiveruserId", recvUserId.String()))
		return apperror.
			New(apperror.InternalServerError).
			Describe("error while opening chat")
	}

	return nil
}
