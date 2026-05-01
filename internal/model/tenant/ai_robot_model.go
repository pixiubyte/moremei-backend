package tenant

import (
	"context"
	commonConst "moremei/ai-saas/pkg/consts"
	"moremei/ai-saas/pkg/pagination"

	"gorm.io/gorm"
)

var _ AiRobotModel = (*customAiRobotModel)(nil)

const (
	AiRobotCategoryAssistant = "assistant"
	AiRobotCategoryProfessor = "professor"
	AiRobotCategoryTool      = "tool"

	AiRobotLlmTypeOpenai  = "openai"
	AiRobotLlmTypeChatGLM = "chatglm"

	AiRobotSceneDefault     = "chats"
	AiRobotSceneHealthGroup = "health_group"
	AiRobotSceneInstore     = "in_store"
)

type (
	// AiRobotModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiRobotModel.
	AiRobotModel interface {
		aiRobotModel
		FindDefaultOne(ctx context.Context) (*AiRobot, error)
		FindByStatus(ctx context.Context, status uint64) ([]*AiRobot, error)
		FindByIds(ctx context.Context, ids []uint64) ([]*AiRobot, error)
		FindEarliestOneByCategory(ctx context.Context, category string) (*AiRobot, error)
		PageFind(ctx context.Context, pageReq *pagination.PageRequest) (pageResp *pagination.PageResponse, err error)
		GetAll(ctx context.Context) ([]*AiRobot, error)
		FindByScene(ctx context.Context, scene string) ([]*AiRobot, error)
		FindBySceneDefault(ctx context.Context, scene string) (*AiRobot, error)
	}

	customAiRobotModel struct {
		*defaultAiRobotModel
	}
)

// NewAiRobotModel returns a model for the database table.
func NewAiRobotModel(conn *gorm.DB) AiRobotModel {
	return &customAiRobotModel{
		defaultAiRobotModel: newAiRobotModel(conn),
	}
}

func (m *customAiRobotModel) FindDefaultOne(ctx context.Context) (*AiRobot, error) {
	var resp AiRobot
	err := m.conn.WithContext(ctx).Model(&AiRobot{}).
		Where("is_default = ?", commonConst.FlagTrue).
		First(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customAiRobotModel) FindByStatus(ctx context.Context, status uint64) ([]*AiRobot, error) {
	var resp []*AiRobot
	err := m.conn.WithContext(ctx).Model(&AiRobot{}).
		Where("status = ?", status).
		Find(&resp).Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *customAiRobotModel) FindByIds(ctx context.Context, ids []uint64) ([]*AiRobot, error) {
	var resp []*AiRobot
	err := m.conn.WithContext(ctx).Model(&AiRobot{}).
		Where("id in (?)", ids).
		Find(&resp).Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *customAiRobotModel) FindEarliestOneByCategory(ctx context.Context, category string) (*AiRobot, error) {
	var resp AiRobot
	err := m.conn.WithContext(ctx).Model(&AiRobot{}).Where("category = ?", category).Order("id ASC").First(&resp).Error
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (m *customAiRobotModel) PageFind(ctx context.Context, pageReq *pagination.PageRequest) (pageResp *pagination.PageResponse, err error) {
	var (
		total  int64
		models []*AiRobot
	)

	query := m.conn.WithContext(ctx).Model(&AiRobot{}).Order("id ASC")

	query.Count(&total)

	err = query.Scopes(pagination.Paginate(pageReq)).Find(&models).Error

	if err != nil {
		return nil, err
	}

	return pagination.NewPageResponse(total, pageReq, models), nil
}

func (m *customAiRobotModel) GetAll(ctx context.Context) (res []*AiRobot, err error) {
	err = m.conn.WithContext(ctx).Model(&AiRobot{}).Find(&res).Error
	return
}

func (m *customAiRobotModel) FindByScene(ctx context.Context, scene string) (res []*AiRobot, err error) {
	err = m.conn.WithContext(ctx).Model(&AiRobot{}).Where("status =? AND scene =?", commonConst.FlagTrue, scene).Find(&res).Error
	return
}

func (m *customAiRobotModel) FindBySceneDefault(ctx context.Context, scene string) (res *AiRobot, err error) {
	err = m.conn.WithContext(ctx).Model(&AiRobot{}).Where("status =? AND scene =? AND is_default = ?", commonConst.FlagTrue, scene, commonConst.FlagTrue).First(&res).Error
	return
}
