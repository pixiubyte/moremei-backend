package messages

type (
	// 文本消息 type: 1
	Text string

	// 图文消息 type: 102
	ImageText struct {
		ImgUrls []string `json:"img_urls"` // 图片链接
		Text    string   `json:"text"`     // 文本内容
	}

	// 文件消息 type: 109
	Files struct {
		FileUrls []string `json:"file_urls"` // 文件链接
		Text     string   `json:"text"`      // 文本内容
	}

	ListItem struct {
		Id    uint64 `json:"id"`          // 资源id
		Uuid  string `json:"uuid"`        // 资源uuid
		Title string `json:"title"`       // 标题
		Desc  string `json:"description"` // 描述
		Image string `json:"image"`       // 图片链接
		Path  string `json:"path"`        // 跳转路径
		Url   string `json:"url"`
	}

	// 列表消息 type: 110
	List struct {
		Text     string      `json:"text"`      // 文本内容
		More     string      `json:"more"`      // 更多链接
		MoreTo   string      `json:"more_to"`   // 更多跳转链接
		ViewType string      `json:"view_type"` // 列表视图类型 simple|detailed
		Items    []*ListItem `json:"items"`     // 列表项
	}

	// 多列表消息 type: 111
	MultiList struct {
		Text     string  `json:"text"`      // 文本内容
		ViewType string  `json:"view_type"` // 列表视图类型
		List     []*List `json:"list"`      // 多列表
	}

	// 卡片消息 type: 112
	Cards struct {
		Text     string   `json:"text"`      // 文本内容
		ViewType string   `json:"view_type"` // 列表视图类型
		Item     ListItem `json:"item"`      // 卡片项
	}

	ChoiceOption struct {
		Label    string          `json:"label"`    // 选项标签
		Score    float64         `json:"score"`    // 选项得分
		Text     string          `json:"text"`     // 选项文本
		Value    string          `json:"value"`    // 选项值
		Children []*ChoiceOption `json:"children"` // 子选项
	}

	// 选择题消息 type: 113
	Choices struct {
		Uuid     string          `json:"uuid"`     // 问题uuid
		Question string          `json:"question"` // 问题
		Multiple bool            `json:"multiple"` // 是否多选
		Tips     string          `json:"tips"`     // 提示
		Options  []*ChoiceOption `json:"options"`  // 选项
	}

	// 问卷消息 type: 114
	Questionnaire struct {
		Uuid      string    `json:"uuid"`       // 问卷uuid
		Text      string    `json:"text"`       // 文本内容
		Title     string    `json:"title"`      // 标题
		Category  string    `json:"category"`   // 分类
		QuestType string    `json:"quest_type"` // 问卷类型：func 功能型，chat 聊天型
		Questions []Choices `json:"questions"`  // 问题列表
	}

	ChoiceAnswer struct {
		Question string   `json:"question"` // 问题
		Answers  []string `json:"answers"`  // 回答
	}

	// 问卷回答消息 type: 115
	QuestionnaireAnswer struct {
		Uuid      string          `json:"uuid"`       // 问卷uuid
		Title     string          `json:"title"`      // 标题
		Category  string          `json:"category"`   // 分类
		QuestType string          `json:"quest_type"` // 问卷类型：func 功能型，chat 聊天型
		Questions []*ChoiceAnswer `json:"questions"`  // 回答列表
	}

	// 系统通知消息 type: 400
	SystemNotice struct {
		Type    int64  `json:"type"`    // 通知类型
		Title   string `json:"title"`   // 标题
		Content string `json:"content"` // 内容
	}

	// 通话交互消息 type: 401
	CallingCtrl struct {
		SubType     int64     `json:"sub_type"`     // 交互类型
		Mode        string    `json:"mode"`         // 模式
		Content     string    `json:"content"`      // 内容
		ContentUuid string    `json:"content_uuid"` // 内容uuid
		Questions   []Choices `json:"questions"`    // 问题列表
	}
)

// 按dify消息格式定义，包含query(用户输入)、inputs(输入参数)和fiels(文件输入)
type AiMessage struct {
	Inputs   map[string]interface{} `json:"inputs"`    // 输入参数
	Query    string                 `json:"query"`     // 用户输入
	Files    []string               `json:"files"`     // 文件列表
	FileType string                 `json:"file_type"` // 文件类型, 每次只支持一种类型的文件，image|document
}

// 流式返回格式
type ChatStreamMessage struct {
	Type           MessageType `json:"type"`                 // 消息类型
	Id             string      `json:"id"`                   // 消息id(前端使用)
	Content        string      `json:"content"`              // 消息内容
	ConversationId string      `json:"conversation_id"`      // 会话id
	RawMsgId       string      `json:"raw_msg_id,omitempty"` // 原始消息id(来自dify)
	CreatedTime    int64       `json:"created_time"`         // 创建时间
}

// 消息
type Message struct {
	Type    MessageType `json:"type"`    // 消息类型
	Content any         `json:"content"` // 消息内容
}
