package chat

import (
	"context"
	"encoding/json"
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
	"moremei/ai-saas/pkg/messages"
	httpx "moremei/ai-saas/pkg/x/http"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"gorm.io/gorm"
)

type SendRobotMessageSSELogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	w      http.ResponseWriter
}

func NewSendRobotMessageSSELogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter) *SendRobotMessageSSELogic {
	return &SendRobotMessageSSELogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
	}
}

func (l *SendRobotMessageSSELogic) SendRobotMessageSSE(req *types.SendRobotMessageRequest) error {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return err
	}
	// 查询用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		return errors.Errorf("find user %d err:%v", authUser.UserId, err)
	}
	// 查询机器人
	aiRobot, err := l.svcCtx.AiRobotModel.FindOneByUuid(l.ctx, req.AiRobotId)
	if err != nil {
		return errors.Errorf("find ai %s err:%v", req.AiRobotId, err)
	}
	// 获取会话，如果没有则创建
	conversation, err := l.getConversation(user, aiRobot, req)
	if err != nil {
		return err
	}

	// 消息类型处理
	dmsg := &dify.ChatMessageRequest{
		User:           user.Uuid,
		ConversationID: conversation.Cid,
		Inputs:         req.Inputs,
	}
	// 转换消息
	req_msg, err := messages.ParseMsgToAiType(req.ContentType, req.Content)
	if err != nil {
		return errors.Errorf("convert message `%v` to dify message err: %v", req, err)
	}
	// 按dify格式设置参数
	setMessageToDifyRequest(req_msg, dmsg)

	// 注册SSE服务
	dataChan := make(chan []byte, 1)

	threading.GoSafeCtx(l.ctx, func() {
		httpx.RegisterSSE(l.ctx, l.w, dataChan, func() {
			close(dataChan)
		})
	})

	// 初始化dify客户端
	c := dify.NewClient(aiRobot.ServerHost, aiRobot.ApiKey, dify.WithThinking(true))

	// 发送消息
	ch, err := c.ChatMessagesStream(l.ctx, *dmsg)
	if err != nil {
		return errors.Errorf("user %d chat %d message stream error: %v", authUser.UserId, conversation.Id, err)
	}

	var (
		cid         string
		answerStrBd strings.Builder
		msgType     messages.MessageType // Ai返回的消息类型, 一次请求只返回一种类型的消息
	)
	// 提前生成消息id，用于流式返回使用
	userMsgUuid := l.svcCtx.Snowflake.Generate()
	aiMsgUuid := l.svcCtx.Snowflake.Generate()

	for {
		select {
		case <-l.ctx.Done():
			if answerStrBd.Len() > 0 {
				goto M
			} else {
				return errors.Errorf("user %d chat %d canceled but empty response", authUser.UserId, conversation.Id)
			}
		case r, isOpen := <-ch:
			if err = r.Err; err != nil {
				// 已经有返回数据时出错
				if answerStrBd.Len() > 0 {
					goto M
				} else {
					return err
				}
			}
			// 结束
			if !isOpen {
				if answerStrBd.Len() > 0 {
					// 手动发送关闭消息
					dataChan <- []byte(httpx.DoneData)

					goto M
				} else {
					return errors.Errorf("user %d chat %d message stream is closed but empty response", authUser.UserId, conversation.Id)
				}
			}

			cid = r.ConversationID

			// 解析消息类型
			msg := messages.FromRawContent(r.Answer)
			if !(msgType > messages.MessageTypeUnknown) {
				// 未定义
				msgType = msg.Type
			}
			answerStrBd.WriteString(msg.ContentString())

			answer := &messages.ChatStreamMessage{
				Type:           msgType,
				Id:             aiMsgUuid,
				Content:        msg.ContentString(),
				ConversationId: r.ConversationID,
				RawMsgId:       r.ID, // 原始消息id(来自dify)
				CreatedTime:    time.Now().Unix(),
			}
			b, _ := json.Marshal(answer)
			dataChan <- b
		}
	}

M:
	// 解析用户消息类型
	umsgContent := req.Content
	umsg := messages.FromRawContent(req.Content)
	if umsg.Type == messages.MessageType(req.ContentType) {
		umsgContent = umsg.ContentString()
	}

	parentMessage := &tenant.Message{
		Uuid:           userMsgUuid,
		ConversationId: conversation.Id,
		ParentId:       0,
		AuthorType:     tenant.MessageAuthorTypeUser,
		AuthorId:       authUser.UserId,
		Source:         req.MsgType,
		Type:           req.ContentType,
		Content:        umsgContent,
	}
	message := tenant.Message{
		Uuid:           aiMsgUuid,
		ConversationId: conversation.Id,
		AuthorType:     tenant.MessageAuthorTypeAI,
		AuthorId:       aiRobot.Id,
		Source:         req.MsgType,
		Type:           int64(msgType),
		Content:        answerStrBd.String(),
	}
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {

		if conversation.Cid == "" && cid != "" {
			conversation.Cid = cid
			err = l.svcCtx.ConversationModel.Update(l.ctx, conversation, tx)
			if err != nil {
				return errors.Errorf("update conversation %d cid %s error: %s", conversation.Id, conversation.Cid, err)
			}
		}

		err = l.svcCtx.MessageModel.Insert(l.ctx, parentMessage, tx)
		if err != nil {
			return errors.Errorf("user %d chat question message insert err:%v", authUser.UserId, err)
		}
		message.ParentId = parentMessage.Id
		err = l.svcCtx.MessageModel.Insert(l.ctx, &message, tx)
		if err != nil {
			return errors.Errorf("user %d chat answer message insert err:%v", authUser.UserId, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (l *SendRobotMessageSSELogic) getConversation(user *tenant.User, aiRobot *tenant.AiRobot, req *types.SendRobotMessageRequest) (conversation *tenant.Conversation, err error) {

	conversationModel := l.svcCtx.ConversationModel
	conversation, err = conversationModel.FindByUserIdAiRobotId(l.ctx, user.Id, aiRobot.Id)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return nil, errors.Errorf("user %d find ai %d conversation err: %v", user.Id, aiRobot.Id, err)
		}
		if !commonConst.IsFlagTrue(aiRobot.IsPublic) {
			_, err = l.svcCtx.AiRobotWhitelistModel.FindOneByAiRobotIdUserId(l.ctx, aiRobot.Id, user.Id)
			if err != nil {
				if !errors.Is(err, gormc.ErrNotFound) {
					return nil, err
				}
				return nil, errs.AiRobotNotFoundError
			}
		}
		conversation = &tenant.Conversation{
			Uuid:      l.svcCtx.Snowflake.Generate(),
			UserId:    user.Id,
			AiRobotId: aiRobot.Id,
		}
		err = l.svcCtx.ConversationModel.Insert(l.ctx, conversation, nil)
		if err != nil {
			return nil, errors.Errorf("user %d create ai %d conversation err: %v", user.Id, aiRobot.Id, err)
		}
	} else if conversation.UserId != user.Id {
		return nil, errs.ConversationNotFoundError
	}
	return conversation, nil
}

// 将消息转换为dify请求格式
func setMessageToDifyRequest(req *messages.AiMessage, dmsg *dify.ChatMessageRequest) {
	dmsg.Query = req.Query
	dmsg.Inputs = req.Inputs
	if dmsg.Files == nil {
		dmsg.Files = make([]*dify.ChatFile, 0)
	}
	for _, f := range req.Files {
		dmsg.Files = append(dmsg.Files, &dify.ChatFile{
			Type:           req.FileType,
			Url:            f,
			TransferMethod: dify.ChatFileTransferMethodRemote,
		})
	}
}
