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

type GetConversationMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationMessagesLogic {
	return &GetConversationMessagesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetConversationMessagesLogic) GetConversationMessages(req *types.GetConversationMessagesRequest) (*types.GetConversationMessagesResponse, error) {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	conversation, err := l.svcCtx.ConversationModel.FindOneByUuid(l.ctx, req.ConversationId)
	if err != nil {
		return nil, err
	}
	queryParams := map[string]interface{}{"conversation_id": conversation.Id, "source": tenant.MessageSourceWeb}
	messageModel := l.svcCtx.MessageModel
	if req.LatestNode != "" {
		latestMessage, err := messageModel.FindOneByUuid(l.ctx, req.LatestNode)
		if err != nil {
			if !errors.Is(err, gormc.ErrNotFound) {
				return nil, err
			}
			return &types.GetConversationMessagesResponse{CurrentNode: "", Messages: []*types.Message{}}, nil
		}
		queryParams["latest_id"] = latestMessage.Id
	}
	response, err := messageModel.PageFind(l.ctx, pagination.NewPageRequest(1, req.PageSize, queryParams))
	if err != nil {
		return nil, err
	}
	models := response.List.([]*tenant.Message)
	if len(models) == 0 {
		return &types.GetConversationMessagesResponse{CurrentNode: "", Messages: []*types.Message{}}, nil
	}
	pmIds := make([]uint64, 0)
	messageIds := make([]uint64, 0, len(models))
	for _, m := range models {
		if m.ParentId > 0 {
			pmIds = append(pmIds, m.ParentId)
		}
		messageIds = append(messageIds, m.Id)
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
	return &types.GetConversationMessagesResponse{
		CurrentNode: models[0].Uuid,
		Messages: lo.Map(models, func(m *tenant.Message, _ int) *types.Message {
			opinion, ok := moMap[m.Id]
			if !ok {
				opinion = tenant.MessageOpinionStatusCancel
			}
			msg := &types.Message{Id: m.Uuid, ConversationId: conversation.Uuid, AuthorType: m.AuthorType, Content: m.Content, Type: m.Type, Source: m.Source, Opinion: opinion, CreatedTime: m.CreatedAt.Unix()}
			if m.ParentId > 0 {
				if parentMessage, ok := pmMap[m.ParentId]; ok {
					msg.ParentId = parentMessage.Uuid
					msg.ParentContent = parentMessage.Content
				}
			}
			return msg
		}),
		}, nil
}
