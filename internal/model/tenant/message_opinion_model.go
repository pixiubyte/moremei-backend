package tenant

import (
	"context"

	"gorm.io/gorm"
)

const (
	MessageOpinionStatusThumbsDown = -1
	MessageOpinionStatusCancel     = 0
	MessageOpinionStatusThumbsUp   = 1
)

var _ MessageOpinionModel = (*customMessageOpinionModel)(nil)

type (
	// MessageOpinionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageOpinionModel.
	MessageOpinionModel interface {
		messageOpinionModel
		FindByUserIdMsgIds(ctx context.Context, userId uint64, msgIds []uint64) ([]*MessageOpinion, error)
		FindByUserId(ctx context.Context, userId uint64) ([]*MessageOpinion, error)
		UpdateByUserId(ctx context.Context, old uint64, new uint64) error
		FindByMessageIds(ctx context.Context, messageIds []uint64) ([]*MessageOpinion, error)
		FindByMessageIdAndUserId(ctx context.Context, messageId uint64, userId uint64) (*MessageOpinion, error)
	}

	customMessageOpinionModel struct {
		*defaultMessageOpinionModel
	}
)

// NewMessageOpinionModel returns a model for the database table.
func NewMessageOpinionModel(conn *gorm.DB) MessageOpinionModel {
	return &customMessageOpinionModel{
		defaultMessageOpinionModel: newMessageOpinionModel(conn),
	}
}

func (m *customMessageOpinionModel) FindByUserIdMsgIds(ctx context.Context, userId uint64, msgIds []uint64) ([]*MessageOpinion, error) {
	var resp []*MessageOpinion

	err := m.conn.WithContext(ctx).Model(&MessageOpinion{}).
		Where("user_id = ?", userId).
		Where("message_id in (?)", msgIds).
		Find(&resp).Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *customMessageOpinionModel) FindByUserId(ctx context.Context, userId uint64) (res []*MessageOpinion, err error) {
	err = m.conn.WithContext(ctx).Model(&MessageOpinion{}).Where("user_id = ?", userId).Find(&res).Error
	return
}

func (m *customMessageOpinionModel) UpdateByUserId(ctx context.Context, old uint64, new uint64) error {
	return m.conn.WithContext(ctx).Model(&MessageOpinion{}).
		Where("user_id = ?", old).
		Update("user_id", new).Error
}

func (m *customMessageOpinionModel) FindByMessageIds(ctx context.Context, messageIds []uint64) ([]*MessageOpinion, error) {
	var resp []*MessageOpinion

	err := m.conn.WithContext(ctx).Model(&MessageOpinion{}).
		Where("message_id in (?)", messageIds).
		Find(&resp).Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *customMessageOpinionModel) FindByMessageIdAndUserId(ctx context.Context, messageId uint64, userId uint64) (resp *MessageOpinion, err error) {
	err = m.conn.WithContext(ctx).Model(&MessageOpinion{}).
		Where("message_id = ? AND user_id = ?", messageId, userId).
		First(&resp).Error
	return
}
