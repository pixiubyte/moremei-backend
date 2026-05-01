package messages

type MessageType int64

// 消息类型
const (
	MessageTypeUnknown             MessageType = 0 // 未知类型
	MessageTypeText                MessageType = 1 // 文本
	MessageTypeImage               MessageType = 2
	MessageTypeAudio               MessageType = 3
	MessageTypeVideo               MessageType = 5
	MessageTypeMiniProgram         MessageType = 7
	MessageTypeLink                MessageType = 8
	MessageTypeFile                MessageType = 9
	MessageTypeQuote               MessageType = 15  // 引用
	MessageTypeImageText           MessageType = 102 // 图文
	MessageTypeFiles               MessageType = 109 // 文件
	MessageTypeList                MessageType = 110 // 列表
	MessageTypeMultiList           MessageType = 111 // 多列表
	MessageTypeCards               MessageType = 112 // 卡片
	MessageTypeChoices             MessageType = 113 // 选择题
	MessageTypeQuestionnaire       MessageType = 114 // 问卷
	MessageTypeQuestionnaireAnswer MessageType = 115 // 问卷回答
	MessageTypeSystemNotice        MessageType = 400 // 系统通知
	MessageTypeCallingCtl          MessageType = 401 // 通话控制
)

const (
	MessageFileTypeImage    = "image"
	MessageFileTypeDocument = "document"
)
