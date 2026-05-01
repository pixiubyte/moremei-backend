package tenant

import (
	"gorm.io/gorm"
)

var _ ParamsStorageModel = (*customParamsStorageModel)(nil)

type (
	// ParamsStorageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customParamsStorageModel.
	ParamsStorageModel interface {
		paramsStorageModel
	}

	customParamsStorageModel struct {
		*defaultParamsStorageModel
	}
)

// NewParamsStorageModel returns a model for the database table.
func NewParamsStorageModel(conn *gorm.DB) ParamsStorageModel {
	return &customParamsStorageModel{
		defaultParamsStorageModel: newParamsStorageModel(conn),
	}
}
