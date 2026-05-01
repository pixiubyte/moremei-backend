package tenant

import (
	"gorm.io/gorm"
)

var _ TokenManageModel = (*customTokenManageModel)(nil)

type (
	// TokenManageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTokenManageModel.
	TokenManageModel interface {
		tokenManageModel
	}

	customTokenManageModel struct {
		*defaultTokenManageModel
	}
)

// NewTokenManageModel returns a model for the database table.
func NewTokenManageModel(conn *gorm.DB) TokenManageModel {
	return &customTokenManageModel{
		defaultTokenManageModel: newTokenManageModel(conn),
	}
}
