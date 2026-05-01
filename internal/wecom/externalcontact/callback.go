package externalcontact

import (
	"encoding/xml"

	"moremei/ai-saas/internal/wecom/util"
)

// 原始回调消息内容
type callbackOriginMessage struct {
	ToUserName string // 企业微信的CorpID，当为第三方套件回调事件时，CorpID的内容为suiteid
	AgentID    string // 接收的应用id，可在应用的设置页面获取
	Encrypt    string // 消息结构体加密后的字符串
}

// EventCallbackMessage 微信客户联系回调消息
// https://developer.work.weixin.qq.com/document/path/92130
type EventCallbackMessage struct {
	ToUserName     string `json:"to_user_name"`
	FromUserName   string `json:"from_user_name"`
	CreateTime     int64  `json:"create_time"`
	MsgType        string `json:"msg_type"`
	Event          string `json:"event"`
	ChangeType     string `json:"change_type"`
	UserID         string `json:"user_id,omitempty"`
	ExternalUserID string `json:"external_user_id,omitempty"`
	State          string `json:"state,omitempty"`
	WelcomeCode    string `json:"welcome_code,omitempty"`
	Source         string `json:"source,omitempty"`
	FailReason     string `json:"fail_reason,omitempty"`
	ChatId         string `json:"chat_id,omitempty"`
	UpdateDetail   string `json:"update_detail,omitempty"`
	JoinScene      int64  `json:"join_scene,omitempty"`
	QuitScene      int64  `json:"quit_scene,omitempty"`
	MemChangeCnt   int64  `json:"mem_change_cnt,omitempty"`
	Id             string `json:"id,omitempty"`
	TagType        string `json:"tag_type,omitempty"`
	StrategyId     string `json:"strategy_id,omitempty"`
}

// GetCallbackMessage 获取联系客户回调事件中的消息内容
func (r *Client) GetCallbackMessage(encryptedMsg []byte) (msg EventCallbackMessage, err error) {
	var origin callbackOriginMessage
	if err = xml.Unmarshal(encryptedMsg, &origin); err != nil {
		return
	}
	_, bData, err := util.DecryptMsg(r.CorpID, origin.Encrypt, r.ExternalContact.EncodingAESKey)
	if err != nil {
		return
	}
	if err = xml.Unmarshal(bData, &msg); err != nil {
		return
	}
	return
}
