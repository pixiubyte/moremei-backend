package chat

import (
	"context"

	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	commonConst "moremei/ai-saas/pkg/consts"
	"moremei/ai-saas/pkg/dify"
	"moremei/ai-saas/pkg/gormc"
	errorx "moremei/ai-saas/pkg/x/error"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetRobotOpeningStatementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRobotOpeningStatementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRobotOpeningStatementLogic {
	return &GetRobotOpeningStatementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRobotOpeningStatementLogic) GetRobotOpeningStatement(req *types.GetRobotOpeningStatementRequest) (resp *types.GetRobotOpeningStatementResponse, err error) {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		return nil, err
	}
	// 查询机器人
	aiRobot, err := l.svcCtx.AiRobotModel.FindOneByUuid(l.ctx, req.RobotId)
	if err != nil {
		if errors.Is(err, gormc.ErrNotFound) { // 未找到机器人
			return nil, errorx.WithCause(errorx.InvalidParamsError, "ai robot not found")
		}
		return nil, err
	}
	// 检查机器人是否可用
	if !commonConst.IsFlagTrue(aiRobot.IsPublic) {
		_, err = l.svcCtx.AiRobotWhitelistModel.FindOneByAiRobotIdUserId(l.ctx, aiRobot.Id, authUser.UserId)
		if err != nil {
			if !errors.Is(err, gormc.ErrNotFound) {
				return nil, err
			}
			return nil, errs.AiRobotNotFoundError
		}
	}
	// 查询机器人的 opening statement
	c := dify.NewClient(aiRobot.ServerHost, aiRobot.ApiKey)

	res, err := c.Parameters(l.ctx, dify.ParametersRequest{
		User: user.Uuid,
	})
	if err != nil {
		return nil, errors.New("获取机器人 opening statement 失败")
	}
	qures := lo.Map(res.SuggestedQuestions, func(item any, index int) string {
		return item.(string)
	})
	resp = &types.GetRobotOpeningStatementResponse{
		OpeningStatement: res.OpeningStatement,
		Suggested:        qures,
	}

	return resp, nil
}
