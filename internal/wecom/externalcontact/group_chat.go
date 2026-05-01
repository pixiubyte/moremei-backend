package externalcontact

import (
	"fmt"

	"moremei/ai-saas/internal/wecom/util"
)

const (
	// groupChatURL 客户群
	groupChatURL       = "https://qyapi.weixin.qq.com/cgi-bin/externalcontact/groupchat"
	openGidToChatIdURL = "https://qyapi.weixin.qq.com/cgi-bin/externalcontact/opengid_to_chatid"
)

type (
	// AddJoinWayRequest 添加群配置请求参数
	AddJoinWayRequest struct {
		Scene          int      `json:"scene"`            // 必填 1 - 群的小程序插件,2 - 群的二维码插件
		Remark         string   `json:"remark"`           //非必填	联系方式的备注信息，用于助记，超过30个字符将被截断
		AutoCreateRoom int      `json:"auto_create_room"` //非必填	当群满了后，是否自动新建群。0-否；1-是。 默认为1
		RoomBaseName   string   `json:"room_base_name"`   //非必填	自动建群的群名前缀，当auto_create_room为1时有效。最长40个utf8字符
		RoomBaseID     int      `json:"room_base_id"`     //非必填	自动建群的群起始序号，当auto_create_room为1时有效
		ChatIDList     []string `json:"chat_id_list"`     //必填	使用该配置的客户群ID列表，支持5个。见客户群ID获取方法
		State          string   `json:"state"`            //非必填	企业自定义的state参数，用于区分不同的入群渠道。不超过30个UTF-8字符
	}

	// AddJoinWayResponse 添加群配置返回值
	AddJoinWayResponse struct {
		util.CommonError
		ConfigID string `json:"config_id"`
	}
)

// AddJoinWay 加入群聊
// @see https://developer.work.weixin.qq.com/document/path/92229
func (r *Client) AddJoinWay(req *AddJoinWayRequest) (*AddJoinWayResponse, error) {
	var (
		accessToken string
		err         error
		response    []byte
	)
	if accessToken, err = r.GetAccessToken(); err != nil {
		return nil, err
	}
	response, err = util.PostJSON(fmt.Sprintf("%s/add_join_way?access_token=%s", groupChatURL, accessToken), req)
	if err != nil {
		return nil, err
	}
	result := &AddJoinWayResponse{}
	if err = util.DecodeWithError(response, result, "AddJoinWay"); err != nil {
		return nil, err
	}
	return result, nil
}

type (
	//JoinWayConfigRequest 获取或删除群配置的请求参数
	JoinWayConfigRequest struct {
		ConfigID string `json:"config_id"`
	}

	//JoinWay 群配置
	JoinWay struct {
		ConfigID       string   `json:"config_id"`
		Scene          int      `json:"scene"`
		Remark         string   `json:"remark"`
		AutoCreateRoom int      `json:"auto_create_room"`
		RoomBaseName   string   `json:"room_base_name"`
		RoomBaseID     int      `json:"room_base_id"`
		ChatIDList     []string `json:"chat_id_list"`
		QrCode         string   `json:"qr_code"`
		State          string   `json:"state"`
	}
	//GetJoinWayResponse 获取群配置的返回值
	GetJoinWayResponse struct {
		util.CommonError
		JoinWay JoinWay `json:"join_way"`
	}
)

// GetJoinWay 获取客户群进群方式配置
// @see https://developer.work.weixin.qq.com/document/path/92229
func (r *Client) GetJoinWay(req *JoinWayConfigRequest) (*GetJoinWayResponse, error) {
	var (
		accessToken string
		err         error
		response    []byte
	)
	if accessToken, err = r.GetAccessToken(); err != nil {
		return nil, err
	}
	response, err = util.PostJSON(fmt.Sprintf("%s/get_join_way?access_token=%s", groupChatURL, accessToken), req)
	if err != nil {
		return nil, err
	}
	result := &GetJoinWayResponse{}
	if err = util.DecodeWithError(response, result, "GetJoinWay"); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateJoinWayRequest 更新群配置的请求参数
type UpdateJoinWayRequest struct {
	ConfigID       string   `json:"config_id"`
	Scene          int      `json:"scene"`            // 必填 1 - 群的小程序插件,2 - 群的二维码插件
	Remark         string   `json:"remark"`           //非必填	联系方式的备注信息，用于助记，超过30个字符将被截断
	AutoCreateRoom int      `json:"auto_create_room"` //非必填	当群满了后，是否自动新建群。0-否；1-是。 默认为1
	RoomBaseName   string   `json:"room_base_name"`   //非必填	自动建群的群名前缀，当auto_create_room为1时有效。最长40个utf8字符
	RoomBaseID     int      `json:"room_base_id"`     //非必填	自动建群的群起始序号，当auto_create_room为1时有效
	ChatIDList     []string `json:"chat_id_list"`     //必填	使用该配置的客户群ID列表，支持5个。见客户群ID获取方法
	State          string   `json:"state"`            //非必填	企业自定义的state参数，用于区分不同的入群渠道。不超过30个UTF-8字符
}

// UpdateJoinWay 更新客户群进群方式配置
// @see https://developer.work.weixin.qq.com/document/path/92229
func (r *Client) UpdateJoinWay(req *UpdateJoinWayRequest) error {
	var (
		accessToken string
		err         error
		response    []byte
	)
	if accessToken, err = r.GetAccessToken(); err != nil {
		return err
	}
	response, err = util.PostJSON(fmt.Sprintf("%s/update_join_way?access_token=%s", groupChatURL, accessToken), req)
	if err != nil {
		return err
	}
	return util.DecodeWithCommonError(response, "UpdateJoinWay")
}

// DelJoinWay 删除客户群进群方式配置
// @see https://developer.work.weixin.qq.com/document/path/92229
func (r *Client) DelJoinWay(req *JoinWayConfigRequest) error {
	var (
		accessToken string
		err         error
		response    []byte
	)
	if accessToken, err = r.GetAccessToken(); err != nil {
		return err
	}
	response, err = util.PostJSON(fmt.Sprintf("%s/del_join_way?access_token=%s", groupChatURL, accessToken), req)
	if err != nil {
		return err
	}
	return util.DecodeWithCommonError(response, "DelJoinWay")
}

type ListGroupChatRequest struct {
	StatusFilter int                       `json:"status_filter,omitempty"` // 客户群跟进状态过滤：0 - 所有列表(即不过滤)；1 - 离职待继承；2 - 离职继承中；3 - 离职继承完成。
	OwnerFilter  *ListGroupChatOwnerFilter `json:"owner_filter,omitempty"`  // 群主过滤。如果不填，表示获取应用可见范围内全部群主的数据（但是不建议这么用，如果可见范围人数超过1000人，为了防止数据包过大，会报错 81017）
	Cursor       string                    `json:"cursor,omitempty"`        // 用于分页查询的游标，字符串类型，由上一次调用返回，首次调用不填
	Limit        int                       `json:"limit"`                   // 分页，预期请求的数据量，取值范围 1 ~ 1000
}

type ListGroupChatOwnerFilter struct {
	UseridList []string `json:"userid_list"` // 用户ID列表。最多100个
}

type ListGroupChatResponse struct {
	util.CommonError
	GroupChatList []GroupChatListResponseChat `json:"group_chat_list"` // 客户群列表
	NextCursor    string                      `json:"next_cursor"`     // 分页游标，下次请求时填写以获取之后分页的记录。如果该字段返回空则表示已没有更多数据
}

type GroupChatListResponseChat struct {
	ChatId string `json:"chat_id"` // 客户群ID
	Status int    `json:"status"`  // 客户群跟进状态：0 - 跟进人正常；1 - 跟进人离职；2 - 离职继承中；3 - 离职继承完成。
}

// ListGroupChat 获取客户群列表
// @see https://developer.work.weixin.qq.com/document/path/92120
func (r *Client) ListGroupChat(req *ListGroupChatRequest) (*ListGroupChatResponse, error) {
	var (
		accessToken string
		err         error
		response    []byte
	)
	if accessToken, err = r.GetAccessToken(); err != nil {
		return nil, err
	}
	response, err = util.PostJSON(fmt.Sprintf("%s/list?access_token=%s", groupChatURL, accessToken), req)
	if err != nil {
		return nil, err
	}
	result := &ListGroupChatResponse{}
	if err = util.DecodeWithError(response, result, "GroupChatList"); err != nil {
		return nil, err
	}
	return result, nil
}

type GetGroupChatRequest struct {
	ChatId   string `json:"chat_id"`             // 客户群ID
	NeedName int64  `json:"need_name,omitempty"` // 是否需要返回群成员的名字group_chat.member_list.name。0-不返回；1-返回。默认不返回
}
type GetGroupChatResponse struct {
	util.CommonError
	GroupChat GroupChat `json:"group_chat"` // 客户群详情
}
type GroupChat struct {
	ChatId     string            `json:"chat_id"`     // 客户群ID
	Name       string            `json:"name"`        // 群名
	Owner      string            `json:"owner"`       // 群主ID
	CreateTime int64             `json:"create_time"` // 群的创建时间
	Notice     string            `json:"notice"`      // 群公告
	MemberList []GroupChatMember `json:"member_list"` // 群成员列表
	AdminList  []GroupChatAdmin  `json:"admin_list"`  // 群管理员列表
}
type GroupChatMember struct {
	UserId        string            `json:"userid"`            // 群成员id
	Type          int64             `json:"type"`              // 成员类型：1 - 企业成员；2 - 外部联系人。
	UnionId       string            `json:"unionid,omitempty"` // 外部联系人在微信开放平台的唯一身份标识（微信unionid），通过此字段企业可将外部联系人与公众号/小程序用户关联起来。仅当群成员类型是微信用户（包括企业成员未添加好友），且企业绑定了微信开发者ID有此字段。第三方不可获取，上游企业不可获取下游企业客户的unionid字段
	JoinTime      int64             `json:"join_time"`         // 入群时间
	JoinScene     int64             `json:"join_scene"`        // 入群方式：1 - 由群成员邀请入群（直接邀请入群）；2 - 由群成员邀请入群（通过邀请链接入群）；3 - 通过扫描群二维码入群。
	Invitor       *GroupChatInvitor `json:"invitor,omitempty"` // 邀请者。目前仅当是由本企业内部成员邀请入群时会返回该值
	GroupNickname string            `json:"group_nickname"`    // 在群里的昵称
	Name          string            `json:"name,omitempty"`    // 名字。仅当 need_name = 1 时返回：如果是微信用户，则返回其在微信中设置的名字；如果是企业微信联系人，则返回其设置对外展示的别名或实名。
}
type GroupChatAdmin struct {
	UserId string `json:"userid"` // 群管理员userid
}
type GroupChatInvitor struct {
	UserId string `json:"userid"` // 邀请者的userid
}

// GetGroupChat 获取客户群详情
// @see https://developer.work.weixin.qq.com/document/path/92122
func (r *Client) GetGroupChat(req *GetGroupChatRequest) (*GetGroupChatResponse, error) {
	var (
		accessToken string
		err         error
		response    []byte
	)
	if accessToken, err = r.GetAccessToken(); err != nil {
		return nil, err
	}
	response, err = util.PostJSON(fmt.Sprintf("%s/get?access_token=%s", groupChatURL, accessToken), req)
	if err != nil {
		return nil, err
	}
	result := &GetGroupChatResponse{}
	if err = util.DecodeWithError(response, result, "GroupChatGet"); err != nil {
		return nil, err
	}
	return result, nil
}

type OpenGidToChatIdRequest struct {
	OpenGid string `json:"opengid"`
}
type OpenGidToChatIdResponse struct {
	util.CommonError
	ChatId string `json:"chat_id"`
}

// OpenGidToChatId 客户群opengid转换
// @see https://developer.work.weixin.qq.com/document/path/94822
func (r *Client) OpenGidToChatId(req *OpenGidToChatIdRequest) (*OpenGidToChatIdResponse, error) {
	var (
		accessToken string
		err         error
		response    []byte
	)
	if accessToken, err = r.GetAccessToken(); err != nil {
		return nil, err
	}
	response, err = util.PostJSON(fmt.Sprintf("%s?access_token=%s", openGidToChatIdURL, accessToken), req)
	if err != nil {
		return nil, err
	}
	result := &OpenGidToChatIdResponse{}
	if err = util.DecodeWithError(response, result, "OpenGidToChatId"); err != nil {
		return nil, err
	}
	return result, nil
}
