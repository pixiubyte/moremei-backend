package chat

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/jsonx"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/pkg/gormc"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShareRobotChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShareRobotChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShareRobotChatHistoryLogic {
	return &GetShareRobotChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShareRobotChatHistoryLogic) GetShareRobotChatHistory(req *types.GetShareRobotChatHistoryRequest) (resp *types.GetShareRobotChatHistoryResponse, err error) {
	// 获取当前登录用户的租户信息
	if _, err = utils.GetAuthUserCtx(l.ctx); err != nil {
		return nil, err
	}

	// 根据分享ID查询记录
	shareHistory, err := l.svcCtx.SharedChatHistoryModel.FindByUUID(l.ctx, req.Id)
	if err != nil {
		if errors.Is(err, gormc.ErrNotFound) {
			return nil, err
		}
		return nil, err
	}

	// 解析消息列表
	var messages []types.ChatHistoryMessageResponse
	err = jsonx.UnmarshalFromString(shareHistory.Messages, &messages)
	if err != nil {
		return nil, err
	}

	return &types.GetShareRobotChatHistoryResponse{
		Messages: messages,
		ShareInfo: types.ShareInfo{
			CreatedAt:  shareHistory.CreatedAt.Unix(),
			SharerName: shareHistory.UserName,
			SharerId:   shareHistory.Uuid,
		},
	}, nil

}
