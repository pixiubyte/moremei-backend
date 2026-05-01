package tenant

import (
	"context"
	"moremei/ai-saas/pkg/pagination"

	"gorm.io/gorm"
)

var _ MiniUserModel = (*customMiniUserModel)(nil)

type (
	// MiniUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMiniUserModel.
	MiniUserModel interface {
		miniUserModel
		FindOneByUnionid(ctx context.Context, unionid string) (*MiniUser, error)
		GetAll(ctx context.Context) ([]*MiniUser, error)
		FindByOfficalOpenId(ctx context.Context) ([]*MiniUser, error)
		UpdateUserId(ctx context.Context, old, new uint64, tx *gorm.DB) error
		PageFind(ctx context.Context, pageReq *pagination.PageRequest) (*pagination.PageResponse, error)
	}

	customMiniUserModel struct {
		*defaultMiniUserModel
	}
)

// NewMiniUserModel returns a model for the database table.
func NewMiniUserModel(conn *gorm.DB) MiniUserModel {
	return &customMiniUserModel{
		defaultMiniUserModel: newMiniUserModel(conn),
	}
}

func (m *customMiniUserModel) FindOneByUnionid(ctx context.Context, unionid string) (resp *MiniUser, err error) {
	err = m.conn.WithContext(ctx).Model(&MiniUser{}).Where("unionid = ?", unionid).First(&resp).Error
	return
}

func (m *customMiniUserModel) GetAll(ctx context.Context) (resp []*MiniUser, err error) {
	err = m.conn.WithContext(ctx).Model(&MiniUser{}).Find(&resp).Error
	return
}

func (m *customMiniUserModel) FindByOfficalOpenId(ctx context.Context) (resp []*MiniUser, err error) {
	err = m.conn.WithContext(ctx).Model(&MiniUser{}).Where("official_openid != ?", "").Find(&resp).Error
	return
}

func (m *customMiniUserModel) UpdateUserId(ctx context.Context, old, new uint64, tx *gorm.DB) error {
	conn := m.conn
	if tx != nil {
		conn = tx
	}
	return conn.WithContext(ctx).Model(&MiniUser{}).Where("user_id =?", old).Update("user_id", new).Error
}

func (m *customMiniUserModel) PageFind(ctx context.Context, pageReq *pagination.PageRequest) (*pagination.PageResponse, error) {
	var total int64
	var data []*MiniUser
	query := m.conn.WithContext(ctx).Model(&MiniUser{}).Order("`id` asc")

	query.Count(&total)

	query = query.Scopes(pagination.Paginate(pageReq))
	err := query.Find(&data).Error
	return pagination.NewPageResponse(total, pageReq, data), err
}
