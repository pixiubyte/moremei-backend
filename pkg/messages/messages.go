package messages

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/jsonx"
)

// 将用户消息按类型解析为dify消息格式
func ParseMsgToAiType(msgType int64, content string) (msg *AiMessage, err error) {

	msg = &AiMessage{}

	switch MessageType(msgType) {
	case MessageTypeText: // 文本消息
		msg.Query = content

	case MessageTypeImageText: // 图文消息
		imgText, err := unmarshalMessage[ImageText](content)
		if err == nil {
			msg.Query = imgText.Text
			msg.Files = imgText.ImgUrls
			msg.FileType = MessageFileTypeImage
		} else {
			return nil, fmt.Errorf("unmarshal type `%v` message `%v` failed: %v", msg, content, err)
		}

	case MessageTypeFiles: // 文件消息
		files, err := unmarshalMessage[Files](content)
		if err == nil {
			msg.Files = files.FileUrls
			msg.FileType = MessageFileTypeDocument
		} else {
			return nil, fmt.Errorf("unmarshal type `%v` message `%v` failed: %v", msg, content, err)
		}

	case MessageTypeQuestionnaireAnswer: // 问卷回答消息
		qa, err := unmarshalMessage[QuestionnaireAnswer](content)
		if err == nil {
			// 兼容旧版消息格式（消息体中添加了type）
			if qa.Category == "" && qa.Uuid == "" && len(qa.Questions) == 0 {
				// 消息实际未解析成功
				msg.Query = content
			} else {
				typedMsg := Message{
					Type:    MessageTypeQuestionnaireAnswer,
					Content: qa,
				}
				tmstr, _ := jsonx.MarshalToString(typedMsg)
				msg.Query = tmstr
			}
		} else {
			return nil, fmt.Errorf("unmarshal type `%v` message `%v` failed: %v", msg, content, err)
		}
	case MessageTypeCallingCtl: // 通话控制消息
		ctl, err := unmarshalMessage[CallingCtrl](content)
		if err == nil {
			typedMsg := Message{
				Type:    MessageTypeCallingCtl,
				Content: ctl,
			}
			tmstr, _ := jsonx.MarshalToString(typedMsg)
			msg.Query = tmstr
		} else {
			return nil, fmt.Errorf("unmarshal type `%v` message `%v` failed: %v", msg, content, err)
		}

	default:
		return nil, fmt.Errorf("unsupported message type for user send: %d", msgType)
	}

	return msg, nil
}

func unmarshalMessage[T any](content string) (*T, error) {
	var msg T
	err := jsonx.UnmarshalFromString(content, &msg)
	return &msg, err
}

// 解析原始消息。多数情况为Ai返回的消息
func FromRawContent(content string) (msg *Message) {

	if len(content) > 0 && content[0] == '{' {
		// 尝试解析为json
		msg = &Message{}
		err := jsonx.UnmarshalFromString(content, msg)
		if err == nil {
			return msg
		}
	}
	// 其他情况，视为文本消息
	return &Message{
		Type:    MessageTypeText,
		Content: content,
	}
}

// ContentString 消息内容转换为字符串
func (m *Message) ContentString() string {
	if m.Type == MessageTypeText {
		cot := m.Content.(string)
		return cot
	}
	cstr, err := jsonx.MarshalToString(m.Content)
	if err != nil {
		return ""
	}
	return cstr
}
