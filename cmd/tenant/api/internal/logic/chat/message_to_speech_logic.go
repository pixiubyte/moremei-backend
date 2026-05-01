package chat

import (
	"context"
	"fmt"
	"net/http"

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

type MessageToSpeechLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	w      http.ResponseWriter
}

func NewMessageToSpeechLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter) *MessageToSpeechLogic {
	return &MessageToSpeechLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
	}
}

func (l *MessageToSpeechLogic) MessageToSpeech(req *types.MessageToSpeechRequest) error {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return err
	}

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

	conversation, err := l.svcCtx.ConversationModel.FindOne(l.ctx, message.ConversationId)
	if err != nil {
		return errors.Errorf("find conversation %d err:%v", message.ConversationId, err)
	}

	aiRobot, err := l.svcCtx.AiRobotModel.FindOne(l.ctx, conversation.AiRobotId)
	if err != nil {
		return errors.Errorf("find ai robot %d err:%v", conversation.AiRobotId, err)
	}

	// 查询user uuid
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		return err
	}

	// 修复text中的特殊字符
	text := fixSpecialChars(message.Content)

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
