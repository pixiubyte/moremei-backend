package tenant

import (
	"context"

	"gorm.io/gorm"
)

var _ AuthProviderUserModel = (*customAuthProviderUserModel)(nil)

type (
	// AuthProviderUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuthProviderUserModel.
	AuthProviderUserModel interface {
		authProviderUserModel
		FindOneByProviderUserid(ctx context.Context, provider string, userid uint64) (*AuthProviderUser, error)
		FindOneByProviderAppUnionid(ctx context.Context, provider, app, unionid string) (*AuthProviderUser, error)
	}

	customAuthProviderUserModel struct {
		*defaultAuthProviderUserModel
	}
)

// NewAuthProviderUserModel returns a model for the database table.
func NewAuthProviderUserModel(conn *gorm.DB) AuthProviderUserModel {
	return &customAuthProviderUserModel{
		defaultAuthProviderUserModel: newAuthProviderUserModel(conn),
	}
}

func (m *customAuthProviderUserModel) FindOneByProviderUserid(ctx context.Context, provider string, userid uint64) (*AuthProviderUser, error) {
	var resp AuthProviderUser
	err := m.conn.WithContext(ctx).Model(&AuthProviderUser{}).Where("`provider` = ? and `userid` = ?", provider, userid).Take(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customAuthProviderUserModel) FindOneByProviderAppUnionid(ctx context.Context, provider, app, unionid string) (*AuthProviderUser, error) {
	var resp AuthProviderUser
	err := m.conn.WithContext(ctx).Model(&AuthProviderUser{}).Where("`provider` = ? and `app` = ? and `unionid` = ?", provider, app, unionid).Take(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
