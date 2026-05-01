package third_chat

import (
	"fmt"
	"time"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

// 实现基于手机号的第三方认证提供商
type Provider struct {
	name          string
	miniProgramId string
}

// 创建一个新的手机号认证提供商
func NewThirdChatProvider(name, miniProgramId string) goth.Provider {
	return &Provider{
		name:          name,
		miniProgramId: miniProgramId,
	}
}

// Name 实现goth.Provider接口
func (p *Provider) Name() string {
	return p.name
}

// SetName 实现goth.Provider接口
func (p *Provider) SetName(name string) {
	p.name = name
}

// BeginAuth 实现goth.Provider接口
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	session := &Session{
		AuthURL: state,
	}
	return session, nil
}

// FetchUser 实现goth.Provider接口
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	s := session.(*Session)
	return goth.User{
		Provider: p.Name(),
		UserID:   s.Mobile,
		Name:     fmt.Sprintf("用户%s", s.Mobile[len(s.Mobile)-4:]),
		NickName: fmt.Sprintf("用户%s", s.Mobile[len(s.Mobile)-4:]),
		RawData: map[string]interface{}{
			"unionid":         s.Mobile,
			"mini_program_id": p.miniProgramId,
			"register_time":   time.Now().Unix(),
		},
	}, nil
}

// Debug 实现goth.Provider接口
func (p *Provider) Debug(debug bool) {}

// RefreshToken 实现goth.Provider接口
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return nil, nil
}

// RefreshTokenAvailable 实现goth.Provider接口
func (p *Provider) RefreshTokenAvailable() bool {
	return false
}
