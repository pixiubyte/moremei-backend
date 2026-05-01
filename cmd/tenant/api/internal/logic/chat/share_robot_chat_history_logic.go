package chat

import (
	"context"
	"fmt"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/jsonx"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareRobotChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareRobotChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareRobotChatHistoryLogic {
	return &ShareRobotChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareRobotChatHistoryLogic) ShareRobotChatHistory(req *types.ShareRobotChatHistoryRequest) (resp *types.ShareRobotChatHistoryResponse, err error) {
	if len(req.Messages) == 0 {
		return nil, errors.New("messages cannot be empty")
	}

	// 获取当前登录用户
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth user")
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to find user: %d", authUser.UserId))
	}

	// 获取指定的机器人
	robot, err := l.svcCtx.AiRobotModel.FindOneByUuid(l.ctx, req.Id)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to find robot: %s", req.Id))
	}

	formattedMessages := make([]types.ChatHistoryMessage, 0, len(req.Messages))
	for _, message := range req.Messages {
		newMsg := types.ChatHistoryMessage{
			Role:        message.Role,
			Content:     message.Content,
			ContentType: message.ContentType,
			CreatedAt:   message.CreatedAt,
		}

		switch message.Role {
		case "USER":
			newMsg.UUID = user.Uuid
			newMsg.Name = user.Nickname
			newMsg.Avatar = user.Avatar
		case "AI":
			newMsg.UUID = robot.Uuid
			newMsg.Name = robot.Name
			newMsg.Avatar = robot.Avator
		}
		formattedMessages = append(formattedMessages, newMsg)
	}

	jsonData, err := jsonx.Marshal(formattedMessages)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal messages")
	}

	// 生成分享ID
	shareId := l.svcCtx.Snowflake.Generate()

	// 创建分享记录
	shareHistory := &tenant.SharedChatHistory{
		Uuid:     shareId,
		UserId:   authUser.UserId,
		UserName: user.Nickname,
		Messages: string(jsonData),
	}

	// 保存到数据库
	if err := l.svcCtx.SharedChatHistoryModel.Insert(l.ctx, shareHistory, nil); err != nil {
		return nil, errors.Wrap(err, "failed to insert share history")
	}

	return &types.ShareRobotChatHistoryResponse{
		ShareId: shareId,
	}, nil
}
