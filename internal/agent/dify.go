package agent

import (
	"context"
	"fmt"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/dify"

	"github.com/zeromicro/go-zero/core/jsonx"
)

func RunDifyWorkflow[T any](ctx context.Context, conf *tenant.AiAgentConfig, user string, input map[string]any) (*T, error) {
	client := dify.NewClient(conf.ServerHost, conf.ApiKey)
	resp, err := client.WorkflowsRun(ctx, dify.WorkflowRunRequest{
		Inputs: input,
		User:   user,
	})
	if err != nil {
		return nil, err
	}
	if resp.Data.Status != dify.StatusSucceeded {
		return nil, fmt.Errorf("RunDifyWorkflow [%s/%s]运行状态不成功：%s", resp.WorkflowRunId, resp.TaskID, resp.Data.Error)
	}
	// 通过序列化和反序列化方式解析输出数据为指定类型
	output, err := jsonx.MarshalToString(resp.Data.Outputs)
	if err != nil {
		return nil, fmt.Errorf("RunDifyWorkflow [%s/%s]输出数据解析失败：%s", resp.WorkflowRunId, resp.TaskID, err.Error())
	}
	var result T
	err = jsonx.UnmarshalFromString(output, &result)
	if err != nil {
		return nil, fmt.Errorf("RunDifyWorkflow [%s/%s]输出数据解析失败：%s", resp.WorkflowRunId, resp.TaskID, err.Error())
	}
	return &result, nil
}
