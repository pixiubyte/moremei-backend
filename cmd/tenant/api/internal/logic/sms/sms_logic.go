package sms

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/alisms"
	commonConst "moremei/ai-saas/pkg/consts"
	tencentSms "moremei/ai-saas/pkg/sms"
	commonUtils "moremei/ai-saas/pkg/utils"

	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	codeKeyFmt         = "%s:%s:sms:%s:%s:code"
	countKeyFmt        = "%s:%s:sms:%s:%s:count"
	repeatKeyFmt       = "%s:%s:sms:%s:%s:repeat"
	verifyMobileScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1], KEYS[2])
else
    return 0
end`
)

type SmsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SmsLogic {
	return &SmsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SmsLogic) Send(method, mobile, clientIp string) error {
	codeKey := fmt.Sprintf(codeKeyFmt, l.svcCtx.Config.Name, "default", method, mobile)
	repeatKey := fmt.Sprintf(repeatKeyFmt, l.svcCtx.Config.Name, "default", method, mobile)
	countKey := fmt.Sprintf(countKeyFmt, l.svcCtx.Config.Name, "default", method, mobile)

	//发送前检测
	existsRepeatKey, err := l.svcCtx.RedisClient.ExistsCtx(l.ctx, repeatKey)
	if err != nil {
		return err
	} else if existsRepeatKey {
		return errs.SmsSendRepeatError
	}

	var count int
	existCountKey, err := l.svcCtx.RedisClient.Exists(countKey)
	if err != nil {
		return err
	}
	if existCountKey {
		countStr, err := l.svcCtx.RedisClient.GetCtx(l.ctx, countKey)
		if err != nil {
			return err
		}
		count, err = strconv.Atoi(countStr)
		if err != nil {
			return err
		}
		if count >= l.svcCtx.Config.Sms.ErrorLimit {
			return errs.SmsSendToManyTimesError
		}
	}

	//生成短信验证码
	var code string
	if l.svcCtx.Config.Mode == service.DevMode && len(l.svcCtx.Config.Sms.DevCode) == consts.SmsCodeLength {
		code = l.svcCtx.Config.Sms.DevCode
	} else {
		code = commonUtils.RandCode(consts.SmsCodeLength)
	}

	content, err := jsonx.MarshalToString([]string{code, strconv.Itoa(consts.SmsCodeExpireSeconds / 60)})
	if err != nil {
		return err
	}
	sms := tenant.Sms{
		PhoneNum: "+86" + mobile,
		BizType:  method,
		Content:  content,
		Status:   commonConst.StatusNew,
		ClientIp: clientIp,
	}
	err = l.svcCtx.SmsModel.Insert(l.ctx, &sms, l.svcCtx.DB)
	if err != nil {
		return err
	}

	currentEndTime, err := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02")+" 23:59:59", time.Local)
	if err != nil {
		return err
	}
	countKeyExpire := time.Duration(currentEndTime.Unix()-time.Now().Unix()) * time.Second

	//设置缓存数据
	err = l.svcCtx.RedisClient.PipelinedCtx(l.ctx, func(p redis.Pipeliner) error {
		if err := p.SetEX(l.ctx, codeKey, code, time.Second*consts.SmsCodeExpireSeconds).Err(); err != nil {
			return err
		}
		if err := p.SetEX(l.ctx, repeatKey, code, time.Second*consts.SmsCodeRepeatExpireSeconds).Err(); err != nil {
			return err
		}
		if count > 0 {
			return p.Incr(l.ctx, countKey).Err()
		} else {
			return p.SetEX(l.ctx, countKey, "1", countKeyExpire).Err()
		}
	})
	if err != nil {
		return err
	}

	if l.svcCtx.Config.Mode != service.DevMode {
		if err := l.sendProviderSms(method, mobile, code); err != nil {
			return err
		}
	}

	return nil
}

func (l *SmsLogic) sendProviderSms(method, mobile, code string) error {
	smsConfig, err := tenant.GetConfigByKey[tenant.DefaultSmsConfig](l.ctx, l.svcCtx.SystemConfigModel, tenant.SystemConfigKeyEihSmsConfig)
	if err != nil {
		return fmt.Errorf("get sms config: %w", err)
	}
	if smsConfig == nil {
		return errors.New("sms config is empty")
	}

	config, template, err := selectSmsConfig(smsConfig, method)
	if err != nil {
		return err
	}
	phoneNumber := "+86" + mobile
	expireMinutes := strconv.Itoa(consts.SmsCodeExpireSeconds / 60)

	if config.SecretId != "" && config.SecretKey != "" {
		endpoint := "sms.tencentcloudapi.com"
		if config.Region == "" {
			config.Region = "ap-nanjing"
		}
		client := tencentSms.MustNewClient(&tencentSms.Config{
			SecretId:  config.SecretId,
			SecretKey: config.SecretKey,
			AppId:     config.AppId,
			Endpoint:  endpoint,
			Region:    config.Region,
		})
		responses, err := client.Send(&tencentSms.Request{
			SignName:     template.SignName,
			TemplateId:   template.TemplateId,
			PhoneNumbers: []string{phoneNumber},
			Params:       []string{code, expireMinutes},
		})
		if err != nil {
			return err
		}
		response, ok := responses[phoneNumber]
		if !ok || !response.Result {
			return fmt.Errorf("send tencent sms failed: %s", response.Message)
		}
		return nil
	}

	if config.AccesskeyId != "" && config.AccesskeySecret != "" {
		region := config.Region
		if region == "" {
			region = "cn-hangzhou"
		}
		client := alisms.MustNewClient(&alisms.Config{
			AccessKeyId:     config.AccesskeyId,
			AccessKeySecret: config.AccesskeySecret,
			Region:          region,
		})
		responses, err := client.Send(&alisms.Request{
			SignName:     template.SignName,
			TemplateId:   template.TemplateId,
			PhoneNumbers: []string{mobile},
			Params: map[string]interface{}{
				"code":   code,
				"minute": expireMinutes,
			},
		})
		if err != nil {
			return err
		}
		response, ok := responses[mobile]
		if !ok || !response.Result {
			return fmt.Errorf("send ali sms failed: %s", response.Message)
		}
		return nil
	}

	return errors.New("sms provider credentials are empty")
}

func selectSmsConfig(defaultConfig *tenant.DefaultSmsConfig, method string) (*tenant.SmsConfigDetail, *tenant.DefaultBizType, error) {
	for _, config := range defaultConfig.Configs {
		if defaultConfig.DefaultConfig != "" && config.Name != defaultConfig.DefaultConfig {
			continue
		}
		if template := selectSmsTemplate(config.Config.BizType, method); template != nil {
			return &config.Config, template, nil
		}
	}
	for _, config := range defaultConfig.Configs {
		if template := selectSmsTemplate(config.Config.BizType, method); template != nil {
			return &config.Config, template, nil
		}
	}
	return nil, nil, fmt.Errorf("sms template not found for method %s", method)
}

func selectSmsTemplate(templates []tenant.DefaultBizType, method string) *tenant.DefaultBizType {
	for idx := range templates {
		for _, name := range templates[idx].Name {
			if name == method {
				return &templates[idx]
			}
		}
	}
	return nil
}

func (l *SmsLogic) Verify(mobile, code, method string) error {
	if l.svcCtx.Config.Mode == service.DevMode {
		return nil
	}
	// 放通测试账号
	for _, ta := range l.svcCtx.Config.TestAccounts {
		if ta.Mobile == mobile && ta.Code == code {
			return nil
		}
	}

	codeKey := fmt.Sprintf(codeKeyFmt, l.svcCtx.Config.Name, "default", method, mobile)
	countKey := fmt.Sprintf(countKeyFmt, l.svcCtx.Config.Name, "default", method, mobile)

	eval, err := l.svcCtx.RedisClient.Eval(verifyMobileScript, []string{codeKey, countKey}, []string{code})
	if err != nil {
		return err
	}

	reply, ok := eval.(int64)
	if !ok || reply != 2 {
		return errs.SmsCodeError
	}

	return nil
}
