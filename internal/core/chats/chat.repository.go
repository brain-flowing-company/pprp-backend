package chats

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllChats(*[]models.ChatsResponses, uuid.UUID) error
	GetMessagesInChat(*[]models.Messages, uuid.UUID, uuid.UUID, int, int) error
	CreateChatStatus(uuid.UUID, uuid.UUID) error
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
	unreadQuery := repo.db.
		Select("messages.sender_id, SUM((NOT messages.read)::INT) AS unread_count").
		Table("chat_status").
		Joins("LEFT JOIN messages ON chat_status.receiver_id = messages.receiver_id AND chat_status.sender_id = messages.sender_id").
		Group("messages.sender_id").
		Where("messages.receiver_id = ? AND messages.created_at >= chat_status.last_active_at", userId)

	latestMessages := repo.db.
		Select("DISTINCT ON (messages.sender_id) sender_id, messages.CONTENT, messages.created_at").
		Table("messages").
		Order("sender_id, created_at DESC")

	return repo.db.
		Select("*").
		Table("(?) AS a", unreadQuery).
		Joins("LEFT JOIN (?) AS b ON a.sender_id = b.sender_id", latestMessages).
		Order("created_at DESC").
		Find(results).Error
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

func (repo *repositoryImpl) CreateChatStatus(sendUserId uuid.UUID, recvUserId uuid.UUID) error {
	now := time.Now()
	status := []models.ChatStatus{
		{
			SenderId:     sendUserId,
			ReceiverId:   recvUserId,
			LastActiveAt: now,
		},
		{
			SenderId:     recvUserId,
			ReceiverId:   sendUserId,
			LastActiveAt: now,
		},
	}
	return repo.db.CreateInBatches(status, 2).Error
}

func (repo *repositoryImpl) CreateMessages(msg *models.Messages) error {
	return repo.db.Model(&models.Messages{}).Create(msg).Error
}
