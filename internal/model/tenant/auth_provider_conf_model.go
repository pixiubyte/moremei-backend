package tenant

import (
	"context"

	"gorm.io/gorm"
)

var _ AuthProviderConfModel = (*customAuthProviderConfModel)(nil)

const (
	// 第三方配置类型, 如登录、云存储、短信等第三方功能
	ConfigTypeAuth = "auth" // 第三方登录
)

type (
	// AuthProviderConfModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuthProviderConfModel.
	AuthProviderConfModel interface {
		authProviderConfModel
		FindAll(ctx context.Context) ([]*AuthProviderConf, error)
	}

	customAuthProviderConfModel struct {
		*defaultAuthProviderConfModel
	}
)

// NewAuthProviderConfModel returns a model for the database table.
func NewAuthProviderConfModel(conn *gorm.DB) AuthProviderConfModel {
	return &customAuthProviderConfModel{
		defaultAuthProviderConfModel: newAuthProviderConfModel(conn),
	}
}

func (m *customAuthProviderConfModel) FindAll(ctx context.Context) ([]*AuthProviderConf, error) {
	confs := []*AuthProviderConf{}
	err := m.conn.WithContext(ctx).Model(&AuthProviderConf{}).Find(&confs).Error
	if err != nil {
		return nil, err
	}
	return confs, nil
}
