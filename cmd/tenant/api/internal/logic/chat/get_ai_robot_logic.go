package chat

import (
	"context"
	"errors"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/pkg/consts"
	"moremei/ai-saas/pkg/gormc"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAiRobotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAiRobotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAiRobotLogic {
	return &GetAiRobotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAiRobotLogic) GetAiRobot(req *types.GetAiRobotRequest) (resp *types.AiRobot, err error) {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	robot, err := l.svcCtx.AiRobotModel.FindOneByUuid(l.ctx, req.RobotId)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return nil, err
		}
		return nil, nil // 找不到机器人也返回 nil 数据
	}
	if robot.Status == consts.FlagFalse {
		return nil, nil // 机器人已禁用, 返回 nil 数据
	}
	if robot.IsPublic == consts.FlagFalse {
		// 检查是否在白名单中
		whitelist, err := l.svcCtx.AiRobotWhitelistModel.FindAiRobotIdsByUserId(l.ctx, authUser.UserId)
		if err != nil {
			return nil, err
		}
		_, ok := lo.Find(whitelist, func(item uint64) bool {
			return item == robot.Id
		})
		if !ok {
			return nil, nil // 不在白名单中, 返回 nil 数据
		}
	}

	return &types.AiRobot{
		Id:        robot.Uuid,
		Name:      robot.Name,
		Desc:      robot.Description,
		Category:  robot.Category,
		IsDefault: consts.IsFlagTrue(robot.IsDefault),
		Avator:    robot.Avator,
	}, nil
}
