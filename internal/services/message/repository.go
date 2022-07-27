package message

import (
	"context"
	"dialog/internal/database/entities"
	"dialog/pb"
	"dialog/pb/messages"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Repository struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func (r Repository) GetAllMessages() (*[]entities.Message, error) {
	var messages []entities.Message
	results := r.DB.Find(&messages)

	return &messages, results.Error
}

func (r Repository) GetMessageById(id string) (*entities.Message, error) {
	var message entities.Message
	result := r.DB.Take(&message, "uid = ?", id)

	return &message, result.Error
}

func (r Repository) CreateMessage(message *entities.Message) (err error) {
	err = r.DB.Create(message).Error

	if err != nil {
		return
	}

	notification, err := json.Marshal(messages.MessageMutateEvent{
		Body: message.Model(),
		Type: pb.MutateEventType_CREATE,
		Id:   message.Uid,
	})

	if err != nil {
		return
	}

	err = r.RedisClient.Publish(context.Background(), MutateChannel, notification).Err()

	return
}

func (r Repository) UpdateMessage(message *entities.Message) (err error) {
	err = r.DB.Updates(message).Error

	if err != nil {
		return
	}

	notification, err := json.Marshal(messages.MessageMutateEvent{
		Body: message.Model(),
		Type: pb.MutateEventType_UPDATE,
		Id:   message.Uid,
	})

	if err != nil {
		return
	}

	err = r.RedisClient.Publish(context.Background(), MutateChannel, notification).Err()

	return
}

func (r Repository) DeleteMessage(message *entities.Message) (err error) {
	err = r.DB.Delete(message).Error

	if err != nil {
		return
	}

	notification, err := json.Marshal(messages.MessageMutateEvent{
		Type: pb.MutateEventType_DESTROY,
		Id:   message.Uid,
	})

	if err != nil {
		return
	}

	err = r.RedisClient.Publish(context.Background(), MutateChannel, notification).Err()

	return
}
