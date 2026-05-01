package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"moremei/ai-saas/internal/auth/third_chat"
	"moremei/ai-saas/internal/model/tenant"
	"strings"
	"sync"
	"time"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/wechat"
)

type UserInfo struct {
	Id       string
	Name     string
	NickName string
	UnionId  string
	Email    string
	Phone    string
	Avatar   string
	Extra    map[string]interface{}
}

type ThirdPartyAuthManager interface {
	GetProviderNames() []string
	GetProvider(name, app string) (Provider, error)
	ClearProviders()
}

type ResetProviderCallback func(*provider) error

// 初始化三方登录
func NewAuthManager(authModel tenant.AuthProviderConfModel) ThirdPartyAuthManager {
	// 初始化三方登录配置
	return &authManager{
		authModel: authModel,
	}
}

type authManager struct {
	lock sync.Mutex

	authModel tenant.AuthProviderConfModel
}

func (a *authManager) GetProviderNames() (names []string) {
	names = []string{}
	pds, err := a.authModel.FindAll(context.Background())
	if err != nil {
		return names
	}

	for _, p := range pds {
		names = append(names, p.Provider)
	}
	return names
}

// 清空所有的 providers. 可以在配置更新时使用
func (a *authManager) ClearProviders() {
	a.lock.Lock()
	defer a.lock.Unlock()

	goth.ClearProviders()
}

// 生成一个预创建的 provider, 用于三方登录
func (a *authManager) GetProvider(name, app string) (Provider, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	// 从数据库中获取 provider 的配置
	pconf, err := a.authModel.FindOneByProviderTypeApp(context.Background(), name, tenant.ConfigTypeAuth, app)
	if err != nil {
		return nil, fmt.Errorf("no provider for %s exists", name)
	}

	// 重置 provider 的 回调
	reset := func(p *provider) error {
		a.lock.Lock()
		defer a.lock.Unlock()
		// 重置 provider
		gp, err := a.newProvider(pconf.Provider, pconf.Config, pconf.App)
		if err != nil {
			return err
		}
		goth.UseProviders(gp)
		p.Provider = gp
		return nil
	}

	provider_name := genProviderName(app, name)

	gp, _ := goth.GetProvider(provider_name)
	if gp != nil {
		sess, _ := gp.UnmarshalSession("{}") // 生成一个临时session
		return &provider{
			Provider: gp,
			Session:  sess,
			config:   pconf.Config,
			reset:    reset,
		}, nil
	}

	gp, err = a.newProvider(pconf.Provider, pconf.Config, pconf.App)
	if err != nil {
		return nil, err
	}

	sess, _ := gp.UnmarshalSession("{}") // 生成一个临时session
	pr := &provider{
		Provider: gp,
		Session:  sess,
		config:   pconf.Config,
		reset:    reset,
	}
	// 注册 provider
	goth.UseProviders(gp)

	return pr, nil
}

func (a *authManager) newProvider(provider, conf, app string) (gp goth.Provider, err error) {

	// 生成一个唯一的 provider_name, 用于区分不同的租户和应用
	provider_name := genProviderName(app, provider)

	switch provider {
	case ProviderApple: // apple 登录

		config := &AppleConf{}
		err = json.Unmarshal([]byte(conf), config)
		if err != nil {
			return nil, err
		}
		secret, err := makeAppleSecret(config)
		if err != nil {
			return nil, err
		}
		gp = apple.New(config.ClientId, secret, config.RedirectUrl, nil)

	case ProviderWechat: // 微信登录

		config := &WechatConf{}
		err = json.Unmarshal([]byte(conf), config)
		if err != nil {
			return nil, err
		}
		gp = wechat.New(config.AppId, config.AppSecret, config.RedirectUrl, wechat.WECHAT_LANG_CN)

	case ProviderThirdChat:
		gp = third_chat.NewThirdChatProvider(provider_name, app)
	default:
		return nil, fmt.Errorf("no provider for %s exists", provider)
	}

	gp.SetName(provider_name)

	return gp, nil
}

func ParseProviderName(provider_name string) (tenantId, app, provider string) {
	parts := strings.Split(provider_name, ":")
	if len(parts) != 2 {
		return "", "", ""
	}
	return "", parts[0], parts[1]
}

// 生成一个唯一的 provider_name, 用于区分不同的应用
func genProviderName(app, provider string) string {
	return app + ":" + provider
}

func makeAppleSecret(apple_c *AppleConf) (string, error) {
	iat := time.Now().Unix()
	exp := iat + 60*60*24*30 // 30day
	sp := apple.SecretParams{
		ClientId:        apple_c.ClientId,
		TeamId:          apple_c.TeamId,
		KeyId:           apple_c.KeyId,
		PKCS8PrivateKey: apple_c.PrivateKey,
		Iat:             int(iat),
		Exp:             int(exp),
	}
	secret, err := apple.MakeSecret(sp)
	if err != nil {
		return "", err
	}
	return *secret, nil
}
