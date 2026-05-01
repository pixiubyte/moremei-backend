package chat

import (
	"context"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationsLogic) GetConversations() (resp *types.GetConversationsResponse, err error) {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	userConversations, err := l.svcCtx.ConversationModel.FindByUserId(l.ctx, authUser.UserId)
	if err != nil {
		return nil, err
	}

	aiRobotIds := lo.Map[*tenant.UserConversation, uint64](userConversations, func(m *tenant.UserConversation, index int) uint64 {
		return m.AiRobotId
	})

	aiRobots, err := l.svcCtx.AiRobotModel.FindByIds(l.ctx, aiRobotIds)
	if err != nil {
		return nil, err
	}
	aiRobotMap := lo.SliceToMap[*tenant.AiRobot, uint64, *tenant.AiRobot](aiRobots, func(m *tenant.AiRobot) (uint64, *tenant.AiRobot) {
		return m.Id, m
	})

	return &types.GetConversationsResponse{
		List: lo.Map[*tenant.UserConversation, *types.Conversation](userConversations, func(m *tenant.UserConversation, index int) *types.Conversation {
			aiRobot := aiRobotMap[m.AiRobotId]
			return &types.Conversation{
				Id:          m.Uuid,
				CreatedTime: m.CreatedAt.Unix(),
				AiRobot: types.AiRobot{
					Id:       aiRobot.Uuid,
					Name:     aiRobot.Name,
					Desc:     aiRobot.Description,
					Category: aiRobot.Category,
				},
			}
		}),
	}, nil
}
