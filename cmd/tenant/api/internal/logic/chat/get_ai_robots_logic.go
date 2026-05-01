package chat

import (
	"context"
	"sort"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	commonConst "moremei/ai-saas/pkg/consts"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAiRobotsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAiRobotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAiRobotsLogic {
	return &GetAiRobotsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAiRobotsLogic) GetAiRobots(req *types.GetAiRobotsRequest) (*types.GetAiRobotsResponse, error) {
	var authUser *utils.AuthUser
	authUser, _ = utils.GetAuthUserCtx(l.ctx)
	if req.Scene == "" {
		req.Scene = tenant.AiRobotSceneDefault
	}
	aiRobots, err := l.svcCtx.AiRobotModel.FindByScene(l.ctx, req.Scene)
	if err != nil {
		return nil, err
	}
	var aiRobotIds []uint64
	if authUser != nil {
		aiRobotIds, err = l.svcCtx.AiRobotWhitelistModel.FindAiRobotIdsByUserId(l.ctx, authUser.UserId)
		if err != nil {
			return nil, err
		}
	}
	aiRobotIdsMap := lo.SliceToMap(aiRobotIds, func(item uint64) (uint64, uint64) { return item, item })
	sort.SliceStable(aiRobots, func(i, j int) bool { return aiRobots[i].Sequence > aiRobots[j].Sequence })
	list := make([]*types.AiRobot, 0)
	for _, aiRobot := range aiRobots {
		if !commonConst.IsFlagTrue(aiRobot.IsPublic) {
			if _, ok := aiRobotIdsMap[aiRobot.Id]; !ok {
				continue
			}
		}
		list = append(list, &types.AiRobot{
			Id:        aiRobot.Uuid,
			Name:      aiRobot.Name,
			Desc:      aiRobot.Description,
			Category:  aiRobot.Category,
			IsDefault: commonConst.IsFlagTrue(aiRobot.IsDefault),
			Avator:    aiRobot.Avator,
		})
	}
	return &types.GetAiRobotsResponse{List: list}, nil
}
