package tenant

import (
	"context"

	"gorm.io/gorm"
)

var _ AiRobotWhitelistModel = (*customAiRobotWhitelistModel)(nil)

type (
	// AiRobotWhitelistModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiRobotWhitelistModel.
	AiRobotWhitelistModel interface {
		aiRobotWhitelistModel
		FindAiRobotIdsByUserId(ctx context.Context, userId uint64) ([]uint64, error)
		UpdateUserId(ctx context.Context, old, new uint64, tx *gorm.DB) error
	}

	customAiRobotWhitelistModel struct {
		*defaultAiRobotWhitelistModel
	}
)

// NewAiRobotWhitelistModel returns a model for the database table.
func NewAiRobotWhitelistModel(conn *gorm.DB) AiRobotWhitelistModel {
	return &customAiRobotWhitelistModel{
		defaultAiRobotWhitelistModel: newAiRobotWhitelistModel(conn),
	}
}

func (m *customAiRobotWhitelistModel) FindAiRobotIdsByUserId(ctx context.Context, userId uint64) ([]uint64, error) {
	var resp []uint64
	err := m.conn.WithContext(ctx).Model(&AiRobotWhitelist{}).
		Where("user_id = ?", userId).
		Select("ai_robot_id").
		Pluck("ai_robot_id", &resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customAiRobotWhitelistModel) UpdateUserId(ctx context.Context, old, new uint64, tx *gorm.DB) error {
	conn := m.conn
	if tx != nil {
		conn = tx
	}
	return conn.WithContext(ctx).Model(&AiRobotWhitelist{}).Where("user_id =?", old).Update("user_id", new).Error
}
