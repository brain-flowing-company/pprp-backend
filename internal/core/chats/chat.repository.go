package chats

import (
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllChats(*[]models.ChatsResponses, uuid.UUID) error
	GetMessagesInChat(*[]models.Messages, uuid.UUID, uuid.UUID, int, int) error
	CreateMessages(msg *models.Messages) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetAllChats(results *[]models.ChatsResponses, userId uuid.UUID) error {
	return nil
}

func (repo *repositoryImpl) GetMessagesInChat(msgs *[]models.Messages, sendUserId uuid.UUID, recvUserId uuid.UUID, offset int, limit int) error {
	subQuery := repo.db.Model(&models.Messages{}).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			sendUserId, recvUserId, recvUserId, sendUserId)

	return repo.db.Select("*").
		Table("(?) AS a", subQuery).
		Order("created_at ASC").
		Find(msgs).Error
}

func (repo *repositoryImpl) CreateMessages(msg *models.Messages) error {
	return repo.db.Model(&models.Messages{}).Create(msg).Error
}
