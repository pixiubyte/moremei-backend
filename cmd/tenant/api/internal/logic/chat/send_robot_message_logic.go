package chat

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	commonConst "moremei/ai-saas/pkg/consts"
	"moremei/ai-saas/pkg/dify"
	"moremei/ai-saas/pkg/gormc"
	errorx "moremei/ai-saas/pkg/x/error"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendRobotMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	w      http.ResponseWriter
}

func NewSendRobotMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter) *SendRobotMessageLogic {
	return &SendRobotMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
	}
}

func (l *SendRobotMessageLogic) SendRobotMessage(req *types.SendRobotMessageRequest) error {
	var cid string
	flusher, ok := l.w.(http.Flusher)
	if !ok {
		return errorx.InternalServerError
	}
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return err
	}

	aiRobot, err := l.svcCtx.AiRobotModel.FindOneByUuid(l.ctx, req.AiRobotId)
	if err != nil {
		return err
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		return err
	}

	conversationModel := l.svcCtx.ConversationModel
	conversation, err := conversationModel.FindByUserIdAiRobotId(l.ctx, authUser.UserId, aiRobot.Id)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return err
		}
		if !commonConst.IsFlagTrue(aiRobot.IsPublic) {
			_, err = l.svcCtx.AiRobotWhitelistModel.FindOneByAiRobotIdUserId(l.ctx, aiRobot.Id, authUser.UserId)
			if err != nil {
				if !errors.Is(err, gormc.ErrNotFound) {
					return err
				}
				return errs.AiRobotNotFoundError
			}
		}
		conversation = &tenant.Conversation{
			Uuid:      l.svcCtx.Snowflake.Generate(),
			UserId:    user.Id,
			AiRobotId: aiRobot.Id,
		}
		err = l.svcCtx.ConversationModel.Insert(l.ctx, conversation, nil)
		if err != nil {
			return err
		}
	} else if conversation.UserId != authUser.UserId {
		return errs.ConversationNotFoundError
	}

	parentMessage := &tenant.Message{
		Uuid:           l.svcCtx.Snowflake.Generate(),
		ConversationId: conversation.Id,
		ParentId:       0,
		AuthorType:     tenant.MessageAuthorTypeUser,
		AuthorId:       authUser.UserId,
		Source:         req.MsgType,
		Type:           tenant.WorktoolQaTextTypeText,
		Content:        req.Content,
	}
	err = l.svcCtx.MessageModel.Insert(l.ctx, parentMessage, nil)
	if err != nil {
		l.Logger.Errorf("message insert err:%v", err)
		return err
	}

	l.w.Header().Set("Content-Type", "application/octet-stream")
	l.w.Header().Set("Cache-Control", "no-cache")
	l.w.Header().Set("Connection", "keep-alive")
	l.w.WriteHeader(http.StatusOK)
	flusher.Flush()

	c := dify.NewClient(aiRobot.ServerHost, aiRobot.ApiKey, dify.WithThinking(true))
	ch, err := c.ChatMessagesStream(l.ctx, dify.ChatMessageRequest{
		Query:          req.Content,
		User:           user.Uuid,
		ConversationID: conversation.Cid,
		Inputs:         req.Inputs,
	})
	if err != nil {
		return err
	}

	var strBuilder strings.Builder
	firstChunk := false
	message := tenant.Message{
		Uuid:           l.svcCtx.Snowflake.Generate(),
		ConversationId: conversation.Id,
		ParentId:       parentMessage.Id,
		AuthorType:     tenant.MessageAuthorTypeAI,
		AuthorId:       aiRobot.Id,
		Source:         req.MsgType,
		Type:           tenant.WorktoolQaTextTypeText,
	}
	answer := types.RobotIdChatMessage{
		Content:     "",
		CreatedTime: time.Now().Unix(),
	}
	for {
		select {
		case <-l.ctx.Done():
			goto M
		case r, isOpen := <-ch:
			if err = r.Err; err != nil {
				return err
			}
			if !isOpen {
				goto M
			}
			cid = r.ConversationID
			strBuilder.WriteString(r.Answer)
			answer.Content = r.Answer
			b, err := json.Marshal(answer)
			if err != nil {
				return err
			}
			if firstChunk {
				l.w.Write(b)
				firstChunk = true
			} else {
				l.w.Write([]byte("\n" + string(b)))
			}
			flusher.Flush()
		}
	}

M:
	message.Content = strBuilder.String()
	if cid == "" || message.Content == "" {
		answer.Content = "请求失败，请稍后再试"
		b, err := json.Marshal(answer)
		if err != nil {
			return err
		}
		l.w.Write(b)
		flusher.Flush()
		message.Content = answer.Content
	}
	if conversation.Cid == "" {
		conversation.Cid = cid
		err = conversationModel.Update(l.ctx, conversation, nil)
		if err != nil {
			return err
		}
	}
	err = l.svcCtx.MessageModel.Insert(l.ctx, &message, nil)
	if err != nil {
		l.Logger.Errorf("message insert err:%v", err)
		return err
	}

	return nil
}
