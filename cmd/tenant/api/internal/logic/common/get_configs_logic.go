package common

import (
	"context"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/internal/model/tenant"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GetConfigsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConfigsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigsLogic {
	return &GetConfigsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConfigsLogic) GetConfigs() (resp *types.GetConfigsResponse, err error) {
	var (
		aiRobot             *tenant.AiRobot
		appConfigsMapper    map[string]any
		memberConfigsMapper map[string]string
	)
	err = mr.Finish(func() error {
		appConfigs, err := l.svcCtx.SystemConfigModel.FindByGroup(l.ctx, tenant.SystemConfigGroupApp)
		if err != nil {
			return err
		}
		appConfigsMapper = lo.SliceToMap[*tenant.SystemConfig, string, any](appConfigs, func(item *tenant.SystemConfig) (string, any) {
			return item.Key, item.GetValue()
		})
		return nil
	}, func() error {
		memberConfigs, err := l.svcCtx.SystemConfigModel.FindByGroup(l.ctx, tenant.SystemConfigGroupMember)
		if err != nil {
			return err
		}
		memberConfigsMapper = lo.SliceToMap[*tenant.SystemConfig, string, string](memberConfigs, func(item *tenant.SystemConfig) (string, string) {
			return item.Key, item.Value
		})
		return nil
	}, func() error {
		aiRobot, err = l.svcCtx.AiRobotModel.FindDefaultOne(l.ctx)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &types.GetConfigsResponse{
		App:   appConfigsMapper,
		Wecom: types.WecomConfig{WecomCorpId: ""},
		DefaultAiRobot: types.DefaultAiRobot{
			Id:             aiRobot.Uuid,
			Name:           aiRobot.Name,
			Description:    aiRobot.Description,
			ContractImgUrl: "",
		},
		MemberConfig: types.MemberConfig{
			MemberBuyUrl: memberConfigsMapper[tenant.SystemConfigKeyMemberBuyUrl],
		},
	}, nil
}
