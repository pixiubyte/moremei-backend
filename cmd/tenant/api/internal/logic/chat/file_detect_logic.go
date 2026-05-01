package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/dify"

	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileDetectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type FileDetectResult struct {
	Category  string   `json:"category"`  // 分类
	Questions []string `json:"questions"` // 推荐问题列表
}

func NewFileDetectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDetectLogic {
	return &FileDetectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDetectLogic) FileDetect(req *types.FileDetectRequest) (resp *types.FileDetectResponse, err error) {
	if _, err = utils.GetAuthUserCtx(l.ctx); err != nil {
		return nil, err
	}

	fdr, err := l.fileDetectByAi(req.Url, req.Type)
	if err != nil {
		l.Errorf("file detect by ai error: %v", err)
		return &types.FileDetectResponse{
			Result:    false,
			Msg:       "文件不可用",
			Questions: []string{},
		}, nil
	}

	return &types.FileDetectResponse{
		Result:    true,
		Category:  fdr.Category,
		Questions: fdr.Questions,
	}, nil
}

func (l *FileDetectLogic) fileDetectByAi(url string, fileType string) (*FileDetectResult, error) {

	systemConfig, err := l.svcCtx.SystemConfigModel.FindByKey(l.ctx, tenant.SystemConfigKeyAiFileDetect)
	if err != nil {
		return nil, fmt.Errorf("sysconf file dectect key:%s, err:%v", tenant.SystemConfigKeyAiFileDetect, err)
	}
	var aiAgentConfig tenant.SysConfAgentApi
	err = json.Unmarshal([]byte(systemConfig.Value), &aiAgentConfig)
	if err != nil {
		return nil, fmt.Errorf("sysconf key:%s, config `%v` err:%v", tenant.SystemConfigKeyAiFileDetect, systemConfig, err)
	}

	caClient := dify.NewClient(aiAgentConfig.ServerHost, aiAgentConfig.ApiKey)

	// 文件传入格式
	fl := map[string]string{
		"type":            fileType,
		"transfer_method": dify.ChatFileTransferMethodRemote,
		"url":             url,
	}

	resp, err := caClient.WorkflowsRun(l.ctx, dify.WorkflowRunRequest{
		Inputs: map[string]interface{}{
			"user_file": fl, // 文件
		},
		User: "default",
	})
	if err != nil {
		return nil, fmt.Errorf("workflow 请求失败: %v", err)
	}

	if resp.Data.Status != dify.StatusSucceeded {
		return nil, fmt.Errorf("workflow 状态异常: %v", resp.Data.Status)
	}

	data, ok := resp.Data.Outputs.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%s", "数据异常")
	}
	var fdr FileDetectResult
	err = mapstructure.Decode(data, &fdr)
	if err != nil {
		return nil, fmt.Errorf("数据解析失败: %v", err)
	}
	return &fdr, nil
}
