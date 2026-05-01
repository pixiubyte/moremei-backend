package feishu

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type TableInitializer struct {
	client    *Client
	appToken  string
	tableID   string
	logger    logx.Logger
	sheetType int
}

func NewTableInitializer(client *Client, appToken, tableID string, logger logx.Logger, sheetType int) *TableInitializer {
	return &TableInitializer{
		client:    client,
		appToken:  appToken,
		tableID:   tableID,
		logger:    logger,
		sheetType: sheetType,
	}
}

func InitTable(client *Client, appToken, tableID string, sheetType int) error {
	initializer := NewTableInitializer(client, appToken, tableID, logx.WithContext(context.Background()), sheetType)
	return initializer.Initialize()
}

func (t *TableInitializer) Initialize() error {
	if t.sheetType == 1 {
		if err := t.cleanExistingFields(); err != nil {
			return errors.Wrap(err, "清理现有字段失败")
		}

		if err := t.setupFields(); err != nil {
			return errors.Wrap(err, "设置字段失败")
		}
	} else if t.sheetType == 2 {
		if err := t.cleanExistingFieldsMei(); err != nil {
			return errors.Wrap(err, "清理现有字段失败")
		}
		if err := t.setupFieldsMei(); err != nil {
			return errors.Wrap(err, "设置字段失败")
		}
	}
	if err := t.cleanExistingRecords(); err != nil {
		return errors.Wrap(err, "清理现有记录失败")
	}

	return nil
}

func (t *TableInitializer) cleanExistingFields() error {
	list, err := t.client.ListFields(t.appToken, t.tableID, 20, "")
	if err != nil {
		return errors.Wrap(err, "获取字段列表失败")
	}

	for _, item := range list.Data.Items {
		if item.IsPrimary {
			if err := t.updatePrimaryField(item.FieldID); err != nil {
				return err
			}
			continue
		}

		if err := t.client.DeleteField(t.appToken, t.tableID, item.FieldID); err != nil {
			return errors.Wrapf(err, "删除字段[%s]失败", item.FieldName)
		}
	}
	return nil
}

func (t *TableInitializer) updatePrimaryField(fieldID string) error {
	updateReq := UpdateFieldRequest{
		FieldName: "会话ID",
		Type:      2,
		Property: map[string]interface{}{
			"formatter": "0",
		},
	}

	if err := t.client.UpdateField(t.appToken, t.tableID, fieldID, updateReq); err != nil {
		return errors.Wrap(err, "更新主键字段失败")
	}
	return nil
}

func (t *TableInitializer) setupFields() error {
	fields := []CreateFieldRequest{
		{FieldName: "会话人", Type: 1},
		{FieldName: "手机号", Type: 1},
		{FieldName: "会话内容", Type: 1},
		{FieldName: "会话时间", Type: 1},
		{FieldName: "会话日期", Type: 5},
	}

	for _, field := range fields {
		if err := t.client.CreateField(t.appToken, t.tableID, field); err != nil {
			return errors.Wrapf(err, "创建字段[%s]失败", field.FieldName)
		}
	}

	return nil
}

func (t *TableInitializer) cleanExistingRecords() error {
	resp, err := t.client.QueryRecords(t.appToken, t.tableID, QueryRecordsRequest{}, "", 20)
	if err != nil {
		return errors.Wrap(err, "查询现有记录失败")
	}

	var recordIds []string
	for _, v := range resp.Data.Items {
		if recordId, ok := v["record_id"]; ok {
			recordIds = append(recordIds, recordId.(string))
		}
	}

	if len(recordIds) > 0 {
		_, err = t.client.BatchDeleteRecords(
			t.appToken,
			t.tableID,
			BatchDeleteRecordsRequest{
				Records: recordIds,
			},
		)
		if err != nil {
			return errors.Wrap(err, "批量删除记录失败")
		}
	}

	return nil
}

func CreateBitable(name string, client *Client) (*CreateBitableResponse, error) {
	resp, err := client.CreateBitable(name, "")
	if err != nil {
		logx.Errorf("Error getting batch id: %v", err)
		return nil, err
	}
	return resp, nil
}

func (t *TableInitializer) cleanExistingFieldsMei() error {
	list, err := t.client.ListFields(t.appToken, t.tableID, 20, "")
	if err != nil {
		return errors.Wrap(err, "获取字段列表失败")
	}

	for _, item := range list.Data.Items {
		if item.IsPrimary {
			if err := t.updatePrimaryFieldMei(item.FieldID); err != nil {
				return err
			}
			continue
		}

		if err := t.client.DeleteField(t.appToken, t.tableID, item.FieldID); err != nil {
			return errors.Wrapf(err, "删除字段[%s]失败", item.FieldName)
		}
	}
	return nil
}

func (t *TableInitializer) updatePrimaryFieldMei(fieldID string) error {
	updateReq := UpdateFieldRequest{
		FieldName: "用户ID",
		Type:      1,
	}

	if err := t.client.UpdateField(t.appToken, t.tableID, fieldID, updateReq); err != nil {
		return errors.Wrap(err, "更新主键字段失败")
	}
	return nil
}

func (t *TableInitializer) setupFieldsMei() error {
	fields := []CreateFieldRequest{
		{FieldName: "用户名称", Type: 1},
		{FieldName: "手机号", Type: 1},
		{FieldName: "酶好天使ID", Type: 1},
		{FieldName: "酶好天使", Type: 1},
		{FieldName: "类型", Type: 1},
		{FieldName: "设备SN", Type: 1},
		{FieldName: "血糖数值", Type: 2},
		{FieldName: "进餐时机", Type: 1},
		{FieldName: "测量时间", Type: 1},
		{FieldName: "测量日期", Type: 5},
	}

	for _, field := range fields {
		if err := t.client.CreateField(t.appToken, t.tableID, field); err != nil {
			return errors.Wrapf(err, "创建字段[%s]失败", field.FieldName)
		}
	}

	return nil
}
