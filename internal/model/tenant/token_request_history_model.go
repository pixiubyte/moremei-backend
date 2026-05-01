package tenant

import (
	"gorm.io/gorm"
)

var _ TokenRequestHistoryModel = (*customTokenRequestHistoryModel)(nil)

type (
	// TokenRequestHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTokenRequestHistoryModel.
	TokenRequestHistoryModel interface {
		tokenRequestHistoryModel
	}

	customTokenRequestHistoryModel struct {
		*defaultTokenRequestHistoryModel
	}
)

// NewTokenRequestHistoryModel returns a model for the database table.
func NewTokenRequestHistoryModel(conn *gorm.DB) TokenRequestHistoryModel {
	return &customTokenRequestHistoryModel{
		defaultTokenRequestHistoryModel: newTokenRequestHistoryModel(conn),
	}
}
