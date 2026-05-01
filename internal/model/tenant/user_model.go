package tenant

import (
	"context"
	"time"

	commonConst "moremei/ai-saas/pkg/consts"
	"moremei/ai-saas/pkg/pagination"

	"gorm.io/gorm"
)

var (
	_ UserModel = (*customUserModel)(nil)

	UserGenders = map[int64]string{
		UserGenderUnspecified: "未知",
		UserGenderMale:        "男",
		UserGenderFemale:      "女",
	}
)

const (
	UserGenderUnspecified = 0
	UserGenderMale        = 1
	UserGenderFemale      = 2
)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		FindByIds(ctx context.Context, ids []uint64) ([]*User, error)
		PageFind(ctx context.Context, param *pagination.PageRequest) (*pagination.PageResponse, error)
		FindUserWeb(ctx context.Context, userId uint64) (resp *UserWeb, err error)
		FindByPage(ctx context.Context, param *pagination.PageRequest) (*pagination.PageResponse, error)
		FindByStatus(ctx context.Context, status uint64) ([]*User, error)
		Count(ctx context.Context, status uint64) (int64, error)
	}

	customUserModel struct {
		*defaultUserModel
	}

	UserWeb struct {
		Id         uint64    `gorm:"column:id"`
		Nickname   string    `gorm:"column:nickname"`
		Avatar     string    `gorm:"column:avatar"`
		Status     int64     `gorm:"column:status"`
		Mobile     string    `gorm:"column:mobile"`
		RegisterAt time.Time `gorm:"column:register_at"`
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn *gorm.DB) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *customUserModel) FindByIds(ctx context.Context, ids []uint64) ([]*User, error) {
	var resp []*User
	err := m.conn.WithContext(ctx).Model(&User{}).Where("`id` IN (?)", ids).Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customUserModel) PageFind(ctx context.Context, pageReq *pagination.PageRequest) (*pagination.PageResponse, error) {
	var (
		total  int64
		models []*UserWeb
	)

	query := m.conn.WithContext(ctx).Table("`user`").
		Joins("left join `web_user` on `user`.`id` = `web_user`.`user_id` and `web_user`.`deleted` = ?", commonConst.FlagFalse).
		Select("`user`.`id` as id, `user`.`nickname`, `user`.`avatar` as avatar, `user`.`status`, `web_user`.`mobile` as mobile, `user`.`register_at`").
		Order("`user`.`id` desc")

	if pageReq.QueryParams["mobile"] != "" {
		query.Where("`web_user`.`mobile` = ?", pageReq.QueryParams["mobile"])
	}
	if userId, ok := pageReq.QueryParams["userId"].([]uint64); ok && len(userId) != 0 {
		query.Where("`user`.`id` = (?)", userId)
	}
	if nickname, ok := pageReq.QueryParams["nickname"].(string); ok && nickname != "" {
		query.Where("`user`.`nickname` like ?", "%"+nickname+"%")
	}
	if status, ok := pageReq.QueryParams["status"].(int64); ok && status >= 0 {
		query.Where("`user`.`status` = ?", status)
	}

	query.Count(&total)

	err := query.Scopes(pagination.Paginate(pageReq)).Find(&models).Error

	if err != nil {
		return nil, err
	}

	return pagination.NewPageResponse(total, pageReq, models), nil
}

func (m *customUserModel) FindUserWeb(ctx context.Context, userId uint64) (resp *UserWeb, err error) {
	err = m.conn.WithContext(ctx).Table("`user`").
		Joins("left join `web_user` on `user`.`id` = `web_user`.`user_id`and `web_user`.`deleted` = ?", commonConst.FlagFalse).
		Select("`user`.`id` as id, `user`.`nickname`, `user`.`avatar` as avatar, `user`.`status`, `web_user`.`mobile` as mobile, `register_at`").
		Where("`user`.`id` = ?", userId).Take(&resp).Error
	return
}

func (m *customUserModel) FindByPage(ctx context.Context, pageReq *pagination.PageRequest) (*pagination.PageResponse, error) {
	var total int64
	var data []*User

	query := m.conn.WithContext(ctx).Model(&User{}).Count(&total)
	query = query.Scopes(pagination.Paginate(pageReq))

	if status, ok := pageReq.QueryParams["status"]; ok {
		query = query.Where("`status` = ?", status)
	}
	query = query.Order("id DESC")

	err := query.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return pagination.NewPageResponse(total, pageReq, data), nil
}

func (m *customUserModel) FindByStatus(ctx context.Context, status uint64) ([]*User, error) {
	var users []*User
	err := m.conn.WithContext(ctx).Where("status = ?", status).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (m *customUserModel) Count(ctx context.Context, status uint64) (int64, error) {
	var total int64
	err := m.conn.WithContext(ctx).Model(&User{}).Where("status = ?", status).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
