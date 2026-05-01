package chat

import (
	"context"
	commonConst "moremei/ai-saas/pkg/consts"

	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/gormc"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddConversationLogic {
	return &AddConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddConversationLogic) AddConversation(req *types.AddConversationRequest) (resp *types.Conversation, err error) {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	aiRobot, err := l.svcCtx.AiRobotModel.FindOneByUuid(l.ctx, req.AiRobotId)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return nil, err
		}
		return nil, errs.AiRobotNotFoundError
	}

	if !commonConst.IsFlagTrue(aiRobot.IsPublic) {
		_, err = l.svcCtx.AiRobotWhitelistModel.FindOneByAiRobotIdUserId(l.ctx, aiRobot.Id, authUser.UserId)
		if err != nil {
			if !errors.Is(err, gormc.ErrNotFound) {
				return nil, err
			}
			return nil, errs.AiRobotNotFoundError
		}
	}

	conversationModel := l.svcCtx.ConversationModel

	conversation, err := conversationModel.FindOneByUserIdAiRobotId(l.ctx, authUser.UserId, aiRobot.Id)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return nil, err
		}
		conversation = &tenant.Conversation{
			Uuid:      l.svcCtx.Snowflake.Generate(),
			UserId:    authUser.UserId,
			AiRobotId: aiRobot.Id,
		}
		err = conversationModel.Insert(l.ctx, conversation, nil)
		if err != nil {
			return nil, err
		}
	}

	return &types.Conversation{
		Id: conversation.Uuid,
		AiRobot: types.AiRobot{
			Id:       aiRobot.Uuid,
			Name:     aiRobot.Name,
			Desc:     aiRobot.Description,
			Category: aiRobot.Category,
		},
		CreatedTime: conversation.CreatedAt.Unix(),
	}, nil
}
