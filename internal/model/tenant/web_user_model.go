package tenant

import (
	"context"

	commonConst "moremei/ai-saas/pkg/consts"

	"gorm.io/gorm"
)

var FromMap = map[string]int64{
	"maomei_miniprogram": 1,
}

var _ WebUserModel = (*customWebUserModel)(nil)

type (
	// WebUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWebUserModel.
	WebUserModel interface {
		webUserModel
		FindOneByMobile(ctx context.Context, mobile string) (*WebUser, error)
		FindOneByUserId(ctx context.Context, userId uint64) (*WebUser, error)
		UpdateByUserId(ctx context.Context, userId uint64, data WebUser) error
		FindByMobile(ctx context.Context, mobile []string) ([]*WebUser, error)
		FindOneByUserIdNoDel(ctx context.Context, userId uint64) (*WebUser, error)
		FindByUserIds(ctx context.Context, userIds []uint64) ([]*WebUser, error)
	}

	customWebUserModel struct {
		*defaultWebUserModel
	}
)

// NewWebUserModel returns a model for the database table.
func NewWebUserModel(conn *gorm.DB) WebUserModel {
	return &customWebUserModel{
		defaultWebUserModel: newWebUserModel(conn),
	}
}

func (m *customWebUserModel) FindOneByMobile(ctx context.Context, mobile string) (*WebUser, error) {
	var resp WebUser
	err := m.conn.WithContext(ctx).Model(&WebUser{}).Where("mobile = ?", mobile).
		Where("deleted = ?", commonConst.FlagFalse).First(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customWebUserModel) FindOneByUserId(ctx context.Context, userId uint64) (*WebUser, error) {
	var resp WebUser
	err := m.conn.WithContext(ctx).Model(&WebUser{}).Where("user_id = ?", userId).
		Where("deleted = ?", commonConst.FlagFalse).First(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customWebUserModel) UpdateByUserId(ctx context.Context, userId uint64, data WebUser) error {
	err := m.conn.WithContext(ctx).Model(&WebUser{}).Where("user_id = ? ", userId).Updates(data).Error
	return err
}

func (m *customWebUserModel) FindByMobile(ctx context.Context, mobile []string) (res []*WebUser, err error) {
	err = m.conn.WithContext(ctx).Where("mobile in (?)", mobile).Find(&res).Error
	return
}

func (m *customWebUserModel) FindOneByUserIdNoDel(ctx context.Context, userId uint64) (*WebUser, error) {
	var resp WebUser
	err := m.conn.WithContext(ctx).Model(&WebUser{}).Where("user_id = ?", userId).
		First(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customWebUserModel) FindByUserIds(ctx context.Context, userIds []uint64) (resp []*WebUser, err error) {
	err = m.conn.WithContext(ctx).Model(&WebUser{}).Where("user_id in (?)", userIds).Find(&resp).Error
	return
}
