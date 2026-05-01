package tenant

import (
	"context"

	"gorm.io/gorm"
)

var _ SharedChatHistoryModel = (*customSharedChatHistoryModel)(nil)

type (
	// SharedChatHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSharedChatHistoryModel.
	SharedChatHistoryModel interface {
		sharedChatHistoryModel
		FindByUUID(ctx context.Context, uuid string) (*SharedChatHistory, error)
	}

	customSharedChatHistoryModel struct {
		*defaultSharedChatHistoryModel
	}
)

// NewSharedChatHistoryModel returns a model for the database table.
func NewSharedChatHistoryModel(conn *gorm.DB) SharedChatHistoryModel {
	return &customSharedChatHistoryModel{
		defaultSharedChatHistoryModel: newSharedChatHistoryModel(conn),
	}
}

func (m *customSharedChatHistoryModel) FindByUUID(ctx context.Context, uuid string) (resp *SharedChatHistory, err error) {
	err = m.conn.WithContext(ctx).Where("uuid = ?", uuid).First(&resp).Error
	return
}
