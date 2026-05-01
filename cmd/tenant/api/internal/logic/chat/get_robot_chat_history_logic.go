package chat

import (
	"context"
	"errors"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/gormc"
	"moremei/ai-saas/pkg/pagination"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetRobotChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRobotChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRobotChatHistoryLogic {
	return &GetRobotChatHistoryLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetRobotChatHistoryLogic) GetRobotChatHistory(req *types.GetRobotChatHistoryRequest) (*types.GetRobotChatHistoryResponse, error) {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	aiRobot, err := l.svcCtx.AiRobotModel.FindOneByUuid(l.ctx, req.AiRobotId)
	if err != nil {
		return nil, err
	}
	conversation, err := l.svcCtx.ConversationModel.FindByUserIdAiRobotId(l.ctx, authUser.UserId, aiRobot.Id)
	if err != nil {
		if errors.Is(err, gormc.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	queryParams := map[string]interface{}{"conversation_id": conversation.Id, "source": req.MsgType}
	messageModel := l.svcCtx.MessageModel
	if req.LatestNode != "" {
		latestMessage, err := messageModel.FindOneByUuid(l.ctx, req.LatestNode)
		if err != nil {
			if !errors.Is(err, gormc.ErrNotFound) {
				return nil, err
			}
			return &types.GetRobotChatHistoryResponse{CurrentNode: "", Messages: []*types.RobotIdChatHistory{}}, nil
		}
		queryParams["latest_id"] = latestMessage.Id
	}
	response, err := messageModel.PageFind(l.ctx, pagination.NewPageRequest(1, req.PageSize, queryParams))
	if err != nil {
		return nil, err
	}
	models := response.List.([]*tenant.Message)
	if len(models) == 0 {
		return &types.GetRobotChatHistoryResponse{CurrentNode: "", Messages: []*types.RobotIdChatHistory{}}, nil
	}
	pmIds := make([]uint64, 0)
	messageIds := make([]uint64, 0)
	for _, m := range models {
		if m.ParentId > 0 {
			pmIds = append(pmIds, m.ParentId)
			messageIds = append(messageIds, m.Id)
		}
	}
	pmMap := map[uint64]*tenant.Message{}
	if len(pmIds) > 0 {
		parentMessages, err := messageModel.FindByIds(l.ctx, pmIds)
		if err != nil {
			return nil, err
		}
		pmMap = lo.SliceToMap(parentMessages, func(m *tenant.Message) (uint64, *tenant.Message) { return m.Id, m })
	}
	moMap := map[uint64]int64{}
	if len(messageIds) > 0 {
		messageOpinions, err := l.svcCtx.MessageOpinionModel.FindByUserIdMsgIds(l.ctx, authUser.UserId, messageIds)
		if err != nil {
			return nil, err
		}
		moMap = lo.SliceToMap(messageOpinions, func(m *tenant.MessageOpinion) (uint64, int64) { return m.MessageId, m.Status })
	}
	messages := make([]*types.RobotIdChatHistory, 0)
	for _, v := range models {
		if v.ParentId == 0 {
			continue
		}
		opinion, ok := moMap[v.Id]
		if !ok {
			opinion = tenant.MessageOpinionStatusCancel
		}
		messages = append(messages, &types.RobotIdChatHistory{Id: v.Uuid, AiRobotId: req.AiRobotId, ParentId: pmMap[v.ParentId].Uuid, ParentContent: pmMap[v.ParentId].Content, ParentCreatedTime: pmMap[v.ParentId].CreatedAt.Unix(), AuthorType: v.AuthorType, Content: v.Content, Type: v.Type, CreatedTime: v.CreatedAt.Unix(), Opinion: opinion})
	}
	return &types.GetRobotChatHistoryResponse{CurrentNode: models[len(models)-1].Uuid, Messages: messages}, nil
}
