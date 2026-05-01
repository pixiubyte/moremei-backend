package chat

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/dify"
	"moremei/ai-saas/pkg/gormc"
	errx "moremei/ai-saas/pkg/x/error"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type TextToSpeechLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	w      http.ResponseWriter
}

func NewTextToSpeechLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter) *TextToSpeechLogic {
	return &TextToSpeechLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
	}
}

func (l *TextToSpeechLogic) TextToSpeech(req *types.TextToSpeechRequest) error {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return err
	}

	var aiRobot *tenant.AiRobot

	text := req.Text // 默认使用 text

	if req.AiRobotId != "" { // 存在 robotId 则从 robot 中获取
		aiRobot, err = l.svcCtx.AiRobotModel.FindOneByUuid(l.ctx, req.AiRobotId)
		if err != nil {
			return err
		}
	}

	if req.MessageId != "" { // 存在 messageId 则从 message 中获取

		message, err := l.svcCtx.MessageModel.FindOneByUuid(l.ctx, req.MessageId)
		if err != nil {
			if !errors.Is(err, gormc.ErrNotFound) {
				return err
			}
			return errx.DataNotFoundError
		}
		if message.AuthorType != tenant.MessageAuthorTypeAI {
			return errx.ForbiddenError
		}

		text = message.Content

		if len(req.AiRobotId) == 0 { // 不存在 robotId 则从 message 中获取
			conversation, err := l.svcCtx.ConversationModel.FindOne(l.ctx, message.ConversationId)
			if err != nil {
				return errors.Errorf("find conversation %d err:%v", message.ConversationId, err)
			}

			aiRobot, err = l.svcCtx.AiRobotModel.FindOne(l.ctx, conversation.AiRobotId)
			if err != nil {
				return errors.Errorf("find ai robot %d err:%v", conversation.AiRobotId, err)
			}
		}
	}

	if aiRobot == nil { // 都不存在则使用默认 robot
		aiRobot, err = l.svcCtx.AiRobotModel.FindDefaultOne(l.ctx)
		if err != nil {
			return err
		}
	}

	if text == "" {
		return errx.WithCause(errx.ForbiddenError, "text or message_id is empty")
	}

	// 修复text中的特殊字符
	text = fixSpecialChars(text)

	// 查询user uuid
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		return err
	}

	client := dify.NewClient(aiRobot.ServerHost, aiRobot.ApiKey)
	if req.Stream {
		ch, err := client.TextToAudioStream(l.ctx, dify.TextToAudioRequest{Text: text, User: user.Uuid, Stream: req.Stream})
		if err != nil {
			return err
		}
		writeAudioHeader(l.w)
		flusher, ok := l.w.(http.Flusher)
		for {
			select {
			case <-l.ctx.Done():
				return nil
			case data, isOpen := <-ch:
				if len(data) > 0 {
					_, err = l.w.Write(data)
					if err != nil {
						l.Errorf("TextToAudioStream %s write data err:%v", req.MessageId, err)
					} else if ok {
						flusher.Flush()
					}
				}
				if !isOpen {
					return nil
				}
			}
		}
	} else {
		data, err := client.TextToAudio(l.ctx, dify.TextToAudioRequest{Text: text, User: user.Uuid, Stream: req.Stream})
		if err != nil {
			l.Errorf("TextToAudio %s err:%v", req.MessageId, err)
			return err
		}
		if len(data) == 0 {
			return fmt.Errorf("TextToAudio failed, data is empty. by user `%v`", authUser.UserId)
		}
		writeAudioHeader(l.w)
		_, err = l.w.Write(data)
		if err != nil {
			l.w.WriteHeader(http.StatusInternalServerError)
			l.Errorf("TextToAudioStream %s write data err:%v", req.MessageId, err)
		}
		return nil
	}
}

func writeAudioHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")
}

// 修复text中的特殊字符
func fixSpecialChars(text string) string {
	text = strings.ReplaceAll(text, "—", "-") // 中文破折号会被语音模型忽略

	// 处理Markdown链接格式：移除URL部分，只保留[]中的内容
	// 匹配格式 [文本](URL) 并替换为 "文本"
	re := regexp.MustCompile(`\[(.*?)\]\([^\)]+\)`)
	text = re.ReplaceAllString(text, "$1")

	return text
}
