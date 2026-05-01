package tenant

import (
	"context"

	"gorm.io/gorm"
)

var _ AiMessageModel = (*customAiMessageModel)(nil)

const (
	ThemeMessage                  = "theme_message"         // 主题消息
	SkinTestUrl                   = "skin_test_url"         // 测肤链接
	KnowLedgeMessage              = "knowledge_message"     // 小知识
	InviteMessage                 = "invite_message"        // 邀请信息
	ProductMessage                = "product_message"       // 产品信息
	AiMessageQuestionnaireMessage = "questionnaire_message" // 问卷消息
	AiMessageSkinTestMessage      = "skin_test_message"     // 测肤消息
	AiMessageQuestionnaireUrl     = "questionnaire_url"     // 问卷链接
	MsgKeyword                    = "message_keyword"       // 聊天关键词
	HoneybeeSkinTestKey           = "honeybee_skin_test"    // 设备测肤
)

type (
	// AiMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiMessageModel.
	AiMessageModel interface {
		aiMessageModel
		GetByStatus(ctx context.Context, status []uint64) ([]*AiMessage, error)
		GetByTypeStatus(ctx context.Context, types string, status uint64) ([]*AiMessage, error)
	}

	customAiMessageModel struct {
		*defaultAiMessageModel
	}
)

// NewAiMessageModel returns a model for the database table.
func NewAiMessageModel(conn *gorm.DB) AiMessageModel {
	return &customAiMessageModel{
		defaultAiMessageModel: newAiMessageModel(conn),
	}
}

func (m *customAiMessageModel) GetByStatus(ctx context.Context, status []uint64) (res []*AiMessage, err error) {
	err = m.conn.WithContext(ctx).Model(&AiMessage{}).Where("status in (?)", status).Find(&res).Error
	return
}

func (m *customAiMessageModel) GetByTypeStatus(ctx context.Context, types string, status uint64) (resp []*AiMessage, err error) {
	err = m.conn.WithContext(ctx).Model(&AiMessage{}).Where("type = ? and status = ?", types, status).Find(&resp).Error
	return
}
