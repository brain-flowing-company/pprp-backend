package chats

import (
	"errors"

	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/brain-flowing-company/pprp-backend/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetAllChats(*[]models.ChatsResponses, uuid.UUID) *apperror.AppError
	GetMessagesInChat(*[]models.Messages, uuid.UUID, uuid.UUID, int, int) *apperror.AppError
	JoinChat(uuid.UUID, uuid.UUID) *apperror.AppError
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

func (s *serviceImpl) GetAllChats(chats *[]models.ChatsResponses, userId uuid.UUID) *apperror.AppError {
	err := s.repo.GetAllChats(chats, userId)
	if err != nil {
		s.logger.Error("Could not get all chats", zap.Error(err), zap.String("userId", userId.String()))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not get all chats")
	}

	return nil
}

func (s *serviceImpl) GetMessagesInChat(msgs *[]models.Messages, sendUserId uuid.UUID, recvUserId uuid.UUID, offset int, limit int) *apperror.AppError {
	limit = utils.Min(limit, 50)

	err := s.repo.GetMessagesInChat(msgs, sendUserId, recvUserId, offset, limit)
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

func (s *serviceImpl) JoinChat(sendUserId uuid.UUID, recvUserId uuid.UUID) *apperror.AppError {
	if sendUserId == recvUserId {
		return apperror.
			New(apperror.BadRequest).
			Describe("Could not chat with yourself")
	}

	err := s.repo.CreateChatStatus(sendUserId, recvUserId)
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return apperror.
			New(apperror.BadRequest).
			Describe("Could not find the specified recvUserId")
	} else if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		s.logger.Error("Could not create chat status",
			zap.Error(err),
			zap.String("senderUserId", sendUserId.String()),
			zap.String("receiveruserId", recvUserId.String()))
		return apperror.
			New(apperror.InternalServerError).
			Describe("Could not join chat")
	}

	return nil
}
