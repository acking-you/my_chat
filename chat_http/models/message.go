package models

import "gorm.io/gorm"

// TextMessage 文字消息
type TextMessage struct {
	gorm.Model
	Sender   int64
	Receiver int64
	Text     string
}

type TextMessageDAO struct {
}

var textDAO TextMessageDAO

func NewTextMessageDAO() TextMessageDAO {
	return textDAO
}

func (t TextMessage) AddMessage(sender, receiver int64, text string) error {
	err := DB.Create(&TextMessage{
		Sender:   sender,
		Receiver: receiver,
		Text:     text,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (t TextMessage) RemoveMessage(message TextMessage) error {
	err := DB.Delete(&message).Error
	if err != nil {
		return err
	}
	return nil
}

func (t TextMessage) GetMessageFlow(userId, userToId int64, limit int) (*[]TextMessage, error) {
	ret := make([]TextMessage, limit)
	err := DB.Model(&TextMessage{}).
		Where("(sender = ? and receiver = ?) or (sender=? and receiver=?)", userId, userToId, userToId, userId).
		Limit(limit).
		Find(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
