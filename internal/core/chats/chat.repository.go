package chats

import (
	"database/sql"
	"time"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllChats(*[]models.ChatPreviews, uuid.UUID) error
	GetMessagesInChat(*[]models.Messages, uuid.UUID, uuid.UUID, int, int) error
	SaveMessages(msg *models.Messages) error
	ReadMessages(sendUserId uuid.UUID, recvUserId uuid.UUID) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		db,
	}
}

func (repo *repositoryImpl) GetAllChats(results *[]models.ChatPreviews, userId uuid.UUID) error {
	_ = repo.db.Model(&models.Messages{}).
		Raw(`
		SELECT *
			FROM (
				SELECT
					DISTINCT ON (user_id) a.user_id,
					COALESCE(unread_messages, 0) AS unread_messages,
					content, sent_at
				FROM (
					SELECT 
						CASE
							WHEN sender_id = @user_id THEN receiver_id
							WHEN receiver_id = @user_id THEN sender_id 
						END AS user_id, message_id, content, read_at, sent_at
					FROM messages
					WHERE sender_id = @user_id OR receiver_id = @user_id
					ORDER BY sent_at DESC
				) AS a
				LEFT JOIN (
					SELECT
						CASE
							WHEN sender_id = @user_id THEN receiver_id
							WHEN receiver_id = @user_id THEN sender_id
						END AS user_id, count(*) AS unread_messages
					FROM messages
					WHERE receiver_id = @user_id AND read_at IS NULL
					GROUP BY user_id
				) AS b
				ON a.user_id = b.user_id
			) AS c
		ORDER BY unread_messages DESC, sent_at DESC
		`, sql.Named("user_id", userId)).
		Scan(results).Error
	return nil
}

func (repo *repositoryImpl) GetMessagesInChat(msgs *[]models.Messages, sendUserId uuid.UUID, recvUserId uuid.UUID, offset int, limit int) error {
	subQuery := repo.db.Model(&models.Messages{}).
		Order("sent_at DESC").
		Offset(offset).Limit(limit).
		Where("(sender_id = @sender_id AND receiver_id = @receiver_id) OR (sender_id = @receiver_id AND receiver_id = @sender_id)",
			sql.Named("sender_id", sendUserId),
			sql.Named("receiver_id", recvUserId))

	return repo.db.Select("*").
		Table("(?) AS a", subQuery).
		Order("sent_at ASC").
		Find(msgs).Error
}

func (repo *repositoryImpl) SaveMessages(msg *models.Messages) error {
	return repo.db.Model(&models.Messages{}).Create(msg).Error
}

func (repo *repositoryImpl) ReadMessages(sendUserId uuid.UUID, recvUserId uuid.UUID) error {
	return repo.db.Model(&models.Messages{}).
		Where("sender_id = ? AND receiver_id = ?", sendUserId, recvUserId).
		Update("read_at", time.Now()).Error
}
