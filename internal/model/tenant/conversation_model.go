package tenant

import (
	"context"
	commonConst "moremei/ai-saas/pkg/consts"

	"gorm.io/gorm"
)

var _ ConversationModel = (*customConversationModel)(nil)

type (
	// ConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationModel.
	ConversationModel interface {
		conversationModel
		FindByUserId(ctx context.Context, userId uint64) ([]*UserConversation, error)
		GetAllByUserId(ctx context.Context, userId uint64) ([]*Conversation, error)
		FindByUserIdAiRobotIdIn(ctx context.Context, userId uint64, aiRobotId []uint64) ([]*Conversation, error)
		GetAll(ctx context.Context) (res []*Conversation, err error)
		FindByUserIdAiRobotId(ctx context.Context, userId uint64, aiRobotId uint64) (*Conversation, error)
		FindByIds(ctx context.Context, ids []uint64) ([]*Conversation, error)
	}

	customConversationModel struct {
		*defaultConversationModel
	}

	UserConversation struct {
		Conversation
		AiRobotUuid string
	}
)

// NewConversationModel returns a model for the database table.
func NewConversationModel(conn *gorm.DB) ConversationModel {
	return &customConversationModel{
		defaultConversationModel: newConversationModel(conn),
	}
}

func (m *customConversationModel) FindByUserId(ctx context.Context, userId uint64) ([]*UserConversation, error) {
	var resp []*UserConversation
	err := m.conn.WithContext(ctx).Model(&Conversation{}).
		Joins("INNER JOIN `ai_robot` ON `ai_robot`.`id` = `conversation`.`ai_robot_id`").
		Where("`conversation`.`user_id` = ?", userId).
		Select([]string{
			"`conversation`.`id`",
			"`conversation`.`uuid`",
			"`conversation`.`user_id`",
			"`conversation`.`ai_robot_id`",
			"`conversation`.`created_at`",
			"`conversation`.`updated_at`",
			"`ai_robot`.`uuid` AS `ai_robot_uuid`",
		}).
		Find(&resp).
		Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customConversationModel) GetAllByUserId(ctx context.Context, userId uint64) ([]*Conversation, error) {
	var resp []*Conversation
	err := m.conn.WithContext(ctx).Model(&Conversation{}).Where("user_id = ?", userId).Find(&resp).Error
	return resp, err
}

func (m *customConversationModel) FindByUserIdAiRobotIdIn(ctx context.Context, userId uint64, aiRobotId []uint64) (resp []*Conversation, err error) {
	err = m.conn.WithContext(ctx).Model(&Conversation{}).
		Where("`user_id` = ? and `ai_robot_id` in (?)", userId, aiRobotId).
		Order("id desc").
		Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customConversationModel) GetAll(ctx context.Context) (res []*Conversation, err error) {
	err = m.conn.WithContext(ctx).Model(&Conversation{}).
		Joins("LEFT JOIN `wecom_ai_robot` ON `wecom_ai_robot`.`ai_robot_id` = `conversation`.`ai_robot_id`").
		Joins("LEFT JOIN `wecom_external_user` ON `wecom_external_user`.`user_id` = `conversation`.`user_id` AND `wecom_external_user`.`wecom_ai_robot_id` = `wecom_ai_robot`.`id`").
		Where("`wecom_external_user`.`deleted` = ?", commonConst.FlagFalse).
		Select("`conversation`.*").
		Find(&res).Error
	return
}

func (m *customConversationModel) FindByUserIdAiRobotId(ctx context.Context, userId uint64, aiRobotId uint64) (resp *Conversation, err error) {
	err = m.conn.WithContext(ctx).Model(&Conversation{}).
		Where("`user_id` = ? and `ai_robot_id` = ?", userId, aiRobotId).Take(&resp).Error
	return
}

func (m *customConversationModel) FindByIds(ctx context.Context, ids []uint64) ([]*Conversation, error) {
	var resp []*Conversation
	err := m.conn.WithContext(ctx).Model(&Conversation{}).
		Where("`id` in (?)", ids).Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}
