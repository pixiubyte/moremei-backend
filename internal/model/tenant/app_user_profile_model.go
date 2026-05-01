package tenant

import (
	"gorm.io/gorm"
)

var _ AppUserProfileModel = (*customAppUserProfileModel)(nil)

type (
	// AppUserProfileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAppUserProfileModel.
	AppUserProfileModel interface {
		appUserProfileModel
	}

	customAppUserProfileModel struct {
		*defaultAppUserProfileModel
	}
)

// NewAppUserProfileModel returns a model for the database table.
func NewAppUserProfileModel(conn *gorm.DB) AppUserProfileModel {
	return &customAppUserProfileModel{
		defaultAppUserProfileModel: newAppUserProfileModel(conn),
	}
}
