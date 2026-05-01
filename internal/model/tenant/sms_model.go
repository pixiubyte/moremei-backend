package tenant

import (
	"gorm.io/gorm"
)

const (
	SmsBizTypeLogin = "login"
	SmsBizTypeBind  = "bind"
)

var _ SmsModel = (*customSmsModel)(nil)

type (
	// SmsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSmsModel.
	SmsModel interface {
		smsModel
	}

	customSmsModel struct {
		*defaultSmsModel
	}
)

// NewSmsModel returns a model for the database table.
func NewSmsModel(conn *gorm.DB) SmsModel {
	return &customSmsModel{
		defaultSmsModel: newSmsModel(conn),
	}
}
