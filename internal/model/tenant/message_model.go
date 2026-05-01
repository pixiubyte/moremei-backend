package tenant

import (
	"context"
	"fmt"
	commonConst "moremei/ai-saas/pkg/consts"
	"moremei/ai-saas/pkg/pagination"

	"gorm.io/gorm"
)

const (
	MessageAuthorTypeUser    string = "User"
	MessageAuthorTypeAI      string = "AI"
	MessageSourceWeb         int64  = 0
	MessageSourceWecom       int64  = 1
	MessageSourceMiniProgram int64  = 3
	WorktoolQaTextTypeText   int64  = 1
	WorktoolQaTextTypeVoice  int64  = 2
)

var _ MessageModel = (*customMessageModel)(nil)

type (
	// MessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageModel.
	MessageModel interface {
		messageModel
		FindByIds(ctx context.Context, ids []uint64) ([]*Message, error)
		FindOneByConversationIdParentId(ctx context.Context, cid, pid uint64) (*Message, error)
		PageFind(ctx context.Context, pageReq *pagination.PageRequest) (pageResp *pagination.PageResponse, err error)
		GetAllUserLatestMessage(ctx context.Context, limit, offset int) ([]*Message, error)
		GetLatestByConversationId(ctx context.Context, cid uint64) (res *Message, err error)
		UpdateByConversationId(ctx context.Context, cid []uint64, newCid uint64) error
		FindByDate(ctx context.Context, startDate, endDate string) ([]*Message, error)
		FindBeforeOne(ctx context.Context, cid, mid uint64, source int64) (*Message, error)
		FindAfterOne(ctx context.Context, cid, mid uint64, source int64) (*Message, error)
		UpdateConversationId(ctx context.Context, cid uint64, newCid uint64) error
		BatchInsert(ctx context.Context, data []*Message, tx *gorm.DB) error
		CountByConversationIdAuthorType(ctx context.Context, beforeId, cid uint64, atype string) (int64, error)
		GetLatestByConversationIdAndType(ctx context.Context, cid uint64, atype string) (res *Message, err error)
		FindAuthorIdsByMonth(ctx context.Context, date string) ([]uint64, error)
		FindAuthorIdsByWeek(ctx context.Context, date string) ([]uint64, error)
		FindAuthorIdsByDaily(ctx context.Context, date string) ([]uint64, error)
		FindAiNoReplyIdsBySourceDate(ctx context.Context, source int64, startDate, endDate string) ([]uint64, error)
		FindByCidAndId(ctx context.Context, cid, id uint64, limit int) ([]*Message, error)
		FindByConvIdIntervalDaySource(ctx context.Context, convId uint64, intervalDay, source int64, limit int) ([]*UserAiMessage, error)
		FindConvUsersByIntervalDaySource(ctx context.Context, intervalDay, source int64) ([]*ConvUser, error)
		FindOneByCidAuthTypeContentCreatedAt(ctx context.Context, cid uint64, authType string, start, end string) (*Message, error)
		FindByCid(ctx context.Context, cid uint64) ([]*Message, error)
		FindByAuthorId(ctx context.Context, authorId uint64) ([]*Message, error)
		FindUserId(ctx context.Context) ([]uint64, error)
		FindConversationIdByAuthorIds(ctx context.Context, authorIds []uint64) ([]*Message, error)
		GetLatestMessageByAuthorId(ctx context.Context, authorId uint64) (*Message, error)
		FindUserConversationIds(ctx context.Context, userId uint64) ([]uint64, error)
		FindByConversations(ctx context.Context, conversations []uint64) ([]*Message, error)
		FindConversationIdsByAuthorIdAndIdRange(ctx context.Context, authorId uint64, start, end uint64) ([]uint64, error)
		FindByConversationsAndIdRange(ctx context.Context, conversations []uint64, start, end uint64) ([]*Message, error)
	}

	customMessageModel struct {
		*defaultMessageModel
	}

	UserAiMessage struct {
		Query  string
		Answer string
	}

	ConvUser struct {
		UserId         uint64
		ConversationId uint64
		LatestDate     string
	}
)

// NewMessageModel returns a model for the database table.
func NewMessageModel(conn *gorm.DB) MessageModel {
	return &customMessageModel{
		defaultMessageModel: newMessageModel(conn),
	}
}

func (m *customMessageModel) FindByIds(ctx context.Context, ids []uint64) ([]*Message, error) {
	var resp []*Message
	err := m.conn.WithContext(ctx).Model(&Message{}).
		Where("id in (?)", ids).Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customMessageModel) FindOneByConversationIdParentId(ctx context.Context, cid, pid uint64) (*Message, error) {
	var resp Message
	err := m.conn.WithContext(ctx).Model(&Message{}).
		Where("conversation_id = ?", cid).
		Where("parent_id = ?", pid).
		First(&resp).
		Error
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (m *customMessageModel) PageFind(ctx context.Context, pageReq *pagination.PageRequest) (pageResp *pagination.PageResponse, err error) {
	var (
		total  int64
		models []*Message
	)
	query := m.conn.WithContext(ctx).Model(&Message{}).
		Where("conversation_id = ?", pageReq.QueryParams["conversation_id"].(uint64)).
		Order("id DESC")

	if s, ok := pageReq.QueryParams["source"].(int64); ok {
		query = query.Where("source = ?", s)
	}

	if q, ok := pageReq.QueryParams["latest_id"]; ok {
		query = query.Where("id < ?", q)
	}

	if keyword, ok := pageReq.QueryParams["keyword"].(string); ok && keyword != "" {
		query = query.Where("content like ?", fmt.Sprintf("%%%v%%", keyword))
	}

	if createdAt, ok := pageReq.QueryParams["start_at"].(string); ok && createdAt != "" {
		query = query.Where("created_at >= ?", createdAt)
	}

	if createdAt, ok := pageReq.QueryParams["end_at"].(string); ok && createdAt != "" {
		query = query.Where("created_at <= ?", createdAt)
	}
	query.Count(&total)

	err = query.Scopes(pagination.Paginate(pageReq)).Find(&models).Error
	if err != nil {
		return nil, err
	}

	return pagination.NewPageResponse(total, pageReq, models), nil
}

func (m *customMessageModel) GetAllUserLatestMessage(ctx context.Context, limit, offset int) (res []*Message, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Select("id,author_id,conversation_id, max(created_at) as created_at").
		Where("`author_type` = ? and `source` = ? ", MessageAuthorTypeUser, 1).
		Group("conversation_id,author_id").
		Order("conversation_id ASC").
		Limit(limit).
		Offset(offset).
		Find(&res).Error
	return
}

func (m *customMessageModel) GetLatestByConversationId(ctx context.Context, cid uint64) (res *Message, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).
		Where("conversation_id = ?", cid).
		Order("id desc").
		First(&res).Error
	return
}

func (m *customMessageModel) UpdateByConversationId(ctx context.Context, cid []uint64, newCid uint64) error {
	return m.conn.WithContext(ctx).Model(&Message{}).
		Where("conversation_id in ?", cid).
		Update("conversation_id", newCid).Error
}

func (m *customMessageModel) FindByDate(ctx context.Context, startDate, endDate string) ([]*Message, error) {
	var resp []*Message
	err := m.conn.WithContext(ctx).Model(&Message{}).
		Where("`type` IN (?)", []int64{WorktoolQaTextTypeText, WorktoolQaTextTypeVoice}).
		Where("`author_type` = ?", MessageAuthorTypeUser).
		Where("`created_at` >= ?", startDate).
		Where("`created_at` < ?", endDate).
		Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customMessageModel) FindBeforeOne(ctx context.Context, cid, mid uint64, source int64) (*Message, error) {
	var resp Message
	err := m.conn.WithContext(ctx).Model(&Message{}).
		Where("`conversation_id` = ?", cid).
		Where("`author_type` = ?", MessageAuthorTypeUser).
		Where("`id` < ?", mid).
		Order("`id` DESC").
		First(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customMessageModel) FindAfterOne(ctx context.Context, cid, mid uint64, source int64) (*Message, error) {
	var resp Message
	err := m.conn.WithContext(ctx).Model(&Message{}).
		Where("`conversation_id` = ?", cid).
		Where("`author_type` = ?", MessageAuthorTypeUser).
		Where("`id` > ?", mid).
		Order("`id` ASC").
		First(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customMessageModel) UpdateConversationId(ctx context.Context, cid uint64, newCid uint64) error {
	return m.conn.WithContext(ctx).Model(&Message{}).
		Where("conversation_id = ?", cid).
		Update("conversation_id", newCid).Error
}

func (m *customMessageModel) BatchInsert(ctx context.Context, data []*Message, tx *gorm.DB) error {
	if tx != nil {
		return tx.WithContext(ctx).Model(&Message{}).CreateInBatches(data, 100).Error
	}

	return m.conn.WithContext(ctx).Model(&Message{}).CreateInBatches(data, 100).Error
}

func (m *customMessageModel) CountByConversationIdAuthorType(ctx context.Context, beforeId, cid uint64, atype string) (res int64, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).
		Where("id < ? and conversation_id = ? and author_type = ?", beforeId, cid, atype).
		Count(&res).Error
	return
}

func (m *customMessageModel) GetLatestByConversationIdAndType(ctx context.Context, cid uint64, atype string) (res *Message, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).
		Where("conversation_id = ? and author_type = ?", cid, atype).Order("id desc").First(&res).Error
	return
}

func (m *customMessageModel) FindAuthorIdsByMonth(ctx context.Context, date string) ([]uint64, error) {
	var resp []uint64
	err := m.conn.WithContext(ctx).Model(&Message{}).
		Where("author_type = ?", MessageAuthorTypeUser).
		Where("created_at >= DATE_FORMAT(DATE_SUB(?, INTERVAL 1 MONTH), '%Y-%m-01')", date).
		Where("created_at < DATE_FORMAT(?, '%Y-%m-01')", date).
		Group("author_id, month(created_at)").
		Distinct("author_id").
		Pluck("author_id", &resp).Error

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customMessageModel) FindAuthorIdsByWeek(ctx context.Context, date string) ([]uint64, error) {
	var resp []uint64
	conn := m.conn.WithContext(ctx)
	q1 := conn.Model(&Message{}).
		Where("author_type = ?", MessageAuthorTypeUser).
		Where("created_at >= DATE_FORMAT(DATE_SUB(?, INTERVAL 1 MONTH), '%Y-%m-01')", date).
		Where("created_at < DATE_FORMAT(?, '%Y-%m-01')", date).
		Group("author_id, week(created_at, 1)")

	err := conn.Table("(?) as t", q1).Group("author_id").
		Having(
			"count(*) = (week(DATE_FORMAT(DATE_SUB(DATE_FORMAT(?, '%Y-%m-01'), INTERVAL 1 DAY), '%Y-%m-%D'), 1) - week(DATE_FORMAT(DATE_SUB(?, INTERVAL 1 MONTH), '%Y-%m-01'), 1))",
			date,
			date,
		).
		Distinct("author_id").
		Pluck("author_id", &resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customMessageModel) FindAuthorIdsByDaily(ctx context.Context, date string) ([]uint64, error) {
	var resp []uint64
	conn := m.conn.WithContext(ctx)
	q1 := conn.Model(&Message{}).
		Where("author_type = ?", MessageAuthorTypeUser).
		Where("created_at >= DATE_FORMAT(DATE_SUB(?, INTERVAL 1 MONTH), '%Y-%m-01')", date).
		Where("created_at < DATE_FORMAT(?, '%Y-%m-01')", date).
		Group("author_id, date(created_at)")
	err := conn.Table("(?) as t", q1).
		Joins("LEFT JOIN user ON user.id = t.author_id").
		Group("author_id").
		Having("count(*) = day(last_day(DATE_SUB(?, INTERVAL 1 MONTH)))", date).
		Error

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customMessageModel) FindAiNoReplyIdsBySourceDate(ctx context.Context, source int64, startDate, endDate string) ([]uint64, error) {
	var resp []uint64
	err := m.conn.WithContext(ctx).Model(&Message{}).Table("message as m1").
		Joins("LEFT JOIN message as m2 ON m2.parent_id = m1.id and m2.author_type = ?", MessageAuthorTypeAI).
		Where("m1.author_type = ?", MessageAuthorTypeUser).
		Where("m1.source = ?", source).
		Where("m1.created_at >= ?", startDate).
		Where("m1.created_at < ?", endDate).
		Where("m2.id is null").
		Distinct("m1.id").
		Pluck("m1.id", &resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customMessageModel) FindByCidAndId(ctx context.Context, cid, id uint64, limit int) (resp []*Message, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Where("conversation_id = ? and id <= ? and source = ?", cid, id, MessageSourceWecom).
		Order("created_at desc").Limit(limit).Find(&resp).Error
	return
}

func (m *customMessageModel) FindByConvIdIntervalDaySource(ctx context.Context, convId uint64, intervalDay, source int64, limit int) ([]*UserAiMessage, error) {
	var resp []*UserAiMessage

	err := m.conn.WithContext(ctx).Table("`message` as `m1`").
		Joins("LEFT JOIN `message` AS `m2` ON `m1`.`id` = `m2`.`parent_id` AND `m2`.`author_type` = ?", MessageAuthorTypeAI).
		Where("`m1`.`author_type` = ?", MessageAuthorTypeUser).
		Where("`m1`.`conversation_id` = ?", convId).
		Where("`m1`.`source` = ?", source).
		//Where("`m1`.`created_at` >= date_format(now() - INTERVAL ? DAY, '%Y-%m-%d')", intervalDay).
		Select("`m1`.`content` as `query`, `m2`.`content` as `answer`").
		Order("`m1`.`id` DESC").
		Limit(limit).
		Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customMessageModel) FindConvUsersByIntervalDaySource(ctx context.Context, intervalDay, source int64) ([]*ConvUser, error) {
	var resp []*ConvUser
	err := m.conn.WithContext(ctx).Model(&Message{}).
		Joins("LEFT JOIN `conversation` ON `conversation`.`id` = `message`.`conversation_id`").
		Joins("LEFT JOIN `wecom_ai_robot` ON `wecom_ai_robot`.`ai_robot_id` = `conversation`.`ai_robot_id`").
		Joins("LEFT JOIN `wecom_external_user` ON `wecom_external_user`.`user_id` = `message`.`author_id` AND `wecom_external_user`.`wecom_ai_robot_id` = `wecom_ai_robot`.`id`").
		Where("`message`.`author_type` = ?", MessageAuthorTypeUser).
		Where("`message`.`source` = ?", source).
		Where("`message`.`created_at` >= date_format(now() - INTERVAL ? DAY, '%Y-%m-%d')", intervalDay).
		Where("`wecom_external_user`.`deleted` = ?", commonConst.FlagFalse).
		Group("`message`.`author_id`, `message`.`conversation_id`").
		Select("`message`.`author_id` as `user_id`, `message`.`conversation_id`, date_format(max(`message`.`created_at`), '%Y-%m-%d') as `latest_date`").
		Having("`latest_date` = date_format(now() - INTERVAL ? DAY, '%Y-%m-%d')", intervalDay).
		Find(&resp).
		Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customMessageModel) FindOneByCidAuthTypeContentCreatedAt(ctx context.Context, cid uint64, authType string, start, end string) (model *Message, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Where("conversation_id = ? and author_type = ? and created_at >= ? and created_at <= ?", cid, authType, start, end).
		First(&model).Error
	return
}

func (m *customMessageModel) FindByCid(ctx context.Context, cid uint64) (models []*Message, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Where("conversation_id = ?", cid).Order("id asc").Limit(100).Find(&models).Error
	return
}

func (m *customMessageModel) FindByAuthorId(ctx context.Context, authorId uint64) (models []*Message, err error) {
	// 直接查询该用户的所有消息
	err = m.conn.WithContext(ctx).Model(&Message{}).
		Where("author_id = ?", authorId).
		Order("created_at ASC").
		Find(&models).Error
	return
}

func (m *customMessageModel) FindUserId(ctx context.Context) (userIds []uint64, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Select("author_id").Group("author_id").Find(&userIds).Error
	return
}

func (m *customMessageModel) FindConversationIdByAuthorIds(ctx context.Context, authorIds []uint64) (resp []*Message, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Where("author_id IN (?) AND author_type = ?", authorIds, "User").Find(&resp).Error
	return
}

func (m *customMessageModel) GetLatestMessageByAuthorId(ctx context.Context, authorId uint64) (resp *Message, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Where("author_id = ? AND author_type = ?", authorId, "User").Order("id DESC").First(&resp).Error
	return
}

func (m *customMessageModel) FindUserConversationIds(ctx context.Context, userId uint64) (resp []uint64, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Where("author_id = ? AND author_type = ?", userId, "User").Group("conversation_id").Pluck("conversation_id", &resp).Error
	return
}

func (m *customMessageModel) FindByConversations(ctx context.Context, conversations []uint64) (resp []*Message, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Where("conversation_id IN (?)", conversations).Find(&resp).Error
	return
}

func (m *customMessageModel) FindConversationIdsByAuthorIdAndIdRange(ctx context.Context, authorId uint64, start, end uint64) (resp []uint64, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Where("author_id = ? AND author_type = ? AND id > ? AND id <= ?", authorId, "User", start, end).Group("conversation_id").Pluck("conversation_id", &resp).Error
	return
}

func (m *customMessageModel) FindByConversationsAndIdRange(ctx context.Context, conversations []uint64, start, end uint64) (resp []*Message, err error) {
	err = m.conn.WithContext(ctx).Model(&Message{}).Where("conversation_id IN (?) AND id > ? AND id <= ?", conversations, start, end).Find(&resp).Error
	return
}
