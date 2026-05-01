package chat

import (
	"context"

	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/gormc"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type OpinionConversationMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpinionConversationMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpinionConversationMessageLogic {
	return &OpinionConversationMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpinionConversationMessageLogic) OpinionConversationMessage(req *types.OpinionConversationMessageRequest) error {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return err
	}

	messageModel := l.svcCtx.MessageModel
	message, err := messageModel.FindOneByUuid(l.ctx, req.MessageId)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return err
		}
		return errs.MessageNotFoundError
	}
	if message.AuthorType != tenant.MessageAuthorTypeAI {
		return errs.MessageNotFoundError
	}

	conversation, err := l.svcCtx.ConversationModel.FindOne(l.ctx, message.ConversationId)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return err
		}
		return errs.MessageNotFoundError
	}

	if conversation.UserId != authUser.UserId {
		return errs.MessageNotFoundError
	}

	var opinionStatus int64
	switch req.Action {
	case consts.MessageOpinionActionThumbsUp:
		opinionStatus = tenant.MessageOpinionStatusThumbsUp
	case consts.MessageOpinionActionThumbsDown:
		opinionStatus = tenant.MessageOpinionStatusThumbsDown
	case consts.MessageOpinionActionCancel:
		opinionStatus = tenant.MessageOpinionStatusCancel
	}

	messageOpinionModel := l.svcCtx.MessageOpinionModel

	messageOpinion, err := messageOpinionModel.FindOneByMessageIdUserId(l.ctx, message.Id, authUser.UserId)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return err
		}

		messageOpinion = &tenant.MessageOpinion{
			MessageId: message.Id,
			UserId:    authUser.UserId,
			Status:    opinionStatus,
		}

		return messageOpinionModel.Insert(l.ctx, messageOpinion, nil)
	}

	messageOpinion.Status = opinionStatus
	return messageOpinionModel.Update(l.ctx, messageOpinion, nil)
}
